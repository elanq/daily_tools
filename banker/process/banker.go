package process

import (
	"context"
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
// this struct reference needed dependencies
func NewBanker() *Banker {
	ctx := context.Background()
	dbName := os.Getenv("DB_NAME")
	collectionName := "collection_banker"
	reader := parser.NewBankReader()
	mongoDriver := db.NewMongoDriver(dbName, collectionName)
	sheetDriver, err := db.NewSheetDriver(ctx)
	if err != nil {
		panic(err)
	}
	bankerHandler := bankerhttp.NewHandler(reader, mongoDriver, sheetDriver)

	return &Banker{
		BankerHandler: bankerHandler,
		Reader:        reader,
		MongoDriver:   mongoDriver,
		SheetDriver:   sheetDriver,
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
	router.Post("/banker/report/backup", bankerHandler.Backup)
	router.Get("/banker/report/monthly", bankerHandler.MonthlyReport)
	router.Get("/banker/report/yearly", bankerHandler.YearlyReport)
	// TODO
	// monthly report endpoint should be only naratively describes current financial status
	// make new endpoint to generate fancy charts for your financial data
	// if possible, CSV upload could be scheduled properly

	return router
}
