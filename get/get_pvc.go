package get

import (
	"fmt"
	"k8s_tools/lib"
	"reflect"
	"time"

	v1 "k8s.io/api/core/v1"
)

func GetPvInfo(pvList *v1.PersistentVolumeList, pvName string) (string, string) {
	var pvType string
	var comments string
	for _, pv := range pvList.Items {
		if pv.Name != pvName {
			continue
		}
		vPvSource := reflect.ValueOf(pv.Spec.PersistentVolumeSource)
		tPvSource := reflect.TypeOf(pv.Spec.PersistentVolumeSource)
		fildeNum := vPvSource.NumField()
		for i := 0; i < fildeNum; i++ {
			if vPvSource.Field(i).IsNil() {
				continue
			}
			pvType = tPvSource.Field(i).Name
			if vPvSource.Field(i).Kind() == reflect.Ptr {
				switch iPvSource := vPvSource.Field(i).Interface().(type) {
				case *v1.HostPathVolumeSource:
					comments = iPvSource.Path
				case *v1.LocalVolumeSource:
					comments = iPvSource.Path
				case *v1.CSIPersistentVolumeSource:
					comments = iPvSource.VolumeHandle
				}
			}
			break // only one type is using for pv, other type is nil
		}
		break // pv name is unique in k8s
	}
	return pvType, comments
}

func ShowPVCList(pvcList *v1.PersistentVolumeClaimList, pvList *v1.PersistentVolumeList, podsList *v1.PodList, pvcNameList []string) bool {
	pvcHeader := []string{"NAME", "STATUS", "PV_NAME", "CAPA", "STORAGECLASS", "PV_TYPE", "POD_NAME", "COMMENTS", "AGE"}
	var data [][]string
	for _, pvc := range pvcList.Items {
		name := pvc.GetName()
		if !lib.MatchResName(name, pvcNameList) {
			continue
		}
		pvcStatus := pvc.Status
		status := fmt.Sprintf("%v", pvcStatus.Phase)
		pvName := pvc.Spec.VolumeName
		capa := fmt.Sprintf("%v", pvcStatus.Capacity.Name("storage", "BinarySI"))
		sc := *pvc.Spec.StorageClassName
		var podName string
		pvType, comments := GetPvInfo(pvList, pvc.Spec.VolumeName)
		if pvType == "HostPath" || pvType == "Local" {
			for _, pod := range podsList.Items {
				for _, vol := range pod.Spec.Volumes {
					volPVC := vol.VolumeSource.PersistentVolumeClaim
					if volPVC != nil && volPVC.ClaimName == name {
						comments = pod.Spec.NodeName + ":" + comments
						podName = pod.Name
						goto breakHere
					}
				}
			}
		}
	breakHere:
		// age
		creationTime := pvc.GetCreationTimestamp()
		age := time.Since(creationTime.Time).Round(time.Second)
		ageS := lib.HoursConverter(age.String())
		data = append(data, []string{name, status, pvName, capa, sc, pvType, podName, comments, ageS})
	}
	if len(data) == 0 {
		fmt.Println("[ERROR] not found PVC: ", pvcNameList)
		return false
	}
	lib.SortDataListWithIndex0(data)
	return lib.FormatPrint(pvcHeader, data)
}
