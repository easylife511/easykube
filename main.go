package main

import (
	"fmt"
	"k8s_tools/process"
	"os"
)

func main() {
	if !process.Launch() {
		fmt.Println("process failure !!!")
		os.Exit(0)
	}

	// fmt.Println("main done ...")
}
