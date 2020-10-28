package main

import (
	"copy-basta/internal/common/log"
)

func main() {
	logger := log.NewLogger()
	lvl, err := log.ToLevel(log.Debug)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(lvl)

	msg := "Hello Logging!"
	logger.DebugWithData(msg, log.Data{"With": "Data", "Seventy Five": 75})
	logger.Info(msg)
	logger.WarnWithData(msg, log.Data{"Your Last": "Warning!", "True?": true})
	logger.Error(msg)
}
