package main

import (
	"ksauth/pkg/crdadaptor"

	casbin "github.com/casbin/casbin/v2"
)

func main() {
	adaptor, err := crdadaptor.NewK8sCRDAdaptorByYamlDefinition("policy", "crd_example.yaml", crdadaptor.EXTERNAL_CLIENT)
	if err != nil {
		panic(err.Error())
	}

	enforcer, err := casbin.NewEnforcer("external_ip.conf", adaptor)
	if err != nil {
		panic(err.Error())
	}
	ok, err := enforcer.AddPolicy("default", "1.1.1.1", "allow")
	if err != nil {
		panic(err.Error())
	}
	ok, err = enforcer.AddPolicy("default", "1.1.1.2", "allow")
	if err != nil {
		panic(err.Error())
	}
	ok, err = enforcer.AddPolicy("default", "1.1.1.3", "allow")
	if err != nil {
		panic(err.Error())
	}

	ok = enforcer.HasPolicy("default", "1.1.1.1", "allow")
	if !ok {
		panic("policy not found")
	}
	ok = enforcer.HasPolicy("default", "1.1.1.2", "allow")
	if !ok {
		panic("policy not found")
	}
	ok = enforcer.HasPolicy("default", "1.1.1.3", "allow")
	if !ok {
		panic("policy not found")
	}
	ok, err = enforcer.RemovePolicy("default", "1.1.1.2", "allow")
	if err != nil {
		panic(err.Error())
	}
	ok = enforcer.HasPolicy("default", "1.1.1.2", "allow")
	if ok {
		panic("policy should have been deleted")
	}
	ok, err = enforcer.UpdatePolicy([]string{"default", "1.1.1.3", "allow"}, []string{"default", "1.1.1.4", "allow"})
	if err != nil {
		panic(err.Error())
	}
	ok = enforcer.HasPolicy("default", "1.1.1.3", "allow")
	if ok {
		panic("policy should have been modified")
	}
	ok = enforcer.HasPolicy("default", "1.1.1.4", "allow")
	if !ok {
		panic("policy not found")
	}

}
