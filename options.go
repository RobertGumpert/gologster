package logger

// Для пользовательского интерфейса используется
// паттерн программирования 'Functional options'.
// Опции ('IsOption') выстыпают в качестве функций,
// которые возвращают функции 'Mode'. Функции 'Mode',
// в свою очередь вызывают метод 'add' у конкретного
// логгера реализующего интерфейс 'iLogger' (объект хранится
// в глобальном пользовательском объекте 'Logger').
// Обычно этот паттерн используется для создания новых объектов,
// но в данном случае только выбора вывода лога.
//
// The user interface uses the 'Functional options' programming pattern.
// Options ('IsOption') pop out as functions that return 'Mode' functions.
// The 'Mode' functions, in turn, call the 'add' method of a specific logger
// that implements the 'iLogger' interface
// (the object is stored in the global user object 'Logger').
// Usually this pattern is used to filledTemplate new objects,
// but in this case only to select the callingMode output.

// Mode : в теле содержит вызов метода 'add' конкретного логгера.
//        in the body contains a call to the 'add' method of a particular logger.
//
type Mode func(logger *Logger, log *logData)

// IsOption : возвращает 'Mode' соответствующий  выбранной пользователем опции.
//          returns 'Mode' corresponding to the option selected by the user.
//
type Option func(param ...string) Mode

// OptionConsole : возвращает 'Mode' соответствующий 'loggerConsoleSimple'.
//           Вызов в том же потоке.
//           returns 'Mode' corresponding to 'loggerConsoleSimple'.
//           Call on the same thread.
//
func OptionConsole(param ...string) Mode {
	return func(logger *Logger, log *logData) {
		logger.modeConsole.add(log, param...)
	}
}

// OptionFileMulti : возвращает 'Mode' соответствующий 'loggerFileMultithreading'.
//             Вызов в том же потоке.
//             returns 'Mode' corresponding to 'loggerFileMultithreading'.
//             Call on the same thread.
//
func OptionFileMulti(param ...string) Mode {
	return func(logger *Logger, log *logData) {
		if logger.modeFileMulti == nil {
			logger.modeConsole.add(log, param...)
		}
		logger.modeFileMulti.add(log, param...)
	}
}

// OptionFileMulti : возвращает 'Mode' соответствующий 'loggerFileMutex'.
//             Вызов в том же потоке.
//             returns 'Mode' corresponding to 'loggerFileMutex'.
//             Call on the same thread.
//
func OptionFileMutex(param ...string) Mode {
	return func(logger *Logger, log *logData) {
		if logger.modeFileMutex == nil {
			logger.modeConsole.add(log, param...)
		}
		logger.modeFileMutex.add(log, param...)
	}
}

// GoOptionConsole : возвращает 'Mode' соответствующий 'loggerConsoleSimple'.
//             Вызов в отдельном потоке.
//             returns 'Mode' corresponding to 'loggerConsoleSimple'.
//             Call in a separate thread.
//
func GoOptionConsole(param ...string) Mode {
	return func(logger *Logger, log *logData) {
		go logger.modeConsole.add(log, param...)
	}
}

// GoOptionFileMulti : возвращает 'Mode' соответствующий 'loggerFileMultithreading'.
//               Вызов в отдельном потоке.
//               returns 'Mode' corresponding to 'loggerFileMultithreading'.
//               Call in a separate thread.
//
func GoOptionFileMulti(param ...string) Mode {
	return func(logger *Logger, log *logData) {
		if logger.modeFileMulti == nil {
			go logger.modeConsole.add(log, param...)
		}
		go logger.modeFileMulti.add(log, param...)
	}
}

// GoOptionFileMutex : возвращает 'Mode' соответствующий 'loggerFileMutex'.
//               Вызов в отдельном потоке.
//               returns 'Mode' corresponding to 'loggerFileMutex'.
//               Call in a separate thread.
//
func GoOptionFileMutex(param ...string) Mode {
	return func(logger *Logger, log *logData) {
		if logger.modeFileMutex == nil {
			go logger.modeConsole.add(log, param...)
		}
		go logger.modeFileMutex.add(log, param...)
	}
}
