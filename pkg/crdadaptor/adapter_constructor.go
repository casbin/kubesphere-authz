package crdadaptor

import (
	"fmt"
	"io/ioutil"

	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewK8sCRDAdaptor(group, version, namespace, policyNameKind, policyNamePlural string, mode ClientType) (*K8sCRDAdaptor, error) {
	var res = K8sCRDAdaptor{
		group:            group,
		version:          version,
		namespace:        namespace,
		policyNameKind:   policyNameKind,
		policyNamePlural: policyNamePlural,
		mode:             mode,
	}
	switch mode {
	case EXTERNAL_CLIENT:
		err := res.establishExternalClient()
		if err != nil {
			return nil, err
		}
	case INTERNAL_CLIENT:
		err := res.establishInternalClient()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("error: Invalid mode %s. mode should be either %s or %s", mode, EXTERNAL_CLIENT, INTERNAL_CLIENT)
	}
	return &res, nil
}

//warning: if multiple versions are specified in yaml definition, only the 1st element will be used.
func NewK8sCRDAdaptorByYamlDefinition(namespace string, yamlDefinitionPath string, mode ClientType) (*K8sCRDAdaptor, error) {
	var definition apiextensions.CustomResourceDefinition
	fileData, err := ioutil.ReadFile(yamlDefinitionPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(fileData, &definition)
	if err != nil {
		return nil, err
	}

	if len(definition.Spec.Versions) == 0 {
		return nil, fmt.Errorf("no versions information provided")
	}
	//TODO: remove the hard code index 0
	return NewK8sCRDAdaptor(
		definition.Spec.Group,
		definition.Spec.Versions[0].Name,
		namespace,
		definition.Spec.Names.Kind,
		definition.Spec.Names.Plural,
		mode)
}

func (k *K8sCRDAdaptor) establishInternalClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *K8sCRDAdaptor) establishExternalClient() error {
	home := homedir.HomeDir()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset
	return nil
}
