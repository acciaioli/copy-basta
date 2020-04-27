package main

import (
	"copy-basta/cmd/copy-basta/common/log"
)

func main() {
	logger := log.NewLogger()
	logger.SetLevel(log.Debug)

	msg := "Hello Logging!"
	logger.DebugWithData(msg, log.Data{"With": "Data", "Seventy Five": 75})
	logger.Info(msg)
	logger.WarnWithData(msg, log.Data{"Your Last": "Warning!", "True?": true})
	logger.Error(msg)
}
