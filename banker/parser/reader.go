package parser

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/elanq/daily_tools/banker/model"
)

//csv reader built specificly on klikbca format
type BankReader struct {
	Filepath   string
	RawContent string
}

//init reader
func NewBankReader() *BankReader {
	return &BankReader{
		Filepath:   "",
		RawContent: "",
	}
}

//remove unwanted content
//and save the sanitize content to struct
func (p *BankReader) sanitizeContent(raw []byte) {
	rawContent := string(raw)
	re := regexp.MustCompile("(?m)[\r\n]+^.*Mata Uang.*|Nama.*|No. Rekening.*|Saldo Awal.*|Kredit.*|Debet.*|Saldo Akhir.*")
	p.RawContent = re.ReplaceAllString(rawContent, "")
}

//return -1 if debit transaction
//return 1 if credit transaction
//default 0 if somewhat we got weird data. just don't count it
func (p *BankReader) getFactor(factor string) int {
	if factor == "DB" {
		return -1
	} else if factor == "CR" {
		return 1
	}

	return 0
}

//read file based on filepath value and sanitize the content
func (p *BankReader) ReadFile(filepath string) error {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	p.Filepath = filepath
	p.sanitizeContent(raw)

	return nil
}

// read content by slice of byte
func (p *BankReader) ReadBytes(bytes []byte) {
	p.sanitizeContent(bytes)
}

// parse date request into time.Time type
func ParseDate(rawDate string) time.Time {
	rawDate = strings.Replace(rawDate, "'", "", -1)
	// i hate this weird string to time parsing in golang
	date, _ := time.Parse("02/01/06", rawDate)

	return date
}

//the main processor of parsing content. Use this method to parse csv format into model.BankContent type
//should call ReadFile/ ReadBytes first before using this method. That way rawcontent variable won't be empty
func (p *BankReader) ParseContent(year string) ([]*model.BankContent, error) {
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

		p.fillContent(&content, &contents, record, year)
	}
	return contents, nil
}

//set parsed record to model.BankContent
func (p *BankReader) fillContent(content *model.BankContent, contents *[]*model.BankContent, record []string, year string) {
	amount, _ := strconv.ParseFloat(record[3], 32)
	balance, _ := strconv.ParseFloat(record[5], 32)

	content.ID = bson.NewObjectId()
	content.Date = ParseDate(record[0] + "/" + year)
	content.Notes = record[1]
	content.Branch = record[2]

	content.Amount = int(amount)
	content.Factor = p.getFactor(record[4])
	content.Balance = int(balance)

	*contents = append(*contents, content)

}
