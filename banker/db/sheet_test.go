package db_test

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/elanq/daily_tools/banker/db"
	"github.com/elanq/daily_tools/banker/model"
	"github.com/subosito/gotenv"
)

var (
	sDriver     *db.SheetDriver
	bankContent *model.BankContent
)

func TestSheet(t *testing.T) {
	gotenv.Load("../.env")
	ctx := context.TODO()
	driver, err := db.NewSheetDriver(ctx)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	sDriver = driver
	bankContent = &model.BankContent{
		ID:      "1",
		Date:    time.Now(),
		Notes:   "notes",
		Branch:  "branch",
		Amount:  1090000,
		Factor:  -1,
		Balance: 849430485,
	}

	t.Run("test write", testWrite)
	t.Run("test read", testRead)
	t.Run("test batch read", testBatchRead)
}

func testWrite(t *testing.T) {
	updateRange := "TestSheet!A1:G2"
	valueRange := [][]interface{}{
		model.SheetHeader(),
		bankContent.SheetContent(),
	}

	resp, err := sDriver.Write(updateRange, valueRange)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if resp.HTTPStatusCode != 200 {
		log.Printf("http status code return %d", resp.HTTPStatusCode)
		t.FailNow()
	}
}

func testRead(t *testing.T) {
	values, err := sDriver.Read("TestSheet!A1:G2")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if values.HTTPStatusCode != 200 {
		log.Printf("invalid response from sheet API, expected 200 got %d", values.HTTPStatusCode)
		t.FailNow()
	}

	if size := len(values.Values); size < 1 {
		log.Printf("expect record size bigger than 1. got %d", size)
		t.FailNow()
	}

	contents := []model.BankContent{}
	for idx, row := range values.Values {
		//skip header
		if idx == 0 {
			continue
		}
		assignContent(&contents, row)
	}

	if size := len(contents); size < 1 {
		log.Printf("expect record size bigger than 1. got %d", size)
		t.FailNow()
	}
}

func testBatchRead(t *testing.T) {
	values, err := sDriver.BatchRead()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if size := len(values); size < 1 {
		log.Printf("expect record size bigger than 1. got %d", size)
	}

	contents := []model.BankContent{}
	for _, resp := range values {
		for idx, row := range resp.Values {
			//skip header
			if idx == 0 {
				continue
			}
			assignContent(&contents, row)
		}
	}

	if size := len(contents); size < 1 {
		log.Printf("expect record size bigger than 1. got %d", size)
		t.FailNow()
	}
}

func assignContent(contents *[]model.BankContent, row []interface{}) {
	content := model.BankContent{}

	content.Notes = row[2].(string)
	content.Branch = row[3].(string)
	content.Amount = parseNumber(row[4].(string))
	content.Factor = parseNumber(row[5].(string))
	content.Balance = parseNumber(row[6].(string))

	*contents = append(*contents, content)
}

func parseNumber(value string) int {
	val, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return val
}
