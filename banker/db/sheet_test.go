package db_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/elanq/daily_tools/banker/db"
	"github.com/elanq/daily_tools/banker/model"
	"github.com/subosito/gotenv"
	"google.golang.org/api/sheets/v4"
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
	sheets, err := sDriver.Read()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	for _, s := range sheets {
		if s.Properties.Title == "TestSheet" {
			traverse(s.Data)
		}
	}
}

func traverse(grid []*sheets.GridData) {
	for _, sData := range grid {
		traverseRow(sData.RowData)
	}
}

func traverseRow(rows []*sheets.RowData) {
	for _, rData := range rows {
		traverseCell(rData.Values)
	}
}

func traverseCell(cells []*sheets.CellData) {
	for _, cData := range cells {
		log.Println(cData.FormattedValue)
	}
}
