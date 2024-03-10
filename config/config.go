package config

import (
	"fmt"
	"k8s_tools/lib"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/util/homedir"
)

const (
	RecommendHomeDir      = ".ek"
	RecommendConfFileName = "config.yaml"
)

var (
	RecommendConfDir  = filepath.Join(homedir.HomeDir(), RecommendHomeDir)
	RecommendConfFile = filepath.Join(RecommendConfDir, RecommendConfFileName)
)

type BaseArgs struct {
	KubeConfigFilePath string   `yaml:"kubeconfig"`
	Node               []string `yaml:"node"`
	User               string   `yaml:"user"`
	Namespace          []string `yaml:"namespace"`
	Pods               []string `yaml:"pods"`
	Containers         []string `yaml:"containers"`
}

func (ba *BaseArgs) CheckAndCreatConfigFile() bool {
	// create dir
	if _, err := os.Stat(RecommendConfDir); err != nil {
		fmt.Printf("[Warning] Config dir not found: %s\n", RecommendConfDir)
		err := os.MkdirAll(RecommendConfDir, 0755)
		if err != nil {
			fmt.Println("[ERROR] create dir err:", err)
			return false
		}
		fmt.Println("create  dir: ", RecommendConfDir)
	}
	// create file
	if _, err := os.Stat(RecommendConfFile); err != nil {
		fmt.Printf("[Warning] Config file not found: %s\n", RecommendConfFile)
		file, err := os.Create(RecommendConfFile)
		if err != nil {
			fmt.Println("[ERROR] create file err:", err)
			return false
		}
		fmt.Println("create file: ", RecommendConfFile)
		ba.User = "root"
		if !ba.UpdateBaseArgs() {
			return false
		}
		defer file.Close()
	}
	// check and create success
	return true
}

func (ba *BaseArgs) ShowBaseArgs() {
	fmt.Printf("##### config file: %v #####\n", RecommendConfFile)
	fmt.Println("================================")
	fmt.Println("kubeconfig is: ", ba.KubeConfigFilePath)
	fmt.Println("node is: ", ba.Node)
	fmt.Println("user is: ", ba.User)
	fmt.Println("namespace is: ", ba.Namespace)
	fmt.Println("pods is: ", ba.Pods)
	fmt.Println("containers is: ", ba.Containers)
	fmt.Println("================================")
}

func (ba *BaseArgs) ShowVersion() {
	fmt.Println("version: ", Version)
}

func (ba *BaseArgs) ShowUsage() {
	usageContent := strings.Replace(Usage, "RecommendConfFile", RecommendConfFile, -1)
	fmt.Println(usageContent)
}

func (ba *BaseArgs) ReadBaseArgs() bool {
	if !ba.CheckAndCreatConfigFile() {
		return false
	}
	content := lib.ReadFile(RecommendConfFile)
	if content == "" {
		fmt.Printf("[ERROR] Read config file failed, %v\n", RecommendConfFile)
		return false
	}
	if err := yaml.Unmarshal([]byte(content), ba); err != nil {
		fmt.Printf("[ERROR] Yaml Unmarshal [%v] err: %v\n", content, err)
		return false
	}
	return true
}

func (ba *BaseArgs) UpdateBaseArgs() bool {
	yamlByte, err := yaml.Marshal(ba)
	if err != nil {
		fmt.Printf("[ERROR] Yaml Marshal [%v] err: %v\n", *ba, err)
		return false
	}
	if !lib.UpdateFile(string(yamlByte), RecommendConfFile) {
		return false
	}
	return true
}
