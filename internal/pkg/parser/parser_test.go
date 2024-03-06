package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makeStringPtrInPlace(v string) *string { return &v }

func TestParser_ParseTransactionUserInputDataFromTextSuccess(t *testing.T) {
	parser := New()

	testCases := []struct {
		name           string
		text           string
		expectedResult *TransactionUserInputData
	}{
		{
			name: "full string",
			text: "9,5 lunch (grenka, dumplings) I need food!",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     []string{"grenka", "dumplings"},
				comment:  makeStringPtrInPlace("I need food!"),
			},
		},
		{
			name: "no comment, no tags",
			text: "9,5 lunch",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     nil,
				comment:  nil,
			},
		},
		{
			name: "no comment",
			text: "9,5 lunch (grenka, dumplings)",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     []string{"grenka", "dumplings"},
				comment:  nil,
			},
		},
		{
			name: "no tags",
			text: "9,5 lunch I need food!",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     nil,
				comment:  makeStringPtrInPlace("I need food!"),
			},
		},
		{
			name: "dot amount separator",
			text: "9.5 lunch",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     nil,
				comment:  nil,
			},
		},
		{
			name: "category normalization",
			text: "9,5 Lunch ",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     nil,
				comment:  nil,
			},
		},
		{
			name: "tags normalization",
			text: "9,5 lunch (grenkA,   Dumplings,)",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     []string{"grenka", "dumplings"},
				comment:  nil,
			},
		},
		{
			name: "comment normalization",
			text: "9,5 lunch I need food!   ",
			expectedResult: &TransactionUserInputData{
				amount:   9.5,
				category: "lunch",
				tags:     nil,
				comment:  makeStringPtrInPlace("I need food!"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.ParseTransactionUserInputDataFromText(tc.text)
			require.Nil(t, err)
			require.Equal(t, tc.expectedResult, result, "TransactionUserInputDatas aren't equal")
		})
	}
}

func TestParser_ParseTransactionUserInputDataFromTextValidationError(t *testing.T) {
	parser := New()

	testCases := []struct {
		name          string
		text          string
	}{
		{
			name: "no amount",
			text: "lunch",
		},
		{
			name: "not valid amount format",
			text: "9;5 lunch",
		},
		{
			name: "no category",
			text: "9,5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.ParseTransactionUserInputDataFromText(tc.text)
			require.Error(t, err)
		})
	}
}
