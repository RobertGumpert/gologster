package gologster

import (
	"text/template"
)

// loggerBaseConsole : определяет базовое поведение логгера в консоль| defines the base behavior of the logger to the console
//
type loggerBaseConsole struct {
	// Объект базового логгера, со стандартным поведением.
	// Basic logger object, with standard behavior.
	base     *loggerBase
	tmpl     *template.Template
	packages map[string]struct{}
}

// newBaseConsole : constructor
//
func newBaseConsole(base *loggerBase, packages map[string]struct{}, tmpl *template.Template) *loggerBaseConsole {
	logger := new(loggerBaseConsole)
	logger.base = base
	logger.tmpl = tmpl
	logger.packages = packages
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
func (logger *loggerBaseConsole) add(log *logData, param ...string) {
	return
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
func (logger *loggerBaseConsole) createOutputString(log *logData, param ...string) (*string, error) {
	out := log.filledTemplate(logger.tmpl)
	return out, nil
}

// output : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBase'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the base logger 'loggerBase'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBaseConsole) output(out *string, param ...string) error {
	err := logger.base.output(out)
	return err
}

// errorOutput : implement iLogger interface
//
// Поведение определенно базовым логгером  'loggerBase'.
// Типы, которые встраивают в себя данный тип, могут самостоятельно
// определять поведение.
//
// The behavior is defined by the base logger 'logger Basic'.
// Types that embed a given type can define behavior on their own.
//
func (logger *loggerBaseConsole) errorOutput(out *string, err error) {
	logger.base.errorOutput(out, err)
}
