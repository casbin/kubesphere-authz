package rule

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/admission/v1"
	core "k8s.io/api/core/v1"
)

func TestRequiredAnnotations(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	podObject.Spec.Containers = append(podObject.Spec.Containers, container)
	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "pods"
	data, err := json.Marshal(podObject)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	res := rule.RequiredAnnotations(&review, "../../../example/required_annotations/required_annotations.conf", "file://../../../example/required_annotations/required_annotations.csv")
	if res == nil {
		t.Error("should be rejected")
		return
	}
}

func TestRequiredAnnotations2(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	podObject.Spec.Containers = append(podObject.Spec.Containers, container)

	podObject.SetAnnotations(map[string]string{
		"a8r.io/owner": "test-100",
	})
	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "pods"
	data, err := json.Marshal(podObject)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	res := rule.RequiredAnnotations(&review, "../../../example/required_annotations/required_annotations.conf", "file://../../../example/required_annotations/required_annotations.csv")
	if res != nil {
		t.Error("should not be rejected")
		return
	}
}

func TestRequiredAnnotations3(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	podObject.Spec.Containers = append(podObject.Spec.Containers, container)

	podObject.SetAnnotations(map[string]string{
		"a8r.io/owner": "tdddest-100",
	})
	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "pods"
	data, err := json.Marshal(podObject)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	res := rule.RequiredAnnotations(&review, "../../../example/required_annotations/required_annotations.conf", "file://../../../example/required_annotations/required_annotations.csv")
	if res == nil {
		t.Error("should  be rejected")
		return
	}
}

func TestRequiredAnnotations4(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	podObject.Spec.Containers = append(podObject.Spec.Containers, container)

	podObject.SetAnnotations(map[string]string{
		"a8r.io/owner":  "test-100",
		"a8r.dio/owner": "tdddest-100",
	})
	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "pods"
	data, err := json.Marshal(podObject)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	res := rule.RequiredAnnotations(&review, "../../../example/required_annotations/required_annotations.conf", "file://../../../example/required_annotations/required_annotations.csv")
	if res != nil {
		t.Error("should not be rejected")
		return
	}
}

func TestRequiredAnnotations5(t *testing.T) {
	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "deployments"

	review.Request.Operation = "DELETE"

	var rule *Rules
	res := rule.ReplicaLimits(&review, "../../../example/replica_limits/replica_limits.conf", "file://../../../example/replica_limits/replica_limits.csv")
	if res != nil {
		t.Error("should not be rejected ")
		return
	}

}
