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
package etcd

import (
	"context"
	"fmt"

	"github.com/k0sproject/k0s/pkg/constant"
	"github.com/k0sproject/k0s/pkg/etcd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func etcdListCmd(k0sVars constant.CfgVars) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member-list",
		Short: "Returns etcd cluster members list",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			etcdClient, err := etcd.NewClient(k0sVars.CertRootDir, k0sVars.EtcdCertDir)
			if err != nil {
				return fmt.Errorf("can't list etcd cluster members: %v", err)
			}
			members, err := etcdClient.ListMembers(ctx)
			if err != nil {
				return fmt.Errorf("can't list etcd cluster members: %v", err)
			}
			l := logrus.New()
			l.SetFormatter(&logrus.JSONFormatter{})

			l.WithField("members", members).
				Info("done")
			return nil
		},
	}
	return cmd
}
