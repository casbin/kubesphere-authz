package rule

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

func TestRequiredProbesForPod(t *testing.T) {
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
	res := rule.RequiredProbes(&review, "../../../example/required_probes/required_probes.conf", "../../../example/required_probes/required_probes.csv")
	if res == nil {
		t.Error("should be rejected")
		return
	}
}

func TestRequiredProbesForPod2(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"

	container.ReadinessProbe = &core.Probe{}
	container.ReadinessProbe.TCPSocket = &core.TCPSocketAction{}

	container.LivenessProbe = &core.Probe{}
	container.LivenessProbe.TCPSocket = &core.TCPSocketAction{}

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
	res := rule.RequiredProbes(&review, "../../../example/required_probes/required_probes.conf", "../../../example/required_probes/required_probes.csv")
	if res != nil {
		t.Error("should not be rejected")
		return
	}
}

func TestRequiredProbesForPod3(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"

	container.ReadinessProbe = &core.Probe{}
	container.ReadinessProbe.Exec = &core.ExecAction{}

	container.LivenessProbe = &core.Probe{}
	container.LivenessProbe.TCPSocket = &core.TCPSocketAction{}

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
	res := rule.RequiredProbes(&review, "../../../example/required_probes/required_probes.conf", "../../../example/required_probes/required_probes.csv")
	if res == nil {
		t.Error("should be rejected")
		return
	}
}

func TestRequiredProbesForDeployment(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"

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

	res := rule.RequiredProbes(&review, "../../../example/required_probes/required_probes.conf", "../../../example/required_probes/required_probes.csv")
	if res == nil {
		t.Error("should be rejected")
	}
}

func TestRequiredProbesForDeployment2(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"

	container.ReadinessProbe = &core.Probe{}
	container.ReadinessProbe.TCPSocket = &core.TCPSocketAction{}

	container.LivenessProbe = &core.Probe{}
	container.LivenessProbe.TCPSocket = &core.TCPSocketAction{}

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

	res := rule.RequiredProbes(&review, "../../../example/required_probes/required_probes.conf", "../../../example/required_probes/required_probes.csv")
	if res != nil {
		t.Error("should not be rejected")
	}
}

func TestRequiredProbesForDeployment3(t *testing.T) {
	var rule *Rules
	//var review v1.AdmissionReview
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"

	container.ReadinessProbe = &core.Probe{}
	container.ReadinessProbe.Exec = &core.ExecAction{}

	container.LivenessProbe = &core.Probe{}
	container.LivenessProbe.TCPSocket = &core.TCPSocketAction{}

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

	res := rule.RequiredProbes(&review, "../../../example/required_probes/required_probes.conf", "../../../example/required_probes/required_probes.csv")
	if res == nil {
		t.Error("should  be rejected")
	}
}
