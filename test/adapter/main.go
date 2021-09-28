package main

import (
	casbin "github.com/casbin/casbin/v2"
	"ksauth/pkg/crdadapter"
)

func main() {
	adapter, err := crdadapter.NewK8sCRDAdapterByYamlDefinition("policy", "crd_example.yaml", crdadapter.EXTERNAL_CLIENT)
	if err != nil {
		panic(err.Error())
	}

	enforcer, err := casbin.NewEnforcer("external_ip.conf", adapter)
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
