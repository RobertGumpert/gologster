package logger

import (
	"errors"
	"log"
	"os"
	"strings"
)

// loggerFileMutex : логгер в файл с использованием мьютекса. | logger to file using mutex.
//
// Работает со стандартным пакетом 'callingMode', используя 'callingMode.SetOutput(io.Writer)',
// с переданным в качестве параметра, поток ввода, открытого файла.
// Запись в файл производится в той же гоурутине (потоке), где был вызван 'loggerFileMutex.add()'.
//
// Works with the standard package 'callingMode' using 'callingMode.SetOutput(io.Writer)',
// with an open file input stream passed as a parameter.
// Writing to the file is done in the same goroutine (stream) where 'loggerFileMutex.add()' was called.
//
type loggerFileMutex struct {
	// Базовый объект работы с файлом.
	// Basic file worker.
	baseFile *loggerBaseFile

	// "key" -> fileAgent
	// EXAMPLE: "sql" -> fileAgent : { path: "./log_sql" }
	config map[string]fileAgent
}

// newLoggerFileMutex : constructor
//
func newLoggerFileMutex(baseFile *loggerBaseFile) *loggerFileMutex {
	logger := new(loggerFileMutex)
	logger.config = make(map[string]fileAgent, 0)
	logger.baseFile = baseFile
	for key, path := range baseFile.config {
		logger.newFile(key, path)
	}
	return logger
}

func (logger *loggerFileMutex) newFile(key, path string) {
	if _, exist := logger.config[key]; !exist {
		logger.baseFile.config[key] = path
	}
	file := fileAgent{
		path: path,
	}
	logger.config[key] = file
}

// add : implement iLogger interface
//
func (logger *loggerFileMutex) add(log *logData, param ...string) {
	var (
		performOutput = func(log *logData, logger *loggerFileMutex, key string) {
			out, err := logger.createOutputString(log)
			if err != nil {
				logger.baseFile.errorOutput(out, err)
				return
			}
			if file, exist := logger.config[key]; exist {
				err := logger.output(out, file.path)
				if err != nil {
					logger.errorOutput(out, err)
				}
				return
			} else {
				logger.errorOutput(out, errors.New("loggerFileMutex.add isn't exist by key : '"+key+"'"))
			}
		}
		fileKey = ""
	)
	if log.IsOption {
		err, key := logger.baseFile.getParams(log, param...)
		if err != nil {
			return
		}
		fileKey = key
	} else {
		for pckg := range logger.baseFile.config {
			if strings.Contains(log.Package, pckg) {
				fileKey = pckg
				break
			}
		}
	}
	performOutput(log, logger, fileKey)
}

// createOutputString : implement iLogger interface
//
// Исрользуется / используйте 'baseFile.createOutputString()'
//
// Use 'baseFile.createOutputString()'
//
func (logger *loggerFileMutex) createOutputString(log *logData, param ...string) (*string, error) {
	return logger.baseFile.createOutputString(log, param...)
}

// output : implement iLogger interface
//
func (logger *loggerFileMutex) output(out *string, param ...string) error {
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
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.Print(*out)
	return closeFile(file, nil)
}

// errorOutput : implement iLogger interface
//
// Исрользуется / используйте 'baseFile.errorOutput()'
//
// Use 'baseFile.errorOutput()'
//
func (logger *loggerFileMutex) errorOutput(out *string, err error) {
	logger.baseFile.errorOutput(out, err)
}
