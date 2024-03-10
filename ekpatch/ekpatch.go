package ekpatch

import (
	"encoding/json"
	"fmt"
	"k8s_tools/config"
	"k8s_tools/get"
	"strings"
)

type PatchData struct {
	Op    string `json:"op"`   // add, remove, replace
	Path  string `json:"path"` // /metadata/labels/test_Label ,  /spec/template/spec/containers/0/image
	Value string `json:"value"`
}

func ProcessPathArgsError(ekAction string, argsList []string) bool {
	fmt.Printf("[ERROR] invalid args: ek %v %v\n", ekAction, strings.Join(argsList, " "))
	return false
}

func ProcessPatch(action string, args *config.AllArgs, k8s *config.K8s) bool {
	node := args.Node
	remainingArgs := args.Remaining

	switch {
	case len(remainingArgs) == 0:
		fmt.Printf("[ERROR] plesae specify which resource want to [%v], please run 'ek' to see [ek %v] usage\n", action, action)
		return false
	case len(remainingArgs) != 3: // only support 3 remaining args, e.g: ek add [node label aaa=bbb]
		return ProcessPathArgsError(action, remainingArgs)
	}
	// ek add [node label aaa=bbb]
	switch remainingArgs[0] { // node
	case "node":
		switch remainingArgs[1] { // label
		case "lab", "label", "labels":
			content := remainingArgs[2] // aaa=bbb
			if action != "remove" && !strings.Contains(content, "=") {
				fmt.Printf("[ERROR] action is: %v, but not find '=' in :%v\n", action, content)
				return false
			}
			var key string
			var value string
			if action == "remove" {
				key = content
			} else {
				key = strings.Split(content, "=")[0]
				value = strings.Split(content, "=")[1]
			}
			labelPath := []PatchData{{
				Op:    action,
				Path:  "/metadata/labels/" + key,
				Value: value,
			}}
			payloadBytes, _ := json.Marshal(labelPath)
			for _, nodeName := range get.GetFullResNameList(k8s.GetNodesList(), node) {
				fmt.Printf("################# [node: %v] #################\n", nodeName)
				if action == "remove" {
					if _, ok := k8s.GetNode(nodeName).Labels[key]; !ok {
						fmt.Printf("[ERROR] not find key: %v\n", key)
						return false
					}
				}
				newNode := k8s.PatchNode(nodeName, payloadBytes) // update patch label
				if action == "remove" {
					if _, ok := newNode.Labels[key]; !ok {
						fmt.Printf("%v %v label [ %v ]successful\n", action, remainingArgs[0], key)
					} else {
						fmt.Printf("%v %v label [ %v ]failed\n", action, remainingArgs[0], key)
					}
					continue
				}
				if newNode.Labels[key] == value {
					fmt.Printf("%v %v label successful, key: %v value: %v\n", action, remainingArgs[0], key, value)
				} else {
					fmt.Printf("[ERROR] label key: %v expect value: %v, actual value: %v\n", key, value, newNode.Labels[key])
					return false
				}
			}
		default:
			return ProcessPathArgsError(action, remainingArgs)
		}
	default:
		return ProcessPathArgsError(action, remainingArgs)
	}

	return true
}
