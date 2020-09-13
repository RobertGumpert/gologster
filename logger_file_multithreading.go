package logger

import (
	"bufio"
	"errors"
	"os"
	"runtime"
	"strings"
)

// loggerFileMultithreading : логгер в файл с использованием очереди на запись. | logger to file using write queue.
//
// Идея заключается в том, что для каждого из файлов создаётся
// буфферизированный канал на 1000 элементов (строк, которые надо записать в файл).
// Существует только одна горутина-читатель 'loggerFileMultithreading.receiver()',
// для этого канала, которая имеет право вызвать функцию записи в файл, что
// гарантирует то, что не возникнет ситуация гонки.
// Буффер на 1000 элементов, позволит повысить производительность, так как
// вероятность того, что этот буффер переполнится, теоритически, очень мала,
// как следствие все горутины-писатели, после записи в канал, завершают работу
// и не занимают место в очереди 'sendq', тем самым не занимая место в памяти.
//
// The idea is that for each of the files
// a buffered channel is created with 1000 elements (lines to be written to the file).
// There is only one goroutine-reader 'loggerFileMultithreading.receiver()',
// for this channel, which has the right to call the function of writing to the file,
// which ensures that there is no race situation.
// A buffer for 1000 elements will improve performance,
// since the probability that this buffer will overflow is theoretically very small,
// as a result, all goroutines-writers, after writing to the channel,
// exit and do not take up space in the 'sendq' queue, thereby without taking up memory space.
//
type loggerFileMultithreading struct {
	baseFile *loggerBaseFile
	config   map[string]fileAgent
}

// newLoggerFileMultithreading : constructor
//
func newLoggerFileMultithreading(baseFile *loggerBaseFile) *loggerFileMultithreading {
	logger := new(loggerFileMultithreading)
	logger.config = make(map[string]fileAgent, 0)
	logger.baseFile = baseFile
	for key, path := range baseFile.config {
		logger.newFile(key, path)
	}
	return logger
}

func (logger *loggerFileMultithreading) newFile(key, path string) {
	if _, exist := logger.config[key]; !exist {
		logger.baseFile.config[key] = path
	}
	file := fileAgent{
		path:    path,
		channel: make(chan *string, 1000),
	}
	logger.config[key] = file
	go logger.receiver(file)
}

// receiver : итерируется по каналу, созданному для конкретного файла.
//			  iterates over the pipe created for the specific file.
//
func (logger *loggerFileMultithreading) receiver(file fileAgent) {
	for outputString := range file.channel {
		runtime.Gosched()
		err := logger.output(outputString, file.path)
		if err != nil {
			logger.errorOutput(outputString, err)
		}
	}
}

// add : implement iLogger interface
//
// После создания строки лога, она записывается в канал,
// созданный для конкретного файла.
// Следом, горутина, записавшая эту строку,
// завершает работу (так как, каналы буферизированные и их буфер равен 1000 элементам,
// вероятность, того, что горутина-писатель встанет в очередь sendq,
// занимая место в памати, очень мала).
// Эту строку считает горутина 'loggerFileMultithreading.receiver()',
// запущенная для конкретного файла, в теле которой находится цикл,
// итерирующийся по каналу, вызвая функцию записи в файл.
//
// After creating a callingMode line, it is written to a channel
// created for a specific file.
// Next, the goroutine that wrote this line exits
// (since the channels are buffered and their buffer is 1000 elements,
// the probability that the writer goroutine will be queued sendq,
// taking up memory space is very small).
// This line is considered by the goroutine 'loggerFileMultithreading.receiver()',
// launched for a specific file, in the body of which there is a loop,
// iterating over the channel, calling the function to write to the file.
//
func (logger *loggerFileMultithreading) add(log *logData, param ...string) {
	var (
		performOutput = func(log *logData, logger *loggerFileMultithreading, key string) {
			out, err := logger.createOutputString(log)
			if err != nil {
				logger.baseFile.errorOutput(out, err)
				return
			}
			if file, exist := logger.config[key]; exist {
				file.channel <- out
				return
			} else {
				logger.errorOutput(out, errors.New("loggerFileMultithreading.add isn't exist by key : '"+key+"'"))
			}
		}
		fileKey = ""
	)
	if log.IsOption {
		err, key := logger.baseFile.getParams(log, param...)
		if err != nil {
			return
		}
		fileKey = key
	} else {
		for pckg := range logger.baseFile.config {
			if strings.Contains(log.Package, pckg) {
				fileKey = pckg
				break
			}
		}
	}
	performOutput(log, logger, fileKey)
}

// errorOutput : implement iLogger interface
//
// Исрользуется / используйте 'baseFile.createOutputString()'
//
// Use 'baseFile.createOutputString()'
//
func (logger *loggerFileMultithreading) createOutputString(log *logData, param ...string) (*string, error) {
	return logger.baseFile.createOutputString(log, param...)
}

// output : implement iLogger interface
//
// После открытия файла, создаёт временный буффер,
// в который будет выполняться запись.
// Следом, вызывается системный вызов fsync(),
// для сборса буферов файловой системы на диск.
// После записи содержимого в буффер,
// данные сбрасываются в файл через '(*io.Writer).Flush()'.
//
// After opening the file, it creates a temporary buffer to write to.
// Next, the fsync() system call is called
// to collect the file system buffers to disk.
// After writing the content to the buffer,
// the data is flushed to the file via '(*io.Writer).Flush()'.
//
func (logger *loggerFileMultithreading) output(out *string, param ...string) error {
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
		return closeFile(file, err)
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

// errorOutput : implement iLogger interface
//
// Исрользуется / используйте 'baseFile.errorOutput()'
//
// Use 'baseFile.errorOutput()'
//
func (logger *loggerFileMultithreading) errorOutput(out *string, err error) {
	logger.baseFile.errorOutput(out, err)
}
