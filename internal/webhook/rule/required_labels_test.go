package rule

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/admission/v1"
	core "k8s.io/api/core/v1"
)

func TestRequiredLabel(t *testing.T) {
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
	res := rule.RequiredLabel(&review, "../../../example/required_labels/required_labels.conf", "file://../../../example/required_labels/required_labels.csv")
	if res == nil {
		t.Error("should be rejected")
		return
	}
}

func TestRequiredLabel2(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	podObject.Spec.Containers = append(podObject.Spec.Containers, container)

	podObject.SetLabels(map[string]string{
		"owner": "test-100",
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
	res := rule.RequiredLabel(&review, "../../../example/required_labels/required_labels.conf", "file://../../../example/required_labels/required_labels.csv")
	if res != nil {
		t.Error("should not be rejected")
		return
	}
}

func TestRequiredLabel3(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	podObject.Spec.Containers = append(podObject.Spec.Containers, container)

	podObject.SetLabels(map[string]string{
		"owner": "tdddest-100",
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
	res := rule.RequiredLabel(&review, "../../../example/required_labels/required_labels.conf", "file://../../../example/required_labels/required_labels.csv")
	if res == nil {
		t.Error("should  be rejected")
		return
	}
}

func TestRequiredLabel4(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	podObject.Spec.Containers = append(podObject.Spec.Containers, container)

	podObject.SetLabels(map[string]string{
		"owner":  "test-100",
		"owner2": "tdddest-100",
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
	res := rule.RequiredLabel(&review, "../../../example/required_labels/required_labels.conf", "file://../../../example/required_labels/required_labels.csv")
	if res != nil {
		t.Errorf("should not be rejected, but got %s", res.Error())
		return
	}
}

func TestRequiredLabel5(t *testing.T) {
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
