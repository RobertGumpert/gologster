package logger

// Для пользовательского интерфейса используется
// паттерн программирования 'Functional options'.
// Опции ('Option') выстыпают в качестве функций,
// которые возвращают функции 'Mode'. Функции 'Mode',
// в свою очередь вызывают метод 'add' у конкретного
// логгера реализующего интерфейс 'iLogger' (объект хранится
// в глобальном пользовательском объекте 'LogInterface').
// Обычно этот паттерн используется для создания новых объектов,
// но в данном случае только выбора вывода лога.
//
// The user interface uses the 'Functional options' programming pattern.
// Options ('Option') pop out as functions that return 'Mode' functions.
// The 'Mode' functions, in turn, call the 'add' method of a specific logger
// that implements the 'iLogger' interface
// (the object is stored in the global user object 'LogInterface').
// Usually this pattern is used to create new objects,
// but in this case only to select the log output.

// Mode : в теле содержит вызов метода 'add' конкретного логгера.
//        in the body contains a call to the 'add' method of a particular logger.
//
type Mode func(logger *LogInterface, value interface{}, lvl level, date string)

// Option : возвращает 'Mode' соответствующий  выбранной пользователем опции.
//          returns 'Mode' corresponding to the option selected by the user.
//
type Option func(param ...string)

// Console : возвращает 'Mode' соответствующий 'loggerConsoleSinglethreading'.
//           Вызов в том же потоке.
//           returns 'Mode' corresponding to 'loggerConsoleSinglethreading'.
//           Call on the same thread.
//
func Console(param ...string) Mode {
	return func(logger *LogInterface, value interface{}, lvl level, date string) {
		logger.modeConsole.add(value, lvl, date, param...)
	}
}

// FileMulti : возвращает 'Mode' соответствующий 'loggerFileMultithreading'.
//             Вызов в том же потоке.
//             returns 'Mode' corresponding to 'loggerFileMultithreading'.
//             Call on the same thread.
//
func FileMulti(param ...string) Mode {
	return func(logger *LogInterface, value interface{}, lvl level, date string) {
		if logger.modeFileMulti == nil {
			logger.modeConsole.add(value, lvl, date, param...)
		}
		logger.modeFileMulti.add(value, lvl, date, param...)
	}
}

// FileMulti : возвращает 'Mode' соответствующий 'loggerFileMutex'.
//             Вызов в том же потоке.
//             returns 'Mode' corresponding to 'loggerFileMutex'.
//             Call on the same thread.
//
func FileMutex(param ...string) Mode {
	return func(logger *LogInterface, value interface{}, lvl level, date string) {
		if logger.modeFileMutex == nil {
			logger.modeConsole.add(value, lvl, date, param...)
		}
		logger.modeFileMutex.add(value, lvl, date, param...)
	}
}

// GoConsole : возвращает 'Mode' соответствующий 'loggerConsoleSinglethreading'.
//             Вызов в отдельном потоке.
//             returns 'Mode' corresponding to 'loggerConsoleSinglethreading'.
//             Call in a separate thread.
//
func GoConsole(param ...string) Mode {
	return func(logger *LogInterface, value interface{}, lvl level, date string) {
		go logger.modeConsole.add(value, lvl, date, param...)
	}
}

// GoFileMulti : возвращает 'Mode' соответствующий 'loggerFileMultithreading'.
//               Вызов в отдельном потоке.
//               returns 'Mode' corresponding to 'loggerFileMultithreading'.
//               Call in a separate thread.
//
func GoFileMulti(param ...string) Mode {
	return func(logger *LogInterface, value interface{}, lvl level, date string) {
		if logger.modeFileMulti == nil {
			go logger.modeConsole.add(value, lvl, date, param...)
		}
		go logger.modeFileMulti.add(value, lvl, date, param...)
	}
}

// GoFileMutex : возвращает 'Mode' соответствующий 'loggerFileMutex'.
//               Вызов в отдельном потоке.
//               returns 'Mode' corresponding to 'loggerFileMutex'.
//               Call in a separate thread.
//
func GoFileMutex(param ...string) Mode {
	return func(logger *LogInterface, value interface{}, lvl level, date string) {
		if logger.modeFileMutex == nil {
			go logger.modeConsole.add(value, lvl, date, param...)
		}
		go logger.modeFileMutex.add(value, lvl, date, param...)
	}
}
