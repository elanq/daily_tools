package utility_test

import (
	"testing"
	"time"

	"github.com/elanq/daily_tools/banker/model"
	"github.com/elanq/daily_tools/banker/utility"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSummary(t *testing.T) {
	var contents []model.BankContent

	emptySummary := utility.GenerateSummary(contents)
	assert.NotNil(t, emptySummary, "should never return nil")

	fillData(&contents)
	summary := utility.GenerateSummary(contents)
	assert.NotNil(t, summary, "should not return nil")
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
