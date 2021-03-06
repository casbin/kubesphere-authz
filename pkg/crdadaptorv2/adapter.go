package crdadaptorv2

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	model "github.com/casbin/casbin/v2/model"
	persist "github.com/casbin/casbin/v2/persist"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ClientType string

const (
	EXTERNAL_CLIENT ClientType = "external_client"
	INTERNAL_CLIENT ClientType = "internal_client"
)

var gvr = schema.GroupVersionResource{
	Group:    "auth.casbin.org",
	Version:  "v1",
	Resource: "universalpolicies",
}

//crd adapter specially for universal policy crd generated by kubebuilder
type K8sCRDAdaptor struct {
	namespace  string
	modelLabel string
	mode       ClientType
	clientset  dynamic.Interface
}

func NewK8sCRDAdaptor(namespace, modelLabel string, mode ClientType) (*K8sCRDAdaptor, error) {
	var res = K8sCRDAdaptor{
		namespace:  namespace,
		modelLabel: modelLabel,
		mode:       mode,
	}
	switch mode {
	case EXTERNAL_CLIENT:
		err := res.establishExternalClient()
		if err != nil {
			return nil, err
		}
	case INTERNAL_CLIENT:
		err := res.establishInternalClient()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("error: Invalid mode %s. mode should be either %s or %s", mode, EXTERNAL_CLIENT, INTERNAL_CLIENT)
	}
	return &res, nil

}

func (k *K8sCRDAdaptor) establishInternalClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *K8sCRDAdaptor) establishExternalClient() error {
	home := homedir.HomeDir()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *K8sCRDAdaptor) LoadPolicy(model model.Model) error {
	policyObjectList, err := k.getAllPolicyObjects()
	if err != nil {
		return err
	}
	for _, policyObject := range policyObjectList.Items {
		//multiple line of policies will not be accepted
		policy := removeStringAndLineBreaks(policyObject.Spec.PolicyItem)
		persist.LoadPolicyLine(policy, model)
	}
	return nil
}

func (k *K8sCRDAdaptor) SavePolicy(model model.Model) error {
	//before we save the policy, we must confirm which policies should be deleted.
	//create a set for all policies
	newPolicyMap := make(map[string]bool)
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			policy := policyToString(ptype, rule)
			newPolicyMap[policy] = false
		}
	}
	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			policy := policyToString(ptype, rule)
			newPolicyMap[policy] = false
		}
	}

	oldPolicyObjectList, err := k.getAllPolicyObjects()
	if err != nil {
		return err
	}
	for _, oldPolicyObject := range oldPolicyObjectList.Items {
		oldPolicy := removeStringAndLineBreaks(oldPolicyObject.Spec.PolicyItem)
		if _, ok := newPolicyMap[oldPolicy]; !ok {
			//this crd object should be removed
			if err := k.deletePolicyObject(oldPolicyObject.Name); err != nil {
				return err
			}
			newPolicyMap[oldPolicy] = true
		}
	}
	//then add into the missing data
	for newPolicy, exist := range newPolicyMap {
		if !exist {
			var name string
			err := k.insertPolicyObject(name, newPolicy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (k *K8sCRDAdaptor) UpdatePolicy(sec string, ptype string, oldRule, newPolicy []string) error {
	//FIXME: should ptype be put into string?
	oldRuleString := removeStringAndLineBreaks(strings.Join(append([]string{ptype}, oldRule...), ","))
	newRuleString := removeStringAndLineBreaks(strings.Join(append([]string{ptype}, newPolicy...), ","))

	//find out the old policy
	policyObjectList, err := k.getAllPolicyObjects()
	if err != nil {
		return err
	}
	for _, policyObject := range policyObjectList.Items {
		if removeStringAndLineBreaks(policyObject.Spec.PolicyItem) == oldRuleString {
			//this object is going to be modified
			err := k.updatePolicyObject(policyObject.Name, newRuleString)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *K8sCRDAdaptor) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error {
	if len(oldRules) != len(newRules) {
		return fmt.Errorf("Adaptor::UpdatePolicies: parameter oldRule and newRules don't have the same length")
	}
	for i, _ := range oldRules {
		err := k.UpdatePolicy(sec, ptype, oldRules[i], newRules[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *K8sCRDAdaptor) UpdateFilteredPolicies(sec string, ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return nil, errors.New("not implemented")
}

// AddPolicy adds a policy rule to the storage.
func (k *K8sCRDAdaptor) AddPolicy(sec string, ptype string, rule []string) error {
	//confirm that there is no same policy
	//find out the old policy
	newRuleString := removeStringAndLineBreaks(strings.Join(append([]string{ptype}, rule...), ","))

	policyObjectList, err := k.getAllPolicyObjects()
	if err != nil {
		return err
	}
	for _, policyObject := range policyObjectList.Items {
		if removeStringAndLineBreaks(policyObject.Spec.PolicyItem) == newRuleString {
			return nil
		}
	}
	err = k.insertPolicyObject(generatePolicyName(k.modelLabel, newRuleString), newRuleString)
	if err != nil {
		return err
	}
	return err
}

// AddPolicies adds policy rules to the storage.
func (k *K8sCRDAdaptor) AddPolicies(sec string, ptype string, rules [][]string) error {
	for i, _ := range rules {
		err := k.AddPolicy(sec, ptype, rules[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// RemovePolicy removes a policy rule from the storage.
func (k *K8sCRDAdaptor) RemovePolicy(sec string, ptype string, rule []string) error {
	ruleString := removeStringAndLineBreaks(strings.Join(append([]string{ptype}, rule...), ","))
	policyObjectList, err := k.getAllPolicyObjects()
	if err != nil {
		return err
	}
	for _, policyObject := range policyObjectList.Items {
		if removeStringAndLineBreaks(policyObject.Spec.PolicyItem) == ruleString {
			//this rule should be removed
			err := k.deletePolicyObject(policyObject.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RemovePolicies removes policy rules from the storage.
func (k *K8sCRDAdaptor) RemovePolicies(sec string, ptype string, rules [][]string) error {
	for i, _ := range rules {
		err := k.RemovePolicy(sec, ptype, rules[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (k *K8sCRDAdaptor) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
