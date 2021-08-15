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
)

func (g *Rules) AllowedRepos(review v1.AdmissionReview, model string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	switch resourceKind {
	case "pods":
		return g.allowedReposForPod(review, model, policy)
	case "deployments":
		return g.allowedReposForDeployment(review, model, policy)
	default:
		return nil
	}
}

func (g *Rules) allowedReposForPod(review v1.AdmissionReview, model string, policy string) error {
	enforcer, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		return err
	}

	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		return nil
	}
	var podObject core.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		return err
	}
	containers := podObject.Spec.Containers
	for _, container := range containers {
		var image = container.Image
		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			image,
		)

		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("casbin rejects the untrusted image %s", image)
		}

	}
	return nil
}

func (g *Rules) allowedReposForDeployment(review v1.AdmissionReview, model string, policy string) error {
	enforcer, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		return err
	}
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		return nil
	}
	var podObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		return err
	}
	containers := podObject.Spec.Template.Spec.Containers
	for _, container := range containers {
		var image = container.Image

		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			image,
		)

		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("casbin rejects the untrusted image %s", image)
		}

	}
	return nil
}
