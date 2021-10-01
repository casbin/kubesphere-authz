package rule

import (
	"encoding/json"
	"fmt"
	"log"

	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
)

func (g *Rules) RequiredLabel(review *v1.AdmissionReview, modelUrl string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	adaptor,err:=getAdaptorObject(policy)
	if err != nil {
		log.Printf("RequiredLabel: %s %s:%s rejected due to error:%s", resourceKind, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	model,err:=getModelObject(modelUrl)
	if err != nil {
		log.Printf("RequiredLabel: %s %s:%s rejected due to error:%s", resourceKind, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)

	if err != nil {
		log.Printf("RequiredLabel: %s %s:%s rejected due to error:%s", resourceKind, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("RequiredLabel: %s %s:%s approved", resourceKind, review.Request.Namespace, review.Request.Name)
		return nil
	}

	var object map[string]interface{}
	if err := json.Unmarshal(review.Request.Object.Raw, &object); err != nil {
		log.Printf("RequiredLabel: %s %s:%s rejected due to error:%s", resourceKind, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	metadata, ok := object["metadata"].(map[string]interface{})
	if !ok {
		log.Printf("RequiredLabel: %s %s:%s rejected due to error: type assertion failure for metadata", resourceKind, review.Request.Namespace, review.Request.Name)
		return fmt.Errorf("RequiredLabel: %s %s:%s rejected due to error: type assertion failure for metadata", resourceKind, review.Request.Namespace, review.Request.Name)
	}

	annotations, ok := metadata["labels"].(map[string]interface{})
	if !ok {
		log.Printf("RequiredLabel: %s %s:%s rejected due to error: type assertion failure for labels", resourceKind, review.Request.Namespace, review.Request.Name)
		return fmt.Errorf("RequiredLabel: %s %s:%s rejected due to error: type assertion failure for labels", resourceKind, review.Request.Namespace, review.Request.Name)
	}
	var passedCount = 0
	for k, v := range annotations {
		ok, err := enforcer.Enforce(k, v)
		if err != nil {
			log.Printf("RequiredLabel: %s %s:%s rejected due to error: %s", resourceKind, review.Request.Namespace, review.Request.Name, err.Error())
			return err
		}
		if ok {
			passedCount++
		}
	}
	requiredCount := len(enforcer.GetModel()["p"]["p"].Policy)
	fmt.Println(passedCount, requiredCount)
	if passedCount < requiredCount {
		log.Printf("RequiredLabel: %s %s:%s rejected due to: %d labels required but %d qualified labels detected", resourceKind, review.Request.Namespace, review.Request.Name, passedCount, requiredCount)
		return fmt.Errorf("RequiredLabel: %s %s:%s rejected due to: %d labels required but %d qualified labels detected", resourceKind, review.Request.Namespace, review.Request.Name, passedCount, requiredCount)
	}
	log.Printf("RequiredLabel: %s %s:%s approved", resourceKind, review.Request.Namespace, review.Request.Name)
	return nil

}
