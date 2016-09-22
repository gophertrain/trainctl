// Copyright © 2016 Brian Ketelsen <me@brianketelsen.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gophertrain/trainctl/templates"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search modules by metadata",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := checkSearchParams(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		mods, err := search(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, m := range mods {
			fmt.Println(m)
		}
		return
	},
}

func checkSearchParams(cmd *cobra.Command) error {

	name, err := cmd.PersistentFlags().GetString("name")
	if err != nil {
		return errors.Wrap(err, "Check parameters: name")
	}
	topic, err := cmd.PersistentFlags().GetString("topic")
	if err != nil {
		return errors.Wrap(err, "Check parameters: topic")
	}

	level, err := cmd.PersistentFlags().GetString("level")
	if err != nil {
		return errors.Wrap(err, "Check parameters: level")
	}
	var found bool
	if name != "" {
		found = true
	}
	if topic != "" {
		found = true
	}

	if level != "" {
		found = true
	}
	if !found {

		return errors.New("At least one search parameter required: name, topic, level")
	}
	return nil
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().String("name", "", "Module name")
	searchCmd.PersistentFlags().String("topic", "", "Module Topic {Go,Kubernetes}")
	searchCmd.PersistentFlags().String("level", "", "{beginner,intermediate,advanced,expert}")
	searchCmd.PersistentFlags().String("description", "", "Module description")
}

func search(cmd *cobra.Command) ([]templates.Module, error) {
	var results []templates.Module

	dir, err := os.Open(ProjectPath())
	if err != nil {
		return results, errors.Wrap(err, "open project directory")
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return results, errors.Wrap(err, "list project directory")
	}

	name, err := cmd.PersistentFlags().GetString("name")
	if err != nil {
		return results, errors.Wrap(err, "Check parameters: name")
	}
	topic, err := cmd.PersistentFlags().GetString("topic")
	if err != nil {
		return results, errors.Wrap(err, "Check parameters: topic")
	}

	level, err := cmd.PersistentFlags().GetString("level")
	if err != nil {
		return results, errors.Wrap(err, "Check parameters: level")
	}
	for _, f := range files {
		if f.IsDir() {
			var hit bool
			manifestPath := filepath.Join(ProjectPath(), f.Name(), f.Name()+".json")
			_, err := os.Stat(manifestPath)
			if err != nil {
				continue
			}

			module, err := getManifest(f.Name())
			if err != nil {
				return results, errors.Wrap(err, "get module")
			}
			if name != "" {
				if strings.Contains(module.ShortName, name) {
					hit = true
				}
			}

			if topic != "" {
				if strings.Contains(module.Topic, topic) {
					hit = true
				}
			}
			if level != "" {
				if strings.Contains(string(module.Level), level) {
					hit = true
				}
			}

			if hit {
				results = append(results, module)
			}

		}
	}

	return results, nil

}
