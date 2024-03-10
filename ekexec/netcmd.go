package ekexec

import (
	"fmt"
	"k8s_tools/config"
	"k8s_tools/get"
	"k8s_tools/lib"
	"os/exec"
	"strings"
)

func ProcessNetCmd(args *config.AllArgs, k8s *config.K8s) bool {
	node := args.Node
	user := args.User
	ns := args.Namespace
	pods := args.Pods
	containers := args.Containers
	remainingArgs := args.Remaining
	if len(remainingArgs) == 0 {
		fmt.Println("[ERROR] please input which net cmd you want to execute")
		return false
	}
	fullNodeNameList := get.GetFullResNameList(k8s.GetNodesList(), node)
	for _, nodeName := range fullNodeNameList { // node
		findPod := false
		fmt.Printf("################# [node: %v] #################\n", nodeName)
		fullNsNameList := get.GetFullResNameList(k8s.GetNamespacesList(), ns)
		for _, nsName := range fullNsNameList { // ns
			printNs := true
			podsList := k8s.GetPodsList(nsName)
			for _, pod := range podsList.Items { // pod
				if pod.Spec.NodeName != nodeName {
					continue
				}
				if !lib.MatchResName(pod.Name, pods) {
					continue
				}
				findPod = true
				if printNs {
					fmt.Printf("***************** [  ns: %v ] *****************\n", nsName)
					printNs = false
				}
				fmt.Printf("================= [ pod: %v ] ================\n", pod.Name)
				sshHead := "ssh " + user + "@" + nodeName
				runtime := "docker"
				err := exec.Command("bash", "-c", sshHead+" which docker").Run()
				if err != nil {
					runtime = "crictl"
				}
				findContainer := false
				for _, container := range pod.Status.ContainerStatuses { // container
					if !lib.MatchResName(container.Name, containers) {
						continue
					}
					findContainer = true
					fmt.Printf("----------- [ container: %v ] -----------\n", container.Name)
					switch {
					case container.State.Running != nil:
						var sshCmd string
						if runtime == "docker" {
							sshCmd = fmt.Sprintf("%v %v ps | grep %v_%v | cut -f1 -d ' '", sshHead, runtime, container.Name, pod.Name)
						} else {
							sshCmd = fmt.Sprintf("%v sudo %v ps | grep ' %v ' | grep %v | cut -f1 -d ' '", sshHead, runtime, container.Name, pod.Name)
						}
						out, err := exec.Command("bash", "-c", sshCmd).Output()
						if err != nil {
							fmt.Printf("[ERROR] exec ssh cmd [%v] failed, err: %v\v", sshCmd, err)
							return false
						}
						if len(out) == 0 {
							fmt.Println("[ERROR] not get container ID with ssh cmd: ", sshCmd)
							return false
						}
						containerId := strings.TrimSpace(string(out))
						sshCmd = fmt.Sprintf("%v sudo %v inspect %v | grep -iE pid.?: | cut -d : -f2 | cut -d , -f1", sshHead, runtime, containerId)
						out, err = exec.Command("bash", "-c", sshCmd).Output()
						if err != nil {
							fmt.Printf("[ERROR] exec ssh cmd [%v] failed, err: %v\v", sshCmd, err)
							return false
						}
						pid := strings.TrimSpace(string(out))
						sshCmd = fmt.Sprintf("%v sudo nsenter -t %v -n %v", sshHead, pid, strings.Join(args.Remaining, " "))
						out, err = exec.Command("bash", "-c", sshCmd).Output()
						if err != nil {
							fmt.Printf("[ERROR] exec ssh cmd [%v] failed, err: %v\v", sshCmd, err)
							return false
						}
						fmt.Println(string(out)) // nsenter output
					case container.State.Waiting != nil:
						fmt.Println("[ERROR] container Waiting: ", container.State.Waiting.Reason)
					case container.State.Terminated != nil:
						fmt.Println("[ERROR] container Terminated: ", container.State.Terminated.Reason)
					}
					fmt.Println()
				}
				if !findContainer {
					fmt.Printf("[ERROR] not find Containers: %v under Pods: %v\n\n", args.Containers, pod.Name)
				}
			}
		}
		if !findPod {
			fmt.Printf("[ERROR] not find pods: %v under ns: %v in node: %v\n", pods, ns, nodeName)
		}
		fmt.Println()
	}
	return true
}
