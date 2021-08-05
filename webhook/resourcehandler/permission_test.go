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
package resourcehandler
import (
	v1 "k8s.io/api/admission/v1"
	"testing"
)
func TestCheckPermission(t *testing.T){
	var review v1.AdmissionReview
	review.Request=&v1.AdmissionRequest{}
	review.Request.Namespace="default"
	review.Request.Name="my-nginx-svc"
	review.Request.Resource.Resource="services"
	review.Request.Operation="DELETE"
	res:=CheckPermission(review,"../casbinconfig/permission_model.conf","../casbinconfig/permission_policy.csv")
	if res==nil{
		t.Error("should be rejected")
	}
}

func TestCheckPermission2(t *testing.T){
	var review v1.AdmissionReview
	review.Request=&v1.AdmissionRequest{}
	review.Request.Namespace="default"
	review.Request.Name="my-nginx-svc"
	review.Request.Resource.Resource="services"
	review.Request.Operation="UPDATE"
	res:=CheckPermission(review,"../casbinconfig/permission_model.conf","../casbinconfig/permission_policy.csv")
	if res!=nil{
		t.Error(res)
	}
}