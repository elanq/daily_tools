package process

import (
	"net/http"
	"os"
	"time"

	bankerhttp "github.com/elanq/daily_tools/banker/http"
	"github.com/elanq/daily_tools/banker/mongo"
	"github.com/elanq/daily_tools/banker/parser"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type Banker struct {
	BankerHandler *bankerhttp.Handler
	Reader        *parser.BankReader
	MongoDriver   *mongo.MongoDriver
	Router        http.Handler
}

func NewBanker() *Banker {
	dbName := os.Getenv("DB_NAME")
	collectionName := "collection_banker"
	reader := parser.NewBankReader()
	bankerHandler := bankerhttp.NewHandler(reader)
	mongoDriver := mongo.NewMongoDriver(dbName, collectionName)
	return &Banker{
		BankerHandler: bankerHandler,
		Reader:        reader,
		MongoDriver:   mongoDriver,
		Router:        setRouter(bankerHandler),
	}
}

func setRouter(bankerHandler *bankerhttp.Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Post("/banker/upload", bankerHandler.FileUpload)

	return router
}