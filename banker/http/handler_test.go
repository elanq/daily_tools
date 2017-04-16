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

	banker "github.com/elanq/daily_tools/banker/http"
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

func TestHttpSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (h *HttpSuite) initHTTP() {
	router := chi.NewRouter()
	router.Post("/banker/upload", h.bankerHandler.FileUpload)

	h.banker.Router = router
}

func (h *HttpSuite) SetupSuite() {
	gotenv.Load("../env.sample")
	h.reader = parser.NewBankReader()
	h.bankerHandler = banker.NewHandler(h.reader)
	h.banker = process.NewBanker()
	h.initHTTP()
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
	recorder := httptest.NewRecorder()
	correctDir, err := filepath.Abs("../test/test_files/bank_sample.csv")
	h.Assert().Nil(err, "Should not error")

	correctRequest := h.prepareMultipartRequest(correctDir, h.bankerHandler.CSVKey)

	h.banker.Router.ServeHTTP(recorder, correctRequest)
	response := recorder.Result()
	h.Assert().Equal(200, response.StatusCode, "Should return 200 OK")

	// incorrect dir
	invalidDirRecorder := httptest.NewRecorder()
	invalidDir, err := filepath.Abs("../test/test_files/invalid_bank_sample.csv")
	h.Assert().Nil(err, "Should not error")

	invalidDirRequest := h.prepareMultipartRequest(invalidDir, h.bankerHandler.CSVKey)

	h.banker.Router.ServeHTTP(invalidDirRecorder, invalidDirRequest)
	invalidResponse := invalidDirRecorder.Result()
	h.Assert().Equal(400, invalidResponse.StatusCode, "Should return 400 bad request ")
}
