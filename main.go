package main

import (
	"fmt"
	"go-market-crawler/lib"
	"go-market-crawler/models"
)

func main() {
	http := new(lib.HttpClient)
	http.SetDebugMode(true)

	url := "https://play.google.com/store/apps/details?id=com.cashslide"
	html := http.Get(url)

	app := new(models.App)
	app.Parsing(html)
	fmt.Println("-----------------------------------")
	fmt.Println(app.ToJson())
	fmt.Println("-----------------------------------")

	apps := make(models.Apps, 1)
	apps = append(apps, *app)
	apps = append(apps, *app)
	fmt.Println(apps.ToJson())

	fmt.Println("-----------------------------------")

	if false {
		url = "https://play.google.com/store/apps/category/SOCIAL/collection/topselling_free?authuser=0"
		params := lib.HttpParams{
		// "message": "HELLO",
		// "key": "178",
		}

		html = http.Post(url, params)
		fmt.Println(s[:10])
	}
}
