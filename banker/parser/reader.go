package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/elanq/banker/model"
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
	// lineCounter := 0
	// for i, bit := range raw {
	// 	if lineCounter == 3 {
	// 		break
	// 	}
	// 	//if reached endline
	// 	if bit == 10 {
	// 		lineCounter++
	// 	}
	// 	// delete this elemnt
	// 	raw = append(raw[:i], raw[i+1:]...)
	//
	// }
	rawContent := string(raw)
	re := regexp.MustCompile("(?m)[\r\n]+^.*Mata Uang.*|Nama.*|No. Rekening.*|Saldo Awal.*|Kredit.*|Debet.*|Saldo Akhir.*")
	res := re.ReplaceAllString(rawContent, "")
	// regex := bytes.Replace(raw, []byte("No. Rekening"), []byte(""), 1)

	return res
}

func (p *BankReader) ReadFile(filepath string) error {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	p.Filepath = filepath
	p.RawContent = p.sanitizeContent(raw)

	csvReader := csv.NewReader(strings.NewReader(p.RawContent))
	for {
		record, csvErr := csvReader.Read()
		if csvErr == io.EOF {
			break
		}
		if csvErr != nil {
			return csvErr
		}
		fmt.Println(record)
	}
	return nil
}

func (p *BankReader) ParseContent() ([]*model.BankContent, error) {
	return nil, nil
}
