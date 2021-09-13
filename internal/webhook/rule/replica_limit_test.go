package rule

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

func TestReplicaLimitTest(t *testing.T) {
	var deployment app.Deployment
	var tmp int32 = 4
	deployment.Spec.Replicas = &tmp
	var container core.Container
	container.Image = "nginx:1.14.2"
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "deployments"

	data, err := json.Marshal(deployment)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data

	var rule *Rules
	res := rule.ReplicaLimits(&review, "../../../example/replica_limits/replica_limits.conf", "file://../../../example/replica_limits/replica_limits.csv")

	if res != nil {
		t.Errorf("should not be rejected, but got %s", err.Error())
	}

}

func TestReplicaLimitTest2(t *testing.T) {
	var deployment app.Deployment
	var tmp int32 = 400
	deployment.Spec.Replicas = &tmp
	var container core.Container
	container.Image = "nginx:1.14.2"
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "deployments"

	data, err := json.Marshal(deployment)
	if err != nil {
		t.Error(err)
		return
	}
	review.Request.Object.Raw = data

	var rule *Rules
	res := rule.ReplicaLimits(&review, "../../../example/replica_limits/replica_limits.conf", "file://../../../example/replica_limits/replica_limits.csv")

	if res == nil {
		t.Error("should be rejected ")
		return
	}

}

func TestReplicaLimitTest3(t *testing.T) {
	var deployment app.Deployment

	var container core.Container
	container.Image = "nginx:1.14.2"
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "deployments"

	data, err := json.Marshal(deployment)
	if err != nil {
		t.Error(err)
		return
	}
	review.Request.Object.Raw = data

	var rule *Rules
	res := rule.ReplicaLimits(&review, "../../../example/replica_limits/replica_limits.conf", "file://../../../example/replica_limits/replica_limits.csv")
	if res == nil {
		t.Error("should be rejected ")
		return
	}

}

func TestReplicaLimitTest4(t *testing.T) {
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
