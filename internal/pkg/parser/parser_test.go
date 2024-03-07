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
				Amount:   9.5,
				Category: "lunch",
				Tags:     []string{"grenka", "dumplings"},
				Comment:  makeStringPtrInPlace("I need food!"),
			},
		},
		{
			name: "no Comment, no Tags",
			text: "9,5 lunch",
			expectedResult: &TransactionUserInputData{
				Amount:   9.5,
				Category: "lunch",
				Tags:     nil,
				Comment:  nil,
			},
		},
		{
			name: "no Comment",
			text: "9,5 lunch (grenka, dumplings)",
			expectedResult: &TransactionUserInputData{
				Amount:   9.5,
				Category: "lunch",
				Tags:     []string{"grenka", "dumplings"},
				Comment:  nil,
			},
		},
		{
			name: "no Tags",
			text: "9,5 lunch I need food!",
			expectedResult: &TransactionUserInputData{
				Amount:   9.5,
				Category: "lunch",
				Tags:     nil,
				Comment:  makeStringPtrInPlace("I need food!"),
			},
		},
		{
			name: "dot Amount separator",
			text: "9.5 lunch",
			expectedResult: &TransactionUserInputData{
				Amount:   9.5,
				Category: "lunch",
				Tags:     nil,
				Comment:  nil,
			},
		},
		{
			name: "integer Amount",
			text: "9 lunch",
			expectedResult: &TransactionUserInputData{
				Amount:   9,
				Category: "lunch",
				Tags:     nil,
				Comment:  nil,
			},
		},
		{
			name: "Category normalization",
			text: "9,5 Lunch ",
			expectedResult: &TransactionUserInputData{
				Amount:   9.5,
				Category: "lunch",
				Tags:     nil,
				Comment:  nil,
			},
		},
		{
			name: "Tags normalization",
			text: "9,5 lunch (grenkA,   Dumplings,)",
			expectedResult: &TransactionUserInputData{
				Amount:   9.5,
				Category: "lunch",
				Tags:     []string{"grenka", "dumplings"},
				Comment:  nil,
			},
		},
		{
			name: "Comment normalization",
			text: "9,5 lunch I need food!   ",
			expectedResult: &TransactionUserInputData{
				Amount:   9.5,
				Category: "lunch",
				Tags:     nil,
				Comment:  makeStringPtrInPlace("I need food!"),
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
			name: "no Amount",
			text: "lunch",
		},
		{
			name: "not valid Amount format",
			text: "9;5 lunch",
		},
		{
			name: "no Category",
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
