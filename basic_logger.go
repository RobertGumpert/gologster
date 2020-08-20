package logger

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"
)

// loggerBasic : определяет базовое поведение любого логгера | defines the basic behavior of any logger
//
// Для вывода использует стандартный пакет 'log' в консоль.
//
// Uses the standard 'log' package for output to the console.
//
type loggerBasic struct{}

// newBasic() : constructor
//
func newBasic() *loggerBasic {
	return new(loggerBasic)
}

// add : implement iLogger interface
//
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBasic) add(value interface{}, lvl level, date, fn string, param ...string) {
	return
}

// createOutputString : implement iLogger interface
//
// Не использует параметров.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// Doesn't use parameters.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBasic) createOutputString(value interface{}, lvl level, date, fn string, param ...string) (*outputString, error) {
	if value == nil {
		if date == "" {
			err := strings.Join([]string{
				"empty 'date' and 'value' log :: time='",
				time.Now().Format("Mon Jan _2 15:04:05 2006"),
				"'",
			}, "")
			return newOutputString([]byte(""), date, fn, lvl), errors.New(err)
		}
		return newOutputString([]byte(""), date, fn, lvl), errors.New("empty log 'value'")
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		e := strings.Join([]string{
			"error in json.Marshal(value)='",
			err.Error(),
			"'",
		}, "")
		return newOutputString([]byte("empty log 'value'"), date, fn, lvl), errors.New(e)
	}
	return newOutputString(bytes, date, fn, lvl), nil
}

// output : implement iLogger interface
//
// Не использует параметров.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// Doesn't use parameters.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBasic) output(out *outputString, param ...string) error {
	log.SetOutput(os.Stdout)
	log.Println(*out)
	return nil
}

// errorOutput : implement iLogger interface
//
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBasic) errorOutput(out *outputString, err error) {
	update := out.addError(err)
	log.Println(update)
}
