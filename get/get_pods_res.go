package get

import (
	"fmt"
	"k8s_tools/lib"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func ShowPodsResource(podsList *v1.PodList, podMetricsList *v1beta1.PodMetricsList, podNameList []string, ctrNameList []string) bool {
	resourceHeader := []string{"NAME", "Requests", "Limits", "Using"}
	havePod := false
	for _, pod := range podsList.Items {
		if !lib.MatchResName(pod.Name, podNameList) {
			continue
		}
		havePod = true
		fmt.Printf("================= [ pod: %v ] ================\n", pod.Name)
		dataMap := map[string][][]string{}
		for _, container := range pod.Spec.Containers {
			containerName := container.Name
			if !lib.MatchResName(containerName, ctrNameList) {
				continue
			}
			dataMap[containerName] = [][]string{}
			if len(container.Resources.Requests) != 0 {
				for k, v := range container.Resources.Requests {
					dataMap[containerName] = append(dataMap[containerName], []string{k.String(), v.String()})
				}
			} else {
				dataMap[containerName] = append(dataMap[containerName], []string{"cpu", "-"})
				dataMap[containerName] = append(dataMap[containerName], []string{"memory", "-"})
			}
			// update Limits
			if len(container.Resources.Limits) != 0 {
				for k, v := range container.Resources.Limits {
					for index := range dataMap[containerName] {
						if dataMap[containerName][index][0] == k.String() {
							dataMap[containerName][index] = append(dataMap[containerName][index], v.String())
						}
					}
				}
			} else {
				for index := range dataMap[containerName] {
					dataMap[containerName][index] = append(dataMap[containerName][index], "-")
				}
			}
			// update metrics
			for _, podInPodMetricsList := range podMetricsList.Items {
				if podInPodMetricsList.Name == pod.Name {
					for _, containerInPod := range podInPodMetricsList.Containers {
						cup := containerInPod.Usage.Cpu()
						cores := resource.NewMilliQuantity(cup.MilliValue(), cup.Format).String()
						memory := containerInPod.Usage.Memory()
						memorySize := resource.NewQuantity(memory.Value(), memory.Format).String()
						memorySize = Ki2Mi(memorySize)
						if containerInPod.Name == containerName {
							for index := range dataMap[containerName] {
								switch dataMap[containerName][index][0] {
								case "cpu":
									dataMap[containerName][index] = append(dataMap[containerName][index], cores)
								case "memory":
									dataMap[containerName][index] = append(dataMap[containerName][index], memorySize)
								default:
									dataMap[containerName][index] = append(dataMap[containerName][index], "-")
								}
							}
						}
					}
					break
				}
			}
		}
		if len(dataMap) == 0 {
			fmt.Println("[ERROR] not found container:", ctrNameList)
			return false
		}
		keys := lib.SortMapKeys(dataMap)
		for _, key := range keys {
			if len(dataMap[key]) == 0 {
				fmt.Println("[ERROR] not found resource under container:", key)
				return false
			}
			fmt.Printf("----------- [ container: %v ] -----------\n", key)
			lib.SortDataListWithIndex0(dataMap[key])
			if !lib.FormatPrint(resourceHeader, dataMap[key]) {
				return false
			}
			fmt.Println()
		}
	}
	if !havePod {
		fmt.Println("[ERROR] not found pods:", podNameList)
		return false
	}
	return true
}
