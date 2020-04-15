package main

import (
	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

func main() {
	logger, err := common.NewLogger(common.WithLevel(common.Debug))
	if err != nil {
		panic(err)
	}

	msg := "Hello Logging!"
	logger.DebugWithData(msg, common.LoggerData{"With": "Data", "Seventy Five": 75})
	logger.Info(msg)
	logger.WarnWithData(msg, common.LoggerData{"Your Last": "Warning!", "True?": true})
	logger.Error(msg)
}
