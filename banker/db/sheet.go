package db

import (
	"errors"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	SheetsDefaultScope       = "https://www.googleapis.com/auth/spreadsheets"
	SheetsDefaultValueOption = "USER_ENTERED"
)

type SheetDriver struct {
	sheetService *sheets.Service
	sheetID      string
}

func NewSheetDriver(ctx context.Context) (*SheetDriver, error) {
	sheetID, err := exportCredential()
	if err != nil {
		return nil, err
	}

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
		sheetID:      sheetID,
	}, nil
}

func (s *SheetDriver) Backup(ctx context.Context, sheetName string, contents [][]interface{}) (int64, error) {
	var sheetID int64
	updateRange := sheetName + "!A:Z"

	resp, err := s.CreateSheet(ctx, sheetName)
	if err != nil {
		return 0, err
	}

	for _, rep := range resp.Replies {
		//assume replies only contains 1 entity
		sheetID = rep.AddSheet.Properties.SheetId
		break
	}

	if _, err := s.Write(ctx, updateRange, contents); err != nil {
		return 0, err
	}

	return sheetID, nil
}

func (s *SheetDriver) DeleteSheet(ctx context.Context, sheetID int64) (*sheets.BatchUpdateSpreadsheetResponse, error) {
	deleteRequest := &sheets.Request{
		DeleteSheet: &sheets.DeleteSheetRequest{
			SheetId: sheetID,
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{deleteRequest},
	}

	resp, err := s.sheetService.Spreadsheets.BatchUpdate(s.sheetID, batchUpdateRequest).Context(ctx).Do()

	return resp, err
}

func (s *SheetDriver) CreateSheet(ctx context.Context, sheetName string) (*sheets.BatchUpdateSpreadsheetResponse, error) {
	addRequest := &sheets.AddSheetRequest{
		Properties: &sheets.SheetProperties{
			Title: sheetName,
		},
	}

	request := &sheets.Request{
		AddSheet: addRequest,
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{request},
	}

	res, err := s.sheetService.Spreadsheets.BatchUpdate(s.sheetID, batchUpdateRequest).Context(ctx).Do()

	return res, err
}

func (s *SheetDriver) Write(ctx context.Context, updateRange string, contents [][]interface{}) (*sheets.UpdateValuesResponse, error) {
	request := &sheets.ValueRange{
		Values: contents,
	}
	response, err := s.sheetService.Spreadsheets.Values.Update(s.sheetID, updateRange, request).
		ValueInputOption(SheetsDefaultValueOption).
		Context(ctx).
		Do()
	if err != nil {
		return response, err
	}
	return response, err
}

func (s *SheetDriver) Read(ctx context.Context, updateRange string) (*sheets.ValueRange, error) {
	values, err := s.sheetService.Spreadsheets.Values.Get(s.sheetID, updateRange).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (s *SheetDriver) BatchRead(ctx context.Context) ([]*sheets.ValueRange, error) {
	ranges := []string{"A:Z"}
	values, err := s.sheetService.Spreadsheets.Values.BatchGet(s.sheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return values.ValueRanges, err
}

func exportCredential() (string, error) {
	sheetID := os.Getenv("GOOGLE_SHEET_ID")
	if sheetID == "" {
		return "", errors.New("can't read GOOGLE_SHEET_ID value, have you set it?")
	}

	credentialLocation := os.Getenv("GOOGLE_CREDENTIAL_LOCATION")
	if credentialLocation == "" {
		return "", errors.New("can't read GOOGLE_CREDENTIAL_LOCATION value, have you set it?")
	}

	if _, err := os.Stat(credentialLocation); err != nil {
		return "", err
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialLocation)
	return sheetID, nil
}
