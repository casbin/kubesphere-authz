package crdmodel

import (
	"context"
	"fmt"
	"io/ioutil"
	crdadaptor "ksauth/pkg/crdadaptorv2"

	"encoding/json"

	crd "ksauth/api/v1"

	model "github.com/casbin/casbin/v2/model"
	"gopkg.in/yaml.v2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

/**
1st return value is the object of model, which is a stirng if the model is stored in a file, or a casbin.Model object if model is stored in other ways
if the model is correctly got but this model is not enabled, the 1st and the 3rd value will be nil together.
2nd enabled
*/
func GetModelFromCrd(group, version, namespace, modelResourcePlural, modelName string, clientType crdadaptor.ClientType) (model.Model, bool, error) {
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
		return nil, false, err
	}
	modelText, enabled, err := getModelTextFromCrd(group, version, namespace, modelResourcePlural, modelName, client)
	if err != nil {
		return nil, false, err
	}
	var res model.Model = model.NewModel()

	if enabled {
		err = res.LoadModelFromText(modelText)
		return res, true, err
	} else {
		return nil, false, nil
	}

}

/**
2nd return value is string of the plural form of the associated policy crd
*/
func GetModelFromCrdByYamlDefinition(yamlDefinitionPath string, namespace string, modelName string, clientType crdadaptor.ClientType) (model.Model, bool, error) {
	var definition apiextensions.CustomResourceDefinition
	fileData, err := ioutil.ReadFile(yamlDefinitionPath)
	if err != nil {
		return nil, false, err
	}
	err = yaml.Unmarshal(fileData, &definition)
	if err != nil {
		return nil, false, err
	}

	if len(definition.Spec.Versions) == 0 {
		return nil, false, fmt.Errorf("no versions information provided")
	}
	return GetModelFromCrd(definition.Spec.Group, definition.Spec.Versions[0].Name, namespace, definition.Spec.Names.Plural, modelName, clientType)
}

/**
1st value is the model text
2nd is boolean marking whether this model is enabled
*/
func getModelTextFromCrd(group, version, namespace, modelResourcePlural, modelName string, client dynamic.Interface) (string, bool, error) {
	var gvr = schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: modelResourcePlural,
	}

	unstructured, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), modelName, metav1.GetOptions{})
	if err != nil {
		return "", false, err
	}

	raw, err := unstructured.MarshalJSON()
	if err != nil {
		return "", false, err
	}
	var crdModel crd.CasbinModel
	err = json.Unmarshal(raw, &crdModel)
	if err != nil {
		return "", false, err
	}

	return crdModel.Spec.ModelText, crdModel.Spec.Enabled, nil

}
