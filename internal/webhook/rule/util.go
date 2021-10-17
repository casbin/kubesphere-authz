package rule

import (
	"fmt"
	"ksauth/controllers"
	"ksauth/internal/config"
	"ksauth/pkg/crdadaptor"
	"ksauth/pkg/crdmodel"
	"strings"
)

func getModelAndPolicyObject(modelUrl, policyUrl string) (interface{}, interface{}, error) {
	modelObject, policyPlural, namespace, err := getModelObject(modelUrl)
	if err != nil {
		return nil, nil, err
	}
	if policyUrl == "" {
		if policyPlural == "" {
			//shouldn't reach here
			return nil, nil, fmt.Errorf("No policy specified or associated with model")
		}
		//should obtain policy adaptor associated with model
		adaptor, err := crdadaptor.NewK8sCRDAdaptorByYamlString(namespace, controllers.GeneratePolicyCrdDefinition(policyPlural), config.GetClientMode())
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
2nd return value is string of the namespaced plural form of the associated policy crd  if model is crd form
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
		model, policyPlural, err := crdmodel.GetModelFromCrdByYamlDefinition(yamlPath, namespace, modelName, config.GetClientMode())
		return model, policyPlural, namespace, err
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
		tmp := strings.Split(path, "#")
		if len(tmp) == 0 {
			return nil, fmt.Errorf("invalid syntax for crd url path. correct syntax: <policy crd name plural form>#<namespace>")
		}
		policyPlural := tmp[0]
		namespace := tmp[1]
		adaptor, err := crdadaptor.NewK8sCRDAdaptorByYamlString(namespace, controllers.GeneratePolicyCrdDefinition(policyPlural), config.GetClientMode())
		//adaptor, err := crdadaptor.NewK8sCRDAdaptorByYamlDefinition(namespace, yamlPath, config.GetClientMode())
		if err != nil {
			return nil, err
		}
		return adaptor, nil
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
