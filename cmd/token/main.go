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
package token

import (
	"github.com/spf13/cobra"

	"github.com/k0sproject/k0s/pkg/config"
)

type CmdOpts config.CLIOptions

var (
	kubeConfig  string
	tokenExpiry string
	tokenRole   string
	waitCreate  bool
)

func NewTokenCmd(c CmdOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage join tokens",
	}

	cmd.SilenceUsage = true
	cmd.AddCommand(tokenCreateCmd(c))
	cmd.AddCommand(tokenListCmd(c))
	cmd.AddCommand(tokenInvalidateCmd(c))
	cmd.Flags().AddFlagSet(c.Flagset)
	return cmd
}
