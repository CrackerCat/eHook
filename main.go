package main

import (
	"eHook/user"
	"eHook/controller"
	"eHook/module"
	"os"
	"eHook/utils"
	"fmt"
	"os/signal"
	"syscall"
    _ "github.com/shuLhan/go-bindata" // add for bindata in Makefile
)

func main() {
	stopper := make(chan os.Signal, 1)
    signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	process, err := controller.CreateProcess(user.PackageName)
	if err != nil {
		fmt.Println("Parse process error: ", err)
		os.Exit(1)
	}

	library, err := controller.CreateLibrary(process, user.LibraryName)
	if err != nil {
		fmt.Println("Parse library error: ", err)
		os.Exit(1)
	}
	btfFile := ""
	if !utils.CheckConfig("CONFIG_DEBUG_INFO_BTF=y") {
		btfFile = utils.FindBTFAssets()
	}
    probe := module.CreateProbeHandler(btfFile)
	err = probe.Run(library)
	if err != nil {
		fmt.Printf("Failed to start module: %v", err)
		os.Exit(1)
	}
	fmt.Println("Module started. Use Ctrl+C to Exit.")
    <-stopper
	err = probe.Stop()
	if err != nil {
		fmt.Printf("Failed to stop module: %v", err)
	}
}