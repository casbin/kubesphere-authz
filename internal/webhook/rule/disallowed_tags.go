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
	"log"
	"strings"

	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

/*
[ks-admission-general-disallowed-tags]

<Brief introduction>
Only images with a tag specified in the associated casbin policy can be allowed.
image without declaring a specific tag will also be rejected.

<Coverage Area of Rule>
pods and deployments resources will be checked.
*/
func (g *Rules) DisallowedTags(review *v1.AdmissionReview, model string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	//fmt.Println(resourceKind)
	switch resourceKind {
	case "pods":
		return g.disallowedTagsForPod(review, model, policy)
	case "deployments":
		return g.disallowedTagsForDeployment(review, model, policy)
	default:
		return nil
	}
}

func (g *Rules) disallowedTagsForPod(review *v1.AdmissionReview, modelUrl string, policy string) error {
	adaptor, err := getAdaptorObject(policy)
	if err != nil {
		log.Printf("DisallowedTags: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	model, err := getModelObject(modelUrl)
	if err != nil {
		log.Printf("DisallowedTags: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)
	if err != nil {
		log.Printf("DisallowedTags: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("DisallowedTags: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var podObject core.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		log.Printf("DisallowedTags: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	allContainers := make([]core.Container, len(podObject.Spec.Containers))
	copy(allContainers, podObject.Spec.Containers)
	allContainers = append(allContainers, podObject.Spec.InitContainers...)
	for _, container := range allContainers {
		var image = container.Image

		slice := strings.Split(image, ":")
		if len(slice) <= 1 {
			log.Printf("DisallowedTags: pod %s:%s rejected due to error: no tags attached", review.Request.Namespace, review.Request.Name)
			return fmt.Errorf("DisallowedTags: pod %s:%s rejected due to error: no tags attached", review.Request.Namespace, review.Request.Name)
		}

		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			slice[len(slice)-1],
		)

		if err != nil {
			log.Printf("DisallowedTags: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("DisallowedTags(%s %s::%s): container %s is not allowed", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, container.Image)
			return fmt.Errorf("casbin rejects the  image %s with prohibited tag", image)
		}

	}
	log.Printf("DisallowedTags(%s %s::%s): approved", review.Request.Resource.Resource, podObject.Namespace, podObject.Name)
	return nil
}

func (g *Rules) disallowedTagsForDeployment(review *v1.AdmissionReview, modelUrl string, policy string) error {
	adaptor, err := getAdaptorObject(policy)
	if err != nil {
		log.Printf("DisallowedTags: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	model, err := getModelObject(modelUrl)
	if err != nil {
		log.Printf("DisallowedTags: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)

	if err != nil {
		log.Printf("DisallowedTags: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("DisallowedTags: deployment %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var deploymentObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &deploymentObject); err != nil {
		log.Printf("DisallowedTags: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	allContainers := make([]core.Container, len(deploymentObject.Spec.Template.Spec.Containers))
	copy(allContainers, deploymentObject.Spec.Template.Spec.Containers)
	allContainers = append(allContainers, deploymentObject.Spec.Template.Spec.InitContainers...)
	for _, container := range allContainers {
		var image = container.Image

		slice := strings.Split(image, ":")
		if len(slice) <= 1 {
			log.Printf("DisallowedTags: deployment %s:%s rejected due to error: no tags attached", review.Request.Namespace, review.Request.Name)
			return fmt.Errorf("DisallowedTags: deployment %s:%s rejected due to error: no tags attached", review.Request.Namespace, review.Request.Name)
		}

		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			slice[len(slice)-1],
		)

		if err != nil {
			log.Printf("DisallowedTags: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("DisallowedTags(%s %s::%s): container %s is not allowed", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, container.Image)
			return fmt.Errorf("casbin rejects the image %s with prohibited tag", image)
		}

	}
	log.Printf("DisallowedTags(%s %s::%s): approved", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name)
	return nil
}
