package crdmodel

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)


func  establishInternalClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil,err
	}
	// creates the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil,err
	}
	return clientset,nil
}
var kubeconfig *string=nil

func establishExternalClient() (dynamic.Interface, error) {
	if kubeconfig==nil{
		if home := homedir.HomeDir();home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
	}
	
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil,err
	}

	// create the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil,err
	}
	return clientset,nil
}

