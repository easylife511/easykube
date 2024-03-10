package get

import (
	"fmt"
	"k8s_tools/lib"
	"reflect"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
)

func ShowVol(podsList *v1.PodList, podNameList []string, ctrNameList []string, pvcList *v1.PersistentVolumeClaimList, pvList *v1.PersistentVolumeList) bool {
	podVolHeader := []string{"VOL_NAME", "VOL_TYPE", "COMMENTS", "AGE"}
	havePod := false
	for _, pod := range podsList.Items {
		if !lib.MatchResName(pod.Name, podNameList) {
			continue
		}
		havePod = true
		var data [][]string
		nodeName := pod.Spec.NodeName
		fmt.Printf("================= [ pod: %v ] ================\n", pod.Name)
		// age
		creationTime := pod.GetCreationTimestamp()
		age := time.Since(creationTime.Time).Round(time.Second)
		ageS := lib.HoursConverter(age.String())
		for _, vol := range pod.Spec.Volumes {
			if strings.Contains(vol.Name, "token") || strings.Contains(vol.Name, "kube-api-access") {
				continue
			}
			vVolSource := reflect.ValueOf(vol.VolumeSource)
			tVolSource := reflect.TypeOf(vol.VolumeSource)
			fildeNum := vVolSource.NumField()
			for i := 0; i < fildeNum; i++ {
				var volType string
				var volComments string
				if vVolSource.Field(i).IsNil() {
					continue
				}
				if vVolSource.Field(i).Kind() != reflect.Ptr {
					continue
				}
				volType = tVolSource.Field(i).Name
				switch iVolSource := vVolSource.Field(i).Interface().(type) {
				case *v1.HostPathVolumeSource:
					volComments = nodeName + ":" + iVolSource.Path
				case *v1.EmptyDirVolumeSource:
					if string(iVolSource.Medium) != "" {
						volComments = string(iVolSource.Medium) + ": " + iVolSource.SizeLimit.String()
					}
				case *v1.SecretVolumeSource:
					volComments = "secret_name: " + iVolSource.SecretName
				case *v1.ConfigMapVolumeSource:
					volComments = "cm_name: " + iVolSource.LocalObjectReference.Name
				case *v1.PersistentVolumeClaimVolumeSource:
					volType = "PVC"
					volComments = "pvc_name: " + iVolSource.ClaimName
				case *v1.DownwardAPIVolumeSource:
					for _, item := range iVolSource.Items {
						var fieldStr string
						if item.FieldRef != nil {
							fieldStr = item.FieldRef.FieldPath
						} else if item.ResourceFieldRef != nil {
							fieldStr = item.ResourceFieldRef.Resource
						}
						volComments += fmt.Sprintf("%v:%v,", item.Path, fieldStr)
					}
					volComments = strings.TrimSuffix(volComments, ",")
				}
				data = append(data, []string{vol.Name, volType, volComments, ageS})
			}
		}
		if len(data) == 0 {
			fmt.Println("[ERROR] not found volume")
		} else {
			lib.SortDataListWithIndex0(data)
			lib.FormatPrint(podVolHeader, data)
		}
		// ---------------------------------------------- container ----------------------------------
		containerVolHeader := []string{"MOUNT_PATH", "VOL_NAME", "PV_NAME", "PV_TYPE", "COMMENTS", "SUB_PATH", "SUB_Expr"}
		findContainer := false
		for _, container := range pod.Spec.Containers {
			if !lib.MatchResName(container.Name, ctrNameList) {
				continue
			}
			findContainer = true
			var containerData [][]string
			haveSubPath := false
			haveSubPathExpr := false
			fmt.Printf("----------- [ container: %v ] -----------\n", container.Name)
			for _, volMount := range container.VolumeMounts {
				volName := volMount.Name
				if strings.Contains(volName, "token") || strings.Contains(volName, "kube-api-access") {
					continue
				}
				var pvComments string
				var pvcName string
				var pvName string
				var pvType string
				mountPath := volMount.MountPath
				subPath := volMount.SubPath
				if subPath != "" {
					haveSubPath = true
				}
				subPathExpr := volMount.SubPathExpr
				if subPathExpr != "" {
					haveSubPathExpr = true
				}
				// update container vol comments
				for index := range data {
					switch {
					case data[index][0] == volName && data[index][1] == "HostPath":
						pvComments = data[index][2]
						goto breakHere
					case data[index][0] == volName && data[index][1] == "PVC":
						for _, pvc := range pvcList.Items {
							pvcName = strings.Replace(data[index][2], "pvc_name: ", "", 1)
							if pvc.Name != pvcName {
								continue
							}
							pvName = pvc.Spec.VolumeName
							pvType, pvComments = GetPvInfo(pvList, pvName)
							if pvType == "HostPath" || pvType == "Local" {
								pvComments = nodeName + ":" + pvComments
							}
							break
						}
						goto breakHere
					}
				}
			breakHere:
				containerData = append(containerData, []string{mountPath, volName, pvName, pvType, pvComments, subPath, subPathExpr})
			}
			if len(containerData) == 0 {
				fmt.Println("[ERROR] not found mount volume for container")
			} else {
				// reduce SubPathExpr or SubPath
				var reduceHeader []string
				if !haveSubPathExpr {
					reduceHeader = containerVolHeader[0 : len(containerVolHeader)-1]
					for index := range containerData {
						containerData[index] = containerData[index][0 : len(containerData[index])-1]
					}
					if !haveSubPath {
						reduceHeader = reduceHeader[0 : len(reduceHeader)-1]
						for index := range containerData {
							containerData[index] = containerData[index][0 : len(containerData[index])-1]
						}
					}
				} else {
					reduceHeader = containerVolHeader
				}
				lib.SortDataListWithIndex0(containerData)
				lib.FormatPrint(reduceHeader, containerData)
				fmt.Println()
			}
		}
		if !findContainer {
			fmt.Printf("[ERROR] not find Containers: %v under Pods: %v\n", ctrNameList, pod.Name)
		}
		fmt.Println()
	}
	if !havePod {
		fmt.Println("[ERROR] not found pods:", podNameList)
		return false
	}
	return true
}
