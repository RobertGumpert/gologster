# gologger - описание | description.

Логгер создавался для того, чтобы одним вызовом функции логгирования, можно было писать сразу в разные накопители.
Одновременная запись, одного лога и в файл и в консоль, с возможностью записывать в отдельном потоке, например только в файл,
а в консоль, только в потоке где была вызвана функция логирования, или вообще все записывать в отдельном потоке,
не заботясь о формате вывода, так как логгер сам создаёт строку вывода в нужном формате.

The logger was created so that with one call to the logging function, it was possible to write to different drives (hard disk, console) at once.
Simultaneous recording of one message to a file and to the console, with the ability to write in a separate stream, for example, only to a file, and to the console, only in the stream where the recording function was called, or write everything in a separate stream, without worrying about the output format, so how the logger itself creates the output string in the desired format.

# Примеры | Examples

Объектом, через которое выполняется логирование является 'UserInterface'. Это объект необходимо создать только один раз в приложении и прокидвать указатель в другие пакеты. 

The object through which logging is performed is 'UserInterface'. This object needs to be created only once in the application and passed the pointer to other packages.

## Базовая настройка.

В базовой настройке доступен только вывод в консоль. Все остальные настройки устанавливаются после получения базовой настройки.

In the basic setting, only console output is available. All other settings are set after receiving the basic setting.

**Базовая настройка | Basic setting :**
```go

// Создаём базового логгера, в котором доступен только вывод в консоль.
// We create a basic gologger in which only output to the console is available.
var logs = gologger.Default()

func main() {

	// Выполним логирование, в этом же потоке.
	// Let's perform logging in the same thread.
	logs.Info("App is started!")
	
	logs.Info("App has terminated!")
}

[OUTPUT]:
2020/08/19 23:31:59 level=[INFO];value=["App is started!"];date=[Wed Aug 19 23:31:59 2020];
2020/08/19 23:31:59 level=[INFO];value=["App has terminated!"];date=[Wed Aug 19 23:31:59 2020];

```

**Базовая настройка. Вывод в отдельном потоке. | Basic setting. Basic setup. Output in a separate thread :**
```go

func PrintNumbers(num int) {

	prev := new(number)

	for i := 0; i < num; i++ {

		currentNumber := &number{
			Value: i,
			Prev:  prev,
		}

		// Распечатаем сообщение в отдельном потоке.
		// Let's print the message in a separate thread.
		logs.Info(currentNumber, gologger.GoConsole())

		prev = currentNumber
	}

	logs.Info("Print all numbers : " + strconv.Itoa(num))
}

```

# Особенности | Features.

Для записи в файл существует две реализации:

- через каналы.

- с помощью стандартного пакета 'log'.

**Запись в файл через каналы.**

Особенностью этого решения является то, что для каждого из файлов создаётся буфферизированный канал на 1000 элементов (строк, которые надо записать в файл).
Для каждого такого канала, запускается в отдельном потоке горутина-читатель, которая имеет право вызвать функцию записи в файл,
что гарантирует то, что не возникнет ситуация гонки. Буффер на 1000 элементов теоритически достаточно большой, для того чтобы горутины писатели не вставали в очередь
на запись в буффер канала.

A feature of this solution is that for each of the files, a buffered channel is created for 1000 elements (lines that must be written to the file).
For each such channel, a reader goroutine is launched in a separate thread, which has the right to call the function of writing to the file, which guarantees that a race situation does not arise. The 1000-element buffer is theoretically large enough to prevent writers from queuing up to write to the channel buffer.

*Горутина-читатель | goroutine-reader :*
```go
func (logger *loggerFileMultithreading) receiver(file fileAgent) {
	for outputString := range file.channel {
		runtime.Gosched()
		err := logger.output(outputString, file.path)
		if err != nil {
			logger.errorOutput(outputString, err)
		}
	}
}
```

После открытия файла, создаёт временный буффер,в который будет выполняться запись.Следом, вызывается системный вызов fsync(),для сборса буферов файловой системы на диск.
После записи содержимого в буффер, данные сбрасываются в файл через '(*io.Writer).Flush()'.

After opening the file, it creates a temporary buffer to write to. Next, the fsync() system call is called to collect the file system buffers to disk.
After writing the content to the buffer, the data is flushed to the file via '(*io.Writer).Flush()'.

*Запись в файл | Write to file :*
```go
func (logger *loggerFileMultithreading) output(out *outputString, param ...string) error {
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
		return err
	}
	buffer := bufio.NewWriter(file)
	err = file.Sync()
	if err != nil {
		return closeFile(file, err)
	}
	_, err = buffer.WriteString(string(*out) + "\n")
	if err != nil {
		return closeFile(file, err)
	}
	err = buffer.Flush()
	if err != nil {
		return closeFile(file, err)
	}
	return closeFile(file, nil)
}
```

![alt text](https://github.com/RobertGumpert/gologger/blob/master/examples/channel.png)


**Запись с помощью стандартного пакета 'log'.**

Используется стандартный пакет 'log', который разрешает состояние гонки с помощью мьютексов.

The standard package 'log' is used, which resolves race conditions using mutexes.

*Запись в файл | Write to file :*
```go
func (logger *loggerFileMutex) output(out *outputString, param ...string) error {
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
	log.Print(*out)
	return closeFile(file, nil)
}
```
