package logger

import (
	"runtime"
	"strings"
	"time"
)

// level : уровень логирования | logging level
//
type level int

const (
	levelInfo  level = 100
	levelError level = 200
	levelPanic level = 300
)

// iLogger : интерфейс логгера. | logger interface.
//
type iLogger interface {
	// add : метод, являющийся точкой доступа к логгеру. | method that is the access point to the logger.
	//
	// * value - логируеммые дынные.
	//           logging data.
	//
	// * lvl - уровень логгирования.
	//         logging level.
	//
	// * date - дата логирования, приведенная к строчному представлению.
	//          logging date, converted to line representation.
	//
	// * param - параметры, необходимые для логгирования.
	// 			 Для файла это может быть ключ по которому можно
	// 	         получить полный путь до файла. Четкий регламент не установлен,
	// 	         передавать в качестве параметра можно все-что угодно, все зависит
	// 	         от того что логгер реализующий данный интерфейс, должен делать
	// 	         с этим параметром.
	// 	         parameters required for logging. For a file, this can be a key
	// 	         by which you can get the full path to the file. There is no clear regulation,
	// 	         anything can be passed as a parameter, it all depends on what the logger
	// 	         implementing the given interface should do with this parameter.
	//
	// Через этот метод необходимо обращаться к логгеру
	// реализующему интерфейс 'iLogger'. Регламинтирует порядок
	// создания строки лога, вывода и обработку ошибок.
	//
	// Through this method it is necessary to contact the logger
	// that implements the 'iLogger' interface. Regulates the order
	// of log line creation, output and error handling.
	//
	add(value interface{}, lvl level, date, fn string, param ...string)

	// createOutputString : метод, ответственный за создание строки лога. | method responsible for creating the log line.
	//
	// * value - логируеммые дынные.
	//           logging data.
	//
	// * lvl - уровень логгирования.
	//         logging level.
	//
	// * date - дата логирования, приведенная к строчному представлению.
	//          logging date, converted to line representation.
	//
	// * param - параметры, необходимые для логгирования.
	// 			 Для файла это может быть ключ по которому можно
	// 	         получить полный путь до файла. Четкий регламент не установлен,
	// 	         передавать в качестве параметра можно все-что угодно, все зависит
	// 	         от того что логгер реализующий данный интерфейс, должен делать
	// 	         с этим параметром.
	// 	         parameters required for logging. For a file, this can be a key
	// 	         by which you can get the full path to the file. There is no clear regulation,
	// 	         anything can be passed as a parameter, it all depends on what the logger
	// 	         implementing the given interface should do with this parameter.
	//
	// Создаёт строку лога и возвращает ошибку в случае
	// если переданные данные не валидны или при их маршалинге
	// произошла какая-то ошибка.
	//
	// Creates a log line and returns an error
	// if the passed data is not valid or
	// some error occurred while marshaling.
	//
	createOutputString(value interface{}, lvl level, date, fn string, param ...string) (*outputString, error)

	// output : метод, выполняющий вывод лога. | method that performs log output.
	//
	// * outputString - строка лога.
	//          		logging line.
	//
	// * param - параметры, необходимые для логгирования.
	// 			 Для файла это может быть ключ по которому можно
	// 	         получить полный путь до файла. Четкий регламент не установлен,
	// 	         передавать в качестве параметра можно все-что угодно, все зависит
	// 	         от того что логгер реализующий данный интерфейс, должен делать
	// 	         с этим параметром.
	// 	         parameters required for logging. For a file, this can be a key
	// 	         by which you can get the full path to the file. There is no clear regulation,
	// 	         anything can be passed as a parameter, it all depends on what the logger
	// 	         implementing the given interface should do with this parameter.
	//
	// Выполняет вывод лога, например в файл, если логгер реализующий
	// данный интерфейс работает с файлами. Возвращает ошибку в случае,
	// если в момент вывода произошла ошибка.
	//
	// Performs log output, for example, to a file,
	// if a logger implementing this interface works with files.
	// Returns an error if an error occurred at the time of output.
	//
	output(outputString *outputString, param ...string) error

	// errorOutput : метод, выполняющий принудительный вывод в другой накопитель,
	//               (например: вывод производился в файл, но в случае ошибки вывода в файл,
	//               вывод будет выполнен в консоль) в случае если,
	//               внутри самого логгера произошла ошибка.
	//               a method that performs forced output to another drive
	//               (for example: output was made to a file, but in case of an error in output to a file,
	//               the output will be made to the console)
	//               if an error occurred inside the logger itself.
	//
	// * outputString - строка лога.
	//          		logging line.
	//
	// * err - ошибка вывода в основной накопитель.
	//         error output to the main drive.
	//
	errorOutput(outputString *outputString, err error)
}

// LogInterface : основной объект. Является пользовательским интерфейсом.
//                 Хранит все объекты реализующие интерфейс 'iLogger'.
//                 Создаётся один раз, на все приложение.
//
//                 main object. Is the user interface. Stores all objects
//                 that implement the 'iLogger' interface.
//                 Created once, for the entire application.
//
type LogInterface struct {
	basic         *loggerBasic
	fileBasic     *loggerFileBasic
	consoleBasic  *loggerConsoleBasic
	modeConsole   *loggerConsoleSinglethreading
	modeFileMulti *loggerFileMultithreading
	modeFileMutex *loggerFileMutex
}

// Default : создаёт базовый пользовательский интерфейс, с выводом в консоль.
//           create a basic user interface, with output to the console.
//
func Default() *LogInterface {
	logger := new(LogInterface)
	logger.basic = newBasic()
	logger.consoleBasic = newBasicConsole(logger.basic)
	logger.modeConsole = newLoggerConsoleSinglethreading(logger.consoleBasic)
	return logger
}

// ConfigFile : настраивает пользовательский интерфейс, для дальнейшей установки модов вывода в файл.
//              configures the user interface to further install output mods to a file.
//
func (logger *LogInterface) ConfigFile(files map[string]string) *LogInterface {
	if logger.fileBasic == nil {
		logger.fileBasic = newBasicFile(logger.basic, files)
	}
	return logger
}

// SetModeFileMulti : настраивает и добавляет в пользовательский интерфейс мод 'loggerFileMultithreading'.
//                    configures and adds the 'loggerFileMultithreading' mode to the user interface.
//
func (logger *LogInterface) SetModeFileMulti(files ...map[string]string) *LogInterface {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	if len(files) == 0 && logger.fileBasic == nil {
		logger.modeConsole.add("LogInterface message : from 'SetModeFileMulti()' files map len() = 0", levelError, date, "")
		return logger
	}
	if logger.fileBasic == nil {
		logger.fileBasic = newBasicFile(logger.basic, files[0])
	}
	logger.modeFileMulti = newLoggerFileMultithreading(logger.fileBasic)
	logger.modeConsole.add("LogInterface message : from 'SetModeFileMulti()' mode was set", levelInfo, date, "")
	return logger
}

// SetModeFileMutex : настраивает и добавляет в пользовательский интерфейс мод 'loggerFileMutex'.
//                    configures and adds the 'loggerFileMutex' mode to the user interface.
//
func (logger *LogInterface) SetModeFileMutex(files ...map[string]string) *LogInterface {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	if len(files) == 0 && logger.fileBasic == nil {
		logger.modeConsole.add("LogInterface message : from 'SetModeFileMutex()' files map len() = 0", levelError, date, "")
		return logger
	}
	if logger.fileBasic == nil {
		logger.fileBasic = newBasicFile(logger.basic, files[0])
	}
	logger.modeFileMutex = newLoggerFileMutex(logger.fileBasic)
	logger.modeConsole.add("LogInterface message : from 'SetModeFileMutex()' mode was set", levelInfo, date, "")
	return logger
}

// Info : логирование уровня 'info'.
//        logging level 'info'.
//
func (logger *LogInterface) Info(value interface{}, modes ...Mode) {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	fn := logger.getFuncFromRuntime(2)
	if modes == nil {
		logger.modeConsole.add(value, levelInfo, date, fn)
		return
	}
	logger.log(value, levelInfo, date, fn, modes...)
}

// Error : логирование уровня 'error'.
//         logging level 'error'.
//
func (logger *LogInterface) Error(value interface{}, modes ...Mode) {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	fn := logger.getFuncFromRuntime(2)
	if modes == nil {
		logger.modeConsole.add(value, levelError, date, fn)
		return
	}
	logger.log(value, levelError, date, fn, modes...)
}

// Panic : логирование уровня 'panic'.
//         logging level 'panic'.
//
func (logger *LogInterface) Panic(value interface{}, modes ...Mode) {
	date := time.Now().Format("Mon Jan _2 15:04:05 2006")
	fn := logger.getFuncFromRuntime(2)
	if modes == nil {
		logger.modeConsole.add(value, levelPanic, date, fn)
		return
	}
	logger.log(value, levelPanic, date, fn, modes...)
}

func (logger *LogInterface) getFuncFromRuntime(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "undefined func"
	}
	fn := runtime.FuncForPC(pc).Name()
	if strings.Contains(fn, "/") {
		split := strings.Split(fn, "/")
		fn = split[len(split)-1]
	}
	return fn
}

func (logger *LogInterface) log(value interface{}, lvl level, date, fn string, modes ...Mode) {
	for _, mode := range modes {
		mode(logger, value, lvl, date, fn)
	}
}
