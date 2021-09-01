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
	networking "k8s.io/api/networking/v1"
)

func TestHttpsOnly1(t *testing.T) {
	var ingress networking.Ingress
	var ingressTls networking.IngressTLS
	ingressTls.Hosts = append(ingressTls.Hosts, "example-host.example.com")
	ingress.Spec.TLS = append(ingress.Spec.TLS, ingressTls)
	ingress.Annotations = make(map[string]string)
	ingress.Annotations["kubernetes.io/ingress.allow-http"] = "false"

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "ingresses"

	data, err := json.Marshal(ingress)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	review.Request.Operation = "CREATE"

	var rule *Rules
	err = rule.HttpsOnly(&review, "", "")
	if err != nil {
		t.Errorf("should have passed, got" + err.Error())
	}

}

func TestHttpsOnly2(t *testing.T) {
	var ingress networking.Ingress
	var ingressTls networking.IngressTLS
	ingressTls.Hosts = append(ingressTls.Hosts, "example-host.example.com")
	ingress.Spec.TLS = append(ingress.Spec.TLS, ingressTls)
	ingress.Annotations = make(map[string]string)
	//ingress.Annotations["kubernetes.io/ingress.allow-http"]="false"

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "ingresses"

	data, err := json.Marshal(ingress)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	review.Request.Operation = "CREATE"

	var rule *Rules
	err = rule.HttpsOnly(&review, "", "")
	if err == nil {
		t.Errorf("should have failed, ")
	}

}

func TestHttpsOnly3(t *testing.T) {
	var ingress networking.Ingress
	//var ingressTls networking.IngressTLS
	//ingressTls.Hosts = append(ingressTls.Hosts, "example-host.example.com")
	//ingress.Spec.TLS=append(ingress.Spec.TLS,ingressTls)
	ingress.Annotations = make(map[string]string)
	ingress.Annotations["kubernetes.io/ingress.allow-http"] = "false"

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "ingresses"

	data, err := json.Marshal(ingress)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	review.Request.Operation = "CREATE"

	var rule *Rules
	err = rule.HttpsOnly(&review, "", "")
	if err == nil {
		t.Errorf("should have failed, ")
	}

}

func TestHttpsOnly4(t *testing.T) {
	var ingress networking.Ingress

	var review v1.AdmissionReview
	review.Request = &v1.AdmissionRequest{}
	review.Request.Namespace = "default"
	review.Request.Resource.Resource = "ingresses"

	data, err := json.Marshal(ingress)
	if err != nil {
		t.Error(err)
	}
	review.Request.Object.Raw = data
	review.Request.Operation = "DELETE"

	var rule *Rules
	err = rule.HttpsOnly(&review, "", "")
	if err != nil {
		t.Errorf("shouldn't have failed, got" + err.Error())
	}

}
