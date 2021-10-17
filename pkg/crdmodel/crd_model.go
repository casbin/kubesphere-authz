package crdmodel

import (
	"context"
	"fmt"
	"io/ioutil"
	"ksauth/pkg/crdadaptor"

	"encoding/json"

	model "github.com/casbin/casbin/v2/model"
	"gopkg.in/yaml.v2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	crd "ksauth/api/v1"
)

/**
2nd return value is string of the plural form of the associated policy crd
*/
func GetModelFromCrd(group, version, namespace, modelResourcePlural, modelName string, clientType crdadaptor.ClientType) (model.Model, string, error) {
	//establish client
	var client dynamic.Interface
	var err error
	switch clientType {
	case crdadaptor.EXTERNAL_CLIENT:
		client, err = establishExternalClient()
	case crdadaptor.INTERNAL_CLIENT:
		client, err = establishInternalClient()
	}
	if err != nil {
		return nil, "", err
	}
	modelText, policyCRDPlural, err := GetModelTextFromCrd(group, version, namespace, modelResourcePlural, modelName, client)
	if err != nil {
		return nil, "", err
	}
	var res model.Model = model.NewModel()
	err = res.LoadModelFromText(modelText)
	return res, policyCRDPlural, err
}

/**
2nd return value is string of the plural form of the associated policy crd
*/
func GetModelFromCrdByYamlDefinition(yamlDefinitionPath string, namespace string, modelName string, clientType crdadaptor.ClientType) (model.Model, string, error) {
	var definition apiextensions.CustomResourceDefinition
	fileData, err := ioutil.ReadFile(yamlDefinitionPath)
	if err != nil {
		return nil, "", err
	}
	err = yaml.Unmarshal(fileData, &definition)
	if err != nil {
		return nil, "", err
	}

	if len(definition.Spec.Versions) == 0 {
		return nil, "", fmt.Errorf("no versions information provided")
	}
	return GetModelFromCrd(definition.Spec.Group, definition.Spec.Versions[0].Name, namespace, definition.Spec.Names.Plural, modelName, clientType)
}

/**
2nd return value is string of the plural form of the associated policy crd
*/
func GetModelTextFromCrd(group, version, namespace, modelResourcePlural, modelName string, client dynamic.Interface) (string, string, error) {
	var gvr = schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: modelResourcePlural,
	}

	unstructured, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), modelName, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	raw, err := unstructured.MarshalJSON()
	if err != nil {
		return "", "", err
	}
	var crdModel crd.CasbinModel
	err = json.Unmarshal(raw, &crdModel)
	if err != nil {
		return "", "", err
	}

	return crdModel.Spec.ModelText, crdModel.Spec.AssociatedPolicyCrdPlural, nil

}
