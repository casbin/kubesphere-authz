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
	networking "k8s.io/api/networking/v1"
	"log"
)

func (g *Rules) HttpsOnly(review *v1.AdmissionReview, _ string, _ string) error {
	var resourceKind = review.Request.Resource.Resource
	if resourceKind != "ingresses" {
		//only ingresses is checked
		return nil
	}
	if review.Request.Operation == "DELETE" {
		//delete operation have no new object to check
		log.Printf("HttpsOnly: service %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var ingress networking.Ingress
	if err := json.Unmarshal(review.Request.Object.Raw, &ingress); err != nil {
		log.Printf("ExternalIP: ingress %s:%s rejected due to error when unmarshal: %s.", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if len(ingress.Spec.TLS) == 0 {
		log.Printf("ExternalIP: ingress %s:%s rejected, no TLS configuration found", ingress.Namespace, ingress.Name)
		return fmt.Errorf("ExternalIP: ingress %s:%s rejected, no TLS configuration found", ingress.Namespace, ingress.Name)
	}

	if v, ok := ingress.Annotations["kubernetes.io/ingress.allow-http"]; !ok || v != "false" {
		log.Printf("ExternalIP: ingress %s:%s rejected, because no annotation 'kubernetes.io/ingress.allow-http=false' found", ingress.Namespace, ingress.Name)
		return fmt.Errorf("ExternalIP: ingress %s:%s rejected, because no annotation 'kubernetes.io/ingress.allow-http=false' found", ingress.Namespace, ingress.Name)
	}
	log.Printf("ExternalIP: ingress %s:%s approved", ingress.Namespace, ingress.Name)
	return nil
}
