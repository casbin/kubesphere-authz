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
	"fmt"
	"ksauth/internal/webhook"
	"os"
)

//please use ${workspaceDir} as cwd, since all unittest and default config depends on this prerequisite
func main() {
	var configPath = ""
	if len(os.Args) > 2 {
		fmt.Println("stating webhook needs exactly 1 parameter: path of config file.")
		return
	} else if len(os.Args) == 1 {
		//use default config path: ${workspaceDir}/config/config/config.json
		configPath = "config/config/config.json"
	} else {
		configPath = os.Args[1]
	}
	webhook.Initconfig(configPath)
	webhook.GetAdmissionWebhook().RunTLS(":8080", "config/certificate/server.crt", "config/certificate/server.key")
}
