package server

import (
	"fmt"
	"sync/atomic"

	"github.com/sirupsen/logrus"

	config "github.com/k0sproject/k0s/pkg/apis/v1beta1"
	k0sv1beta1 "github.com/k0sproject/k0s/pkg/apis/v1beta1"
	"github.com/k0sproject/k0s/pkg/constant"

	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/tools/clientcmd"
)

// ClusterConfigCustomResource is the component that manages the cluster config resource in the API
// ClusterConfigCustomResource is the component that manages the cluster config resource in the API
type ClusterConfigCustomResource struct {
	ClusterConfig *config.ClusterConfig
	Logger        *logrus.Entry

	clientSet     *apiextension.Clientset
	leaderElector LeaderElector
	stopCh        chan struct{}
	kubeConfig    string
}

// NewEndpointReconciler creates new endpoint reconciler
func NewClusterConfigCR(c *k0sv1beta1.ClusterConfig, leaderElector LeaderElector, k0sVars constant.CfgVars) *ClusterConfigCustomResource {
	d := atomic.Value{}
	d.Store(true)
	return &ClusterConfigCustomResource{
		ClusterConfig: c,
		leaderElector: leaderElector,
		stopCh:        make(chan struct{}),
		kubeConfig:    k0sVars.AdminKubeConfigPath,
		Logger:        logrus.WithFields(logrus.Fields{"component": "clusterConfigCR"}),
	}
}

// Init initializes the Config CustomResource
func (c *ClusterConfigCustomResource) Init() error {
	config, err := clientcmd.BuildConfigFromFlags("", c.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to build config for CR operations: %v", err)
	}

	apixClient, err := apiextension.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create API typed Client for CR operations: %v", err)
	}
	c.clientSet = apixClient

	err = k0sv1beta1.CreateConfigCRD(apixClient, *c.ClusterConfig)
	if err != nil {
		return fmt.Errorf("failed to create custom resource of type ClusterConfig: %v", err)
	}
	return nil
}

// Run
func (c *ClusterConfigCustomResource) Run() error {
	return nil
}

// Stop
func (c *ClusterConfigCustomResource) Stop() error {
	close(c.stopCh)
	return nil
}

// Healthy
func (c *ClusterConfigCustomResource) Healthy() error {
	return nil
}
