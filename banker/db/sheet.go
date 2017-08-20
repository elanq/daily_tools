package db

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	SheetsDefaultScope       = "https://www.googleapis.com/auth/spreadsheets"
	SheetsDefaultValueOption = "USER_ENTERED"
)

type SheetDriver struct {
	ctx          context.Context
	sheetService *sheets.Service
}

func NewSheetDriver() (*SheetDriver, error) {
	ctx := context.TODO()
	client, err := google.DefaultClient(ctx, SheetsDefaultScope)
	if err != nil {
		return nil, err
	}
	sheetService, err := sheets.New(client)
	if err != nil {
		return nil, err
	}

	return &SheetDriver{
		sheetService: sheetService,
		ctx:          ctx,
	}, err
}

func (s *SheetDriver) Write(sheetID string, updateRange string, contents [][]interface{}) (*sheets.UpdateValuesResponse, error) {
	request := &sheets.ValueRange{
		Values: contents,
	}
	response, err := s.sheetService.Spreadsheets.Values.Update(sheetID, updateRange, request).ValueInputOption(SheetsDefaultValueOption).Do()
	if err != nil {
		return response, err
	}
	return response, err
}

func (s *SheetDriver) Read(sheetID string) []*sheets.Sheet {
	sheets, err := s.sheetService.Spreadsheets.Get(sheetID).Context(s.ctx).Do()
	if err != nil {
		panic(err)
	}

	return sheets.Sheets
}
