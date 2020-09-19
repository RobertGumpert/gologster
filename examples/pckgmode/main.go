package main

import (
	"./mypackage"
	urep "./repository/user"
	ucase "./usecase/user"
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

func createLogger(root string) *gologger.Logger {
	logger := gologger.Packages(map[string][]gologger.PackageInstaller{
		"main": {
			gologger.PackageConsoleSimple(gologger.BaseLogTemplate, gologger.MultiThreading),
		},
		"mypackage": {
			gologger.PackageConsoleSimple(gologger.BaseLogTemplate, gologger.SingleThreading),
		},
		"/repository/user": {
			gologger.PackageConsoleSimple(gologger.BaseLogTemplate, gologger.MultiThreading),
			gologger.PackageFileMutex(gologger.BaseLogTemplate, gologger.SingleThreading, root + "/logs/file_1.txt"),
		},
		"/usecase/user": {
			gologger.PackageConsoleSimple(gologger.BaseLogTemplate, gologger.MultiThreading),
			gologger.PackageFileMutex(gologger.BaseLogTemplate, gologger.MultiThreading, root + "/logs/file_2.txt"),
		},
	}, )
	return logger
}
