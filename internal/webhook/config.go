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
package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"ksauth/internal/webhook/rule"

	"k8s.io/api/admission/v1"
)

type RuleConfig struct {
	Available bool   `json:"available"`
	Model     string `json:"model"`
	Policy    string `json:"policy"`
	RuleGroup string `json:"-"`
}

//maping check item name to configuration
var webHookConfig map[string]RuleConfig

var valueOfGeneral reflect.Value
var typeOfGeneral reflect.Type

//load the webhook/config.json
func Initconfig(configPath string) {
	fileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("Failed to load "+configPath+" due to" + err.Error())
	}

	err = json.Unmarshal(fileContent, &webHookConfig)
	if err != nil {
		log.Fatal("Failed to load "+configPath+" due to" + err.Error())
	}

	//use reflect to get type and value of general object
	var generalObj *rule.Rules
	valueOfGeneral = reflect.ValueOf(generalObj)
	typeOfGeneral = reflect.TypeOf(generalObj)
}

func enforceGeneralRules(methodName string, review *v1.AdmissionReview, model string, policy string) error {
	args := []reflect.Value{
		reflect.ValueOf(review),
		reflect.ValueOf(model),
		reflect.ValueOf(policy),
	}

	funcValue := valueOfGeneral.MethodByName(methodName)
	if !funcValue.IsValid() {
		return fmt.Errorf("invalid method name %s", methodName)
	}

	res := funcValue.Call(args)
	if len(res) != 1 {
		return fmt.Errorf("invalid method %s which returns %d values", methodName, len(res))
	}
	result := res[0]
	if result.IsNil() {
		return nil
	}
	return result.Interface().(error)

}
