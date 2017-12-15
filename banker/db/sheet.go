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
	ctx          context.Context
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
		ctx:          ctx,
		sheetID:      sheetID,
	}, nil
}

func (s *SheetDriver) Write(updateRange string, contents [][]interface{}) (*sheets.UpdateValuesResponse, error) {
	request := &sheets.ValueRange{
		Values: contents,
	}
	response, err := s.sheetService.Spreadsheets.Values.Update(s.sheetID, updateRange, request).
		ValueInputOption(SheetsDefaultValueOption).
		Do()
	if err != nil {
		return response, err
	}
	return response, err
}

func (s *SheetDriver) Read() ([]*sheets.Sheet, error) {
	sheets, err := s.sheetService.Spreadsheets.Get(s.sheetID).
		Context(s.ctx).
		Do()
	if err != nil {
		return nil, err
	}

	return sheets.Sheets, nil
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
