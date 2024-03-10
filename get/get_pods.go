package get

import (
	"fmt"
	"k8s_tools/lib"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
)

func Ki2Mi(memorySize string) string {
	if strings.Contains(memorySize, "Ki") {
		num, err := strconv.Atoi(strings.Replace(memorySize, "Ki", "", 1))
		if err != nil {
			fmt.Println("[ERROR] transform memorySize to int failed, err: ", err)
			return memorySize
		}
		memorySize = strconv.Itoa(num/1024) + "Mi"
	}
	return memorySize
}

func ShowPodsList(podsList *v1.PodList, podNameList []string) bool {
	podsHeader := []string{"NAME", "READY", "STATUS", "RESTARTS", "NODE", "AGE"}
	var data [][]string
	for _, pod := range podsList.Items {
		name := pod.GetName()
		if !lib.MatchResName(name, podNameList) {
			continue
		}
		// Get the status of each of the pods
		podStatus := pod.Status
		var containerRestarts int32
		var containerReady int
		var totalContainers int
		// If a pod has multiple containers, get the status from all
		for index := range pod.Spec.Containers {
			containerRestarts += podStatus.ContainerStatuses[index].RestartCount
			if podStatus.ContainerStatuses[index].Ready {
				containerReady++
			}
			totalContainers++
		}
		// Get the values from the pod status
		ready := fmt.Sprintf("%v/%v", containerReady, totalContainers)
		status := fmt.Sprintf("%v", podStatus.Phase)
		// status := fmt.Sprintf("%v", podStatus.Message)
		restarts := fmt.Sprintf("%v", containerRestarts)
		// Calculate the age
		creationTime := pod.GetCreationTimestamp()
		age := time.Since(creationTime.Time).Round(time.Second)
		ageS := lib.HoursConverter(age.String())
		node := pod.Spec.NodeName
		data = append(data, []string{name, ready, status, restarts, node, ageS})
	}
	if len(data) == 0 {
		fmt.Println("[ERROR] not found pods:", podNameList)
		return false
	}
	lib.SortDataListWithIndex0(data)
	return lib.FormatPrint(podsHeader, data)
}
