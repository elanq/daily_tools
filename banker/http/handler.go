package http

import (
	"encoding/json"
	"fmt"
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
	file, header, err := r.FormFile(h.CSVKey)
	fmt.Println("incoming request with header ", header.Header)
	if err != nil {
		return
	}

	defer file.Close()

	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	h.Reader.ReadBytes(rawBytes)
	bankContents, err := h.Reader.ParseContent()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bankContents)
}
