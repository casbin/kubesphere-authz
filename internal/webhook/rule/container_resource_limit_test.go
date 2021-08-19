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

func TestContainerResourceForPod1(t *testing.T) {
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceForPod2(t *testing.T) {
	var rule *Rules
	//No memory limit
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceForPod3(t *testing.T) {
	var rule *Rules
	//Exceeded resource limit
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("100Gi")
	container.Resources.Limits["memory"] = resource.MustParse("100Gi")
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceForPod4(t *testing.T) {
	var rule *Rules
	//should pass
	var podObject core.Pod
	var container core.Container
	container.Image = "nginx:1.14.2"
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//shouldn't reject
	if err != nil {
		t.Errorf("container without resource limits shouldn't be rejected")
	}
}

func TestContainerResourceForDeployment1(t *testing.T) {
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceForDeployment2(t *testing.T) {
	var rule *Rules
	//No memory limit
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("100Ki")
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceForDeployment3(t *testing.T) {
	var rule *Rules
	//Exceeded resource limit
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
	container.Resources.Limits["cpu"] = resource.MustParse("100Gi")
	container.Resources.Limits["memory"] = resource.MustParse("100Gi")
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//should reject
	if err == nil {
		t.Errorf("container without resource limits should be rejected")
	}
}

func TestContainerResourceForDeployment4(t *testing.T) {
	var rule *Rules
	//should pass
	var podObject app.Deployment
	var container core.Container
	container.Image = "nginx:1.14.2"
	container.Resources.Limits = make(core.ResourceList)
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

	err = rule.ContainerResourceLimit(&review, "../../../example/casbinconfig/container_resource_limit.conf", "../../../example/casbinconfig/container_resource_limit.csv")
	t.Log(err)
	//shouldn't reject
	if err != nil {
		t.Errorf("container without resource limits shouldn't be rejected")
	}
}
