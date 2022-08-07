package main

import (
	"context"
	"encoding/json"
	"io/ioutil"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getConfigmap() (kickerConfs []KickerConf) {
	// de locatie en naam van de configmap moet ergens bepaald worden
	// volgens een conventie
	fileBytes, err := ioutil.ReadFile("./test.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &kickerConfs)
	if err != nil {
		panic(err)
	}

	return kickerConfs
}

func getPods() *v1.PodList {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	// TODO zorg dat de namespace uit een flag komt
	pods, err := clientset.CoreV1().Pods("anti-aff-test").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// fmt.Println(reflect.TypeOf(pods))
	return pods
}

func getESXNodeofOCPNode(OCPNode string) string {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	// TODO zorg dat de namespace opgehaald wordt voor waar de podkicker draait
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var ESXNode string

	for _, n := range nodes.Items {
		labels := n.ObjectMeta.Labels
		if OCPNode == n.Name {
			// fmt.Printf("VM/OCP node: %v on ESX node: %v\n", n.Name, labels["esx-node"])
			ESXNode = labels["esx-node"]
		}
	}
	return ESXNode
}

func KickPod(Pod string) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	// TODO zorg dat de namespace uit een flag komt
	Err := clientset.CoreV1().Pods("anti-aff-test").Delete(context.TODO(), Pod, metav1.DeleteOptions{})
	if Err != nil {
		panic(Err.Error())
	}
	// fmt.Println(reflect.TypeOf(pods))
}
