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
	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"log"
	"regexp"
)

var imageRegex *regexp.Regexp

func init() {
	imageRegex, _ = regexp.Compile(`@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+`)
}

/*
[ks-admission-general-image-digest]

<Brief introduction>
Only images with a valid digest can be allowed.

<Coverage Area of Rule>
pods and deployments resources will be checked.
*/
func (g *Rules) ImageDigest(review *v1.AdmissionReview, model string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	//fmt.Println(resourceKind)
	switch resourceKind {
	case "pods":
		return g.imageDigestForPod(review, model, policy)
	case "deployments":
		return g.imageDigestForDeployment(review, model, policy)
	default:
		return nil
	}
}

func (g *Rules) imageDigestForPod(review *v1.AdmissionReview, model string, policy string) error {

	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("ImageDigest: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var podObject core.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		log.Printf("ImageDigest: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	allContainers := make([]core.Container, len(podObject.Spec.Containers))
	copy(allContainers, podObject.Spec.Containers)
	allContainers = append(allContainers, podObject.Spec.InitContainers...)

	for _, container := range allContainers {
		var image = container.Image
		ok := imageRegex.MatchString(image)

		if !ok {
			log.Printf("ImageDigest(%s %s::%s):reject due to error no legal digest attached in image %s", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, image)
			return fmt.Errorf("ImageDigest(%s %s::%s):reject due to error no legal digest attached in image %s", review.Request.Resource.Resource, podObject.Namespace, podObject.Name, image)
		}

	}
	log.Printf("ImageDigest(%s %s::%s): approved", review.Request.Resource.Resource, podObject.Namespace, podObject.Name)
	return nil
}

func (g *Rules) imageDigestForDeployment(review *v1.AdmissionReview, model string, policy string) error {

	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("ImageDigest: deployment %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var deploymentObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &deploymentObject); err != nil {
		log.Printf("ImageDigest: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	allContainers := make([]core.Container, len(deploymentObject.Spec.Template.Spec.Containers))
	copy(allContainers, deploymentObject.Spec.Template.Spec.Containers)
	allContainers = append(allContainers, deploymentObject.Spec.Template.Spec.InitContainers...)

	for _, container := range allContainers {
		var image = container.Image
		ok := imageRegex.MatchString(image)

		if !ok {
			log.Printf("ImageDigest(%s %s::%s):reject due to error no legal digest attached in image %s", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, image)
			return fmt.Errorf("ImageDigest(%s %s::%s):reject due to error no legal digest attached in image %s", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, image)
		}

	}
	log.Printf("ImageDigest(%s %s::%s): approved", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name)
	return nil
}
