package crdadaptor

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (k *K8sCRDAdaptor) getAllPolicyObjects() (PolicyList, error) {
	var emptyPolicyList PolicyList
	gvr := GetGroupVersionResource(k.group, k.version, k.policyNamePlural)
	unstructured, err := k.clientset.Resource(gvr).Namespace(k.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return emptyPolicyList, err
	}
	raw, err := unstructured.MarshalJSON()
	if err != nil {
		return emptyPolicyList, err
	}

	var policyObjectList PolicyList
	err = json.Unmarshal(raw, &policyObjectList)
	if err != nil {
		return emptyPolicyList, err
	}
	return policyObjectList, nil
}

func (k *K8sCRDAdaptor) deletePolicyObject(name string) error {
	gvr := GetGroupVersionResource(k.group, k.version, k.policyNamePlural)
	return k.clientset.Resource(gvr).Namespace(k.namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (k *K8sCRDAdaptor) insertPolicyObject(name string, policy string) error {
	gvr := GetGroupVersionResource(k.group, k.version, k.policyNamePlural)
	// oldUnstructured, err := k.clientset.Resource(gvr).Namespace(k.namespace).List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	return err
	// }

	var policyObject Policy
	policyObject.APIVersion = "auth.casbin.org/v1"
	policyObject.Kind = k.policyNameKind
	policyObject.SetName(name)
	policyObject.Spec.PolicyItem = policy
	//policyObject.SetResourceVersion(oldUnstructured.GetResourceVersion())

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
	gvr := GetGroupVersionResource(k.group, k.version, k.policyNamePlural)
	oldUnstructured, err := k.clientset.Resource(gvr).Namespace(k.namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	var policyObject Policy
	policyObject.APIVersion = "auth.casbin.org/v1"
	policyObject.Kind = k.policyNameKind
	policyObject.SetName(name)
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
