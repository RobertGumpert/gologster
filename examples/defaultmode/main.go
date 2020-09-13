package main

import (
	gologger "../../../gologger"
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
	// Выполним логирование, в этом же потоке.
	// Let's perform logging in the same thread.
	logger.Info("App is started!", gologger.OptionConsole(), gologger.OptionFileMutex("log_1"))


	logger.Info("App has terminated!", gologger.OptionConsole(), gologger.GoOptionFileMulti("log_1"))

	time.Sleep(5*time.Second)
}

func createLogger(root string) *gologger.Logger {
	logger := gologger.Default(
		//
		//
		gologger.DefaultConsoleSimple(gologger.BaseLogTemplate),
		//
		//
		gologger.DefaultFileMutex(
			gologger.BaseLogTemplate,
			map[string]string{
				"log_1": root + "/logs/file_1.txt",
			},
		),
		//
		//
		gologger.DefaultFileMulti(
			gologger.BaseLogTemplate,
		),
	)
	return logger
}
