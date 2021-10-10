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
	"fmt"
	"ksauth/pkg/casbinhelper"
	"log"

	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

/*
[ks-admission-general-container-resource-ratios]

<Brief introduction>
When a new container is created via k8s, you must specified the maximum cpu and memory limit for this container,
as well as the cpu and memory requested by this container. It's normal to have some redundancy for the resources so the amount declared in limit should be a little bigger than that in request. However, the redundancy should not be to high. This rule ensures that the redundancy ratio doesn't exceed the limit.

redundancy ratio= (amount of resource declared in limit)/(amount of resource declared in request)

<Coverage Area of Rule>
pods and deployments resources will be checked.
*/
func (g *Rules) ContainerResourceRatio(review *v1.AdmissionReview, model string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	switch resourceKind {
	case "pods":
		return g.containerResourceRatioForPod(review, model, policy)
	case "deployments":
		return g.containerResourceRatioForDeployment(review, model, policy)
	default:
		return nil
	}
}

func (g *Rules) containerResourceRatioForPod(review *v1.AdmissionReview, modelUrl string, policy string) error {
	adaptor, err := getAdaptorObject(policy)
	if err != nil {
		log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	model, err := getModelObject(modelUrl)
	if err != nil {
		log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)
	if err != nil {
		log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer.AddFunction("parseFloat", casbinhelper.ParseFloat)
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("ContainerResourceRatio: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}

	var podObject core.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	allContainers := make([]core.Container, len(podObject.Spec.Containers))
	copy(allContainers, podObject.Spec.Containers)
	allContainers = append(allContainers, podObject.Spec.InitContainers...)

	for _, container := range allContainers {
		resource := container.Resources
		cpuLimit, ok := resource.Limits["cpu"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no cpu limit", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("container %s has no cpu limit", container.Image)
		}
		memoryLimit, ok := resource.Limits["memory"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no memory limit", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("container %s has no memory limit", container.Image)
		}
		cpuRequest, ok := resource.Requests["cpu"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no cpu request", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("container %s has no cpu request", container.Image)
		}
		memoryRequest, ok := resource.Requests["memory"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no memory request", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("container %s has no memory request", container.Image)
		}

		cpuLimitInByte := cpuLimit.Value()
		memoryLimitInByte := memoryLimit.Value()
		cpuRequestInByte := cpuRequest.Value()
		memoryRequestInByte := memoryRequest.Value()

		cpuRedundancyRatio := float64(cpuLimitInByte) / float64(cpuRequestInByte)
		memoryRedundancyRatio := float64(memoryLimitInByte) / float64(memoryRequestInByte)

		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			cpuRedundancyRatio,
			memoryRedundancyRatio,
		)
		if err != nil {
			log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("ContainerResourceLimi(%s %s::%s): container %s resource redundancy ratio [cpu: %f, memory:%f] is higher than the maximum redundancy ratio", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image, cpuRedundancyRatio, memoryRedundancyRatio)
			return fmt.Errorf("container %s resource limit [cpu: %f, memory:%f] is higher than the maximum redundancy ratio", container.Image, cpuRedundancyRatio, memoryRedundancyRatio)
		}
	}
	log.Printf("ContainerResourceLimi(%s %s::%s) approved", review.Request.Resource.Resource, podObject.Namespace, podObject.Name)
	return nil

}

func (g *Rules) containerResourceRatioForDeployment(review *v1.AdmissionReview, modelUrl string, policy string) error {
	adaptor, err := getAdaptorObject(policy)
	if err != nil {
		log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	model, err := getModelObject(modelUrl)
	if err != nil {
		log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)

	if err != nil {
		log.Printf("ContainerResourceRatio: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer.AddFunction("parseFloat", casbinhelper.ParseFloat)
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("ContainerResourceRatio: deployment %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var deploymentObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &deploymentObject); err != nil {
		log.Printf("ContainerResourceRatio: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	allContainers := make([]core.Container, len(deploymentObject.Spec.Template.Spec.Containers))
	copy(allContainers, deploymentObject.Spec.Template.Spec.Containers)
	allContainers = append(allContainers, deploymentObject.Spec.Template.Spec.InitContainers...)

	for _, container := range allContainers {
		resource := container.Resources
		cpuLimit, ok := resource.Limits["cpu"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no cpu limit", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("container %s has no cpu limit", container.Image)
		}
		memoryLimit, ok := resource.Limits["memory"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no memory limit", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("container %s has no memory limit", container.Image)
		}
		cpuRequest, ok := resource.Requests["cpu"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no cpu request", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("container %s has no cpu request", container.Image)
		}
		memoryRequest, ok := resource.Requests["memory"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no memory request", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("container %s has no memory request", container.Image)
		}

		cpuLimitInByte := cpuLimit.Value()
		memoryLimitInByte := memoryLimit.Value()
		cpuRequestInByte := cpuRequest.Value()
		memoryRequestInByte := memoryRequest.Value()

		cpuRedundancyRatio := float64(cpuLimitInByte) / float64(cpuRequestInByte)
		memoryRedundancyRatio := float64(memoryLimitInByte) / float64(memoryRequestInByte)

		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			cpuRedundancyRatio,
			memoryRedundancyRatio,
		)
		if err != nil {
			log.Printf("ContainerResourceRatio: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("ContainerResourceLimi(%s %s::%s): container %s resource redundancy ratio [cpu: %f, memory:%f] is higher than the maximum limit", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image, cpuRedundancyRatio, memoryRedundancyRatio)
			return fmt.Errorf("container %s resource limit [cpu: %f, memory:%f] is higher than the maximum limit", container.Image, cpuRedundancyRatio, memoryRedundancyRatio)
		}

	}
	log.Printf("ContainerResourceLimi(%s %s::%s) approved", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name)
	return nil
}
