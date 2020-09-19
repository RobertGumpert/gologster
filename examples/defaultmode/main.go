package main

import (
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
	// Выполним логирование, в этом же потоке.
	// Let's perform logging in the same thread.
	logger.Info("App is started!", gologster.OptionConsole(), gologster.OptionFileMutex("log_1"))


	logger.Info("App has terminated!", gologster.OptionConsole(), gologster.GoOptionFileMulti("log_1"))

	time.Sleep(5*time.Second)
}

func createLogger(root string) *gologster.Logger {
	logger := gologster.Default(
		//
		//
		gologster.DefaultConsoleSimple(gologster.BaseLogTemplate),
		//
		//
		gologster.DefaultFileMutex(
			gologster.BaseLogTemplate,
			map[string]string{
				"log_1": root + "/logs/file_1.txt",
			},
		),
		//
		//
		gologster.DefaultFileMulti(
			gologster.BaseLogTemplate,
		),
	)
	return logger
}
