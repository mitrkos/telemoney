package main

import "github.com/mitrkos/telemoney/internal/app/telemoney"

func main() {
	err := telemoney.Start()
	if err != nil {
		panic(err)
	}
}
