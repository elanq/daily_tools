package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/subosito/gotenv"
)

var client *http.Client

func main() {
	loadConfig()
	access_token := getAccessToken()
	fmt.Println(access_token)
}

func loadConfig() {
	gotenv.Load()
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func getAccessToken() string {
	uname := os.Getenv("REDDIT_USER")
	pwd := os.Getenv("REDDIT_PWD")
	app_key := os.Getenv("REDDIT_APP_KEY")
	app_secret := os.Getenv("REDDIT_APP_SECRET")

	url := "https://www.reddit.com/api/v1/access_token"
	payload := strings.NewReader("grant_type=password&username=" + uname + "&password=" + pwd)

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	req.SetBasicAuth(app_key, app_secret)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("user-agent", "story_scrapper/0.1 by elanq")

	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if res.StatusCode != 200 {
		return ""
	}

	body, _ := ioutil.ReadAll(res.Body)
	var rawMap map[string]interface{}

	err = json.Unmarshal(body, &rawMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return rawMap["access_token"].(string)

}

func printPersonalInfo() {
}
