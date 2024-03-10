package get

import (
	"fmt"
	"k8s_tools/lib"
	"time"

	v1 "k8s.io/api/core/v1"
)

func GetContainerStatus(state *v1.ContainerState) string {
	stateStr := ""
	switch {
	case state.Running != nil:
		stateStr = "Running"
	case state.Waiting != nil:
		stateStr = state.Waiting.Reason + " --- " + state.Waiting.Message
	case state.Terminated != nil:
		if state.Terminated.Reason == "Completed" {
			stateStr = state.Terminated.Reason
		} else {
			stateStr = state.Terminated.Reason + " --- " + state.Terminated.Message
		}
	}
	return stateStr
}

func GetContainer(pod *v1.Pod) ([]string, [][]string) {
	ctrInfoHeader := []string{"CONTAINER_NAME", "INIT", "READY", "STATUS", "IMAGE", "AGE"}
	var (
		ctrName   string
		isInit    string
		isReady   string
		status    string
		imageName string
	)
	var ctrInfoData [][]string
	// Calculate the age
	creationTime := pod.GetCreationTimestamp()
	age := time.Since(creationTime.Time).Round(time.Second)
	ageS := lib.HoursConverter(age.String())
	allContainers := [2][]v1.Container{pod.Spec.InitContainers, pod.Spec.Containers}
	for _, eachTypeContainer := range allContainers {
		for _, container := range eachTypeContainer {
			ctrName = container.Name
			isInit = ""
			isReady = ""
			status = ""
			imageName = container.Image
			ctrInfoData = append(ctrInfoData, []string{ctrName, isInit, isReady, status, imageName, ageS})
		}
	}
	for _, initCtrStat := range pod.Status.InitContainerStatuses {
		for index := range ctrInfoData {
			if initCtrStat.Name == ctrInfoData[index][0] {
				ctrInfoData[index][1] = "yes"
				ctrInfoData[index][2] = fmt.Sprintf("%v", initCtrStat.Ready)
				ctrInfoData[index][3] = GetContainerStatus(&initCtrStat.State)
				break
			}
		}
	}
	for _, CtrStat := range pod.Status.ContainerStatuses {
		for index := range ctrInfoData {
			if CtrStat.Name == ctrInfoData[index][0] {
				ctrInfoData[index][1] = "no"
				ctrInfoData[index][2] = fmt.Sprintf("%v", CtrStat.Ready)
				ctrInfoData[index][3] = GetContainerStatus(&CtrStat.State)
				break
			}
		}
	}
	return ctrInfoHeader, ctrInfoData
}

func GetImage(pod *v1.Pod) ([]string, [][]string) {
	imageHeader := []string{"CONTAINER_NAME", "INIT", "IMAGE", "PULL_POLICY"}
	var (
		ctrName    string
		isInit     string
		imageName  string
		pullPolicy string
	)
	var imageData [][]string
	for _, initCtr := range pod.Spec.InitContainers {
		ctrName = initCtr.Name
		isInit = "yes"
		imageName = initCtr.Image
		pullPolicy = string(initCtr.ImagePullPolicy)
		imageData = append(imageData, []string{ctrName, isInit, imageName, pullPolicy})
	}
	for _, container := range pod.Spec.Containers {
		ctrName = container.Name
		isInit = "no"
		imageName = container.Image
		pullPolicy = string(container.ImagePullPolicy)
		imageData = append(imageData, []string{ctrName, isInit, imageName, pullPolicy})
	}
	return imageHeader, imageData
}

func GetEnv(pod *v1.Pod, container *v1.Container) ([]string, [][]string) {
	envHeader := []string{"ENV_NAME", "ENV_TYPE", "VALUE", "COMMENTS"}
	var (
		envName     string
		envType     string
		envValue    string
		envComments string
	)
	var envData [][]string
	for _, env := range container.Env {
		envName = env.Name
		envType = ""
		envValue = env.Value
		envComments = ""
		evf := env.ValueFrom
		if evf == nil {
			envType = "chart"
		} else {
			switch {
			case evf.FieldRef != nil:
				envType = "fieldRef"
				envComments = evf.FieldRef.FieldPath
				switch envComments {
				case "metadata.name":
					envValue = pod.Name
				case "metadata.namespace":
					envValue = pod.Namespace
				}
			case evf.ResourceFieldRef != nil:
				envType = "resourceFieldRef"
				envComments = evf.ResourceFieldRef.Resource
			case evf.ConfigMapKeyRef != nil:
				envType = "configMapRef"
				envComments = fmt.Sprintf("cm_name: %v, key: %v", evf.ConfigMapKeyRef.Name, evf.ConfigMapKeyRef.Key)
			case evf.SecretKeyRef != nil:
				envType = "SecretKeyRef"
				envComments = fmt.Sprintf("secret_name: %v, key: %v", evf.SecretKeyRef.Name, evf.SecretKeyRef.Key)
			}
		}
		envData = append(envData, []string{envName, envType, envValue, envComments})
	}
	for _, envFrom := range container.EnvFrom {
		envName = ""
		envValue = ""
		envType = "envFrom"
		envComments = ""
		switch {
		case envFrom.Prefix != "":
			envComments = envFrom.Prefix
		case envFrom.ConfigMapRef != nil:
			envComments = "cm_name: " + envFrom.ConfigMapRef.Name
		case envFrom.SecretRef != nil:
			envComments = "secret_name: " + envFrom.SecretRef.Name
		}
		envData = append(envData, []string{envName, envType, envValue, envComments})
	}
	return envHeader, envData
}

func ShowPodsProfile(podsList *v1.PodList, podNameList []string, ctrNameList []string, profile string) bool {
	result := true
	havePod := false
	for _, pod := range podsList.Items {
		name := pod.GetName()
		if !lib.MatchResName(name, podNameList) {
			continue
		}
		havePod = true
		var header []string
		var data [][]string
		fmt.Printf("================= [ pod: %v ] ================\n", name)
		if profile == "image" || profile == "container" { // pod level

			switch profile {
			case "image":
				header, data = GetImage(&pod)
			case "container":
				header, data = GetContainer(&pod)
			}
		} else { // container level
			findContainer := false
			allContainers := [2][]v1.Container{pod.Spec.InitContainers, pod.Spec.Containers}
			for _, eachTypeContainer := range allContainers {
				for _, container := range eachTypeContainer {
					if !lib.MatchResName(container.Name, ctrNameList) {
						continue
					}
					findContainer = true
					var ctrHeader []string
					var ctrData [][]string
					fmt.Printf("----------- [ container: %v ] -----------\n", container.Name)

					switch profile {
					case "env":
						ctrHeader, ctrData = GetEnv(&pod, &container)
					}
					if len(ctrData) != 0 {
						lib.SortDataListWithIndex0(ctrData)
						if !lib.FormatPrint(ctrHeader, ctrData) {
							result = false
						}
					} else {
						fmt.Println("[ERROR] not find:", profile)
					}
					fmt.Println()
				}
			}
			if !findContainer {
				fmt.Println("[ERROR] not found containers:", ctrNameList)
				result = false
			}
		} // get info end
		if len(data) != 0 {
			lib.SortDataListWithIndex0(data)
			if !lib.FormatPrint(header, data) {
				result = false
			}
		}
		fmt.Println()
	}
	if !havePod {
		fmt.Println("[ERROR] not found pods:", podNameList)
		result = false
	}
	return result
}
