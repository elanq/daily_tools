package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/elanq/daily_tools/banker/model"
	"github.com/elanq/daily_tools/banker/mongo"
	"github.com/elanq/daily_tools/banker/parser"
)

type Handler struct {
	Reader      *parser.BankReader
	CSVKey      string
	MongoDriver *mongo.MongoDriver
}

func NewHandler(reader *parser.BankReader, driver *mongo.MongoDriver) *Handler {
	return &Handler{
		Reader:      reader,
		MongoDriver: driver,
		CSVKey:      os.Getenv("MULTIPART_CSV_KEY"),
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

	h.Reader.ReadBytes(rawBytes)
	bankContents, err := h.Reader.ParseContent()
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
