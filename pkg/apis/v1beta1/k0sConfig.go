package v1beta1

import (
	"context"
	"reflect"

	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CRDPlural   string = "configs"
	CRDGroup    string = "k0sproject.io"
	CRDVersion  string = "v1beta1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

// CreateConfigCRD creates the clusterConfig CRD in the cluster
func CreateConfigCRD(clientset apiextension.Interface, config ClusterConfig) error {
	version := apiextensionv1beta1.CustomResourceDefinitionVersion{Name: CRDVersion}
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{Name: FullCRDName},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group:    CRDGroup,
			Versions: []apiextensionv1beta1.CustomResourceDefinitionVersion{version},
			Scope:    apiextensionv1beta1.NamespaceScoped,
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Plural: CRDPlural,
				Kind:   reflect.TypeOf(config).Name(),
			},
		},
	}

	opts := v1.CreateOptions{}
	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(context.TODO(), crd, opts)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
