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
	// var contents []*model.BankContent
	// path := "/Users/eq/Downloads/sample.csv"
	// reader := parser.NewBankReader()
	// err := reader.ReadFile(path)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println("Can't read file")
	// }
	//
	// contents, err = reader.ParseContent()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// var firstBalance int
	// var lastBalance int
	// for i, content := range contents {
	// 	if i == 0 {
	// 		firstBalance = content.Balance
	// 	}
	// 	if i == len(contents)-1 {
	// 		lastBalance = content.Balance
	// 	}
	// 	fmt.Println(content.Amount)
	// }
	//
	// fmt.Println(firstBalance)
	// fmt.Println(lastBalance)
}
