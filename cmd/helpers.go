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
package cmd

import (
	"github.com/k0sproject/k0s/pkg/apis/v1beta1"
	"github.com/k0sproject/k0s/pkg/config"
)

// ConfigFromYaml returns given k0s config or default config
func ConfigFromYaml(cfgPath string) (clusterConfig *v1beta1.ClusterConfig, err error) {
	cfg, err := config.ValidateYaml(cfgFile, k0sVars.DataDir)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
