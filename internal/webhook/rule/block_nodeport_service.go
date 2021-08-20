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

	v1 "k8s.io/api/admission/v1"
	core "k8s.io/api/core/v1"
)

/*
[ks-admission-general-block-nodeport-services]

<Brief introduction>
Prohibit all NodePort services

<Coverage Area of Rule>
services resources will be checked.
*/
func (g *Rules) BlockNodeportService(review *v1.AdmissionReview, _ string, _ string) error {
	var resourceKind = review.Request.Resource.Resource
	if resourceKind != "services" {
		//only service is checked
		return nil
	}
	if review.Request.Operation == "DELETE" {
		//delete operation have no new object to check
		log.Printf("BlockNodeportService: service %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}

	var serviceObject core.Service
	if err := json.Unmarshal(review.Request.Object.Raw, &serviceObject); err != nil {
		log.Printf("BlockNodeportService: service %s:%s rejected due to error when unmarshal: %s.", serviceObject.Namespace, serviceObject.Name, err.Error())
		return err
	}
	if serviceObject.Spec.Type == core.ServiceTypeNodePort {
		log.Printf("BlockNodeportService: service %s:%s blocked due to type NodePort.", serviceObject.Namespace, serviceObject.Name)
		return fmt.Errorf("user is not allowed to create service of type NodePort")
	}
	log.Printf("BlockNodeportService: service %s:%s approved", serviceObject.Namespace, serviceObject.Name)
	return nil
}
