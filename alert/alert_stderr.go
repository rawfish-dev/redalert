package alert

import (
	"log"
	"os"

	"redalert/common"
)

type StandardError struct {
	log *log.Logger
}

func NewStandardError() StandardError {
	return StandardError{
		log: log.New(os.Stderr, "", log.Ldate|log.Ltime),
	}
}

func (a StandardError) Name() string {
	return "StandardError"
}

func (a StandardError) Trigger(alertPackage *AlertPackage) error {
	a.log.Println(alertPackage.Message)
	alertPackage.AlertLogger.Println(common.White, "Stderr alert successfully triggered.", common.Reset)
	return nil
}
