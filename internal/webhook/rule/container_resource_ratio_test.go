// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package rule

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestContainerResourceRatioForPod1(t *testing.T) {
	var rule *Rules
	//No cpu and memory limit
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceRatioForPod2(t *testing.T) {
	var rule *Rules
	//No request
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("100Ki")
	container.Resources.Limits["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}
func TestContainerResourceRatioForPod3(t *testing.T) {
	var rule *Rules
	//No Limits
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Requests["cpu"] = resource.MustParse("100Ki")
	container.Resources.Requests["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceRatioForPod4(t *testing.T) {
	var rule *Rules
	//Exceeded cpu redundancy ratio
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("300Ki")
	container.Resources.Limits["memory"] = resource.MustParse("100Ki")
	container.Resources.Requests["cpu"] = resource.MustParse("100Ki")
	container.Resources.Requests["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceRatioForPod5(t *testing.T) {
	var rule *Rules
	//should pass
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("101Ki")
	container.Resources.Limits["memory"] = resource.MustParse("101Ki")
	container.Resources.Requests["cpu"] = resource.MustParse("100Ki")
	container.Resources.Requests["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//shouldn't reject
	if err != nil {
		t.Errorf("container without resource limits shouldn't be rejected")
	}
}

func TestContainerResourceRatioForDeployment1(t *testing.T) {
	var rule *Rules
	//No cpu and memory limit
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceRatioForDeployment2(t *testing.T) {
	var rule *Rules
	//No requests
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("100Ki")
	container.Resources.Limits["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceRatioForDeployment3(t *testing.T) {
	var rule *Rules
	//no limits
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Requests["cpu"] = resource.MustParse("100Ki")
	container.Resources.Requests["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceRatioForDeployment4(t *testing.T) {
	var rule *Rules
	//exceeded cpu redundancy
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("300Ki")
	container.Resources.Limits["memory"] = resource.MustParse("100Ki")
	container.Resources.Requests["cpu"] = resource.MustParse("100Ki")
	container.Resources.Requests["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}
func TestContainerResourceRatioForDeployment5(t *testing.T) {
	var rule *Rules
	//should pass
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Requests = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("101Ki")
	container.Resources.Limits["memory"] = resource.MustParse("101Ki")
	container.Resources.Requests["cpu"] = resource.MustParse("100Ki")
	container.Resources.Requests["memory"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceRatio(&review, "file://../../../example/container_resource_ratio/container_resource_ratio.conf", "file://../../../example/container_resource_ratio/container_resource_ratio.csv")
	t.Log(err)
	//shouldn't reject
	if err != nil {
		t.Errorf("container without resource limits shouldn't be rejected")
	}
}
