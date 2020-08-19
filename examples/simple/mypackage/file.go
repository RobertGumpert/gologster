package mypackage

import (
	"../../../logger"
	"strconv"
)

var (
	logs *logger.LogInterface
)

type number struct {
	Value int
	Prev  *number
}

func SetLogs(main *logger.LogInterface) {
	logs = main
}

func PrintNumbers(num int) {

	prev := new(number)

	for i := 0; i < num; i++ {

		currentNumber := &number{
			Value: i,
			Prev:  prev,
		}

		// Распечатаем сообщение в отдельном потоке.
		// Let's print the message in a separate thread.
		logs.Info(currentNumber, logger.GoConsole())

		prev = currentNumber
	}

	logs.Info("Print all numbers : " + strconv.Itoa(num))
}
