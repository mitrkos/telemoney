package parser

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/oriser/regroup"
)

type Parser struct {
	regexp *regroup.ReGroup
}

type TransactionUserInputData struct {
	Amount   float64
	Category string
	Tags     []string
	Comment  *string
}

func New() *Parser {
	regexp := regroup.MustCompile(
		`^(?P<amount>\d+[\.,]?\d*) (?P<category>\w+) ?(?:\((?P<tags>[\w, ]*)\))?(?P<comment>.*$)?`,
	) // to parse "9,5 lunch (grenka, dumplings) I need foood!"
	return &Parser{
		regexp: regexp,
	}
}

func (p *Parser) ParseTransactionUserInputDataFromText(text string) (*TransactionUserInputData, error) {
	result, err := func() (*TransactionUserInputData, error) {
		match, err := p.regexp.Groups(text)
		if err != nil {
			return nil, err
		}

		amountRaw, ok := match["amount"]
		amount, err := parseAndValidateTransactionAmount(amountRaw, ok)
		if err != nil {
			return nil, err
		}

		categoryRaw, ok := match["category"]
		category, err := parseAndValidateTransactionCategory(categoryRaw, ok)
		if err != nil {
			return nil, err
		}

		tagsRaw, ok := match["tags"]
		tags := parseAndValidateTransactionTags(tagsRaw, ok)

		commentRaw, ok := match["comment"]
		comment := parseAndValidateTransactionComment(commentRaw, ok)

		return &TransactionUserInputData{
			Amount:   amount,
			Category: category,
			Tags:     tags,
			Comment:  comment,
		}, nil
	}()

	if err != nil {
		slog.Error("Parsing TransactionUserInputData failed", slog.Any("error", err), slog.Any("input", text))
	} else {
		slog.Info("Parsed TransactionUserInputData", slog.Any("result", result), slog.Any("input", text))
	}

	return result, err
}

func parseAndValidateTransactionAmount(amountRaw string, ok bool) (float64, error) {
	if !ok {
		return 0, errors.New("amount is not found")
	}
	amountRaw = strings.Replace(amountRaw, ",", ".", 1)
	return strconv.ParseFloat(amountRaw, 64)
}

func parseAndValidateTransactionCategory(categoryRaw string, ok bool) (string, error) {
	if !ok {
		return "", errors.New("category is not found")
	}
	category := strings.ToLower(strings.TrimSpace(categoryRaw))
	// TODO: check with allowed categories
	return category, nil
}

func parseAndValidateTransactionTags(tagsRaw string, ok bool) []string {
	if !ok {
		return nil
	}

	var tags []string

	splittedTagsRaw := strings.Split(tagsRaw, ",")
	for _, tagRaw := range splittedTagsRaw {
		normilizedTag := strings.ToLower(strings.TrimSpace(tagRaw))
		if len(normilizedTag) > 0 {
			tags = append(tags, strings.ToLower(strings.TrimSpace(tagRaw)))
		}
	}

	if len(tags) > 0 {
		return tags
	}

	return nil
}

func parseAndValidateTransactionComment(commentRaw string, ok bool) *string {
	if !ok {
		return nil
	}
	comment := strings.TrimSpace(commentRaw)
	if comment != "" {
		return &comment
	}
	return nil
}
