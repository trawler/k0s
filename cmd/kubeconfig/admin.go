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
package kubeconfig

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/k0sproject/k0s/internal/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func kubeConfigAdminCmd(c CmdOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin [command]",
		Short: "Display Admin's Kubeconfig file",
		Long:  "Print kubeconfig for the Admin user to stdout",
		Example: `	$ k0s kubeconfig admin > ~/.kube/config
	$ export KUBECONFIG=~/.kube/config
	$ kubectl get nodes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if util.FileExists(c.K0sVars.AdminKubeConfigPath) {
				content, err := ioutil.ReadFile(c.K0sVars.AdminKubeConfigPath)
				if err != nil {
					log.Fatal(err)
				}

				clusterAPIURL, err := c.getAPIURL()
				if err != nil {
					return errors.Wrap(err, "failed to fetch cluster's API Address: %v.")
				}
				newContent := strings.Replace(string(content), "https://localhost:6443", clusterAPIURL, -1)
				os.Stdout.Write([]byte(newContent))
			} else {
				return errors.Errorf("failed to read admin config, is the control plane initialized on this node?")
			}
			return nil
		},
	}
	return cmd
}
