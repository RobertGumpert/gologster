package mypackage

import (
	"github.com/RobertGumpert/gologster"
	"strconv"
)

var (
	logger *gologster.Logger
)

type number struct {
	Value int
	Prev  *number
}

func SetLogs(main *gologster.Logger) {
	logger = main
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
		logger.Info(currentNumber)

		prev = currentNumber
	}

	logger.Info("Print all numbers : " + strconv.Itoa(num))
}