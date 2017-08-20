package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/elanq/daily_tools/banker/db"
	"github.com/elanq/daily_tools/banker/model"
)

const (
	sheetId = "1eGVdPz1OdPE8btM9J7_y0eRWej-nD7KQTgqyzo3swO8"
)

func TestSheet(t *testing.T) {
	sDriver, err := db.NewSheetDriver()
	if err != nil {
		t.Fail()
	}
	sheets := sDriver.Read(sheetId)
	for _, s := range sheets {
		raw, _ := s.MarshalJSON()
		fmt.Println(string(raw))
	}

	updateRange := "Sheet1!A1:G2"
	content := &model.BankContent{
		ID:      "1",
		Date:    time.Now(),
		Notes:   "notes",
		Branch:  "branch",
		Amount:  1090000,
		Factor:  -1,
		Balance: 849430485,
	}
	valueRange := [][]interface{}{
		model.SheetHeader(),
		content.SheetContent(),
	}
	resp, err := sDriver.Write(sheetId, updateRange, valueRange)
	if err != nil {

		fmt.Println(err)
		t.FailNow()
	}
	raw, _ := resp.MarshalJSON()
	fmt.Println(string(raw))
	t.Fail()
}
