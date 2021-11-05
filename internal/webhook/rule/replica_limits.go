package rule

import (
	"encoding/json"
	"fmt"
	"ksauth/pkg/casbinhelper"
	"log"

	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
)

func (g *Rules) ReplicaLimits(review *v1.AdmissionReview, modelUrl string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	if resourceKind != "deployments" {
		return nil
	}
	model, adaptor, err := getModelAndPolicyObject(modelUrl, policy)
	if err != nil {
		log.Printf("ReplicaLimits: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if model == nil {
		log.Printf("ReplicaLimits: approved due to enable=true")
		return nil
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)

	if err != nil {
		log.Printf("ReplicaLimits: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer.AddFunction("parseInt", casbinhelper.ParseInt)
	enforcer.AddFunction("parseFloat", casbinhelper.ParseFloat)
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("ReplicaLimits: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var deploymentObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &deploymentObject); err != nil {
		log.Printf("ReplicaLimits: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	var replicas int32 = 1
	if deploymentObject.Spec.Replicas != nil {

		replicas = *(deploymentObject.Spec.Replicas)
	}
	fmt.Println(replicas)
	ok, err := enforcer.Enforce(
		review.Request.Namespace,
		int(replicas),
	)
	if err != nil {
		log.Printf("ReplicaLimits: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if !ok {
		log.Printf("ReplicaLimits(%s %s::%s): prohibited replica number %d", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, replicas)
		return fmt.Errorf("ReplicaLimits(%s %s::%s): prohibited replica number %d", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name, replicas)
	}
	log.Printf("ReplicaLimits(%s %s::%s): approved", review.Request.Resource.Resource, deploymentObject.Namespace, deploymentObject.Name)
	return nil
}
