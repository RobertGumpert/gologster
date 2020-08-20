package logger

// loggerConsoleBasic : определяет базовое поведение логгера в консоль| defines the basic behavior of the logger to the console
//
type loggerConsoleBasic struct {
	// Объект базового логгера, со стандартным поведением.
	// Basic logger object, with standard behavior.
	basic *loggerBasic
}

// newBasicConsole : constructor
//
func newBasicConsole(basic *loggerBasic) *loggerConsoleBasic {
	logger := new(loggerConsoleBasic)
	logger.basic = basic
	return logger
}

// add : implement iLogger interface
//
// Поведение определяется самостоятельно типами,
// которые встраивают в себя данный тип.
//
// The behavior is determined independently by the types that
// embed the given type in themselves.
//
func (logger *loggerConsoleBasic) add(value interface{}, lvl level, date, fn string, param ...string) {
	return
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
func (logger *loggerConsoleBasic) createOutputString(value interface{}, lvl level, date, fn string, param ...string) (*outputString, error) {
	return logger.basic.createOutputString(value, lvl, date, fn)
}

// output : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBasic'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the basic logger 'loggerBasic'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerConsoleBasic) output(out *outputString, param ...string) error {
	err := logger.basic.output(out)
	return err
}

// errorOutput : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBasic'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the basic logger 'logger Basic'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerConsoleBasic) errorOutput(out *outputString, err error) {
	logger.basic.errorOutput(out, err)
}
