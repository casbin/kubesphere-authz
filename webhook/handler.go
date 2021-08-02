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
	//"fmt"
	"fmt"
	"io/ioutil"

	"k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	//casbin "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {

	data, _ := ioutil.ReadAll(c.Request.Body)
	var requestBody v1.AdmissionReview
	var decoder runtime.Decoder = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	decoder.Decode(data, nil, &requestBody)
	//for debug only. Todo:remove this block of code

	fmt.Printf("%s\n", requestBody.Request.Resource.Resource)
	if requestBody.Request.Namespace != "default" {
		approve(c, string(requestBody.Request.UID))
		return
	}
	//have all checkitems checked
	for item, config := range webHookConfig {
		if !config.Available {
			continue
		}
		if checkFunc, ok := webHookCheckItem[item]; ok {
			err := checkFunc(requestBody, config.Model, config.Policy)
			if err != nil {
				reject(c, string(requestBody.Request.UID), err.Error())
				return
			}
		}
	}
	approve(c, string(requestBody.Request.UID))

}

func reject(c *gin.Context, uid string, rejectReason string) {
	c.JSON(200, gin.H{
		"apiVersion": "admission.k8s.io/v1",
		"kind":       "AdmissionReview",
		"response": map[string]interface{}{
			"uid":     uid,
			"allowed": false,
			"status": map[string]interface{}{
				"code":    403,
				"message": rejectReason,
			},
		},
	})
}

func approve(c *gin.Context, uid string) {
	c.JSON(200, gin.H{
		"apiVersion": "admission.k8s.io/v1",
		"kind":       "AdmissionReview",
		"response": map[string]interface{}{
			"uid":     uid,
			"allowed": true,
		},
	})
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "pong",
	})
}
