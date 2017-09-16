package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/elanq/daily_tools/banker/db"
	"github.com/elanq/daily_tools/banker/model"
	"github.com/elanq/daily_tools/banker/parser"
	"github.com/elanq/daily_tools/banker/utility"
)

const (
	ErrInternalError         = errors.New("error while processing request")
	ErrContentNotFound       = errors.New("content not found")
	ErrTransactionExists     = errors.New("Transaction already exists")
	ErrYearNotExists         = errors.New("should specify year as param")
	ErrYearAndMonthNotExists = errors.New("should specify month and year as param")
	ErrInvalidForm           = errors.New("Invalid form key")
	ErrInvalidFile           = errors.New("Cannot read specified file")
	ErrInvalidCSV            = errors.New("Invalid csv format")
	ErrorSavingToDb          = errors.New("Error while saving file to db")
)

//Handler type. used to reference csv reader, csv key, year key and mongodriver
type Handler struct {
	Reader      *parser.BankReader
	CSVKey      string
	YearKey     string
	MongoDriver *db.MongoDriver
}

//create new type of request handler
func NewHandler(reader *parser.BankReader, driver *db.MongoDriver) *Handler {
	return &Handler{
		Reader:      reader,
		MongoDriver: driver,
		CSVKey:      os.Getenv("MULTIPART_CSV_KEY"),
		YearKey:     "year",
	}
}

//Insert content to db
func (h *Handler) saveContent(bankContents []*model.BankContent) error {
	if h.checkContent(bankContents[0], bankContents[len(bankContents)-1]) {
		return ErrTransactionExists
	}

	for _, content := range bankContents {
		err := h.MongoDriver.Insert(content)
		if err != nil {
			return err
		}
	}
	return nil
}

//TODO
//1. Handle daily report as charts
//2. Generate brief summary with information like
// -> highest income
// -> highest outcome
// -> most spending in a day
// -> total balance
//3. accept these request types
// -> no type will give raw json format of mutation record
// -> type=monthly_summary to give return as summary
// -> type=chart will give chart of income and expenditure per day. Will presented as bar chart

//provide transaction data per day by given year
// TODO: need to aggregate into monthly data tho
func (h *Handler) YearlyReport(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	var results []model.BankContent

	if year == "" {
		http.Error(w, ErrYearNotExists.Error(), http.StatusBadRequest)
		return
	}
	minTime := parser.ParseDate("01/01/" + year)
	maxTime := minTime.AddDate(1, 0, 0)

	err := h.fetchContent(minTime, maxTime, &results)

	if err != nil {
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	if len(results) < 1 {
		http.Error(w, ErrContentNotFound.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if returnSummary(&w, r, &results) {
		return
	}

	json.NewEncoder(w).Encode(results)
}

//provide transaction data per day by given month
func (h *Handler) MonthlyReport(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	var results []model.BankContent

	if month == "" || year == "" {
		http.Error(w, ErrYearAndMonthNotExists.Error(), http.StatusBadRequest)
		return
	}

	minTime := parser.ParseDate("01/" + month + "/" + year)
	maxTime := minTime.AddDate(0, 1, 0)

	err := h.fetchContent(minTime, maxTime, &results)
	if err != nil {
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	if len(results) < 1 {
		http.Error(w, ErrContentNotFound.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if returnSummary(&w, r, &results) {
		return
	}

	json.NewEncoder(w).Encode(results)
}

func returnSummary(w *http.ResponseWriter, r *http.Request, results *[]model.BankContent) bool {
	contentType := r.URL.Query().Get("type")
	if contentType != "summary" {
		return false
	}
	summary := utility.GenerateSummary(*results)
	json.NewEncoder(*w).Encode(summary)

	return true
}

func (h *Handler) FileUpload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile(h.CSVKey)
	if err != nil {
		http.Error(w, ErrInvalidForm.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()
	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, ErrInvalidFile.Error(), http.StatusBadRequest)
		return
	}

	year := r.FormValue(h.YearKey)
	if year == "" {
		year = "17"
	}

	h.Reader.ReadBytes(rawBytes)
	bankContents, err := h.Reader.ParseContent(year)
	if err != nil {
		http.Error(w, ErrInvalidCSV.Error(), http.StatusBadRequest)
		return
	}

	err = h.saveContent(bankContents)
	if err != nil {
		http.Error(w, ErrorSavingToDb.Error(), http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bankContents)
}

//populate bank content
func (h *Handler) fetchContent(minTime time.Time, maxTime time.Time, results *[]model.BankContent) error {
	query := bson.M{
		"date": bson.M{
			"$gte": minTime,
			"$lt":  maxTime,
		},
	}

	return h.MongoDriver.Find(query, results, "-date")
}

func (h *Handler) checkContent(firstRow *model.BankContent, lastRow *model.BankContent) bool {
	if firstRow == nil && lastRow == nil {
		return false
	}

	var results []model.BankContent

	query := bson.M{
		"date": bson.M{
			"$gte": firstRow.Date,
			"$lte": lastRow.Date,
		},
	}

	h.MongoDriver.Find(query, &results)
	if len(results) > 0 {
		return true
	}
	return false
}
