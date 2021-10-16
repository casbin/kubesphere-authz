package crdmodel
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type CrdModelSpec struct {
	ModelText string `json:"modelText"`
}

type CrdModel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CrdModelSpec `json:"spec"`
}
