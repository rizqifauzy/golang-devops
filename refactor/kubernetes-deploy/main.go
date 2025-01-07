package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var (
		client           *kubernetes.Clientset
		err              error
		deploymentLabels map[string]string
		expectedPods     int32
	)
	ctx := context.Background()

	if client, err = getClient(); err != nil {
		fmt.Printf("Get Client Error : %s", err)
		os.Exit(1)
	}
	//fmt.Printf("Client: %+v\n", client)
	if deploymentLabels, expectedPods, err = deploy(ctx, client); err != nil {
		fmt.Printf("Deployment Error : %s", err)
		os.Exit(1)
	}

	if err = waitForPods(ctx, client, deploymentLabels, expectedPods); err != nil {
		fmt.Printf("Error : %s", err)
		os.Exit(1)
	}
	fmt.Printf("Deploy finished. Did a deploy with labels: %+v\n", deploymentLabels)

}

func getClient() (*kubernetes.Clientset, error) {
	//var kubeconfig *string
	//if home := homedir.HomeDir(); home != "" {
	//kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	//flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func deploy(ctx context.Context, client *kubernetes.Clientset) (map[string]string, int32, error) {
	var deployment *v1.Deployment

	appFile, err := os.ReadFile("app.yaml")
	if err != nil {
		return nil, 0, fmt.Errorf("ReadFile error: %s", err)
	}
	// decode := scheme.Codecs.UniversalDeserializer().Decode
	// obj, groupVersionKind, err := decode(appFile, nil, nil)
	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode(appFile, nil, nil)

	switch obj.(type) {
	case *v1.Deployment:
		deployment = obj.(*v1.Deployment)
	default:
		return nil, 0, fmt.Errorf("Unrecognized type: %s\n", groupVersionKind)
	}

	_, err = client.AppsV1().Deployments("default").Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		deploymentResponse, err := client.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return nil, 0, fmt.Errorf("deployment error: %s", err)
		}
		return deploymentResponse.Spec.Template.Labels, *deploymentResponse.Spec.Replicas, nil
	} else if err != nil && !errors.IsNotFound(err) {
		return nil, 0, fmt.Errorf("deployment get error: %s", err)
	}

	deploymentResponse, err := client.AppsV1().Deployments("default").Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, 0, fmt.Errorf("Update deployment error: %s", err)
	}
	return deploymentResponse.Spec.Template.Labels, *deploymentResponse.Spec.Replicas, nil
}

func waitForPods(ctx context.Context, client *kubernetes.Clientset, deploymentLabels map[string]string, expectedPods int32) error {
	for {
		validatedLabels, err := labels.ValidatedSelectorFromSet(deploymentLabels)
		if err != nil {
			return fmt.Errorf("ValidatedSelectorFromSet: %s", err)
		}

		podList, err := client.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
			LabelSelector: validatedLabels.String(),
		})
		if err != nil {
			return fmt.Errorf("list pods error: %s", err)
		}
		podsRunning := 0
		for _, pods := range podList.Items {
			if pods.Status.Phase == "Running" {
				podsRunning++
			}
		}
		fmt.Printf("waiting for pods to become ready (running %d / %d)\n", podsRunning, len(podList.Items))
		if podsRunning > 0 && podsRunning == len(podList.Items) && podsRunning == int(expectedPods) {
			break
		}
		time.Sleep(5 * time.Second)

	}

	return nil
}
