package logger

import (
	"errors"
	"log"
	"os"
	"strings"
)

// loggerFileMutex : логгер в файл с использованием мьютекса. | logger to file using mutex.
//
// Работает со стандартным пакетом 'log', используя 'log.SetOutput(io.Writer)',
// с переданным в качестве параметра, поток ввода, открытого файла.
// Запись в файл производится в той же гоурутине (потоке), где был вызван 'loggerFileMutex.add()'.
//
// Works with the standard package 'log' using 'log.SetOutput(io.Writer)',
// with an open file input stream passed as a parameter.
// Writing to the file is done in the same goroutine (stream) where 'loggerFileMutex.add()' was called.
//
type loggerFileMutex struct {
	// Базовый объект работы с файлом.
	// Basic file worker.
	basicFileLogger *loggerFileBasic

	// "key" -> fileAgent
	// EXAMPLE: "sql" -> fileAgent : { path: "./log_sql" }
	filesMap map[string]fileAgent
}

// newLoggerFileMutex : constructor
//
func newLoggerFileMutex(file *loggerFileBasic) *loggerFileMutex {
	logger := new(loggerFileMutex)
	logger.filesMap = make(map[string]fileAgent, 0)
	logger.basicFileLogger = file
	for key, path := range file.filesMap {
		file := fileAgent{
			path: path,
		}
		logger.filesMap[key] = file
	}
	return logger
}

// add : implement iLogger interface
//
func (logger *loggerFileMutex) add(value interface{}, lvl level, date string, param ...string) {
	err, key := logger.basicFileLogger.getParams(value, lvl, date, param)
	if err != nil {
		return
	}
	out, err := logger.createOutputString(value, lvl, date)
	if err != nil {
		logger.errorOutput(out, err)
		return
	}
	if file, exist := logger.filesMap[key]; exist {
		err := logger.output(out, file.path)
		if err != nil {
			logger.errorOutput(out, err)
		}
		return
	} else {
		logger.errorOutput(out, errors.New("fileAsync isn't exist by key : '"+key+"'"))
	}
}

// createOutputString : implement iLogger interface
//
// Исрользуется / используйте 'basicFileLogger.createOutputString()'
//
// Use 'basicFileLogger.createOutputString()'
//
func (logger *loggerFileMutex) createOutputString(value interface{}, lvl level, date string, param ...string) (*outputString, error) {
	return logger.basicFileLogger.createOutputString(value, lvl, date)
}

// output : implement iLogger interface
//
func (logger *loggerFileMutex) output(out *outputString, param ...string) error {
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
	log.SetOutput(file)
	log.Print(*out)
	return closeFile(file, nil)
}

// errorOutput : implement iLogger interface
//
// Исрользуется / используйте 'basicFileLogger.errorOutput()'
//
// Use 'basicFileLogger.errorOutput()'
//
func (logger *loggerFileMutex) errorOutput(out *outputString, err error) {
	logger.basicFileLogger.errorOutput(out, err)
}
