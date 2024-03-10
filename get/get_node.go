package get

import (
	"fmt"
	"k8s_tools/lib"
	"os/exec"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
)

func ShowNodesList(nodesList *v1.NodeList, nodeNameList []string, user string) bool {
	nodesHeader := []string{"NAME", "IP", "CPU", "MEMORY", "HPsize", "HPtotal", "HPfree", "STATUS", "AGE"}
	var data [][]string
	for _, node := range nodesList.Items {
		nodeStatus := node.Status
		name := node.GetName()
		if !lib.MatchResName(name, nodeNameList) {
			continue
		}
		address := fmt.Sprintf("%v", nodeStatus.Addresses[0].Address)
		cpu := fmt.Sprintf("%v", nodeStatus.Capacity.Cpu())
		memory := fmt.Sprintf("%v", nodeStatus.Capacity.Memory())
		var (
			HpSize      string
			HpTotal     string
			HpRemaining string
		)
		// Hugepagesize
		sshCmd := "ssh " + user + "@" + name + " cat /proc/meminfo | grep -i Hugepagesize"
		out, err := exec.Command("bash", "-c", sshCmd).Output()
		if err != nil {
			fmt.Printf("[ERROR] exec cmd [%v] failed, err: %v \n", sshCmd, err)
			return false
		}
		if strings.Contains(string(out), "2048 kB") {
			HpSize = "2M"
		} else if strings.Contains(string(out), "1048576 kB") {
			HpSize = "1G"
		} else {
			HpSize = string(out)
		}
		// HugePages_Total
		out, _ = exec.Command("bash", "-c", strings.Replace(sshCmd, "Hugepagesize", "HugePages_Total", 1)).Output()
		if HpSize == "2M" {
			HpNum := strings.TrimSpace(strings.Split(string(out), ":")[1]) // string: 10240
			if HpNum == "0" {
				HpTotal = "0"
			} else {
				Num, _ := strconv.Atoi(HpNum) // int: 10240
				HpTotal = strconv.Itoa(Num/512) + "Gi"
			}
		} else if HpSize == "1G" {
			HpTotal = strings.TrimSpace(strings.Split(string(out), ":")[1]) + "Gi"
		} else {
			HpTotal = string(out)
		}
		// HugePages_Free
		if HpTotal == "0" {
			HpRemaining = "0"
		} else {
			out, _ = exec.Command("bash", "-c", strings.Replace(sshCmd, "Hugepagesize", "HugePages_Free", 1)).Output()
			if HpSize == "2M" {
				HpNum := strings.TrimSpace(strings.Split(string(out), ":")[1]) // string: 10240
				Num, _ := strconv.ParseFloat(HpNum, 64)                        // int: 10240
				HpRemaining = fmt.Sprintf("%0.2f", Num/512) + "Gi"
			} else if HpSize == "1G" {
				HpRemaining = strings.TrimSpace(strings.Split(string(out), ":")[1]) + "Gi"
			} else {
				HpRemaining = string(out)
			}
		}
		status := fmt.Sprintf("%v", nodeStatus.Conditions[len(nodeStatus.Conditions)-1].Type)
		// Calculate the age of the node
		nsCreationTime := node.GetCreationTimestamp()
		age := time.Since(nsCreationTime.Time).Round(time.Second)
		ageS := lib.HoursConverter(age.String())
		data = append(data, []string{name, address, cpu, Ki2Mi(memory), HpSize, HpTotal, HpRemaining, status, ageS})
	}
	if len(data) == 0 {
		fmt.Println("[ERROR] not found node:", nodeNameList)
		return false
	}
	lib.SortDataListWithIndex0(data)
	return lib.FormatPrint(nodesHeader, data)
}

func ShowNodesListSriov(nodeList *v1.NodeList, nodeNameList []string, podList *v1.PodList) bool {
	nodesSriovHeader := []string{"NAME", "Capacity", "Allocatable", "Requests", "Limits"}
	dataMap := map[string][][]string{}
	for _, node := range nodeList.Items {
		nodeStatus := node.Status
		name := node.GetName()
		if !lib.MatchResName(name, nodeNameList) {
			continue
		}
		address := fmt.Sprintf("%v", nodeStatus.Addresses[0].Address)
		status := fmt.Sprintf("%v", nodeStatus.Conditions[len(nodeStatus.Conditions)-1].Type)
		nodeInfo := fmt.Sprintf("%v | %v | %v", name, address, status)
		dataMap[nodeInfo] = [][]string{}
		for k, v := range nodeStatus.Capacity {
			if strings.Contains(k.String(), "sriov") || strings.Contains(k.String(), "dpdk") {
				dataMap[nodeInfo] = append(dataMap[nodeInfo], []string{k.String(), v.String()})
			}
		}
		for k, v := range nodeStatus.Allocatable {
			if !strings.Contains(k.String(), "sriov") && !strings.Contains(k.String(), "dpdk") {
				continue
			}
			for index := range dataMap[nodeInfo] {
				if dataMap[nodeInfo][index][0] == k.String() {
					dataMap[nodeInfo][index] = append(dataMap[nodeInfo][index], v.String(), "0", "0")
				}
			}
		}
		for _, pod := range podList.Items {
			if pod.Spec.NodeName != name {
				continue
			}
			for _, container := range pod.Spec.Containers {
				for k, v := range container.Resources.Requests {
					if !strings.Contains(k.String(), "sriov") && !strings.Contains(k.String(), "dpdk") {
						continue
					}
					for index := range dataMap[nodeInfo] {
						if dataMap[nodeInfo][index][0] == k.String() {
							oldNum, _ := strconv.Atoi(dataMap[nodeInfo][index][3])
							newNum, _ := strconv.Atoi(v.String())
							dataMap[nodeInfo][index][3] = strconv.Itoa(oldNum + newNum)
						}
					}
				}
				for k, v := range container.Resources.Limits {
					if !strings.Contains(k.String(), "sriov") && !strings.Contains(k.String(), "dpdk") {
						continue
					}
					for index := range dataMap[nodeInfo] {
						if dataMap[nodeInfo][index][0] == k.String() {
							oldNum, _ := strconv.Atoi(dataMap[nodeInfo][index][4])
							newNum, _ := strconv.Atoi(v.String())
							dataMap[nodeInfo][index][4] = strconv.Itoa(oldNum + newNum)
						}
					}
				}
			}
		}
	}
	if len(dataMap) == 0 {
		fmt.Println("[ERROR] not found node:", nodeNameList)
		return false
	}
	keys := lib.SortMapKeys(dataMap)
	for _, key := range keys {
		fmt.Printf("---------- [ %v ] ----------\n", key)
		lib.SortDataListWithIndex0(dataMap[key])
		if len(dataMap[key]) == 0 {
			fmt.Println("[Warning] not found sriov resource")
		} else if !lib.FormatPrint(nodesSriovHeader, dataMap[key]) {
			return false
		}
		fmt.Println()
	}
	return true
}
