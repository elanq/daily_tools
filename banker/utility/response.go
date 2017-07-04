package utility

import (
	"strconv"

	"github.com/elanq/daily_tools/banker/model"
)

type Summary struct {
	totalIncome    int
	totalOutcome   int
	currentBalance int
	profit         int
	mostIncome     model.BankContent
	mostOutcome    model.BankContent
}

func PrintSummaryType(contents []model.BankContent) string {
	var response string
	income := 0
	outcome := 0
	summary := &Summary{}

	if len(contents) < 1 {
		response = "No data available"
		return response
	}

	for _, c := range contents {
		if c.Factor < 0 {
			summary.totalOutcome += c.Amount
			if outcome < c.Amount {
				summary.mostOutcome = c
				outcome = c.Amount
			}
		} else {
			summary.totalIncome += c.Amount
			if income < c.Amount {
				summary.mostIncome = c
				income = c.Amount
			}
		}

	}

	summary.currentBalance = contents[len(contents)-1].Balance
	summary.profit = summary.totalIncome - summary.totalOutcome
	response = "your income this month Rp " + strconv.Itoa(summary.totalIncome) +
		"\nyour expenditure this month Rp " + strconv.Itoa(summary.totalOutcome) +
		"\nyour current balance now is Rp " + strconv.Itoa(summary.currentBalance) +
		"\nyour most income for about Rp " + strconv.Itoa(summary.mostIncome.Amount) +
		"\ncomes with detail " + summary.mostIncome.Notes +
		"\nhowever, your most expenditure for about Rp " + strconv.Itoa(summary.mostOutcome.Amount) +
		"\ncomes with detail " + summary.mostOutcome.Notes +
		"\nand this month you saved about Rp " + strconv.Itoa(summary.profit)

	return response
}
