package main

import (
	"fmt"
)

func main() {
	http := new(HttpClient)
	http.set_debug_mode(! true)

	url := "https://play.google.com/store/apps/details?id=com.cashslide"
	html := http.get(url)

	app := new(App)
	app.parsing(html)
	fmt.Println("-----------------------------------")
	fmt.Println(app.to_json())
	fmt.Println("-----------------------------------")


	if false {
		url = "https://play.google.com/store/apps/category/SOCIAL/collection/topselling_free?authuser=0"
		params := HttpParams{
			// "message": "HELLO",
			// "key": "178",
		}

		html = http.post(url, params)
	}


	// fmt.Println(s[:10])
}

