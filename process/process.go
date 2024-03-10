package process

import (
	"fmt"
	"k8s_tools/config"
	"k8s_tools/ekexec"
	"k8s_tools/ekpatch"
	"k8s_tools/get"
	"strings"
)

func StartK8sClient(kubeConfigFilePath string) *config.K8s {
	k8s := &config.K8s{}
	k8s.GetK8sconfig(kubeConfigFilePath)
	k8s.GetK8sClient()
	return k8s
}

func Launch() bool {
	args := config.LoadArgs()
	if args == nil {
		return false
	}
	switch args.Action {
	case "version":
		args.ShowVersion()
	case "config":
		if len(args.Remaining) == 0 {
			args.ShowBaseArgs()
			if args.UpdateConfig && !args.UpdateBaseArgs() {
				return false
			}
		} else {
			fmt.Println("[ERROR] not support: ek config", strings.Join(args.Remaining, " "))
			return false
		}
	case "get":
		if !get.ProcessGet(args, StartK8sClient(args.KubeConfigFilePath)) {
			return false
		}
	case "exec", "vppctl":
		if !ekexec.ProcessExec(args.Action, args, StartK8sClient(args.KubeConfigFilePath)) {
			return false
		}
	case "netcmd":
		if !ekexec.ProcessNetCmd(args, StartK8sClient(args.KubeConfigFilePath)) {
			return false
		}
	case "add", "remove", "replace":
		if !ekpatch.ProcessPatch(args.Action, args, StartK8sClient(args.KubeConfigFilePath)) {
			return false
		}
	default:
		fmt.Printf("[ERROR] invalid args: ek %v, please run 'ek' to see usage\n", args.Action)
		return false
	}
	return true
}
