package config

import (
	"flag"
	"fmt"
	"k8s_tools/lib"
	"os"
	"strings"
)

// filter define and undefine args
func ManuallyProcessArgs(args []string) ([]string, []string) {
	known := []string{}
	unknown := []string{}
	//################################ Define input args  ###################################
	kvargs := []string{"-h", "--help", "--kubeconfig", "-n", "--namespace", "--show-labels", "-p", "-c", "-d", "-u"}
	kvargs_single := []string{"-h", "--help", "--show-labels"}
	//#######################################################################################
	knownArg := func(a string) bool {
		for _, pre := range kvargs {
			if strings.HasPrefix(a, pre+"=") {
				return true
			}
		}
		return false
	}
	isKnown := func(v string) string {
		for _, i := range kvargs {
			if i == v {
				return v
			}
		}
		return ""
	}

	for i := 0; i < len(args); i++ {
		switch a := args[i]; a {
		case "--debug":
			known = append(known, a)
		case isKnown(a):
			known = append(known, a)
			if lib.IsElementExistsInSlice(kvargs_single, a) {
				continue
			}
			i++
			if i < len(args) {
				known = append(known, args[i])
			}
		default:
			if knownArg(a) {
				known = append(known, a)
				continue
			}
			unknown = append(unknown, a)
		}
	}
	return known, unknown
}

func CombineStrArg(argFromInput string, argFromConfig string, argDefaultValue string) string {
	if argFromInput == "" && argFromConfig != "" {
		argFromInput = argFromConfig
	}
	if argFromInput == "reset" { // clean KubeConfigFilePath name in config file
		argFromInput = argDefaultValue
	}
	return argFromInput
}
func CombineListArg(argFromInput []string, argFromConfig []string) []string {
	if len(argFromInput) == 0 && len(argFromConfig) != 0 {
		argFromInput = argFromConfig
	}
	if len(argFromInput) != 0 && argFromInput[0] == "reset" { // clean node name in config file
		argFromInput = []string{}
	}
	return argFromInput
}

type AllArgs struct {
	BaseArgs
	UpdateConfig bool
	Action       string
	Remaining    []string
}

func (aa *AllArgs) SplitStrArg(strArg string, listArg []string) []string {
	if strArg != "" {
		aa.UpdateConfig = true
		listArg = append(listArg, strings.Split(strArg, ",")...)
	}
	return listArg
}
func (aa *AllArgs) ParseArgs(args []string) {
	defineArgs, undefineArgs := ManuallyProcessArgs(args)
	var nodesNameStr string
	var nsNameStr string
	var podsNameStr string
	var ctrsNameStr string
	flag.StringVar(&aa.KubeConfigFilePath, "kubeconfig", "", "--kubeconfig same as kubectl")
	flag.StringVar(&nsNameStr, "n", "", "namespace name")
	flag.StringVar(&nodesNameStr, "d", "", "node name")
	flag.StringVar(&aa.User, "u", "", "user name for ssh worker node")
	flag.StringVar(&podsNameStr, "p", "", "part of pods name")
	flag.StringVar(&ctrsNameStr, "c", "", "part of container name")
	flag.CommandLine.Parse(defineArgs)

	if aa.KubeConfigFilePath != "" || aa.User != "" {
		aa.UpdateConfig = true
	}
	aa.Node = aa.SplitStrArg(nodesNameStr, aa.Node)
	aa.Namespace = aa.SplitStrArg(nsNameStr, aa.Namespace)
	aa.Pods = aa.SplitStrArg(podsNameStr, aa.Pods)
	aa.Containers = aa.SplitStrArg(ctrsNameStr, aa.Containers)
	aa.Action = undefineArgs[0]
	aa.Remaining = undefineArgs[1:]
}

// ----------------------------------------
func LoadArgs() *AllArgs {
	allArgs := &AllArgs{}
	// parse input args
	if len(os.Args) == 1 {
		allArgs.ShowUsage()
		os.Exit(0)
	}
	allArgs.ParseArgs(os.Args[1:])
	if allArgs.Action == "config" && len(allArgs.Remaining) == 1 && allArgs.Remaining[0] == "reset" {
		if err := os.Remove(RecommendConfFile); err != nil {
			fmt.Printf("[ERROR] delete %v failed, err: %v\n", RecommendConfFile, err)
			os.Exit(0)
		}
		allArgs.Remaining = []string{}
	}
	// read args from config file
	baseArgs := &BaseArgs{}
	if !baseArgs.ReadBaseArgs() {
		return nil
	}
	allArgs.KubeConfigFilePath = CombineStrArg(allArgs.KubeConfigFilePath, baseArgs.KubeConfigFilePath, "")
	allArgs.Node = CombineListArg(allArgs.Node, baseArgs.Node)
	allArgs.User = CombineStrArg(allArgs.User, baseArgs.User, "root")
	allArgs.Namespace = CombineListArg(allArgs.Namespace, baseArgs.Namespace)
	allArgs.Pods = CombineListArg(allArgs.Pods, baseArgs.Pods)
	allArgs.Containers = CombineListArg(allArgs.Containers, baseArgs.Containers)
	return allArgs
}
