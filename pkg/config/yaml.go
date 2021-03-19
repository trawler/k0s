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
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/k0sproject/k0s/pkg/apis/v1beta1"
	"github.com/sirupsen/logrus"
)

func GetYamlFromFile(cfgPath string, dataDir string) (clusterConfig *v1beta1.ClusterConfig, err error) {
	if cfgPath == "" {
		// no config file exists, using defaults
		logrus.Info("no config file given, using defaults")
	}
	cfg, err := ValidateYaml(cfgPath, dataDir)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func ValidateYaml(cfgPath string, dataDir string) (clusterConfig *v1beta1.ClusterConfig, err error) {
	if cfgPath == "" {
		// no config file exists, using defaults
		clusterConfig = v1beta1.DefaultClusterConfig()
	} else if isInputFromPipe() {
		clusterConfig, err = v1beta1.FromYamlPipe(os.Stdin)
	} else {
		clusterConfig, err = v1beta1.FromYamlFile(cfgPath)
	}
	if err != nil {
		return nil, err
	}

	if clusterConfig.Spec.Storage.Type == v1beta1.KineStorageType && clusterConfig.Spec.Storage.Kine == nil {
		clusterConfig.Spec.Storage.Kine = v1beta1.DefaultKineConfig(dataDir)
	}
	if clusterConfig.Spec.Install == nil {
		clusterConfig.Spec.Install = v1beta1.DefaultInstallSpec()
	}

	errors := clusterConfig.Validate()
	if len(errors) > 0 {
		messages := make([]string, len(errors))
		for _, e := range errors {
			messages = append(messages, e.Error())
		}
		return nil, fmt.Errorf(strings.Join(messages, "\n"))
	}
	return clusterConfig, nil
}

func isInputFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}
