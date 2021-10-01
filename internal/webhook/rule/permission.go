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
	"fmt"
	"log"

	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
)

func (g *Rules) ResourceOperationPermission(review *v1.AdmissionReview, modelUrl string, policy string) error {
	adaptor,err:=getAdaptorObject(policy)
	if err != nil {
		return err
	}
	model,err:=getModelObject(modelUrl)
	if err != nil {
		return err
	}
	e, err := casbin.NewEnforcer(model, adaptor)
	if err != nil {
		return err
	}
	ok, err := e.Enforce(
		review.Request.Name,
		review.Request.Resource.Resource,
		string(review.Request.Operation),
	)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if !ok {
		return fmt.Errorf("checkPermission rejected this request")
	}
	return nil

}
