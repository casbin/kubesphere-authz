package controllers

import "testing"

func TestCapitalize(t *testing.T) {

	testArray := [][]string{
		{"test1", "Test1"},
		{"", ""},
		{"Test", "Test"},
		{"123avb", "123avb"},
	}
	for _, tuple := range testArray {
		res := capitalize(tuple[0])
		if res != tuple[1] {
			t.Errorf("expected %s, got %s", tuple[1], res)
		}
	}

}

func TestGeneratePolicyCrdDefinition(t *testing.T) {
	var expected = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: allowed-repo.auth.casbin.org
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
    plural: allowed-repo
    # singular name to be used as an alias on the CLI and for display
    singular: allowed-repo
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: AllowedRepo
`
	res := GeneratePolicyCrdDefinition("allowed-repo")
	if res != expected {
		t.Errorf("incorrect answer generated")
	}

}
