package get

import (
	"fmt"
	"k8s_tools/config"
	"k8s_tools/lib"
	"strings"
)

func ShowMeta(metaType string, metaMap map[string]string, resType string, resName string) {
	fmt.Printf("----------------- [ %v ] %v: %v -----------------\n", strings.ToUpper(metaType), resType, resName)
	var separator string
	if metaType == "labels" {
		separator = "="
	} else if metaType == "annotations" {
		separator = " ---> "
	}
	if len(metaMap) == 0 {
		fmt.Println("[WARNING] not find", metaType)
	} else {
		keys := lib.SortMapKeys(metaMap)
		for _, key := range keys {
			fmt.Printf("%v%v%v\n", key, separator, metaMap[key])
		}
	}
}

func ShowResMeta(metaType string, k8s *config.K8s, resType string, resFullNameList []string, podNameList []string) {
	for index, name := range resFullNameList {
		switch resType {
		case "node":
			node := k8s.GetNode(name)
			if metaType == "labels" {
				ShowMeta(metaType, node.Labels, "node", name)
			} else if metaType == "annotations" {
				ShowMeta(metaType, node.Annotations, "node", name)
			}
		case "ns":
			ns := k8s.GetNamespace(name)
			if metaType == "labels" {
				ShowMeta(metaType, ns.Labels, "ns", name)
			} else if metaType == "annotations" {
				ShowMeta(metaType, ns.Annotations, "ns", name)
			}
		case "pod":
			fmt.Printf("***************** [  ns: %v ] *****************\n", name)
			podFullNamelist := GetFullResNameList(k8s.GetPodsList(name), podNameList)
			for i, podName := range podFullNamelist {
				pod := k8s.GetPod(name, podName)
				if metaType == "labels" {
					ShowMeta(metaType, pod.Labels, "pod", podName)
				} else if metaType == "annotations" {
					ShowMeta(metaType, pod.Annotations, "pod", podName)
				}
				if i < len(podFullNamelist)-1 {
					fmt.Println()
				}
			}
		}
		if index < len(resFullNameList)-1 {
			fmt.Println()
		}
	}
}
