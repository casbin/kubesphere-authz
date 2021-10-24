package rule

import (
	"encoding/json"
	"fmt"
	"log"

	casbin "github.com/casbin/casbin/v2"
	v1 "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

func (g *Rules) RequiredProbes(review *v1.AdmissionReview, model string, policy string) error {
	var resourceKind = review.Request.Resource.Resource
	switch resourceKind {
	case "pods":
		return g.requiredProbesForPod(review, model, policy)
	case "deployments":
		return g.requiredProbesForDeployment(review, model, policy)
	default:
		return nil
	}
}

func (g *Rules) requiredProbesForPod(review *v1.AdmissionReview, modelUrl string, policy string) error {
	model, adaptor, err := getModelAndPolicyObject(modelUrl, policy)
	if err != nil {
		log.Printf("RequiredProbes: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if model == nil {
		log.Printf("RequiredProbes approved due to enable==true")
		return nil
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)
	if err != nil {
		log.Printf("RequiredProbes: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("RequiredProbes: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var podObject core.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &podObject); err != nil {
		log.Printf("RequiredProbes: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	for _, container := range podObject.Spec.Containers {
		err = g.checkProbe(review, &container, enforcer)
		if err != nil {
			return err
		}

	}
	log.Printf("RequiredProbes:%s  %s:%s approved", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
	return nil
}

func (g *Rules) requiredProbesForDeployment(review *v1.AdmissionReview, modelUrl string, policy string) error {
	model, adaptor, err := getModelAndPolicyObject(modelUrl, policy)

	if err != nil {
		log.Printf("RequiredProbes: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if model == nil {
		log.Printf("RequiredProbes approved due to enable==true")
		return nil
	}
	enforcer, err := casbin.NewEnforcer(model, adaptor)
	if err != nil {
		log.Printf("RequiredProbes: pod %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if review.Request.Operation == "DELETE" {
		//delete operation have no docker image to check
		log.Printf("RequiredProbes: pod %s:%s approved", review.Request.Namespace, review.Request.Name)
		return nil
	}
	var deploymentObject app.Deployment
	if err := json.Unmarshal(review.Request.Object.Raw, &deploymentObject); err != nil {
		log.Printf("RequiredProbes: deployment %s:%s rejected due to error:%s", review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}

	for _, container := range deploymentObject.Spec.Template.Spec.Containers {
		err = g.checkProbe(review, &container, enforcer)
		if err != nil {
			return err
		}

	}
	log.Printf("RequiredProbes:%s %s:%s approved", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
	return nil
}

func (g *Rules) checkProbe(review *v1.AdmissionReview, container *core.Container, enforcer *casbin.Enforcer) error {
	//only livenessProbe,readinessProbe startupProbe
	//first: livenessProbe:
	ok, err := enforcer.Enforce("livenessProbe")
	if err != nil {
		log.Printf("RequiredProbes: %s %s:%s rejected due to error:%s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	// livnessProbe is required
	if ok {
		if container.LivenessProbe == nil {
			log.Printf("RequiredProbes: %s %s:%s rejected due to error: LivnessProbe is Required but not found", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
			return fmt.Errorf("RequiredProbes: %s %s:%s rejected due to error: LivnessProbe is Required but not found", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
		}
		err := g.checkProbeType(review, "livenessProbe", container.LivenessProbe, enforcer)
		if err != nil {
			return err
		}
	}

	//second: readinessProbe:
	ok, err = enforcer.Enforce("readinessProbe")
	if err != nil {
		log.Printf("RequiredProbes: %s %s:%s rejected due to error:%s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if ok {
		if container.ReadinessProbe == nil {
			log.Printf("RequiredProbes: %s %s:%s rejected due to error: readinessProbe is Required but not found", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
			return fmt.Errorf("RequiredProbes: %s %s:%s rejected due to error: readinessProbe is Required but not found", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
		}
		err := g.checkProbeType(review, "readinessProbe", container.ReadinessProbe, enforcer)
		if err != nil {
			return err
		}

	}

	//third startupProbe
	ok, err = enforcer.Enforce("startupProbe")
	if err != nil {
		log.Printf("RequiredProbes: %s %s:%s rejected due to error:%s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if ok {
		if container.StartupProbe == nil {
			log.Printf("RequiredProbes: %s %s:%s rejected due to error: startupProbe is Required but not found", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
			return fmt.Errorf("RequiredProbes: %s %s:%s rejected due to error: startupProbe is Required but not found", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name)
		}
		err := g.checkProbeType(review, "startupProbe", container.StartupProbe, enforcer)
		if err != nil {
			return err
		}

	}
	return nil
}

func (g *Rules) checkProbeType(review *v1.AdmissionReview, probeName string, probe *core.Probe, enforcer *casbin.Enforcer) error {
	//there are only 3 allowed type of probe: "tcpSocket", "httpGet", "exec"
	enforceContext := casbin.NewEnforceContext("2")
	//tcpSocket
	ok, err := enforcer.Enforce(enforceContext, probeName, "tcpSocket")
	if err != nil {
		log.Printf("RequiredProbes: %s %s:%s rejected due to error:%s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if ok && probe.TCPSocket != nil {
		return nil
	}
	//httpGet
	ok, err = enforcer.Enforce(enforceContext, probeName, "httpGet")
	if err != nil {
		log.Printf("RequiredProbes: %s %s:%s rejected due to error:%s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if ok && probe.HTTPGet != nil {
		return nil
	}
	//exec
	ok, err = enforcer.Enforce(enforceContext, probeName, "exec")
	if err != nil {
		log.Printf("RequiredProbes: %s %s:%s rejected due to error:%s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, err.Error())
		return err
	}
	if ok && probe.Exec != nil {
		return nil
	}
	log.Printf("RequiredProbes: %s %s:%s rejected due to: failed to find a allowed type for probe %s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, probeName)
	return fmt.Errorf("RequiredProbes: %s %s:%s rejected due to: failed to find a allowed type for probe %s", review.Request.Resource.Resource, review.Request.Namespace, review.Request.Name, probeName)

}
