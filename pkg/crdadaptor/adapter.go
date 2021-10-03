package crdadaptor

import (
	"errors"
	"fmt"
	"strings"

	model "github.com/casbin/casbin/v2/model"
	persist "github.com/casbin/casbin/v2/persist"
	"k8s.io/client-go/dynamic"
	//"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ClientType string

const (
	EXTERNAL_CLIENT ClientType = "external_client"
	INTERNAL_CLIENT ClientType = "internal_client"
)

type K8sCRDAdaptor struct {
	group            string
	version          string
	namespace        string     //in which namespace these name are stored
	policyNameKind   string     //the resource name of this policy
	policyNamePlural string     //the plural form of the resource name of this policy
	mode             ClientType //EXTERNAL_CLIENT if this code is run outside the k8s, INTERNAL_CLIENT if else
	clientset        dynamic.Interface
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
	err = k.insertPolicyObject(generatePolicyName(newRuleString), newRuleString)
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
