package main

import (
	"github.com/mitrkos/telemoney/internal/app/telemoney"
	"github.com/mitrkos/telemoney/internal/pkg/logger"
)

func main() {
	logger.SetLogger()
	deps, err := telemoney.PrepareDependencies()
	if err != nil {
		panic(err)
	}

	t := telemoney.New(deps.Config, deps.Api, deps.TransactionStorage, deps.Parser)
	err = t.Start()
	if err != nil {
		panic(err)
	}
}
