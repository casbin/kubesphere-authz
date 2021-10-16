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
)


func GetModelFromCrd(group, version,namespace,modelResourcePlural,modelName string,clientType crdadaptor.ClientType )(model.Model,error){
	//establish client
	var client dynamic.Interface
	var err error
	switch clientType{
	case crdadaptor.EXTERNAL_CLIENT:
		client,err=establishExternalClient()
	case crdadaptor.INTERNAL_CLIENT:
		client,err=establishInternalClient()
	}
	if err!=nil{
		return nil,err
	}
	modelText,err:=GetModelTextFromCrd(group,version,namespace,modelResourcePlural,modelName,client)
	if err!=nil{
		return nil,err
	}
	var res model.Model=model.NewModel()
	err=res.LoadModelFromText(modelText)
	return res,err
}


func GetModelFromCrdByYamlDefinition(yamlDefinitionPath string,namespace string,modelName string,clientType crdadaptor.ClientType)(model.Model,error){
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
	return GetModelFromCrd(definition.Spec.Group,definition.Spec.Versions[0].Name,namespace,definition.Spec.Names.Plural,modelName,clientType)
}



func GetModelTextFromCrd(group, version,namespace,modelResourcePlural,modelName string,client dynamic.Interface)(string,error){
	var gvr=schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: modelResourcePlural,
	}

	unstructured, err :=client.Resource(gvr).Namespace(namespace).Get(context.TODO(),modelName,metav1.GetOptions{})
	if err!=nil{
		return "",err
	}

	raw, err := unstructured.MarshalJSON()
	if err!=nil{
		return "",err
	}
	var crdModel CrdModel
	err=json.Unmarshal(raw,&crdModel)
	if err!=nil{
		return "",err
	}

	return crdModel.Spec.ModelText,nil
	

}