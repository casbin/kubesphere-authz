package rule

import (
	"fmt"
	"ksauth/internal/config"
	"ksauth/pkg/crdadaptor"
	"ksauth/pkg/crdmodel"
	"strings"
)

func getModelObject(url string) (interface{}, error) {
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
			return nil, fmt.Errorf("invalid syntax for crd url path. correct syntax: <yaml path to crd definition>#<namespace>")
		}
		yamlPath := "config/crd/bases/auth.casbin.org_casbinmodels.yaml"
		modelName := tmp[0]
		namespace := tmp[1]
		model, err := crdmodel.GetModelFromCrdByYamlDefinition(yamlPath, namespace, modelName, config.GetClientMode())
		return model, err
	}
	return nil, fmt.Errorf("invalid scheme %s", scheme)
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
			return nil, fmt.Errorf("invalid syntax for crd url path. correct syntax: <yaml path to crd definition>#<namespace>")
		}
		yamlPath := tmp[0]
		namespace := tmp[1]
		adaptor, err := crdadaptor.NewK8sCRDAdaptorByYamlDefinition(namespace, yamlPath, config.GetClientMode())
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
