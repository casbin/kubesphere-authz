package crdadaptorv2

import (
	"context"
	"encoding/json"
	api "ksauth/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

func (k *K8sCRDAdaptor) getAllPolicyObjects() (api.UniversalPolicyList, error) {
	var policyObjectList api.UniversalPolicyList
	var labelMap = make(map[string]string)
	labelMap["model"] = k.modelLabel
	unstructured, err := k.clientset.Resource(gvr).Namespace(k.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.SelectorFromSet(labelMap).String()})
	if err != nil {
		return policyObjectList, err
	}
	raw, err := unstructured.MarshalJSON()
	if err != nil {
		return policyObjectList, err
	}

	err = json.Unmarshal(raw, &policyObjectList)
	if err != nil {
		return policyObjectList, err
	}
	return policyObjectList, nil
}

func (k *K8sCRDAdaptor) deletePolicyObject(name string) error {
	return k.clientset.Resource(gvr).Namespace(k.namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (k *K8sCRDAdaptor) insertPolicyObject(name string, policy string) error {
	var policyObject api.UniversalPolicy
	policyObject.Spec.PolicyItem = policy
	policyObject.SetName(name)
	policyObject.APIVersion = "auth.casbin.org/v1"
	policyObject.Kind = "UniversalPolicy"
	policyObject.ObjectMeta.Labels = make(map[string]string)
	policyObject.ObjectMeta.Labels["model"] = k.modelLabel

	raw, err := json.Marshal(policyObject)
	if err != nil {
		return err
	}
	var unstructured unstructured.Unstructured
	err = unstructured.UnmarshalJSON(raw)
	if err != nil {
		return err
	}

	_, err = k.clientset.Resource(gvr).Namespace(k.namespace).Create(context.TODO(), &unstructured, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (k *K8sCRDAdaptor) updatePolicyObject(name string, policy string) error {
	var labelMap = make(map[string]string)
	labelMap["model"] = k.modelLabel
	oldUnstructured, err := k.clientset.Resource(gvr).Namespace(k.namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	var policyObject api.UniversalPolicy
	policyObject.SetName(name)
	policyObject.APIVersion = "auth.casbin.org/v1"
	policyObject.Kind = "UniversalPolicy"
	policyObject.Spec.PolicyItem = policy
	policyObject.SetResourceVersion(oldUnstructured.GetResourceVersion())
	raw, err := json.Marshal(policyObject)
	if err != nil {
		return err
	}
	var unstructured unstructured.Unstructured
	err = unstructured.UnmarshalJSON(raw)
	if err != nil {
		return err
	}

	_, err = k.clientset.Resource(gvr).Namespace(k.namespace).Update(context.TODO(), &unstructured, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
