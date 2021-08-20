package rule

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

func TestDisallowedTagsForDeployment(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:test"
	podObject.Spec.Template.Spec.Containers = append(podObject.Spec.Template.Spec.Containers, container)

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "deployments"

	data, err := json.Marshal(podObject)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data

	res := rule.DisallowedTags(&review, "../../../example/disallowed_tags/disallowed_tags.conf", "../../../example/disallowed_tags/disallowed_tags.csv")
	if res == nil {
		t.Error("should be rejected")
	}

}
func TestDisallowedTagsForDeployment2(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:latest"
	podObject.Spec.Template.Spec.Containers = append(podObject.Spec.Template.Spec.Containers, container)

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "deployments"

	data, err := json.Marshal(podObject)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data

	res := rule.DisallowedTags(&review, "../../../example/disallowed_tags/disallowed_tags.conf", "../../../example/disallowed_tags/disallowed_tags.csv")
	if res != nil {
		t.Errorf("should not be rejected,%s", res.Error())
	}

}
func TestDisallowedTagsForDeployment3(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx"
	podObject.Spec.Template.Spec.Containers = append(podObject.Spec.Template.Spec.Containers, container)

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "deployments"

	data, err := json.Marshal(podObject)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data

	res := rule.DisallowedTags(&review, "../../../example/disallowed_tags/disallowed_tags.conf", "../../../example/disallowed_tags/disallowed_tags.csv")
	if res == nil {
		t.Errorf("should be rejected")
	}

}

func TestDisallowedTagsForPod(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:test"
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

	res := rule.DisallowedTags(&review, "../../../example/disallowed_tags/disallowed_tags.conf", "../../../example/disallowed_tags/disallowed_tags.csv")
	if res == nil {
		t.Error("should be rejected")
	}

}
func TestDisallowedTagsForPod2(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:latest"
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

	res := rule.DisallowedTags(&review, "../../../example/disallowed_tags/disallowed_tags.conf", "../../../example/disallowed_tags/disallowed_tags.csv")
	if res != nil {
		t.Errorf("should not be rejected,%s", res.Error())
	}

}
func TestDisallowedTagsForPod3(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx"
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

	res := rule.DisallowedTags(&review, "../../../example/disallowed_tags/disallowed_tags.conf", "../../../example/disallowed_tags/disallowed_tags.csv")
	if res == nil {
		t.Errorf("should be rejected")
	}

}
