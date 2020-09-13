package logger

import (
	"errors"
	"log"
	"strings"
	"text/template"
	"time"
)

// Logger : основной объект. Является пользовательским интерфейсом.
//                 Хранит все объекты реализующие интерфейс 'iLogger'.
//                 Создаётся один раз, на все приложение.
//
//                 main object. Is the user interface. Stores all objects
//                 that implement the 'iLogger' interface.
//                 Created once, for the entire application.
//
type Logger struct {
	base          *loggerBase
	baseFile      *loggerBaseFile
	baseConsole   *loggerBaseConsole
	modeConsole   *loggerConsoleSimple
	modeFileMulti *loggerFileMultithreading
	modeFileMutex *loggerFileMutex
	pckgs         map[string][]Option
}

type DefaultInstaller func(logger *Logger) error
type DefaultConfigurator func(templateString string, params ...map[string]string) DefaultInstaller

type PackageInstaller func(logger *Logger, pckg string) error
type PackageConfigurator func(templateString string, isConcurrency concurrency, params ...string) PackageInstaller

type concurrency bool

const SingleThreading concurrency = false
const MultiThreading concurrency = true

func DefaultConsoleSimple(templateString string, params ...map[string]string) DefaultInstaller {
	return func(logger *Logger) error {
		packages := make(map[string]struct{})
		tmpl, err := template.New("console_simple").Parse(templateString)
		if err != nil {
			tmpl, _ = template.New("console_simple").Parse(BaseLogTemplate)
		}
		logger.baseConsole = newBaseConsole(logger.base, packages, tmpl)
		logger.modeConsole = newLoggerConsoleSimple(logger.baseConsole)
		return nil
	}
}

func DefaultFileMutex(templateString string, params ...map[string]string) DefaultInstaller {
	return func(logger *Logger) error {
		if logger.baseFile == nil {
			if len(params) == 0 && len(params[0]) == 0 {
				return errors.New("DefaultFileMutex : File map isn't exist. ")
			}
		}
		tmpl, err := template.New("file_mutex").Parse(templateString)
		if err != nil {
			tmpl, _ = template.New("file_mutex").Parse(BaseLogTemplate)
		}
		if logger.baseFile == nil {
			logger.baseFile = newBaseFile(logger.base, params[0], tmpl)
		}
		logger.modeFileMutex = newLoggerFileMutex(logger.baseFile)
		return nil
	}
}

func DefaultFileMulti(templateString string, params ...map[string]string) DefaultInstaller {
	return func(logger *Logger) error {
		if logger.baseFile == nil {
			if len(params) == 0 && len(params[0]) == 0 {
				return errors.New("DefaultFileMulti : File map isn't exist. ")
			}
		}
		tmpl, err := template.New("file_mutex").Parse(templateString)
		if err != nil {
			tmpl, _ = template.New("file_mutex").Parse(BaseLogTemplate)
		}
		if logger.baseFile == nil {
			logger.baseFile = newBaseFile(logger.base, params[0], tmpl)
		}
		logger.modeFileMulti = newLoggerFileMultithreading(logger.baseFile)
		return nil
	}
}

func PackageConsoleSimple(templateString string, isConcurrency concurrency, params ...string) PackageInstaller {
	return func(logger *Logger, pckg string) error {
		tmpl, err := template.New("console_simple").Parse(templateString)
		if err != nil {
			tmpl, _ = template.New("console_simple").Parse(BaseLogTemplate)
		}
		//
		if logger.modeConsole == nil {
			packages := make(map[string]struct{})
			packages[pckg] = struct{}{}
			logger.baseConsole = newBaseConsole(logger.base, packages, tmpl)
			logger.modeConsole = newLoggerConsoleSimple(logger.baseConsole)
		} else {
			logger.baseConsole.packages[pckg] = struct{}{}
		}
		//
		if isConcurrency {
			logger.pckgs[pckg] = append(logger.pckgs[pckg], GoOptionConsole)
		} else {
			logger.pckgs[pckg] = append(logger.pckgs[pckg], OptionConsole)
		}
		//
		return nil
	}
}

func PackageFileMutex(templateString string, isConcurrency concurrency, params ...string) PackageInstaller {
	return func(logger *Logger, pckg string) error {
		//
		tmpl, err := template.New("console_simple").Parse(templateString)
		if err != nil {
			tmpl, _ = template.New("console_simple").Parse(BaseLogTemplate)
		}
		//
		if logger.modeFileMutex == nil {
			packages := make(map[string]string, 0)
			for _, file := range params {
				packages[pckg] = file
			}
			logger.baseFile = newBaseFile(logger.base, packages, tmpl)
			logger.modeFileMutex = newLoggerFileMutex(logger.baseFile)
		} else {
			for _, file := range params {
				logger.modeFileMutex.newFile(pckg, file)
			}
		}
		//
		if isConcurrency {
			logger.pckgs[pckg] = append(logger.pckgs[pckg], GoOptionFileMutex)
		} else {
			logger.pckgs[pckg] = append(logger.pckgs[pckg], OptionFileMutex)
		}
		//
		return nil
	}
}


func PackageFileMulti(templateString string, isConcurrency concurrency, params ...string) PackageInstaller {
	return func(logger *Logger, pckg string) error {
		//
		tmpl, err := template.New("console_simple").Parse(templateString)
		if err != nil {
			tmpl, _ = template.New("console_simple").Parse(BaseLogTemplate)
		}
		//
		if logger.modeFileMulti == nil {
			packages := make(map[string]string, 0)
			for _, file := range params {
				packages[pckg] = file
			}
			logger.baseFile = newBaseFile(logger.base, packages, tmpl)
			logger.modeFileMulti = newLoggerFileMultithreading(logger.baseFile)
		} else {
			for _, file := range params {
				logger.modeFileMulti.newFile(pckg, file)
			}
		}
		//
		if isConcurrency {
			logger.pckgs[pckg] = append(logger.pckgs[pckg], GoOptionFileMutex)
		} else {
			logger.pckgs[pckg] = append(logger.pckgs[pckg], OptionFileMutex)
		}
		//
		return nil
	}
}


// Default : создаёт базовый пользовательский интерфейс, с выводом в консоль.
//           filledTemplate a base user interface, with output to the console.
//
func Default(installers ...DefaultInstaller) *Logger {
	logger := new(Logger)
	logger.base = newBase()
	for _, mode := range installers {
		err := mode(logger)
		if err != nil {
			log.Println(err)
		}
	}
	return logger
}

func Packages(packages map[string][]PackageInstaller) *Logger {
	logger := new(Logger)
	logger.base = newBase()
	logger.pckgs = make(map[string][]Option, 0)
	for name, installers := range packages {
		for _ , mode := range installers {
			err := mode(logger, name)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return logger
}

// Info : логирование уровня 'info'.
//        logging level 'info'.
//
func (logger *Logger) Info(value interface{}, modes ...Mode) {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	data := newLogData(value, levelInfo, date).setRuntimeInfo(3)
	_ = data.marshal(logger.base)
	if len(modes) != 0 {
		data.IsOption = true
		logger.callingMode(data, modes...)
	} else {
		for pckg, options := range logger.pckgs {
			if strings.Contains(data.Package, pckg) {
				data.Package = pckg
				logger.callingOption(data, options...)
				break
			}
		}
	}
}

// Error : логирование уровня 'error'.
//         logging level 'error'.
//
func (logger *Logger) Error(value interface{}, modes ...Mode) {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	data := newLogData(value, levelError, date).setRuntimeInfo(4)
	_ = data.marshal(logger.base)
	if len(modes) != 0 {
		data.IsOption = true
		logger.callingMode(data, modes...)
	} else {
		for pckg, options := range logger.pckgs {
			if strings.Contains(data.Package, pckg) {
				data.Package = pckg
				logger.callingOption(data, options...)
				break
			}
		}
	}
}

// Panic : логирование уровня 'panic'.
//         logging level 'panic'.
//
func (logger *Logger) Panic(value interface{}, modes ...Mode) {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	data := newLogData(value, levelPanic, date).setRuntimeInfo(4)
	_ = data.marshal(logger.base)
	if len(modes) != 0 {
		data.IsOption = true
		logger.callingMode(data, modes...)
	} else {
		for pckg, options := range logger.pckgs {
			if strings.Contains(data.Package, pckg) {
				data.Package = pckg
				logger.callingOption(data, options...)
				break
			}
		}
	}
}

func (logger *Logger) callingOption(log *logData, options ...Option) {
	for _, option := range options {
		mode := option()
		mode(logger, log)
	}
}

func (logger *Logger) callingMode(log *logData, modes ...Mode) {
	for _, mode := range modes {
		mode(logger, log)
	}
}
