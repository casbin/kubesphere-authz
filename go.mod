module ksauth

go 1.16

require (
	github.com/casbin/casbin/v2 v2.36.0
	github.com/gin-gonic/gin v1.7.3
	github.com/unrolled/secure v1.0.9
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.1
	k8s.io/apiextensions-apiserver v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
)

replace github.com/bouk/monkey v1.0.2 => bou.ke/monkey v1.0.2

replace bou.ke/monkey v1.0.2 => github.com/bouk/monkey v1.0.2
