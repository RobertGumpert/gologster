package user

import (
	"github.com/RobertGumpert/gologster"
)

var (
	logger *gologster.Logger
)

func SetLogs(main *gologster.Logger) {
	logger = main
}

func Log() {
	logger.Info("Hello from usecase")
}
