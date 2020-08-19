# gologger - описание | description.

Логгер создавался для того, чтобы одним вызовом функции логгирования, можно было писать сразу в разные накопители.
Одновременная запись, одного лога и в файл и в консоль, с возможностью записывать в отдельном потоке, например только в файл,
а в консоль, только в потоке где была вызвана функция логирования, или вообще все записывать в отдельном потоке,
не заботясь о формате вывода, так как логгер сам создаёт строку вывода в нужном формате.

The logger was created so that with one call to the logging function, it was possible to write to different drives (hard disk, console) at once.
Simultaneous recording of one message to a file and to the console, with the ability to write in a separate stream, for example, only to a file, and to the console, only in the stream where the recording function was called, or write everything in a separate stream, without worrying about the output format, so how the logger itself creates the output string in the desired format.

# Особенности | Features.

Для записи в файл существует две реализации:

- через каналы.

- с помощью стандартного пакета 'log'.

**Запись в файл через каналы.**

Особенностью этого решения является то, что для каждого из файлов создаётся буфферизированный канал на 1000 элементов (строк, которые надо записать в файл).
Для каждого такого канала, запускается в отдельном потоке горутина-читатель, которая имеет право вызвать функцию записи в файл,
что гарантирует то, что не возникнет ситуация гонки. Буффер на 1000 элементов теоритически достаточно большой, для того чтобы горутины писатели не вставали в очередь
на запись в буффер канала.

Горутины-читатель | goroutine-reader
```
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

A feature of this solution is that for each of the files, a buffered channel is created for 1000 elements (lines that must be written to the file).
For each such channel, a reader goroutine is launched in a separate thread, which has the right to call the function of writing to the file, which guarantees that a race situation does not arise. The 1000-element buffer is theoretically large enough to prevent writers from queuing up to write to the channel buffer.

![alt text](https://github.com/RobertGumpert/gologger/blob/master/examples/channel.png)