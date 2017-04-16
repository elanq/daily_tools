package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/elanq/daily_tools/banker/parser"
)

type Handler struct {
	Reader *parser.BankReader
	CSVKey string
}

func NewHandler(reader *parser.BankReader) *Handler {
	return &Handler{
		Reader: reader,
		CSVKey: os.Getenv("MULTIPART_CSV_KEY"),
	}
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
		w.Write([]byte("400 - Invalid or corrupted file"))
		return
	}

	h.Reader.ReadBytes(rawBytes)
	bankContents, err := h.Reader.ParseContent()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Invalid csv format"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bankContents)
}
