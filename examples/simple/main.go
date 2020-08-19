package main

import (
	"../../logger"
	"./mypackage"
	"time"
)

// Создаём базового логгера, в котором доступен только вывод в консоль.
// We create a basic logger in which only output to the console is available.
var logs = logger.Default()

func main() {

	// Выполним логирование, в этом же потоке.
	// Let's perform logging in the same thread.
	logs.Info("App is started!")

	mypackage.SetLogs(logs)
	mypackage.PrintNumbers(10)
	time.Sleep(5 * time.Second)

	logs.Info("App has terminated!")
}
