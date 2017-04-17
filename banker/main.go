package main

import (
	"net/http"

	"github.com/elanq/daily_tools/banker/process"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	app := process.NewBanker()
	http.ListenAndServe(":12345", app.Router)
}
