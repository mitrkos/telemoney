package main

import (
	"github.com/mitrkos/telemoney/internal/app/telemoney"
	"github.com/mitrkos/telemoney/internal/pkg/logger"
)

func main() {
	logger.SetLogger()
	err := telemoney.Start()
	if err != nil {
		panic(err)
	}
}
