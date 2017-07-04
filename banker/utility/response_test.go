package utility_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/elanq/daily_tools/banker/model"
	"github.com/elanq/daily_tools/banker/utility"
	"github.com/stretchr/testify/assert"
)

func TestPrintSummaryType(t *testing.T) {
	var contents []model.BankContent

	response := utility.PrintSummaryType(contents)
	assert.Equal(t, "No data available", response, "should return same value")
	fillData(&contents)
	validResponse := utility.PrintSummaryType(contents)
	fmt.Println(validResponse)
	assert.NotEqual(t, "", validResponse, "should return summary")
}

func fillData(contents *[]model.BankContent) {
	for i := 0; i < 4; i++ {
		var factor int
		if i%2 == 0 {
			factor = -1
		} else {
			factor = 1
		}
		content := &model.BankContent{
			Date:    time.Now(),
			Notes:   "Notes",
			Branch:  "1234",
			Amount:  1000000,
			Factor:  factor,
			Balance: 2000000,
		}
		*contents = append(*contents, *content)
	}
}
