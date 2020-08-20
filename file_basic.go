package logger

import (
	"errors"
	"os"
	"strings"
)

// fileAgent : определяет параметры, каналы и так далее,
// 			   которые могут быть необходимы типами,
// 			   которые встраивают в себя данный тип (loggerFileBasic), при работе с файлом.
//
//			   defines parameters, channels and so on,
//			   which may be required by types that
//			   embed this type (loggerFileBasic) when working with a file.
//
type fileAgent struct {
	path    string
	channel chan *outputString
}

// loggerFileBasic : определяет базовое поведение логгера в файл| defines the basic behavior of the logger to the file
//
type loggerFileBasic struct {
	// Объект базового логгера, со стандартным поведением.
	// Basic logger object, with standard behavior.
	basic *loggerBasic

	// key -> value : "sql" -> "~/home/dir/log_sql.txt"
	filesMap map[string]string
}

// newBasicFile : constructor
//
func newBasicFile(basic *loggerBasic, filesMap map[string]string) *loggerFileBasic {
	logger := new(loggerFileBasic)
	logger.basic = basic
	logger.filesMap = filesMap
	return logger
}

// getParams : проверяет наличие только одного параметра - ключ файла. | checks for only one parameter - the file key.
//
func (logger *loggerFileBasic) getParams(value interface{}, lvl level, date, fn string, param []string) (error, string) {
	var (
		key = ""
	)
	if len(param) >= 1 {
		key = param[0]
	} else {
		err := errors.New("fileAsync 'key' isn't exist in 'param ...string'")
		out, e := logger.createOutputString(value, lvl, date, fn)
		if e != nil {
			err = errors.New(strings.Join([]string{
				err.Error(),
				e.Error(),
			}, ";"))
		}
		logger.errorOutput(out, err)
		return err, key
	}
	return nil, key
}

// add : implement iLogger interface
//
// Поведение определяется самостоятельно типами,
// которые встраивают в себя данный тип.
//
// The behavior is determined independently by the types that
// embed the given type in themselves.
//
func (logger *loggerFileBasic) add(value interface{}, lvl level, date, fn string, param ...string) {
	return
}

// createOutputString : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBasic'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the basic logger 'logger Basic'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerFileBasic) createOutputString(value interface{}, lvl level, date, fn string, param ...string) (*outputString, error) {
	return logger.basic.createOutputString(value, lvl, date, fn)
}

// output : implement iLogger interface
//
// Передаваемый массив параметров состоит
// только из одного элмента - путь до файла.
// Запись производится без буффера.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The passed parameter array consists
// of only one element - the path to the file.
// Recording is performed without a buffer.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerFileBasic) output(out *outputString, param ...string) error {
	var (
		path      = param[0]
		closeFile = func(file *os.File, err error) error {
			errFile := file.Close()
			if errFile != nil {
				err = errors.New(strings.Join([]string{
					err.Error(),
					errFile.Error(),
				}, "::"))
			}
			return err
		}
	)
	file, err := os.OpenFile(path, os.O_APPEND, 0666)
	if err != nil {
		return closeFile(file, nil)
	}
	_, err = file.WriteString(string(*out))
	if err != nil {
		return closeFile(file, nil)
	}
	return nil
}

// createOutputString : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBasic'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the basic logger 'loggerBasic'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerFileBasic) errorOutput(out *outputString, err error) {
	logger.basic.errorOutput(out, err)
}
