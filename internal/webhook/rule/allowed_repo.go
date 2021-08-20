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
	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"ksauth/pkg/casbinhelper"
	"log"
)

/*
[ks-admission-general-allowed-repos]

<Brief introduction>
Only images specified in the associated casbin policy can be allowed.
Precisely, if the image name starts with any prefix specified by the policy, then this image is allowed.

<Coverage Area of Rule>
pods and deployments resources will be checked.
*/
func (g *Rules) AllowedRepos(review *v1.AdmissionReview, model string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	//fmt.Println(resourceKind)
	switch resourceKind {
	case "pods":
		return g.allowedReposForPod(review, model, policy)
	case "deployments":
		return g.allowedReposForDeployment(review, model, policy)
	default:
		return nil
	}
}

func (g *Rules) allowedReposForPod(review *v1.AdmissionReview, model string, policy string) error {
	enforcer, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		log.Printf("AllowedRepos: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer.AddFunction("hasPrefix", casbinhelper.HasPrefix)

	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("AllowedRepos: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var podObject core.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		log.Printf("AllowedRepos: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	allContainers := make([]core.Container, len(podObject.Spec.Containers))
	copy(allContainers, podObject.Spec.Containers)
	allContainers = append(allContainers, podObject.Spec.InitContainers...)
	for _, container := range allContainers {
		var image = container.Image
		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			image,
		)

		if err != nil {
			log.Printf("AllowedRepos: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("AllowedRepos(%s %s::%s): container %s is not allowed", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("casbin rejects the untrusted image %s", image)
		}

	}
	log.Printf("AllowedRepos(%s %s::%s): approved", review.Request.Resource.Resource, podObject.Namespace, podObject.Name)
	return nil
}

func (g *Rules) allowedReposForDeployment(review *v1.AdmissionReview, model string, policy string) error {
	enforcer, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		log.Printf("AllowedRepos: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer.AddFunction("hasPrefix", casbinhelper.HasPrefix)
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("AllowedRepos: deployment %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var deploymentObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &deploymentObject); err != nil {
		log.Printf("AllowedRepos: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	allContainers := make([]core.Container, len(deploymentObject.Spec.Template.Spec.Containers))
	copy(allContainers, deploymentObject.Spec.Template.Spec.Containers)
	allContainers = append(allContainers, deploymentObject.Spec.Template.Spec.InitContainers...)
	for _, container := range allContainers {
		var image = container.Image
		fmt.Println(image)
		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			image,
		)

		if err != nil {
			log.Printf("AllowedRepos: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("AllowedRepos(%s %s::%s): container %s is not allowed", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("casbin rejects the untrusted image %s", image)
		}

	}
	log.Printf("AllowedRepos(%s %s::%s): approved", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name)
	return nil
}
