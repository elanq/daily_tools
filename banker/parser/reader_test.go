package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/elanq/daily_tools/banker/parser"
	"github.com/stretchr/testify/assert"
)

func TestNewBankReader(t *testing.T) {
	reader := parser.NewBankReader()
	assert.NotNil(t, reader, "Should not return nil")
}

func TestReadFile(t *testing.T) {
	reader := parser.NewBankReader()

	correctDir, err := filepath.Abs("../test/test_files/bank_sample.csv")
	err = reader.ReadFile(correctDir)
	assert.Nil(t, err, "Should not return any error")

	wrongDir, wrongErr := filepath.Abs("../test/test_files/bank_sample.go")
	wrongErr = reader.ReadFile(wrongDir)
	assert.Error(t, wrongErr, "Should error because not existent file")
}

func TestParseContent(t *testing.T) {
	reader := parser.NewBankReader()
	invalidReader := parser.NewBankReader()

	correctDir, err := filepath.Abs("../test/test_files/bank_sample.csv")
	err = reader.ReadFile(correctDir)
	assert.Nil(t, err, "Should not return any error")

	invalidDir, invalidErr := filepath.Abs("../test/test_files/invalid_bank_sample.csv")
	invalidErr = invalidReader.ReadFile(invalidDir)

	records, err := reader.ParseContent()
	assert.Nil(t, err, "Should not return any error")
	for _, record := range records {
		assert.NotNil(t, record.Date, "Should contains value")
		assert.NotNil(t, record.Notes, "Should contains value")
		assert.NotNil(t, record.Branch, "Should contains value")
		assert.NotNil(t, record.Amount, "Should contains value")
		assert.NotNil(t, record.Factor, "Should contains value")
		assert.NotNil(t, record.Balance, "Should contains value")
	}

	_, invalidErr = invalidReader.ParseContent()
	assert.Error(t, invalidErr, "Should return csv error")

	invalidReader.RawContent = ""
	_, invalidErr = invalidReader.ParseContent()
	assert.Error(t, invalidErr, "Should return content error")

}
