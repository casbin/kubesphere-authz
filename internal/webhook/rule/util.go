package rule

import (
	"fmt"
	//"ksauth/controllers"
	"ksauth/internal/config"
	crdadaptor "ksauth/pkg/crdadaptorv2"
	"ksauth/pkg/crdmodel"
	"strings"
)

func getModelAndPolicyObject(modelUrl, policyUrl string) (interface{}, interface{}, error) {
	modelObject, modelName, namespace, err := getModelObject(modelUrl)
	if err != nil {
		return nil, nil, err
	}
	if policyUrl == "" {
		//should obtain universal policy adaptor 
		adaptor, err := crdadaptor.NewK8sCRDAdaptor(namespace, modelName, config.GetClientMode())
		if err != nil {
			return nil, nil, err
		}
		return modelObject, adaptor, nil
	} else {
		adaptor, err := getAdaptorObject(policyUrl)
		if err != nil {
			return nil, nil, err
		}
		return modelObject, adaptor, nil
	}

}

/**
2nd return value is model name
3rd return value is k8s namespace (if have)
*/
func getModelObject(url string) (interface{}, string, string, error) {
	scheme, path, err := splitSchemeAndPath(url)
	if err != nil {
		return nil, "", "", err
	}
	switch scheme {
	case "file":
		return path, "", "", nil
	case "crd":
		tmp := strings.Split(path, "#")
		if len(tmp) == 0 {
			return nil, "", "", fmt.Errorf("invalid syntax for crd url path. correct syntax: <yaml path to crd definition>#<namespace>")
		}
		yamlPath := "config/crd/bases/auth.casbin.org_casbinmodels.yaml"
		modelName := tmp[0]
		namespace := tmp[1]
		model, _, err := crdmodel.GetModelFromCrdByYamlDefinition(yamlPath, namespace, modelName, config.GetClientMode())
		return model, modelName, namespace, err
	}
	return nil, "", "", fmt.Errorf("invalid scheme %s", scheme)
}

func getAdaptorObject(url string) (interface{}, error) {
	scheme, path, err := splitSchemeAndPath(url)
	if err != nil {
		return nil, err
	}
	switch scheme {
	case "file":
		return path, nil
	case "crd":
		// crd adaptor v1 is no longer supported
	}
	return nil, fmt.Errorf("invalid scheme %s", scheme)
}

func splitSchemeAndPath(url string) (scheme, path string, e error) {
	tmp := strings.Split(url, "://")
	scheme = ""
	path = ""
	if len(tmp) != 2 {
		e = fmt.Errorf("invalid url %s", url)
		return
	}
	return tmp[0], tmp[1], nil
}
