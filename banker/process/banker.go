package process

import (
	"net/http"
	"os"
	"time"

	"github.com/elanq/daily_tools/banker/db"
	bankerhttp "github.com/elanq/daily_tools/banker/http"
	"github.com/elanq/daily_tools/banker/parser"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type Banker struct {
	BankerHandler *bankerhttp.Handler
	Reader        *parser.BankReader
	MongoDriver   *db.MongoDriver
	Router        http.Handler
}

// Initiate Banker struct
// thist struct reference needed dependencies
func NewBanker() *Banker {
	dbName := os.Getenv("DB_NAME")
	collectionName := "collection_banker"
	reader := parser.NewBankReader()
	mongoDriver := db.NewMongoDriver(dbName, collectionName)
	bankerHandler := bankerhttp.NewHandler(reader, mongoDriver)

	return &Banker{
		BankerHandler: bankerHandler,
		Reader:        reader,
		MongoDriver:   mongoDriver,
		Router:        setRouter(bankerHandler),
	}
}

// Set route for the service
func setRouter(bankerHandler *bankerhttp.Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Post("/banker/upload", bankerHandler.FileUpload)
	router.Get("/banker/report/monthly", bankerHandler.MonthlyReport)
	router.Get("/banker/report/yearly", bankerHandler.YearlyReport)
	// TODO
	// monthly report endpoint should be only naratively describes current financial status
	// make new endpoint to generate fancy charts for your financial data
	// if possible, CSV upload is should be scheduled properly

	return router
}
