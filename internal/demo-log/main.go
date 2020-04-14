package main

import "github.com/spin14/copy-basta/cmd/copy-basta/log"

func main() {
	logger, err := log.NewLogger(log.WithLevel(log.Debug))
	if err != nil {
		panic(err)
	}

	msg := "Hello Logging!"
	logger.DebugWithData(msg, log.LoggerData{"With": "Data", "Seventy Five": 75})
	logger.Info(msg)
	logger.WarnWithData(msg, log.LoggerData{"Your Last": "Warning!", "True?": true})
	logger.Error(msg)
}
