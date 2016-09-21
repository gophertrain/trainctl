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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gophertrain/trainctl/templates"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about a module",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		module, err := getManifest(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(module)
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)

	infoCmd.PersistentFlags().String("name", "", "Module name")
}

func getManifest(cmd *cobra.Command) (templates.Module, error) {

	var m templates.Module
	name := cmd.Flag("name").Value.String() + ".json"
	path := getPath(cmd, "")

	file, err := ioutil.ReadFile(filepath.Join(path, name))
	if err != nil {
		return m, errors.Wrap(err, "reading manifest")
	}

	err = json.Unmarshal(file, &m)
	return m, err
}