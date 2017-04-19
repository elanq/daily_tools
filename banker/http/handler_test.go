package http_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/mgo.v2/bson"

	banker "github.com/elanq/daily_tools/banker/http"
	"github.com/elanq/daily_tools/banker/model"
	"github.com/elanq/daily_tools/banker/mongo"
	"github.com/elanq/daily_tools/banker/parser"
	"github.com/elanq/daily_tools/banker/process"
	"github.com/pressly/chi"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
)

type HttpSuite struct {
	suite.Suite
	reader        *parser.BankReader
	bankerHandler *banker.Handler
	banker        *process.Banker
}

type TestData struct {
	path         string
	csvKey       string
	expectedCode int
	url          string
}

func TestHttpSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (h *HttpSuite) initHTTP() {
	router := chi.NewRouter()
	router.Post("/banker/upload", h.bankerHandler.FileUpload)
	router.Get("/banker/report", h.bankerHandler.MonthlyReport)
	h.banker.Router = router
}

func (h *HttpSuite) SetupSuite() {
	gotenv.Load("../env.sample")
	dBName := os.Getenv("DB_NAME")
	collectionName := "banker_test_record"

	mongoDriver := mongo.NewMongoDriver(dBName, collectionName)
	h.reader = parser.NewBankReader()
	h.bankerHandler = banker.NewHandler(h.reader, mongoDriver)
	h.banker = process.NewBanker()
	h.initHTTP()
}

func (h *HttpSuite) TearDownSuite() {
	h.cleanupData()
}

func (h *HttpSuite) TestNewHandler() {
	h.Assert().NotNil(h.bankerHandler, "Should not be nil")
	h.Assert().NotNil(h.bankerHandler.Reader, "should not be nil")
	h.Assert().EqualValues(h.bankerHandler.CSVKey, "hulahula", "env value should be same")
}

func (h *HttpSuite) prepareMultipartRequest(dir string, formKey string) *http.Request {
	buffer := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(buffer)
	filewriter, err := bodywriter.CreateFormFile(formKey, dir)

	h.Assert().Nil(err, "Should not error")

	file, err := os.Open(dir)
	h.Assert().Nil(err, "Should not error")

	_, err = io.Copy(filewriter, file)
	h.Assert().Nil(err, "Should not error")

	bodywriter.Close()
	req, err := http.NewRequest("POST", "http://localhost:12345/banker/upload", buffer)
	h.Assert().Nil(err, "should not error")
	req.Header.Add("Content-Type", bodywriter.FormDataContentType()) //this is super essential

	return req
}

func (h *HttpSuite) TestMonthlyReport() {
	var tests []*TestData
	correctTest := &TestData{
		url:          "http://localhost:12345/banker/report?month=10&year=17",
		expectedCode: http.StatusOK,
	}
	tests = append(tests, correctTest)

	paramNotSpecifiedTest := &TestData{
		url:          "http://localhost:12345/banker/report",
		expectedCode: http.StatusBadRequest,
	}
	tests = append(tests, paramNotSpecifiedTest)

	dataNotFoundTest := &TestData{
		url:          "http://localhost:12345/banker/report?month=1&year=10",
		expectedCode: http.StatusNotFound,
	}
	tests = append(tests, dataNotFoundTest)

	h.doMonthlyReportTest(tests)
}

func (h *HttpSuite) doMonthlyReportTest(tests []*TestData) {
	for _, test := range tests {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", test.url, nil)
		h.Assert().Nil(err, "should not error")

		h.banker.Router.ServeHTTP(recorder, req)
		resp := recorder.Result()
		h.Assert().Equal(test.expectedCode, resp.StatusCode, "should satisfy status code")
	}
}

func (h *HttpSuite) TestFileUpload() {
	var tests []*TestData

	correctTest := &TestData{
		path:         "../test/test_files/bank_sample.csv",
		expectedCode: http.StatusOK,
		csvKey:       h.bankerHandler.CSVKey,
	}
	tests = append(tests, correctTest)

	invalidFileTest := &TestData{
		path:         "../test/test_files/invalid_bank_sample.csv",
		expectedCode: http.StatusBadRequest,
		csvKey:       h.bankerHandler.CSVKey,
	}
	tests = append(tests, invalidFileTest)

	invalidCSVTest := &TestData{
		path:         "../test/test_files/invalid_bank_sample.csv",
		expectedCode: http.StatusBadRequest,
		csvKey:       "invalid",
	}
	tests = append(tests, invalidCSVTest)

	h.doFileUploadTest(tests)
}

func (h *HttpSuite) doFileUploadTest(tests []*TestData) {
	for _, test := range tests {
		recorder := httptest.NewRecorder()
		dir, err := filepath.Abs(test.path)
		h.Assert().Nil(err, "Should not error")

		request := h.prepareMultipartRequest(dir, test.csvKey)

		h.banker.Router.ServeHTTP(recorder, request)
		response := recorder.Result()
		h.Assert().Equal(test.expectedCode, response.StatusCode, "Should return what expected")
	}
}

func (h *HttpSuite) cleanupData() {
	var results []model.BankContent
	h.bankerHandler.MongoDriver.Find(bson.M{}, &results)

	for _, result := range results {
		h.bankerHandler.MongoDriver.Remove(result.ID)
	}
}
