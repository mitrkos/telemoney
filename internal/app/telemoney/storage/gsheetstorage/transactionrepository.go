package gsheetstorage

import (
	"strings"

	"github.com/mitrkos/telemoney/internal/model"
	"github.com/mitrkos/telemoney/internal/pkg/gsheetclient"
)

type TransactionRepository struct {
	gsheetclient           *gsheetclient.GSheetsClient
	transactionSheetId     string
	transactionMessageIDScanRange *gsheetclient.A1Range
}

func New(gsheetclient *gsheetclient.GSheetsClient, transactionSheetId string) *TransactionRepository {
	// TODO: move gsheetclient creation to here
	trr := &TransactionRepository{
		gsheetclient:           gsheetclient,
		transactionSheetId:     transactionSheetId,
		transactionMessageIDScanRange: nil,
	}
	trr.transactionMessageIDScanRange = trr.makeTransactionMessageIdScanRange()
	return trr
}

func (trr *TransactionRepository) Insert(transaction *model.Transaction) error {
	trr.gsheetclient.AppendDataToRange(trr.makeRangeForSheet(nil, nil), convertTransactionToDataRow(transaction))
	return nil // TODO: add errors
}

func (trr *TransactionRepository) Update(transaction *model.Transaction) error {
	msgIDLocation, err := trr.gsheetclient.FindValueLocation(trr.transactionMessageIDScanRange, transaction.MessageID)
	if err != nil {
		return err
	}
	if msgIDLocation == nil {
		return nil // TODO: ? not found error
	}

	trr.gsheetclient.UpdateDataRange(trr.makeTransactionRowRangeFromLocation(msgIDLocation), convertTransactionToDataRow(transaction))
	return nil // TODO: add errors
}

func (trr *TransactionRepository) DeleteByMessageId(transactionMessageID string) error {
	msgIDLocation, err := trr.gsheetclient.FindValueLocation(trr.transactionMessageIDScanRange, transactionMessageID)
	if err != nil {
		return err
	}
	if msgIDLocation == nil {
		return nil // TODO: ? not found error
	}

	trr.gsheetclient.ClearRange(trr.makeTransactionRowRangeFromLocation(msgIDLocation))
	return nil // TODO: add errors
}

func (trr *TransactionRepository) makeRangeForSheet(leftTop *gsheetclient.A1Location, rightBottom *gsheetclient.A1Location) *gsheetclient.A1Range {
	return &gsheetclient.A1Range{
		SheetId:     trr.transactionSheetId,
		LeftTop:     leftTop,
		RightBottom: rightBottom,
	}
}

func (trr *TransactionRepository) makeTransactionMessageIdScanRange() *gsheetclient.A1Range {
	return trr.makeRangeForSheet(&gsheetclient.A1Location{
		Column: "B", // TODO: add schema mapping
		Row:    3,
	}, &gsheetclient.A1Location{
		Column: "B",
		Row:    0,
	})
}

func (trr *TransactionRepository) makeTransactionRowRangeFromLocation(location *gsheetclient.A1Location) *gsheetclient.A1Range {
	return trr.makeRangeForSheet(&gsheetclient.A1Location{
		Column: "A", // TODO: add schema mapping
		Row:    location.Row,
	}, &gsheetclient.A1Location{
		Column: "F",
		Row:    location.Row,
	})
}

func convertTransactionToDataRow(transaction *model.Transaction) []interface{} {
	dataRow := make([]interface{}, 6)

	dataRow[0] = transaction.CreatedAt
	dataRow[1] = transaction.MessageID
	dataRow[2] = transaction.Amount
	dataRow[3] = transaction.Category
	if len(transaction.Tags) > 0 {
		tagsStr := strings.Join(transaction.Tags[:], ",")
		dataRow[4] = tagsStr
	}
	if transaction.Comment != nil {
		dataRow[5] = *transaction.Comment
	}

	return dataRow
}
