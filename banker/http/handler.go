package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/elanq/daily_tools/banker/model"
	"github.com/elanq/daily_tools/banker/mongo"
	"github.com/elanq/daily_tools/banker/parser"
)

type Handler struct {
	Reader      *parser.BankReader
	CSVKey      string
	YearKey     string
	MongoDriver *mongo.MongoDriver
}

func NewHandler(reader *parser.BankReader, driver *mongo.MongoDriver) *Handler {
	return &Handler{
		Reader:      reader,
		MongoDriver: driver,
		CSVKey:      os.Getenv("MULTIPART_CSV_KEY"),
		YearKey:     "year",
	}
}

func (h *Handler) saveContent(bankContents []*model.BankContent) error {
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
//4. change this method to dauly reports. as it would be more appropriate

func (h *Handler) MonthlyReport(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	var results []model.BankContent

	if year == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("should specify year as param"))
		return
	}
	minTime := parser.ParseDate("01/01/" + year)
	maxTime := minTime.AddDate(1, 0, 0)

	err := h.contentQuery(minTime, maxTime, &results)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while processing request"))
		return
	}

	if len(results) < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write(([]byte("404 - Content not found")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) DailyReport(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	var results []model.BankContent

	if month == "" || year == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("should specify month and year as param"))
		return
	}

	minTime := parser.ParseDate("01/" + month + "/" + year)
	maxTime := minTime.AddDate(0, 1, 0)

	err := h.contentQuery(minTime, maxTime, &results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while processing request"))
		return
	}

	if len(results) < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write(([]byte("404 - Content not found")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) FileUpload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile(h.CSVKey)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("400 - Invalid form key"))
		return
	}

	defer file.Close()

	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("400 - Cannot read specified file"))
		return
	}

	year := r.FormValue(h.YearKey)
	if year == "" {
		year = "17"
	}

	h.Reader.ReadBytes(rawBytes)
	bankContents, err := h.Reader.ParseContent(year)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Invalid csv format"))
		return
	}

	err = h.saveContent(bankContents)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		w.Write([]byte("304 - Error while saving file to db"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bankContents)
}

func (h *Handler) contentQuery(minTime time.Time, maxTime time.Time, results *[]model.BankContent) error {
	query := bson.M{
		"date": bson.M{
			"$gte": minTime,
			"$lt":  maxTime,
		},
	}

	return h.MongoDriver.Find(query, results)
}
