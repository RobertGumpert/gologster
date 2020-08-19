package mypackage

import (
	gologger "../../.."
	"strconv"
)

var (
	logs *gologger.LogInterface
)

type number struct {
	Value int
	Prev  *number
}

func SetLogs(main *gologger.LogInterface) {
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
		logs.Info(currentNumber, gologger.GoConsole())

		prev = currentNumber
	}

	logs.Info("Print all numbers : " + strconv.Itoa(num))
}
