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

	//driver and content setup
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
	t.Run("test sheet creation", testCreate)
	t.Run("test backup", testBackup)
}

func testWrite(t *testing.T) {
	ctx := context.TODO()
	updateRange := "TestSheet!A1:G2"
	valueRange := [][]interface{}{
		model.SheetHeader(),
		bankContent.SheetContent(),
	}

	resp, err := sDriver.Write(ctx, updateRange, valueRange)
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
	ctx := context.TODO()
	values, err := sDriver.Read(ctx, "TestSheet!A1:G2")
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
	ctx := context.TODO()
	values, err := sDriver.BatchRead(ctx)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if size := len(values); size < 1 {
		log.Printf("expect record size bigger than 1. got %d", size)
		t.FailNow()
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

func testBackup(t *testing.T) {
	ctx := context.TODO()
	valueRange := [][]interface{}{
		model.SheetHeader(),
		bankContent.SheetContent(),
	}
	sheetName := "testBackupSheet"

	sheetID, err := sDriver.Backup(ctx, sheetName, valueRange)

	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if sheetID == 0 {
		log.Println("sheet ID is 0")
		t.FailNow()
	}

	if err = sheetCleanup(ctx, sheetID); err != nil {
		log.Println(err)
		t.FailNow()
	}
}

func testCreate(t *testing.T) {
	var createdSheetID int64
	ctx := context.TODO()
	sheetName := "testBackupSheet"

	res, err := sDriver.CreateSheet(ctx, sheetName)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if res.HTTPStatusCode != 200 {
		log.Println(err)
		t.FailNow()
	}

	for _, rep := range res.Replies {
		createdSheetID = rep.AddSheet.Properties.SheetId
		//assume replies only contains 1 value
		break
	}

	if createdSheetID == 0 {
		log.Println("failed to get created sheet ID")
		log.Println(res)
		t.FailNow()
	}

	if err = sheetCleanup(ctx, createdSheetID); err != nil {
		log.Println(err)
		t.FailNow()
	}
}

func sheetCleanup(ctx context.Context, sheetID int64) error {
	//delete newly created sheet
	_, err := sDriver.DeleteSheet(ctx, sheetID)
	if err != nil {
		return err
	}

	return nil
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
