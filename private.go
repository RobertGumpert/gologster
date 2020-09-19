package gologster

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

const (
	BaseLogTemplate string = "level=[{{.Level}}];func=[name: {{.Func}}, line: {{.Line}}, package:{{.Package}}];value=[{{.Value}}];date=[{{.Date}}];"
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
	// of callingMode line creation, output and error handling.
	//
	add(data *logData, param ...string)

	// createOutputString : метод, ответственный за создание строки лога. | method responsible for creating the callingMode line.
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
	// Creates a callingMode line and returns an error
	// if the passed data is not valid or
	// some error occurred while marshaling.
	//
	createOutputString(data *logData, param ...string) (*string, error)

	// output : метод, выполняющий вывод лога. | method that performs callingMode output.
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
	// Performs callingMode output, for example, to a file,
	// if a logger implementing this interface works with files.
	// Returns an error if an error occurred at the time of output.
	//
	output(out *string, param ...string) error

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
	errorOutput(out *string, err error)
}

type logData struct {
	UserDataOriginal                        interface{}
	Lvl                                     level
	IsOption                                bool
	Error                                   error
	Value, Level, Package, Date, Func, Line string
}

func newLogData(value interface{}, lvl level, date string) *logData {
	log := new(logData)
	log.UserDataOriginal = value
	log.Date = date
	log.Lvl = lvl
	log.Level = toStringLevel(lvl)
	return log
}

func (log *logData) marshal(base *loggerBase) error {
	out, err := base.createOutputString(log)
	if err != nil {
		log.Value = "marshal error"
		return err
	}
	log.Value = *out
	return err
}

func (log *logData) setRuntimeInfo(skip int) *logData {
	function, pckg, line := getRuntimeInfo(skip)
	log.Func = function
	log.Package = pckg
	log.Line = line
	return log
}

func (log *logData) filledTemplate(tmpl *template.Template) *string {
	var (
		out    = ""
		buffer = new(bytes.Buffer)
	)
	err := tmpl.Execute(buffer, log)
	if err != nil {
		baseTemplate := getTextTemplate("base", BaseLogTemplate, BaseLogTemplate)
		_ = baseTemplate.Execute(buffer, log)
	}
	out = buffer.String()
	return &out
}

func getTextTemplate(name, str, alternative string) *template.Template {
	var (
		textTemplate *template.Template
	)
	if str == "" {
		textTemplate, _ = template.New(name).Parse(alternative)
	} else {
		t, err := template.New(name).Parse(str)
		if err != nil {
			textTemplate, _ = template.New("console").Parse(alternative)
		}
		textTemplate = t
	}
	return textTemplate
}

func toStringLevel(lvl level) string {
	switch lvl {
	case levelInfo:
		return "INFO"
	case levelError:
		return "ERROR"
	case levelPanic:
		return "PANIC"
	default:
		return "NON"
	}
}

func getRuntimeInfo(skip int) (string, string, string) {
	var (
		function = "undefined func"
		line     = "-1"
		pckg     = "undefined package"
	)
	pc, _, lineInt, ok := runtime.Caller(skip)
	if !ok {
		return function, pckg, line
	}
	function = runtime.FuncForPC(pc).Name()
	if strings.Contains(function, "/") {
		//
		split := strings.Split(function, "/")
		function = split[len(split)-1]
		//
		functionSplit := strings.Split(function, ".")
		function = functionSplit[len(functionSplit)-1]
		//
		split = split[0:len(split)-1]
		split = append(split, functionSplit[0])
		pckg = strings.Join(split, "/")
	} else {
		if strings.Contains(function, ".") {
			functionSplit := strings.Split(function, ".")
			function = functionSplit[len(functionSplit)-1]
			pckg = functionSplit[0]
		}
	}
	return function, pckg, strconv.Itoa(lineInt)
}
