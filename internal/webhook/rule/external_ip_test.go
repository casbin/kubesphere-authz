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
	core "k8s.io/api/core/v1"
)

func TestExternalIP(t *testing.T) {
	var rule *Rules

	var serviceObj core.Service
	serviceObj.SetName("my-nginx-svc")
	serviceObj.Spec.Type = core.ServiceTypeNodePort

	serviceObj.Spec.ExternalIPs = append(serviceObj.Spec.ExternalIPs, "1.1.1.1")

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "services"

	data, err := json.Marshal(serviceObj)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data

	err = rule.ExternalIP(&review, "../../../example/external_ip/external_ip.conf", "file://../../../example/external_ip/external_ip.csv")
	if err != nil {
		t.Errorf("Should have passed, but get " + err.Error())
	}
}
func TestExternalIP2(t *testing.T) {
	var rule *Rules

	var serviceObj core.Service
	serviceObj.SetName("my-nginx-svc")
	serviceObj.Spec.Type = core.ServiceTypeNodePort
	serviceObj.Spec.ExternalIPs = append(serviceObj.Spec.ExternalIPs, "1.1.1.1")
	serviceObj.Spec.ExternalIPs = append(serviceObj.Spec.ExternalIPs, "2.2.2.2")

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "services"

	data, err := json.Marshal(serviceObj)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data

	err = rule.ExternalIP(&review, "../../../example/external_ip/external_ip.conf", "file://../../../example/external_ip/external_ip.csv")
	if err == nil {
		t.Errorf("Should have failed, nodeport service shouldn't be allowed")
	}
}

func TestExternalIP3(t *testing.T) {
	var rule *Rules

	var serviceObj core.Service
	serviceObj.SetName("my-nginx-svc")
	serviceObj.Spec.Type = core.ServiceTypeNodePort

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "services"

	data, err := json.Marshal(serviceObj)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data

	err = rule.ExternalIP(&review, "../../../example/external_ip/external_ip.conf", "file://../../../example/external_ip/external_ip.csv")
	if err != nil {
		t.Errorf("Should have passed, but get " + err.Error())
	}
}
