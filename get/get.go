package get

import (
	"fmt"
	"k8s_tools/config"
	"k8s_tools/lib"
	"strings"

	v1 "k8s.io/api/core/v1"
)

func ProcessGetArgsError(argsList []string) bool {
	fmt.Println("[ERROR] invalid args: ek get", strings.Join(argsList, " "))
	return false
}

func GetFullResNameList(v1ListObj interface{}, resNameList []string) []string {
	fullResNameList := []string{}
	var resType string
	switch resList := v1ListObj.(type) {
	case *v1.NodeList:
		for _, res := range resList.Items {
			fullResNameList = lib.AppendResName(fullResNameList, res.Name, resNameList)
		}
		resType = "node"
	case *v1.NamespaceList:
		for _, res := range resList.Items {
			fullResNameList = lib.AppendResName(fullResNameList, res.Name, resNameList)
		}
		resType = "ns"
	case *v1.PodList:
		for _, res := range resList.Items {
			fullResNameList = lib.AppendResName(fullResNameList, res.Name, resNameList)
		}
		resType = "pod"
	}
	if len(fullResNameList) == 0 {
		fmt.Printf("[ERROR] not found %v: %v\n\n", resType, resNameList)
	}
	return fullResNameList
}

func ShowResWithMultipleNs(k8s *config.K8s, ns []string, resType string, filter1 []string, filter2 []string) bool {
	result := true
	fullNsNameList := GetFullResNameList(k8s.GetNamespacesList(), ns)
	for index, nsName := range fullNsNameList {
		fmt.Printf("***************** [  ns: %v ] *****************\n", nsName)
		switch resType {
		case "pod":
			if !ShowPodsList(k8s.GetPodsList(nsName), filter1) {
				result = false
			}
		case "res":
			if !ShowPodsResource(k8s.GetPodsList(nsName), k8s.GetPodMetricsList(nsName), filter1, filter2) {
				result = false
			}
		case "pvc":
			if !ShowPVCList(k8s.GetPVCList(nsName), k8s.GetPVList(), k8s.GetPodsList(nsName), filter1) {
				result = false
			}
		case "hostpath":
			if !ShowHostPathList(k8s.GetPodsList(nsName)) {
				result = false
			}
		case "vol":
			if !ShowVol(k8s.GetPodsList(nsName), filter1, filter2, k8s.GetPVCList(nsName), k8s.GetPVList()) {
				result = false
			}
		case "image":
			if !ShowPodsProfile(k8s.GetPodsList(nsName), filter1, filter2, "image") {
				result = false
			}
		case "env":
			if !ShowPodsProfile(k8s.GetPodsList(nsName), filter1, filter2, "env") {
				result = false
			}
		case "container":
			if !ShowPodsProfile(k8s.GetPodsList(nsName), filter1, filter2, "container") {
				result = false
			}
		}
		if index < len(fullNsNameList)-1 {
			fmt.Println()
		}
	}
	return result
}

func ProcessGet(args *config.AllArgs, k8s *config.K8s) bool {
	node := args.Node
	user := args.User
	ns := args.Namespace
	pods := args.Pods
	containers := args.Containers
	remainingArgs := args.Remaining
	if len(remainingArgs) == 0 {
		fmt.Println("[ERROR] plesae specify which resource want to show, please run 'ek' to see [ek get] usage")
		return false
	}
	var emptyList []string
	switch {
	case len(remainingArgs) == 1:
		switch remainingArgs[0] {
		case "node":
			return ShowNodesList(k8s.GetNodesList(), node, user)
		case "ns":
			return ShowNamespacesList(k8s.GetNamespacesList(), ns)
		case "po", "pod", "pods":
			// return ShowPods(k8s, ns, pods)
			return ShowResWithMultipleNs(k8s, ns, "pod", pods, emptyList) // ShowPodsList
		case "res", "resource":
			return ShowResWithMultipleNs(k8s, ns, "res", pods, containers) // ShowPodsResource
		case "pvc", "PVC", "Pvc":
			return ShowResWithMultipleNs(k8s, ns, "pvc", emptyList, emptyList) // ShowPVCList
		case "hostpath", "HostPath", "hostPath", "Hostpath":
			return ShowResWithMultipleNs(k8s, ns, "hostpath", emptyList, emptyList) // ShowHostPathList
		case "vol", "volume", "volumes", "VOL":
			return ShowResWithMultipleNs(k8s, ns, "vol", pods, containers) // ShowVol
		case "img", "image", "images":
			return ShowResWithMultipleNs(k8s, ns, "image", pods, emptyList) // ShowPodsProfile
		case "env", "ENV":
			return ShowResWithMultipleNs(k8s, ns, "env", pods, containers) // ShowPodsProfile
		case "con", "cont", "container":
			return ShowResWithMultipleNs(k8s, ns, "container", pods, containers) // ShowPodsProfile
		default:
			return ProcessGetArgsError(remainingArgs)
		}
	case len(remainingArgs) == 2:
		result := true
		switch remainingArgs[0] {
		case "node":
			switch remainingArgs[1] {
			case "sriov":
				result = ShowNodesListSriov(k8s.GetNodesList(), node, k8s.GetPodsList(""))
			case "lab", "label", "labels":
				ShowResMeta("labels", k8s, "node", GetFullResNameList(k8s.GetNodesList(), node), emptyList)
			case "ann", "anno", "annotation", "annotations":
				ShowResMeta("annotations", k8s, "node", GetFullResNameList(k8s.GetNodesList(), node), emptyList)
			default:
				return ProcessGetArgsError(remainingArgs)
			}
			return result
		case "ns":
			switch remainingArgs[1] {
			case "lab", "label", "labels":
				ShowResMeta("labels", k8s, "ns", GetFullResNameList(k8s.GetNamespacesList(), ns), emptyList)
			case "ann", "anno", "annotation", "annotations":
				ShowResMeta("annotations", k8s, "ns", GetFullResNameList(k8s.GetNamespacesList(), ns), emptyList)
			default:
				return ProcessGetArgsError(remainingArgs)
			}
			return result
		case "po", "pod", "pods":
			switch remainingArgs[1] {
			case "lab", "label", "labels":
				ShowResMeta("labels", k8s, "pod", GetFullResNameList(k8s.GetNamespacesList(), ns), pods)
			case "ann", "anno", "annotation", "annotations":
				ShowResMeta("annotations", k8s, "pod", GetFullResNameList(k8s.GetNamespacesList(), ns), pods)
			default:
				return ProcessGetArgsError(remainingArgs)
			}
			return result
		case "pvc", "PVC", "Pvc":
			return ShowResWithMultipleNs(k8s, ns, "pvc", []string{remainingArgs[1]}, emptyList) // ShowPVCList
		default:
			return ProcessGetArgsError(remainingArgs)
		}
	default:
		return ProcessGetArgsError(remainingArgs)
	}
}
