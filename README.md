# casbin-kubesphere-auth
Casbin-kubesphere-auth is a plugin which apply several security authentication check on kubesphere via [casbin](https://casbin.org/docs/en/overview). This plugin support the following function:
- check whether it is legal to apply an operation on a k8s resource (e.g. perform a 'DELETE' operation on a certain deployment). Illegal request will be intercepted and rejected.
- check whether the docker image you use on any pod/deployment is trusted. If not, request will be intercepted and rejected.

Functions above are implemented via admission webhook of k8s. Webhook service can be built as a docker file, and to support kubesphere better, this webhook service is also packed as a helm application, which can be uploaded to local kubesphere market and easily deployed.

You can use this plugin in kubesphere and a raw k8s.

## Overview
## Structure for this project:
- casbin-kubesphere/ this folder is used for create a helm package.
- k8sconfig/ this folder contains necessary yaml configuration files to deploy this webhook.
- webhook/ this folder contains real code for webhook service. This folder also include a Dockerfile, which means the docker image of the service should be buit base on this folder. 
    - webhook/casbinconfig contains casbin model and policies.
    - webhook/certificate contains certificates, private keys and public keys **ONLY FOR EXAMPLE!** You **MUST NOT** use this keys in any environment except test environment because **EVERYONE CAN GET THE PRIVATE KEY IN THIS FOLDER**. You should generate a set of your own keys via the method metioned by the tutorial below.
    - the others are go codes implementing this service.

## Get Started: How to make this plugin work.
### step 1: have k8s and kubesphere installed.
Install k8s: <https://kubernetes.io/docs/setup/>

Install kubesphere: <https://kubesphere.com.cn/en/docs/quick-start/minimal-kubesphere-on-k8s/>

(When install kubesphere, please choose 'minimal install on kubernets' instead of 'all in one for linux')

### step 2 Enable the ValidatingWebhookConfiguration of your k8s
Briefly, you should add configuration '--enable-admission-plugins=NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook' to k8s apiserver. Specific method to add this configuration varies depending on how you installed your k8s. 

For example, if you use minikube, you are supposed to stop the minikube and restart it via the following command 
```shell
minikube start --extra-config=apiserver.enable-admission-plugins=NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook
```
Or if you used kubeadm to install the k8s, perhaps you need to add '--enable-admission-plugins=NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook' line to the proper position of your api-server configuration file, which is usually under /etc/kubernetes/manifests.

Or perhaps you may be able to use 'kube-apiserver' directly
......


You may find more information from k8s doc. See <https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/>

### Generate a set of certificate and private keys.
K8s requires that any webhook for k8s must use https, not http.

In this step, we shall generate a self-made CA, and use this CA to sign a self-signed certificate for your webhook. If you already have a certificate signed by a real CA, you can skip this step. 

AN EXAMPLE of the output of this step now exists in webhook/certificate. BUT YOU MUST NOT USE IT directly because private key is also exposed in this folder, and a leaked private key makes your connection insecure.

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


### Reexamine yaml configs
To make the webhook into effort, we need to apply some yaml configs to k8s.

There are two yaml files in k8sconfig folder: webhook2.yaml and webhook3.yaml.

Applying webhook3.yaml will make you create a Deployment using the image of this plugin, and a Service which expose the plugin's ip and port.

Applying webhook2.yaml will make you tell k8s that which service admission webhooks services are, and when operation on specified resources are being applied, request should be sent to these services.

In webhook2.yaml you can see 'caBundle' attribute. This is the base64 encoded string of the  certificate of the CA which signed the certificate for your webhook service, because k8s need to know the CA so that they can ensure the certificate your webhook provides is valid. For example, in this project, you can use 'base64 ca.crt' to get the string. It should be noted that anything like '\n',
'\r' must be removed.

In webhook3.yaml, you can find that we used local docker image and set the image policy to 'never pull images from remote'. for the convenience of running this tiny project. If your corporation has a private docker repo, you should modify this part.
### Reexamine casbin configs
casbin model and policies are stored in webhook/casbinconfig. There are 2 sets of model&policy. image_model.conf and image_policy controls whether a image is trusted, and permission.conf and permission.csv control whether an operation on a resource can be applied

### Reexamine webhook configs
In webhook/webhookconfig you can see config.json. Through this file, you can turn on or turn off a check rule, **or modify the parmeter of casbin Enforcer**. For example, you can modify them so that casbin's enforcer can use something like redis or mysql adaptors  so that ploicies can be modified dynamically. In this project we use files as policy to make an exapmle, which is not recommeded because if so, you have to shut down ther service and rebuild the docker image every time you make some changes to policy.  

### Pack this service into helm app 
If you want to make this service became a infrastructure of your organization, you should ack this service into helm app, which is the only format the kubesphere app store supports. 

If you haven't installed helm yet, see <https://helm.sh/docs/> and have it installed.



Run
```
helm create casbin-kubesphere
```
You can see a folder 'casbin-kubesphere' is created. In this project, we have already runned this command and you can see there's already a folder casbin-kubesphere there.

Remove everything under casbin-kubesphere/templates except deployment.yaml. Combine the contents of k8sconfig/webhook2.yaml and k8sconfig/webhook3.yaml ans copy it into deployment.yaml

Wipe out everything in values.yaml

Run 
```
helm package casbin-kubesphere
```
and you will see a file called 'casbin-kubesphere-0.1.0.tgz' created. This is the helm package for this plugin.

To install this plugin in k8s directly, run helm install casbin-kubesphere-0.1.0.tgz.

To upload this app to kubesphere so that everyone can use it, see <https://kubesphere.com.cn/en/docs/workspace-administration/upload-helm-based-application/>







