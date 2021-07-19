package proxyapi

import (
	"errors"
	"fmt"
	"log"

	"github.com/pnocera/dashgo/pkg/age"
	"github.com/pnocera/dashgo/pkg/instances"
	v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	daprEnabledAnnotation     = "dapr.io/enabled"
	daprIDAnnotation          = "dapr.io/app-id"
	daprPushedAppIDAnnotation = "dapr.proxy/app-id"
	daprIsProxyAnnotation     = "dapr.proxy/isproxy"
	daprProxySuffix           = "_proxied"
)

type ProxyAPI interface {
	GetProxy(scope string) instances.Instance
	Proxy(scope string, appID string) bool
	UnProxy(scope string) bool
}

type proxyapi struct {
	kubeClient kubernetes.Interface
}

// NewInstances returns an Instances instance
func NewProxyApi(kubeClient *kubernetes.Clientset) ProxyAPI {
	i := proxyapi{}

	i.kubeClient = kubeClient

	return &i
}

// Gets the proxy deployment
func (i *proxyapi) GetProxy(scope string) instances.Instance {

	d, err := i.getDeploymentByAnnotation(scope, daprIsProxyAnnotation, "true")

	if err != nil {
		log.Println(err)
		return instances.Instance{}
	}

	id := d.Spec.Template.Annotations[daprIDAnnotation]

	return instances.Instance{
		AppID:            id,
		HTTPPort:         3500,
		GRPCPort:         50001,
		Command:          "",
		Age:              age.GetAge(d.CreationTimestamp.Time),
		Created:          d.GetCreationTimestamp().String(),
		PID:              -1,
		Replicas:         int(*d.Spec.Replicas),
		SupportsDeletion: false,
		SupportsLogs:     true,
		Address:          fmt.Sprintf("%s-dapr:80", id),
		Status:           fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, d.Status.Replicas),
		Labels:           "app:" + d.Labels["app"],
		Selector:         "app:" + d.Labels["app"],
		Config:           d.Spec.Template.Annotations["dapr.io/config"],
	}

}

func (i *proxyapi) Proxy(scope string, appID string) bool {
	//get proxy and switch to appID
	inst := i.GetProxy(scope).AppID

	if inst == "" {
		log.Println("Could not get proxy")
		return false
	}

	err := i.patchDeployment(scope, appID, fmt.Sprintf("%s%s", appID, daprProxySuffix))
	if err != nil {
		log.Println(err)
		return false
	}

	err = i.patchDeployment(scope, inst, appID)
	if err != nil {

		log.Println(err)
		return false
	}

	return true
}

func (i *proxyapi) UnProxy(scope string) bool {

	appID := i.GetProxy(scope).AppID
	if appID == "" {
		log.Println("Could not get proxy")
		return false
	}

	err := i.unpatchDeployment(scope)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (i *proxyapi) patchDeployment(scope string, appID string, newID string) error {
	d, err := i.getDeploymentByAnnotation(scope, daprIDAnnotation, appID)
	if err != nil {
		return err
	}

	deploymentsClient := i.kubeClient.AppsV1().Deployments(scope)
	d.Spec.Template.Annotations[daprIDAnnotation] = newID
	d.Spec.Template.Annotations[daprPushedAppIDAnnotation] = appID

	_, err = deploymentsClient.Update(&d)
	if err != nil {
		return err
	}

	return nil
}

func (i *proxyapi) unpatchDeployment(scope string) error {
	deploymentsClient := i.kubeClient.AppsV1().Deployments("")

	prox, err := i.getDeploymentByAnnotation(scope, daprIsProxyAnnotation, "true")
	if err != nil {
		return err
	}

	appid := prox.Spec.Template.Annotations[daprIDAnnotation]
	prox.Spec.Template.Annotations[daprIDAnnotation] = prox.Spec.Template.Annotations[daprPushedAppIDAnnotation]
	prox.Spec.Template.Annotations[daprPushedAppIDAnnotation] = ""

	_, err = deploymentsClient.Update(&prox)
	if err != nil {
		return err
	}

	d, err := i.getDeploymentByAnnotation(scope, daprIDAnnotation, fmt.Sprintf("%s%s", appid, daprProxySuffix))
	if err != nil {
		return err
	}

	d.Spec.Template.Annotations[daprIDAnnotation] = d.Spec.Template.Annotations[daprPushedAppIDAnnotation]
	d.Spec.Template.Annotations[daprPushedAppIDAnnotation] = ""

	_, err = deploymentsClient.Update(&d)
	if err != nil {
		return err
	}

	return nil
}

// func (i *proxyapi) getDeployment0(labelselector string) (v1.Deployment, error) {
// 	options := meta_v1.ListOptions{
// 		LabelSelector: labelselector,
// 	}
// 	resp, err := i.kubeClient.AppsV1().Deployments("").List((options))
// 	if err != nil {
// 		log.Println(err)
// 		return v1.Deployment{}, err
// 	}
// 	for _, d := range resp.Items {
// 		return d, nil
// 	}
// 	return v1.Deployment{}, errors.New("no matches found")
// }

func (i *proxyapi) getDeploymentByAnnotation(scope string, annotation string, value string) (v1.Deployment, error) {
	options := meta_v1.ListOptions{}
	resp, err := i.kubeClient.AppsV1().Deployments(scope).List((options))
	if err != nil {
		log.Println(err)
		return v1.Deployment{}, err
	}
	for _, d := range resp.Items {
		id := d.Spec.Template.Annotations[annotation]
		if id == value {
			return d, nil
		}

	}
	return v1.Deployment{}, errors.New("no matches found")
}
