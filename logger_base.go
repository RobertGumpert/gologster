package logger

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
)

// loggerBase : определяет базовое поведение любого логгера | defines the base behavior of any logger
//
// Для вывода использует стандартный пакет 'callingMode' в консоль.
//
// Uses the standard 'callingMode' package for output to the console.
//
type loggerBase struct {
}

// newBase() : constructor
//
func newBase() *loggerBase {
	return new(loggerBase)
}

// add : implement iLogger interface
//
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBase) add(log *logData, param ...string) {
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
func (logger *loggerBase) createOutputString(log *logData, param ...string) (*string, error) {
	var (
		out = ""
	)
	bytes, err := json.Marshal(log.UserDataOriginal)
	if err != nil {
		e := strings.Join([]string{
			"error in json.Value(value)='",
			err.Error(),
			"'",
		}, "")
		return &out, errors.New(e)
	}
	out = string(bytes)
	return &out, nil
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
func (logger *loggerBase) output(out *string, param ...string) error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
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
func (logger *loggerBase) errorOutput(out *string, err error) {
	update := strings.Join([]string{
		*out,
		"error=[" + err.Error() + "];",
	}, "")
	_ = logger.output(&update)
}
