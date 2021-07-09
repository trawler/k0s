package v1beta1

import "github.com/k0sproject/k0s/pkg/constant"

// SystemUser defines the user to use for each component
type SystemUser struct {
	Etcd          string `json:"etcdUser,omitempty" yaml:"etcdUser,omitempty"`
	Kine          string `json:"kineUser,omitempty" yaml:"kineUser,omitempty"`
	Konnectivity  string `json:"konnectivityUser,omitempty" yaml:"konnectivityUser,omitempty"`
	KubeAPIServer string `json:"kubeAPIserverUser,omitempty" yaml:"kubeAPIserverUser,omitempty"`
	KubeScheduler string `json:"kubeSchedulerUser,omitempty" yaml:"kubeSchedulerUser,omitempty"`
}

// DefaultSystemUsers returns the default system users to be used for the different components
func DefaultSystemUsers() *SystemUser {
	return &SystemUser{
		Etcd:          constant.EtcdUser,
		Kine:          constant.KineUser,
		Konnectivity:  constant.KonnectivityServerUser,
		KubeAPIServer: constant.ApiserverUser,
		KubeScheduler: constant.SchedulerUser,
	}
}

// DefaultInstallSpec ...
func DefaultInstallSpec() *InstallSpec {
	return &InstallSpec{
		SystemUsers: DefaultSystemUsers(),
	}
}
