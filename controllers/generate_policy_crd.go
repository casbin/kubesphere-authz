package controllers

import "strings"

const template string = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: ${plural}.auth.casbin.org
  namespace: kubesphere-authz-system
spec:
  # group name to use for REST API: /apis/<group>/<version>
  group: auth.casbin.org
  # list of versions supported by this CustomResourceDefinition
  versions:
    - name: v1
      # Each version can be enabled/disabled by Served flag.
      served: true
      # One and only one version must be marked as the storage version.
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                policyItem:
                  #this policyItem contains a line of policy. multiple lines of policy is forbiddened
                  type: string
                
  # either Namespaced or Cluster
  scope: Namespaced
  names:
    # plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: ${plural}
    # singular name to be used as an alias on the CLI and for display
    singular: ${singular}
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: ${kind}
`

func GeneratePolicyCrdDefinition(plural string) string {
	split := strings.Split(plural, "-")
	for i, s := range split {
		split[i] = capitalize(s)
	}
	kind := strings.Join(split, "")

	var crd = template
	crd = strings.Replace(crd, "${plural}", plural, -1)
	crd = strings.Replace(crd, "${singular}", plural, -1)
	crd = strings.Replace(crd, "${kind}", kind, -1)
	return crd

}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	var tmp = []byte(s)
	if tmp[0] >= 'a' && tmp[0] <= 'z' {
		tmp[0] = tmp[0] - 'a' + 'A'
	}
	return string(tmp)
}
