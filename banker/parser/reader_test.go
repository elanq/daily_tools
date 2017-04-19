package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/elanq/daily_tools/banker/parser"
	"github.com/stretchr/testify/suite"
)

type ReaderSuite struct {
	suite.Suite
	Reader        *parser.BankReader
	InvalidReader *parser.BankReader
}

func (r *ReaderSuite) SetupSuite() {
	r.Reader = parser.NewBankReader()
	r.InvalidReader = parser.NewBankReader()
}

func TestReaderSuite(t *testing.T) {
	suite.Run(t, new(ReaderSuite))
}

func (r *ReaderSuite) TestNewBankReader() {
	r.Assert().NotNil(r.Reader, "should be not nil")
}

func (r *ReaderSuite) TestReadFile() {
	correctDir, err := filepath.Abs("../test/test_files/bank_sample.csv")
	r.Assert().Nil(err, "Should not return any error")
	err = r.Reader.ReadFile(correctDir)
	r.Assert().Nil(err, "Should not return any error")

	wrongDir, wrongErr := filepath.Abs("../test/test_files/bank_sample.go")
	r.Assert().Nil(wrongErr, "Should not return any error")
	wrongErr = r.Reader.ReadFile(wrongDir)
	r.Assert().Error(wrongErr, "Should error because not existent file")
}

func (r *ReaderSuite) TestParseContent() {
	correctDir, err := filepath.Abs("../test/test_files/bank_sample.csv")
	r.Assert().Nil(err, "Should not return any error")
	err = r.Reader.ReadFile(correctDir)
	r.Assert().Nil(err, "Should not return any error")

	invalidDir, invalidErr := filepath.Abs("../test/test_files/invalid_bank_sample.csv")
	r.Assert().Nil(invalidErr, "Should not return any error")
	invalidErr = r.InvalidReader.ReadFile(invalidDir)
	r.Assert().Nil(invalidErr, "Should not return any error")

	records, err := r.Reader.ParseContent("17")
	r.Assert().Nil(err, "Should not return any error")
	for _, record := range records {
		r.Assert().NotNil(record.Date, "Should contains value")
		r.Assert().NotNil(record.Notes, "Should contains value")
		r.Assert().NotNil(record.Branch, "Should contains value")
		r.Assert().NotNil(record.Amount, "Should contains value")
		r.Assert().NotNil(record.Factor, "Should contains value")
		r.Assert().NotNil(record.Balance, "Should contains value")
	}

	_, invalidErr = r.InvalidReader.ParseContent("17")
	r.Assert().Error(invalidErr, "Should return csv error")

	r.Reader.RawContent = ""
	_, invalidErr = r.Reader.ParseContent("2017")
	r.Assert().Error(invalidErr, "Should return content error")
}

func (r *ReaderSuite) TestParseDate() {
	testDate := "01/12/17"
	date := parser.ParseDate(testDate)

	day := date.Day()
	month := date.Month().String()
	year := date.Year()

	r.Assert().Equal(1, day, "should be day 1")
	r.Assert().Equal("December", month, "should be December")
	r.Assert().Equal(2017, year, "should be year 2017")
}

// TODO : TestReadByte
