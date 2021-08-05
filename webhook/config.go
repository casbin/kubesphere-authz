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
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"k8s.io/api/admission/v1"
	"webhook/resourcehandler"
)

const (
	RESOURCE_OPERATION_PERMISSION string = "resourceOperationPermission"
	POD_IMAGE_CHECK               string = "podImageCheck"
	DEPLOYMENT_IMAGE_CHECK        string = "deploymentImageCheck"
)

type CasbinCheckItem func(review v1.AdmissionReview, model string, policy string) error

type ItemConfig struct {
	Available bool   `json:"available"`
	Model     string `json:"model"`
	Policy    string `json:"policy"`
}

var webHookCheckItem = map[string]CasbinCheckItem{
	RESOURCE_OPERATION_PERMISSION: resourcehandler.CheckPermission,
	POD_IMAGE_CHECK:               resourcehandler.CheckTrustedImageOfPod,
	DEPLOYMENT_IMAGE_CHECK:        resourcehandler.CheckTrustedImageOfDeployment,
}
var webHookConfig map[string]ItemConfig

//load the webhook/config.json
func init() {
	fileContent, err := ioutil.ReadFile("webhookconfig/config.json")
	if err != nil {
		log.Fatal("Failed to load webhookconfig/config.json")
	}
	json.Unmarshal(fileContent, &webHookConfig)
	//todo: check validity of config
}
