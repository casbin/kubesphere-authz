package config

import "ksauth/pkg/crdadaptor"

var DEBUG bool = true
var CLIENT_MODE crdadaptor.ClientType = crdadaptor.EXTERNAL_CLIENT
//var CLIENT_MODE crdadaptor.ClientType = crdadaptor.INTERNAL_CLIENT
var EXCLUDED_NAMESPACE=[]string{"kube-node-lease","kube-public","kube-system","policy"}

