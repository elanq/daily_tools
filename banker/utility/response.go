package utility

import (
	"github.com/elanq/daily_tools/banker/model"
)

func MonthlySummary(contents []model.BankContent) *model.Summary {
	income := 0
	outcome := 0
	summary := &model.Summary{}

	if len(contents) < 1 {
		return summary
	}

	for _, c := range contents {
		if c.Factor < 0 {
			summary.TotalOutcome += c.Amount
			if outcome < c.Amount {
				summary.MostOutcome = c
				outcome = c.Amount
			}
		} else {
			summary.TotalIncome += c.Amount
			if income < c.Amount {
				summary.MostIncome = c
				income = c.Amount
			}
		}
	}

	//because transactions are sorted descending by default
	summary.CurrentBalance = contents[0].Balance
	summary.Profit = summary.TotalIncome - summary.TotalOutcome

	return summary
}
