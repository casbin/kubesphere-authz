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
[ks-admission-general-container-resource-limits]

<Brief introduction>
When a new container is created via k8s, you must specified the maximum cpu and memory limit for this container, and these limits must not exceed a user specified limit

<Coverage Area of Rule>
pods and deployments resources will be checked.
*/
func (g *Rules) ContainerResourceLimit(review *v1.AdmissionReview, model string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	switch resourceKind {
	case "pods":
		return g.containerResourceLimitForPod(review, model, policy)
	case "deployments":
		return g.containerResourceLimitForDeployment(review, model, policy)
	default:
		return nil
	}
}

func (g *Rules) containerResourceLimitForPod(review *v1.AdmissionReview, model string, policy string) error {
	adaptor,err:=getAdaptorObject(policy)
	if err != nil {
		log.Printf("ContainerResourceLimit: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)
	if err != nil {
		log.Printf("ContainerResourceLimit: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer.AddFunction("parseFloat", casbinhelper.ParseFloat)
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("ContainerResourceLimit: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}

	var podObject core.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		log.Printf("ContainerResourceLimit: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	allContainers := make([]core.Container, len(podObject.Spec.Containers))
	copy(allContainers, podObject.Spec.Containers)
	allContainers = append(allContainers, podObject.Spec.InitContainers...)

	for _, container := range allContainers {
		resource := container.Resources
		cpu, ok := resource.Limits["cpu"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no cpu limit", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("container %s has no cpu limit", container.Image)
		}
		memory, ok := resource.Limits["memory"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no memory limit", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("container %s has no memory limit", container.Image)
		}
		cpuInByte := cpu.Value()
		memoryInByte := memory.Value()
		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			cpuInByte,
			memoryInByte,
		)
		if err != nil {
			log.Printf("ContainerResourceLimit: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("ContainerResourceLimi(%s %s::%s): container %s resource limit [cpu: %d, memory:%d] is higher than the maximum limit", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image, cpuInByte, memoryInByte)
			return fmt.Errorf("container %s resource limit [cpu: %d, memory:%d] is higher than the maximum limit", container.Image, cpuInByte, memoryInByte)
		}
	}
	log.Printf("ContainerResourceLimi(%s %s::%s) approved", review.Request.Resource.Resource, podObject.Namespace, podObject.Name)
	return nil

}

func (g *Rules) containerResourceLimitForDeployment(review *v1.AdmissionReview, model string, policy string) error {
	adaptor,err:=getAdaptorObject(policy)
	if err != nil {
		log.Printf("ContainerResourceLimit: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)

	if err != nil {
		log.Printf("ContainerResourceLimit: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer.AddFunction("parseFloat", casbinhelper.ParseFloat)
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("ContainerResourceLimit: deployment %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var deploymentObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &deploymentObject); err != nil {
		log.Printf("ContainerResourceLimit: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	allContainers := make([]core.Container, len(deploymentObject.Spec.Template.Spec.Containers))
	copy(allContainers, deploymentObject.Spec.Template.Spec.Containers)
	allContainers = append(allContainers, deploymentObject.Spec.Template.Spec.InitContainers...)

	for _, container := range allContainers {
		resource := container.Resources
		cpu, ok := resource.Limits["cpu"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no cpu limit", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("container %s has no cpu limit", container.Image)
		}
		memory, ok := resource.Limits["memory"]
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s has no memory limit", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("container %s has no memory limit", container.Image)
		}
		cpuInByte := cpu.Value()
		memoryInByte := memory.Value()
		fmt.Println(cpuInByte, memoryInByte)
		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			cpuInByte,
			memoryInByte,
		)
		if err != nil {
			log.Printf("ContainerResourceLimit: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("ContainerResourceLimit(%s %s::%s): container %s resource limit [cpu: %d, memory:%d] is higher than the maximum limit", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image, cpuInByte, memoryInByte)
			return fmt.Errorf("container %s resource limit [cpu: %d, memory:%d] is higher than the maximum limit", container.Image, cpuInByte, memoryInByte)
		}

	}
	log.Printf("ContainerResourceLimi(%s %s::%s) approved", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name)
	return nil
}
