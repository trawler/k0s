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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/k0sproject/k0s/cmd/api"
	"github.com/k0sproject/k0s/cmd/controller"
	"github.com/k0sproject/k0s/cmd/etcd"
	"github.com/k0sproject/k0s/cmd/install"
	"github.com/k0sproject/k0s/cmd/kubeconfig"
	"github.com/k0sproject/k0s/cmd/kubectl"
	"github.com/k0sproject/k0s/cmd/reset"
	"github.com/k0sproject/k0s/cmd/status"
	"github.com/k0sproject/k0s/cmd/token"
	"github.com/k0sproject/k0s/cmd/validate"
	"github.com/k0sproject/k0s/cmd/worker"

	"github.com/k0sproject/k0s/internal/util"

	"github.com/k0sproject/k0s/pkg/apis/v1beta1"
	"github.com/k0sproject/k0s/pkg/build"
	"github.com/k0sproject/k0s/pkg/config"
	"github.com/k0sproject/k0s/pkg/constant"
)

var (
	cfgFile       string
	cmdLogLevels  map[string]string
	dataDir       string
	debug         bool
	debugListenOn string
	k0sVars       constant.CfgVars
	logging       map[string]string
)

var defaultLogLevels = map[string]string{
	"etcd":                    "info",
	"containerd":              "info",
	"konnectivity-server":     "1",
	"kube-apiserver":          "1",
	"kube-controller-manager": "1",
	"kube-scheduler":          "1",
	"kubelet":                 "1",
	"kube-proxy":              "1",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "Data Directory for k0s (default: /var/lib/k0s). DO NOT CHANGE for an existing setup, things will break!")
	rootCmd.PersistentFlags().StringVar(&debugListenOn, "debugListenOn", ":6060", "Http listenOn for debug pprof handler")

	// Get relevant Vars from constant package
	k0sVars = constant.GetConfig(dataDir)

	opts := config.CLIOptions{
		CfgFile:          cfgFile,
		Debug:            debug,
		DefaultLogLevels: defaultLogLevels,
		Flagset:          getPersistentFlagSet(),
		K0sVars:          k0sVars,
	}

	rootCmd.AddCommand(api.NewApiCmd(api.CmdOpts(opts)))
	rootCmd.AddCommand(controller.NewControllerCmd(controller.CmdOpts(opts)))
	rootCmd.AddCommand(etcd.NewEtcdCmd(etcd.CmdOpts(opts)))
	rootCmd.AddCommand(install.NewInstallCmd(install.CmdOpts(opts)))
	rootCmd.AddCommand(token.NewTokenCmd(token.CmdOpts(opts)))
	rootCmd.AddCommand(worker.NewWorkerCmd(worker.CmdOpts(opts)))
	rootCmd.AddCommand(reset.NewResetCmd(reset.CmdOpts(opts)))
	rootCmd.AddCommand(status.NewStatusCmd(status.CmdOpts(opts)))
	rootCmd.AddCommand(validate.NewValidateCmd(validate.CmdOpts(opts)))
	rootCmd.AddCommand(kubeconfig.NewKubeConfigCmd(kubeconfig.CmdOpts(opts)))
	rootCmd.AddCommand(kubectl.NewK0sKubectlCmd(kubectl.CmdOpts(opts)))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(docs)
	rootCmd.AddCommand(completionCmd)

	// Add persistent Flags
	rootCmd.Flags().AddFlagSet(opts.Flagset)

	rootCmd.DisableAutoGenTag = true
	longDesc = "k0s - The zero friction Kubernetes - https://k0sproject.io"
	if build.EulaNotice != "" {
		longDesc = longDesc + "\n" + build.EulaNotice
	}
	rootCmd.Long = longDesc
}

var (
	longDesc string

	rootCmd = &cobra.Command{
		Use:   "k0s",
		Short: "k0s - Zero Friction Kubernetes",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set DEBUG from env, or from command flag
			if viper.GetString("debug") != "" || debug {
				logrus.SetLevel(logrus.DebugLevel)
				go func() {
					log.Println("starting debug server under", debugListenOn)
					log.Println(http.ListenAndServe(debugListenOn, nil))
				}()
			}

			// Set logging
			logging = util.MapMerge(cmdLogLevels, defaultLogLevels)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the k0s version",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(build.Version)
		},
	}

	docs = &cobra.Command{
		Use:   "docs",
		Short: "Generate Markdown docs for the k0s binary",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := generateDocs()
			if err != nil {
				return err
			}
			return nil
		},
	}
	configCmd = &cobra.Command{
		Use:   "default-config",
		Short: "Output the default k0s configuration yaml to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := buildConfig(); err != nil {
				return err
			}
			return nil
		},
	}

	completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(k0s completion bash)

# To load completions for each session, execute once:
  $ k0s completion bash > /etc/bash_completion.d/k0s

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ k0s completion zsh > "${fpath[1]}/_k0s"

# You will need to start a new shell for this setup to take effect.

Fish:

$ k0s completion fish | source

# To load completions for each session, execute once:
$ k0s completion fish > ~/.config/fish/completions/k0s.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
			return nil
		},
	}
)

func buildConfig() error {
	conf, _ := yaml.Marshal(v1beta1.DefaultClusterConfig())
	fmt.Print(string(conf))
	return nil
}

func generateDocs() error {
	if err := doc.GenMarkdownTree(rootCmd, "./docs/cli"); err != nil {
		return err
	}
	return nil
}

func getPersistentFlagSet() *pflag.FlagSet {
	flagset := &pflag.FlagSet{}
	flagset.StringVarP(&cfgFile, "config", "c", "", "config file (default: ./k0s.yaml)")
	flagset.BoolVarP(&debug, "debug", "d", false, "Debug logging (default: false)")
	return flagset
}

func Execute() {
	// just a hack to trick linter which requires to check for errors
	// cobra itself already prints out all errors that happen in subcommands
	_ = rootCmd.Execute()
}
