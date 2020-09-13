package logger

import "time"

type zipper struct {
	fileBasicLogger *loggerBaseFile
	duration        time.Duration
}
