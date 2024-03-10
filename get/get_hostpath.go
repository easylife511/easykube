package get

import (
	"fmt"
	"k8s_tools/lib"
	"time"

	v1 "k8s.io/api/core/v1"
)

func ShowHostPathList(podsList *v1.PodList) bool {
	hostPathHeader := []string{"VOL_NAME", "POD_NAME", "NODE", "PATH", "TYPE", "AGE"}
	var data [][]string
	for _, pod := range podsList.Items {
		podName := pod.GetName()
		nodeName := pod.Spec.NodeName
		// age
		creationTime := pod.GetCreationTimestamp()
		age := time.Since(creationTime.Time).Round(time.Second)
		ageS := lib.HoursConverter(age.String())
		for _, vol := range pod.Spec.Volumes {
			volHostPath := vol.VolumeSource.HostPath
			if volHostPath != nil {
				data = append(data, []string{vol.Name, podName, nodeName, volHostPath.Path, string(*volHostPath.Type), ageS})
			}
		}
	}
	if len(data) == 0 {
		fmt.Println("[ERROR] not found HostPath")
		return false
	}
	lib.SortDataListWithIndex0(data)
	return lib.FormatPrint(hostPathHeader, data)
}
