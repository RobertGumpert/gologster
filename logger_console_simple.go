package gologster

import "strings"

// loggerConsoleSimple : определяет поведение логгера в консоль в однопоточном режиме | defines the behavior of the logger to the console in single-threaded mode
//
// Для вывода использует стандартный пакет 'callingMode' в консоль.
// Вывод выполняется в той же горутине (потоке), где был вызван 'loggerConsoleSimple.add()'.
//
// Uses the standard 'callingMode' package for output to the console.
// The output is executed in the same goroutine (thread) where 'loggerConsoleSimple.add ()' was called.
//
type loggerConsoleSimple struct {
	// Базовый объект работы с консолью.
	// Basic object of working with the console.
	baseConsole *loggerBaseConsole
}

// newLoggerConsoleSimple : constructor
//
func newLoggerConsoleSimple(baseConsole *loggerBaseConsole) *loggerConsoleSimple {
	logger := new(loggerConsoleSimple)
	logger.baseConsole = baseConsole
	return logger
}

// add : implement iLogger interface
//
func (logger *loggerConsoleSimple) add(log *logData, param ...string) {
	performOutput := func(log *logData, logger *loggerConsoleSimple) {
		out, err := logger.createOutputString(log)
		if err != nil {
			logger.errorOutput(out, err)
			return
		}
		_ = logger.output(out)
	}
	if log.IsOption {
		performOutput(log, logger)
		return
	}
	for pckg, _ := range logger.baseConsole.packages {
		if strings.Contains(log.Package, pckg) {
			performOutput(log, logger)
			break
		}
	}
	//if _, exist := logger.baseConsole.packages[log.Package]; exist {
	//	performOutput(log, logger)
	//	return
	//}
}

// add : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBaseConsole'.
//
// The behavior is defined by the base logger 'loggerBaseConsole'.
//
func (logger *loggerConsoleSimple) createOutputString(log *logData, param ...string) (*string, error) {
	return logger.baseConsole.createOutputString(log, param...)
}

// output : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBaseConsole'.
//
// The behavior is defined by the base logger 'loggerBaseConsole'.
//
func (logger *loggerConsoleSimple) output(out *string, param ...string) error {
	err := logger.baseConsole.output(out)
	return err
}

// errorOutput : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBaseConsole'.
//
// The behavior is defined by the base logger 'loggerBaseConsole'.
//
func (logger *loggerConsoleSimple) errorOutput(out *string, err error) {
	logger.baseConsole.errorOutput(out, err)
}
