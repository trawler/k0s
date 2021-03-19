/*
Copyright 2021 k0s authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package validate

import (
	"fmt"

	"github.com/k0sproject/k0s/pkg/config"
	"github.com/spf13/cobra"
)

type CmdOpts config.CLIOptions

func NewValidateCmd(c CmdOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Helper command for validating the config file",
	}
	cmd.AddCommand(validateConfigCmd(c))
	return cmd
}

func validateConfigCmd(c CmdOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Helper command for validating the config file",
		Long: `Example:
   k0s validate config --config path_to_config.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.GetYamlFromFile(c.CfgFile, c.K0sVars.DataDir)
			if err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}

	// append flags
	cmd.Flags().AddFlagSet(c.Flagset)
	return cmd
}
