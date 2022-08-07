package main

import (
	"fmt"
	"reflect"
	"strconv"
)

// struct om de json te parsen naar een go type
type KickerConf struct {
	Name           string            `json:"name"`
	Minpods        string            `json:"minPodsRunning"`
	MaxpodsESXnode string            `json:"maxPodsPerESXNode"`
	Matchlabels    map[string]string `json:"matchLabels"`
}

type Maps map[string]string

func MapInnerJoin(m1, m2 Maps) map[string]string {
	m3 := make(map[string]string)
	for key, value := range m1 {
		if m2[key] == value {
			m3[key] = value
		}

	}
	return m3
}

func main() {
	pods := getPods()
	kickerConfs := getConfigmap()

	// het doel is om pods te selecteren die gebruiker opgeeft dmv de labels
	// daarna moeten kijken we naar de verdeling van pods over ESX nodes
	// als het teveel pods zijn voor een bepaalde ESX node dan kicken
	// het teveel aan pods. We houden het minimum aan operationele pods in act

	// loop over alle configuraties in mounted json file
	// en loop over alle pod labels van alle pods in een namespace
	// kijk of de pods label maps overeenkomen met de opgegeven labels in de configuratie
	// als dat zo is voeg ze toe aan countMatchedLabels

	// kijk hoeveel pods er per ESX node draaien als dit er teveel zijn dan kick de pod

	// misschien alle pods in de namespace waar dit programma draait constant labelen
	// met de ESX host waar ze op draaien. De kicker hoeft dan alleen nog maar te kicken
	// als het er teveel zijn.

	// Hoe plaatsen we een nieuwe pod, nadat de oude gekicked is, op de juiste ESX node?
	// met een nodeselector?

	// config2Pod := make(map[string]map[string]string)
	// config2Pod[ESXNode] = make(map[string]string)

	// Als een pod niet gekicked mag worden omdat er al teveel te gekicked zijn
	// ivm minPodsRunning, wat dan? Wachten tot de volgende loop van dit programma?

	config2Pod := make(map[int]map[string]string)
	for i := range kickerConfs {
		config2Pod[i] = make(map[string]string)
	}
	for _, pod := range pods.Items {
		// var podName string = pod.Name
		ESXNode := getESXNodeofOCPNode(pod.Spec.NodeName)
		// fmt.Printf("pod %v running on OCP node: %v and ESX node: %v \n", podName, pod.Spec.NodeName, ESXNode)

		for i := range kickerConfs {
			// innerjoin de pod labels en de labels in de configuration
			innerJoinedMap := MapInnerJoin(pod.Labels, kickerConfs[i].Matchlabels)
			// als de innerjoin gelijk is aan configuration gebruik dan deze configuration
			// voor deze pod
			eq := reflect.DeepEqual(innerJoinedMap, kickerConfs[i].Matchlabels)
			// fmt.Printf("Config \"%v\" = %v to pod %v\n", kickerConfs[i].Name, eq, pod.Name)

			// create a list of pods
			if eq {
				config2Pod[i][pod.Name] = ESXNode
			}
		}
	}
	var countPodsonNodes int
	for i, confmap := range config2Pod {
		maxPodESXNode, _ := strconv.Atoi(kickerConfs[i].MaxpodsESXnode)
		countPodsonNodes = 0
		for podname, _ := range confmap {
			countPodsonNodes += 1
			if countPodsonNodes > maxPodESXNode {
				fmt.Printf("Kicking pod %v\n", podname)
				KickPod(podname)
			}
			if countPodsonNodes > maxPodESXNode {
				numberOfPodsToKick := countPodsonNodes - maxPodESXNode
				fmt.Printf("I kicked %d pods according to Configuration: \"%v\"\n", numberOfPodsToKick, kickerConfs[i].Name)
			}
		}
	}
}