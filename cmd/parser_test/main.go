package main

import (
	"github.com/mitrkos/telemoney/internal/pkg/parser"
)

func main() {
	parser := parser.New()

    parser.ParseTransactionUserInputDataFromText("9,5 lunch (grenka, dumplings) I need food!")
}
