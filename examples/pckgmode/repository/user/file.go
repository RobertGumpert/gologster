package user

import (
	gologger "../../../.."
)

var (
	logger *gologger.Logger
)

func SetLogs(main *gologger.Logger) {
	logger = main
}

func Log() {
	logger.Info("Hello from repository")
}

