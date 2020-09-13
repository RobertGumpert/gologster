package logger

import (
	"errors"
	"os"
	"strings"
	"text/template"
)

// fileAgent : определяет параметры, каналы и так далее,
// 			   которые могут быть необходимы типами,
// 			   которые встраивают в себя данный тип (loggerBaseFile), при работе с файлом.
//
//			   defines parameters, channels and so on,
//			   which may be required by types that
//			   embed this type (loggerBaseFile) when working with a file.
//
type fileAgent struct {
	path    string
	channel chan *string
}

// loggerBaseFile : определяет базовое поведение логгера в файл| defines the base behavior of the logger to the file
//
type loggerBaseFile struct {
	// Объект базового логгера, со стандартным поведением.
	// Basic logger object, with standard behavior.
	base *loggerBase

	// key -> value : "sql" -> "~/home/dir/log_sql.txt"
	config map[string]string
	tmpl   *template.Template
}

// newBaseFile : constructor
//
func newBaseFile(base *loggerBase, config map[string]string, tmpl *template.Template) *loggerBaseFile {
	logger := new(loggerBaseFile)
	logger.base = base
	logger.config = config
	logger.tmpl = tmpl
	return logger
}

// getParams : проверяет наличие только одного параметра - ключ файла. | checks for only one parameter - the file key.
//
func (logger *loggerBaseFile) getParams(log *logData, param ...string) (error, string) {
	var (
		key = ""
	)
	if len(param) >= 1 {
		key = param[0]
	} else {
		err := errors.New("loggerBaseFile.getParams 'key' isn't exist in 'param ...string'")
		out, e := logger.createOutputString(log, param...)
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
func (logger *loggerBaseFile) add(log *logData, param ...string) {
	return
}

// createOutputString : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBase'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the base logger 'logger Basic'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBaseFile) createOutputString(log *logData, param ...string) (*string, error) {
	out := log.filledTemplate(logger.tmpl)
	return out, nil
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
func (logger *loggerBaseFile) output(out *string, param ...string) error {
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
// Поведение определенно базовым логгером  'loggerBase'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the base logger 'loggerBase'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBaseFile) errorOutput(out *string, err error) {
	logger.base.errorOutput(out, err)
}
