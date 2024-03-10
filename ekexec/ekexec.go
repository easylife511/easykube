package ekexec

import (
	"fmt"
	"k8s_tools/config"
	"k8s_tools/get"
	"k8s_tools/lib"
	"strings"
)

func ProcessExec(action string, args *config.AllArgs, k8s *config.K8s) bool {
	ns := args.Namespace
	pods := args.Pods
	containers := args.Containers
	remainingArgs := args.Remaining
	if len(remainingArgs) == 0 {
		switch action {
		case "exec", "vppctl":
			fmt.Printf("[ERROR] please input which %v cmd you want to execute\n", action)
		}
		return false
	}
	fullNsNameList := get.GetFullResNameList(k8s.GetNamespacesList(), ns)
	for _, nsName := range fullNsNameList { // ns
		fmt.Printf("***************** [  ns: %v ] *****************\n", nsName)

		fullPodNameList := get.GetFullResNameList(k8s.GetPodsList(nsName), pods)
		for index, podName := range fullPodNameList { // pod
			pod := k8s.GetPod(nsName, podName)
			fmt.Printf("================= [ pod: %v ] ================\n", pod.Name)

			findContainer := false
			for _, container := range pod.Status.ContainerStatuses { // container
				if !lib.MatchResName(container.Name, containers) {
					continue
				}
				findContainer = true
				fmt.Printf("----------- [ container: %v ] -----------\n", container.Name)
				switch {
				case container.State.Waiting != nil:
					fmt.Println("[ERROR] container Waiting: ", container.State.Waiting.Reason)
				case container.State.Terminated != nil:
					fmt.Println("[ERROR] container Terminated: ", container.State.Terminated.Reason)
				case container.State.Running != nil: // container running
					switch {
					case action == "exec": // exec
						execCmd := strings.Join(remainingArgs, " ")
						if stdout, ok := k8s.ExecCmd(execCmd, nsName, podName, container.Name); ok {
							fmt.Println(stdout)
						}
					}
				}
			}
			if !findContainer {
				fmt.Println("[ERROR] not found containers:", containers)
			}
			if index < len(fullPodNameList)-1 {
				fmt.Println()
			}
		}
	}
	return true
}
