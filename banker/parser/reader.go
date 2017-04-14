package parser

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/elanq/daily_tools/banker/model"
)

type BankReader struct {
	Filepath   string
	RawContent string
}

func NewBankReader() *BankReader {
	return &BankReader{
		Filepath:   "",
		RawContent: "",
	}
}

func (p *BankReader) sanitizeContent(raw []byte) string {
	rawContent := string(raw)
	re := regexp.MustCompile("(?m)[\r\n]+^.*Mata Uang.*|Nama.*|No. Rekening.*|Saldo Awal.*|Kredit.*|Debet.*|Saldo Akhir.*")
	res := re.ReplaceAllString(rawContent, "")

	return res
}

func (p *BankReader) getFactor(factor string) int {
	if factor == "DB" {
		return -1
	} else if factor == "CR" {
		return 1
	}

	return 0
}

func (p *BankReader) ReadFile(filepath string) error {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	p.Filepath = filepath
	p.RawContent = p.sanitizeContent(raw)

	return nil
}

func (p *BankReader) ParseContent() ([]*model.BankContent, error) {
	var contents []*model.BankContent

	if p.RawContent == "" {
		return nil, errors.New("Raw content not available")
	}

	csvReader := csv.NewReader(strings.NewReader(p.RawContent))
	skipFirst := true
	for {
		record, csvErr := csvReader.Read()
		if skipFirst {
			skipFirst = false
			continue
		}

		content := model.BankContent{}
		if csvErr == io.EOF {
			break
		}
		if csvErr != nil {
			return nil, csvErr
		}

		amount, _ := strconv.ParseFloat(record[3], 32)
		balance, _ := strconv.ParseFloat(record[5], 32)

		content.Date = record[0]
		content.Notes = record[1]
		content.Branch = record[2]

		content.Amount = int(amount)
		content.Factor = p.getFactor(record[4])
		content.Balance = int(balance)

		contents = append(contents, &content)

	}
	return contents, nil
}
