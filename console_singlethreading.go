package logger

// loggerConsoleSinglethreading : определяет поведение логгера в консоль в однопоточном режиме | defines the behavior of the logger to the console in single-threaded mode
//
// Для вывода использует стандартный пакет 'log' в консоль.
// Вывод выполняется в той же горутине (потоке), где был вызван 'loggerConsoleSinglethreading.add()'.
//
// Uses the standard 'log' package for output to the console.
// The output is executed in the same goroutine (thread) where 'loggerConsoleSinglethreading.add ()' was called.
//
type loggerConsoleSinglethreading struct {
	// Базовый объект работы с консолью.
	// Basic object of working with the console.
	basicConsoleLogger *loggerConsoleBasic
}

// newLoggerConsoleSinglethreading : constructor
//
func newLoggerConsoleSinglethreading(basic *loggerConsoleBasic) *loggerConsoleSinglethreading {
	logger := new(loggerConsoleSinglethreading)
	logger.basicConsoleLogger = basic
	return logger
}

// add : implement iLogger interface
//
func (logger *loggerConsoleSinglethreading) add(value interface{}, lvl level, date string, param ...string) {
	out, err := logger.createOutputString(value, lvl, date)
	if err != nil {
		logger.errorOutput(out, err)
		return
	}
	_ = logger.output(out)
}

// add : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerConsoleBasic'.
//
// The behavior is defined by the basic logger 'loggerConsoleBasic'.
//
func (logger *loggerConsoleSinglethreading) createOutputString(value interface{}, lvl level, date string, param ...string) (*outputString, error) {
	return logger.basicConsoleLogger.createOutputString(value, lvl, date)
}

// output : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerConsoleBasic'.
//
// The behavior is defined by the basic logger 'loggerConsoleBasic'.
//
func (logger *loggerConsoleSinglethreading) output(out *outputString, param ...string) error {
	err := logger.basicConsoleLogger.output(out)
	return err
}

// errorOutput : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerConsoleBasic'.
//
// The behavior is defined by the basic logger 'loggerConsoleBasic'.
//
func (logger *loggerConsoleSinglethreading) errorOutput(out *outputString, err error) {
	logger.basicConsoleLogger.errorOutput(out, err)
}
