package logger

import "strings"

type outputString string

func newOutputString(value []byte, date string, lvl level) *outputString {
	out := outputString("")
	out = out.addLevel(lvl)
	out = out.addValue(value)
	out = out.addDate(date)
	return &out
}

func (out outputString) addError(err error) outputString {
	if strings.Index(string(out), "level=") == 0 && strings.Contains(string(out), "date=") && strings.Contains(string(out), "value=") {
		return outputString(strings.Join([]string{
			string(out),
			toStringError(err)}, ""))
	}
	return out
}

func (out outputString) addDate(date string) outputString {
	if strings.Index(string(out), "level=") == 0 && !strings.Contains(string(out), "date=") {
		return outputString(strings.Join([]string{
			string(out),
			"date=[",
			date,
			"];"}, ""))
	}
	return out
}

func (out outputString) addValue(value []byte) outputString {
	if strings.Index(string(out), "level=") == 0 && !strings.Contains(string(out), "value=") {
		return outputString(strings.Join([]string{
			string(out),
			toStringValue(value)}, ""))
	}
	return out
}

func (out outputString) addLevel(lvl level) outputString {
	if !strings.Contains(string(out), "level=") && strings.Count(string(out), ";") == 0 {
		return outputString(strings.Join([]string{
			string(out),
			toStringLevel(lvl)}, ""))
	}
	return out
}

func toStringLevel(lvl level) string {
	switch lvl {
	case levelInfo:
		return "level=[INFO];"
	case levelError:
		return "level=[ERROR];"
	case levelPanic:
		return "level=[PANIC];"
	default:
		return "level=[NON];"
	}
}

func toStringValue(value []byte) string {
	return strings.Join([]string{
		"value=[",
		string(value),
		"];"}, "")
}

func toStringError(err error) string {
	return strings.Join([]string{
		"error=[",
		err.Error(),
		"];",
	}, "")
}
