package main

import (
	"./mypackage"
	urep "./repository/user"
	ucase "./usecase/user"
	"github.com/RobertGumpert/gologster"
	"path/filepath"
	"runtime"
	"time"
)

func getRoot() string {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Dir(file)
	return root
}

func main() {


	logger := createLogger(getRoot())

	mypackage.SetLogs(logger)
	urep.SetLogs(logger)
	ucase.SetLogs(logger)

	mypackage.PrintNumbers(10)
	urep.Log()
	ucase.Log()

	// Выполним логирование, в этом же потоке.
	// Let's perform logging in the same thread.
	logger.Info("App is started!")

	logger.Info("App has terminated!")

	time.Sleep(5 * time.Second)
}

func createLogger(root string) *gologster.Logger {
	logger := gologster.Packages(map[string][]gologster.PackageInstaller{
		"main": {
			gologster.PackageConsoleSimple(gologster.BaseLogTemplate, gologster.MultiThreading),
		},
		"mypackage": {
			gologster.PackageConsoleSimple(gologster.BaseLogTemplate, gologster.SingleThreading),
		},
		"/repository/user": {
			gologster.PackageConsoleSimple(gologster.BaseLogTemplate, gologster.MultiThreading),
			gologster.PackageFileMutex(gologster.BaseLogTemplate, gologster.SingleThreading, root + "/logs/file_1.txt"),
		},
		"/usecase/user": {
			gologster.PackageConsoleSimple(gologster.BaseLogTemplate, gologster.MultiThreading),
			gologster.PackageFileMutex(gologster.BaseLogTemplate, gologster.MultiThreading, root + "/logs/file_2.txt"),
		},
	}, )
	return logger
}
