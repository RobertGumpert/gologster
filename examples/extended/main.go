package main

import (
	gologger "../.."
	"bufio"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _   = runtime.Caller(0)
	projectRoot  = filepath.Dir(b)
	loggingFiles = map[string]string{
		"file_1": projectRoot + "/log1.txt",
		"file_2": projectRoot + "/log2.txt",
	}
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	logs := settings(loggingFiles)

	logs.Info("Logger wes set!")

	// Выполним логирование, всеми установленными способами, в этом же потоке,
	// кроме вывода в консоль, его выполним в отдельном потоке.
	// Let's perform logging, in all the established ways, in the same thread,
	// except for output to the console, we will execute it in a separate thread.
	logs.Info("App is Started!", gologger.GoConsole(), gologger.FileMulti("file_1"), gologger.FileMutex("file_2"))

	for scanner.Scan() {

		message := scanner.Text()

		if message == "" {

			// Выполним логирование в консоль, в том же потоке, а в файл с помощью каналов (GoFileMulti), выполним в отдельном потоке.
			// Let's log into the console, in the same thread, and to a file using channels (GoFileMulti), we'll execute it in a separate thread.
			logs.Error("Message is empty", gologger.Console(), gologger.GoFileMulti("file_1"))
		} else {

			// Выполним логирование в консоль и в файл, с помощью стандартного логгера (GoFileMutex), выполним в отдельном потоке.
			// Let's log into the console and into a file using a standard logger (GoFileMutex), and execute it in a separate thread.
			logs.Info(message, gologger.GoConsole(), gologger.GoFileMutex("file_2"))
		}
	}
}

func settings(files map[string]string) *gologger.LogInterface {

	// Создаём базового логгера. Добавляем файлы для логирования.
	// We create a basic logger. Add files for logging.
	logs := gologger.Default().ConfigFile(files)

	// Устанавливаем режим запись через каналы.
	// We set the recording mode through channels.
	logs.SetModeFileMulti()

	// Устанавливаем режим записи через логгер из стандартной библиотеки
	// We set the recording mode through the logger from the standard library
	logs.SetModeFileMutex()

	return logs
}
