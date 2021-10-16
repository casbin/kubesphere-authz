# casbin-kubesphere-auth
[TOC]
## 0.Overview
Casbin-kubesphere-auth is a plugin which apply several security authentication check on kubesphere via [casbin](https://casbin.org/docs/en/overview).

In fact, this plugin is an admission webhook of k8s, performing various kinds of checks and applying the security rules on each operation. 

Actually, you can also apply this plugin on a raw k8s.

## 1.Quick Start

1. Pull this repository.

2. For your convenience, save the current directory path by executing `workspace=$(pwd)`

3. If you haven't start the k8s cluster, start it.

Install k8s: <https://kubernetes.io/docs/setup/>

Before start the k8s cluster, make sure the ValidatingAdmissionWebhook is allowed in your cluster.For example:

If you used kubeadm to install the k8s, perhaps you need to add '--enable-admission-plugins=NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook' line to the proper position of your api-server configuration file, which is usually under /etc/kubernetes/manifests.

If you use minikube, you are supposed to stop the minikube and restart it via the following command 
```shell
minikube start --extra-config=apiserver.enable-admission-plugins=NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook
```

4. (optional)Modify config/config/config.json to switch on or switch off some rules, or change the model and policy source.
If you made changes to the config.json, make sure to rebuild the docker image.

5. (optional)Maybe you have made some modifications on this plugin to implement your own customized features.<br/>
If you want to build the docker image for this plugin. execute
```shell
./build.sh
```
And then tag the new image or push it to remote docker image repositories as you wish.

6. Deploy the webhook
Execute the following command to start the admission webhook service
```shell
cd ${workspace}/deployment 
kubectl apply -f webhook_register_internal_step1.yaml
```

After the deployment and the service is correctly started, execute this command to tell k8s that the webhook is the service we started before.
```shell
cd ${workspace}/deployment 
kubectl apply -f webhook_register_internal_step2.yaml
```

7. Apply all example policies
```shell
cd ${workspace}
python3 load_crd.py
```

8. Now you can have a try to test that whether the casbin admission webhook is fully operational now.
For example, 
```shell
cd ${workspace}/example/allowed_repo
kubectl apply -f allowed_repo_rejected.yaml --dry-run=server

```
you will find out that the request is rejected by our webhook.

## 2. Project structure
```
.
├── cmd
│   └── webhook (the entrance point of the webhook)
├── config
│   ├── certificate (storing the private keys and public keys for the webhook)
│   └── config (storing the config file for the webhook)
├── deployments (stroing the scripts)
├── example (storing the example models, policies and configs for each rule)
|           (each subordinate folder contains a model file, an example csv policy file, some example request that will be either approved or 
|           rejected under the example model and policy, and corresponding yaml configuration files for policies in k8s crd format)
│   ├── allowed_repo
│   ├── block_nodeport_service
│   ├── container_resource_limit
│   ├── container_resource_ratio
│   ├── disallowed_tags
│   ├── external_ip
│   ├── https_only
│   ├── image_digest
│   ├── permission
│   ├── replica_limits
│   ├── required_annotations
│   ├── required_labels
│   └── required_probes
├── internal (source code)
│   ├── config (code for resolving the config or )
│   └── webhook (actual webhook code)
│       └── rule (code used for implementing the rules)
├── pkg
│   ├── casbinhelper (some helper functions)
│   └── crdadaptor  (a crd adaptor)
├── test 
│   ├── adaptor (e2e test code for adaptor)
│   └── e2e (e2e test for each rule implemented)
│       ├── pretest
│       ├── testbuild
│       ├── testframework
│       └── testlog
└── util
    └── policy_to_crd (python script converting csv policy files to k8s yaml configuration files)

```
For developers: please fo not change the folder structure of the projectly arbitrarily, because all unit tests and e2e tests use relative file path, and changing folder structure may prevent them from working properly.

## 3 Dive deeper
### 3.1 Reexamination of the certificate.

Let's diver deeper into the certificate. K8s requires that any webhook for k8s must use https, not http.

In our project, we actually generated a self-made CA, and use this CA to sign a self-signed certificate for your webhook. If you already have a certificate signed by a real CA, you can just use them and skip this part. 

AN EXAMPLE of the output of this step now exists in config/certificate. BUT YOU MUST NOT USE IT directly because private key is also exposed in this folder, and a leaked private key makes your connection insecure.

Generate the private key for the fake CA
```
openssl genrsa -des3 -out ca.key 2048
```

Remove the password protection of the private key.
``` 
openssl rsa -in ca.key -out ca.key
```

Generate a private key for the webhook server and remove the password.
```
openssl genrsa -des3 -out server.key 2048
openssl rsa -in server.key  -out server.key 
```

Copy your system's openssl config file for temporary use. You can use `openssl version -a` to find out the location of the config file.

Find the \[req\] paragraph and add the following line: `req_extensions = v3_req`

Find the \[v3_req\] paragraph and add the following line: `subjectAltName = @alt_names`

Append following lines to the file:
```
[alt_names]
DNS.2=casbin-webhook-svc.default.svc
```
The 'casbin-webhook-svc.default.svc' should be replaced with the real service name of your own service (if you decide to modify the service name)

Use the modified config file to generate a certificate request file
```
openssl req -new -nodes -keyout server.key -out server.csr -config openssl.cnf 
```

Use the self-made CA to respond the request and sign the certificate
```
openssl x509 -req -days 3650 -in server.csr -out server.crt -CA ca.crt  -CAkey ca.key -CAcreateserial -extensions v3_req  -extfile openssl.cnf 
```

### 3.2 Reexamination of the k8s yaml configure file.
In 3.1 chapter we mentioned that k8s requires the certificate. So In this part we all going to point out that how to configure this certificate into k8s.

If you have changed the example certificates we provided (ACTUALLY, YOU MUST), then you really should read this chapter.

There are two yaml files in deployment folder: webhook_register_internal_step1.yaml and webhook_register_internal_step2.yaml .

Applying webhook_register_internal_step1.yam will make you create a Deployment using the image of this plugin, and a Service which expose the plugin's ip and port.

Applying webhook_register_internal_step2.yaml will make you tell k8s that which service admission webhooks services are, and when operation on specified resources are being applied, request should be sent to these services.

In webhook_register_internal_step2.yaml you can see 'caBundle' attribute. This is the base64 encoded string of the  certificate of the CA which signed the certificate for your webhook service, because k8s need to know the CA so that they can ensure the certificate your webhook provides is valid. For example, in this project, you can use 'base64 ca.crt' to get the string. It should be noted that anything like '\n',
'\r' must be removed.

### 3.3 Reexamination of policies in crd format and casbin crd adaptor
In order to maintain the consistency of policy on all instances of the k8s cluster, we convert the policies into k8s crd resources, so that the etcd protocol integrated in k8s can be utilized to perfectly achieve the strong cosistency we require. We also implemented a casbin crd adaptor to load policies we generated into casbin enforcer. 

For more information about crd resources see <https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/>.

For the reason why we demand a strong consistency, study more about distributive system.

For more detailed information see <https://github.com/casbin/kubesphere-authz/issues/5>

In util/policy_to_crd folder, we have a python script to automatically generate all crd yaml configuration files you need for every rules based on two templates (crd_template.yaml and policy_template.yaml) and a csv policy file for that rule. For example, we want to generate 
crd yaml files for rule 'allowed_repos'(in example/allowed_repos), just copy allowed_repo.csv into util/policy_to_crd and the execute
```shell
python3 python3 policy_to_crd.py  allowed_repo.csv 
```
you will see some newly generated files
```
crd
├── allowed_repo_definition.yaml
└── policy
    └── policy1.yaml
```

allowed_repo_definition.yaml is the definition of policy crd resource for this rule. In this file you can find that the name of this crd definition is generated based on the filename of csv policy files. Also, we need to point out that this policy is under "policy" namespace. If you want to change it, modify the crd_template.yaml and policy_template.yaml.

Besides, each line of policy is generated as a crd object, defined in the files contained in crd/policy folder. Each object contains ONLY ONE line of policy.

After these files are generated, use kubectl apply -f to make them into effect. 

As for the adaptor, you can find these code in pkg/crdadaptor. If you don't know what a casbin adaptor is, see <https://casbin.org/docs/en/adapters>

### 3.4 Reexamination of the config file.

Currently, the config file locates at config/config/config.json, which is the default position. But acatually we can change the path of the config file. Actually, our webhook program accept exactly one commandline arg, which is the path to the config file. If this arg is not passed, default value is "config/config/config.json".

In this file, each item corresponds to a rule applied to the request, like this:
```json
"ResourceOperationPermission": {
    "available": false,
    "model": "example/permission/permission.conf",
    "policy": "file://example/permission/permission.csv"
},
```
"ResourceOperationPermission" is the name of the rule, (which is also the actual name of the funtion we implemented to apply the rule. Actually, we use reflexion to do this.). 

"available" is either true or false, controlling whether this rule will be taken into effect.

"model" specifies the path of the model conf file.

"policy" specifies where to find the policies. Currently we support 2 kind of policy: traditional csv format policies and our casbin crd adaptor and crd policies(in Chapter 3.3). For csv policy files, the value of "policy" field should follow this syntax: `"file://<path to csv file>"`, for example:  ` "policy": "file://example/permission/permission.csv"`. For crd adaptor, the value of "policy" field should follow this syntax `crd://<path to crd definition file>#<namespace>`, for example: `"policy": "crd://example/allowed_repo/crd/allowed_repo_definition.yaml#policy"`.

allowed_repo_definition.yaml is the definition of policy crd resource for this rule. In this file you can find that the name of this crd definition is generated based on the filename of csv policy files. Also, we need to point out that this policy is under "policy" namespace. If you want to change it, modify the crd_template.yaml and policy_template.yaml.

## 4 For developer
### 4.1 how to debug this webhook?
To toggle a breakpoint, see the output...... You need to run this webhook externally(outside the k8s cluster, but on your local machine). You need to follow the following instructions:

- step 1: modify internal/config/config.go, change DEBUG to true and set  CLIENT_MODE to crdadaptor.EXTERNAL_CLIENT.
```go
//var DEBUG bool = false
var DEBUG bool = true
//var CLIENT_MODE crdadaptor.ClientType = crdadaptor.INTERNAL_CLIENT
var CLIENT_MODE crdadaptor.ClientType = crdadaptor.EXTERNAL_CLIENT
```
- step 2 (optional): modify host. (If you have already configured another registered domain for this webhook, just skip this step.)mpdify your host file to make sure that "webhook.domain.local" is resolved to to your host.

- step3: start your k8s cluster
- step4: register external webhook by applying deployment/webhook_register.yaml to k8s client. If you changed the certificate of your webhook, make sure the caBundle field of webhook_register.yaml is properly modified.(see chapter 3.1 for more information)
```shell
cd ${workspace}/deployment 
kubectl apply -f webhook_register.yaml
```
- step5: modify config/config/config.json if necessary
- step5: under the root directory (`${workspace}`), execute
```shell        
    go run cmd/webhook/main.go
```


