package http_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/elanq/daily_tools/banker/db"
	banker "github.com/elanq/daily_tools/banker/http"
	"github.com/elanq/daily_tools/banker/model"
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
	testName     string
}

func TestHttpSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (h *HttpSuite) initHTTP() {
	router := chi.NewRouter()
	router.Post("/banker/upload", h.bankerHandler.FileUpload)
	router.Get("/banker/report/monthly", h.bankerHandler.MonthlyReport)
	router.Get("/banker/report/yearly", h.bankerHandler.YearlyReport)
	h.banker.Router = router
}

func (h *HttpSuite) SetupSuite() {
	gotenv.Load("../env.sample")
	dBName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("DB_COLLECTION")

	mongoDriver := db.NewMongoDriver(dBName, collectionName)
	h.reader = parser.NewBankReader()
	h.bankerHandler = banker.NewHandler(h.reader, mongoDriver)
	h.banker = process.NewBanker()
	h.initHTTP()
	h.fillTestData()
}

func (h *HttpSuite) TearDownSuite() {
	h.cleanupData()
}

func (h *HttpSuite) fillTestData() {
	date, _ := time.Parse("02/01/06", "01/10/17")
	content := &model.BankContent{
		Date:    date,
		Notes:   "Notes",
		Branch:  "1234",
		Amount:  1000000,
		Factor:  1,
		Balance: 2000000,
	}
	h.bankerHandler.MongoDriver.Insert(content)
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

func (h *HttpSuite) TestFileUpload() {
	var tests []*TestData

	correctTest := &TestData{
		path:         "../test/test_files/bank_sample.csv",
		expectedCode: http.StatusOK,
		csvKey:       h.bankerHandler.CSVKey,
		testName:     "test:upload:correct",
	}
	tests = append(tests, correctTest)

	invalidFileTest := &TestData{
		path:         "../test/test_files/invalid_bank_sample.csv",
		expectedCode: http.StatusBadRequest,
		csvKey:       h.bankerHandler.CSVKey,
		testName:     "test:upload:invalid_file",
	}
	tests = append(tests, invalidFileTest)

	invalidCSVTest := &TestData{
		path:         "../test/test_files/invalid_bank_sample.csv",
		expectedCode: http.StatusBadRequest,
		csvKey:       "invalid",
		testName:     "test:upload:invalid_csv",
	}
	tests = append(tests, invalidCSVTest)

	duplicateTest := &TestData{
		path:         "../test/test_files/bank_sample.csv",
		expectedCode: http.StatusNotModified,
		csvKey:       h.bankerHandler.CSVKey,
	}
	tests = append(tests, duplicateTest)

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

func (h *HttpSuite) TestReport() {
	var tests []*TestData
	monthlyCorrectTest := &TestData{
		url:          "http://localhost:12345/banker/report/monthly?month=10&year=17",
		expectedCode: http.StatusOK,
		testName:     "test:monthly:success",
	}
	tests = append(tests, monthlyCorrectTest)

	monthlyParamNotSpecifiedTest := &TestData{
		url:          "http://localhost:12345/banker/report/monthly",
		expectedCode: http.StatusBadRequest,
		testName:     "test:monthly:unspecified_param",
	}
	tests = append(tests, monthlyParamNotSpecifiedTest)

	monthlyDataNotFoundTest := &TestData{
		url:          "http://localhost:12345/banker/report/monthly?month=1&year=10",
		expectedCode: http.StatusNotFound,
		testName:     "test:monthly:data_not_found",
	}
	tests = append(tests, monthlyDataNotFoundTest)

	yearlyCorrectTest := &TestData{
		url:          "http://localhost:12345/banker/report/yearly?year=17",
		expectedCode: http.StatusOK,
		testName:     "test:yearly:success",
	}
	tests = append(tests, yearlyCorrectTest)

	yearlyParamNotSpecifiedTest := &TestData{
		url:          "http://localhost:12345/banker/report/yearly",
		expectedCode: http.StatusBadRequest,
		testName:     "test:yearly:unspecified_param",
	}
	tests = append(tests, yearlyParamNotSpecifiedTest)

	yearlyDataNotFoundTest := &TestData{
		url:          "http://localhost:12345/banker/report/yearly?year=10",
		expectedCode: http.StatusNotFound,
		testName:     "test:yearly:data_not_found",
	}
	tests = append(tests, yearlyDataNotFoundTest)

	summaryTypeTest := &TestData{
		url:          "http://localhost:12345/banker/report/monthly?month=10&year=17&type=summary",
		expectedCode: http.StatusOK,
		testName:     "test:monthly:summary_type",
	}
	tests = append(tests, summaryTypeTest)

	summaryTypeYearlyTest := &TestData{
		url:          "http://localhost:12345/banker/report/yearly?year=17&type=summary",
		expectedCode: http.StatusOK,
		testName:     "test:yearly:summary_type",
	}
	tests = append(tests, summaryTypeYearlyTest)

	h.doReportTest(tests)
}

func (h *HttpSuite) doReportTest(tests []*TestData) {
	for _, test := range tests {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", test.url, nil)
		h.Assert().Nil(err, "should not error")

		h.banker.Router.ServeHTTP(recorder, req)
		resp := recorder.Result()
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(test.testName, "return : ", string(b))
		h.Assert().Equal(test.expectedCode, resp.StatusCode, "should satisfy status code")
	}
}

func (h *HttpSuite) cleanupData() {
	var results []model.BankContent
	h.bankerHandler.MongoDriver.Find(bson.M{}, &results)

	for _, result := range results {
		h.bankerHandler.MongoDriver.Remove(result.ID)
	}
}
