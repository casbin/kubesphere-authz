package main

import (
	"fmt"
	"ksauth/pkg/crdadaptor"
	"ksauth/pkg/crdmodel"
)
func main(){
	res,err:=crdmodel.GetModelFromCrdByYamlDefinition("pkg/crdmodel/model.yaml","policy","allowed-repos",crdadaptor.EXTERNAL_CLIENT)
	if err!=nil{
		panic(err.Error())
	}
	fmt.Println(res.ToText())
}