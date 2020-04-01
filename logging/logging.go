package logging

import (
	"context"
	"log"
)

type Data map[string]interface{}

// todo: formatter
func Info(ctx context.Context, id string, data *Data) {
	log.Printf("info: [%s] (%v)", id, data)
}

// todo: formatter
func Error(ctx context.Context, id string, err error, data *Data) {
	log.Printf("error: [%s] %s (%v)", id, err.Error(), data)
}
