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

	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
	core "k8s.io/api/core/v1"
)

func (g *Rules) ExternalIP(review *v1.AdmissionReview, modelUrl string, policy string) error {

	var resourceKind = review.Request.Resource.Resource
	if resourceKind != "services" {
		//only service is checked
		return nil
	}
	if review.Request.Operation == "DELETE" {
		//delete operation have no new object to check
		log.Printf("ExternalIP: service %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}

	var serviceObject core.Service
	if err := json.Unmarshal(review.Request.Object.Raw, &serviceObject); err != nil {
		log.Printf("ExternalIP: service %s:%s rejected due to error when unmarshal: %s.", serviceObject.Namespace, serviceObject.Name, err.Error())
		return err
	}

	model, adaptor, err := getModelAndPolicyObject(modelUrl, policy)
	if err != nil {
		log.Printf("ExternalIP: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)
	if err != nil {
		log.Printf("ExternalIP: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	for _, ip := range serviceObject.Spec.ExternalIPs {
		ok, err := enforcer.Enforce(
			review.Request.Namespace,
			ip,
		)
		if err != nil {
			log.Printf("ExternalIP: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if !ok {
			log.Printf("ExternalIP: Service %s:%s is rejected due to using external ip %s", review.Request.Namespace, serviceObject.Name, ip)
			return fmt.Errorf("ExternalIP: Service %s:%s is rejected due to using external ip %s", review.Request.Namespace, serviceObject.Name, ip)
		}
	}
	log.Printf("ExternalIP: Service %s:%s is approved", review.Request.Namespace, serviceObject.Name)
	return nil
}
