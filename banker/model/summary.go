package model

import (
	"strconv"
)

type Summary struct {
	TotalIncome    int
	TotalOutcome   int
	CurrentBalance int
	Profit         int
	MostIncome     BankContent
	MostOutcome    BankContent
}

func (s *Summary) Print() string {
	return "your income this month Rp " + strconv.Itoa(s.TotalIncome) +
		"\nyour expenditure this month Rp " + strconv.Itoa(s.TotalOutcome) +
		"\nyour current balance now is Rp " + strconv.Itoa(s.CurrentBalance) +
		"\nyour most income for about Rp " + strconv.Itoa(s.MostIncome.Amount) +
		"\ncomes with detail " + s.MostIncome.Notes +
		"\nhowever, your most expenditure for about Rp " + strconv.Itoa(s.MostOutcome.Amount) +
		"\ncomes with detail " + s.MostOutcome.Notes +
		"\nand this month you saved about Rp " + strconv.Itoa(s.Profit)
}
