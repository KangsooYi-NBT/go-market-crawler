package main

import (
	"fmt"
	// "sort"
	"go-market-crawler/lib"
	"go-market-crawler/models"
)

func main() {
	http := new(lib.HttpClient)
	http.SetDebugMode(true)

	url := ""
	html := ""

	if !true {
		url = "https://play.google.com/store/apps/details?id=com.cashslide"
		html = http.Get(url)

		app := new(models.App)
		app.PackageId = "com.cashslide"
		app.Parsing(html)
		fmt.Println("-----------------------------------")
		fmt.Println(app.ToJson())
		fmt.Println("-----------------------------------")

		// apps := make(models.Apps, 1)
		// apps = append(apps, *app)
		// apps = append(apps, *app)
		// fmt.Println(apps.ToJson())
		// fmt.Println("-----------------------------------")
		return
	}

	// SYNC
	if !true {
		//인기 소셜 앱 리스트
		url = "https://play.google.com/store/apps/category/SOCIAL/collection/topselling_free?authuser=0"
		params := lib.HttpParams{
		// "message": "HELLO",
		// "key": "178",
		}
		html = http.Post(url, params)
		a := new(models.AppsCategory)
		packageIds := a.Parsing(html)

		// AppsCategory에서 추출한 PackageID를 순서대로 순회해서 Google에서 Fetch 후 Apps에 추가
		apps := make(models.Apps, 0)
		for _, package_id := range packageIds {
			// fmt.Println(package_id)
			url = "https://play.google.com/store/apps/details?id=" + package_id
			html = http.Get(url)

			app := new(models.App)
			app.PackageId = package_id
			app.Parsing(html)
			apps = append(apps, *app)
		}

		fmt.Println(apps.ToJson())
	}

	// ASYNC by Go Routine
	if true {
		//인기 소셜 앱 리스트
		url = "https://play.google.com/store/apps/category/SOCIAL/collection/topselling_free?authuser=0"
		params := lib.HttpParams{
		// "message": "HELLO",
		// "key": "178",
		}

		html = http.Post(url, params)
		a := new(models.AppsCategory)
		packageIds := a.Parsing(html)
		packageCnt := len(packageIds)

		messages := make(chan models.App, packageCnt)
		apps := make(models.Apps, 0)

		// go routine으로 PackageId별 App 정보 수집
		for category_no, package_id := range packageIds {
			go func(package_id string, category_no int) {
				url = "https://play.google.com/store/apps/details?id=" + package_id
				html = http.Get(url)

				app := new(models.App)
				app.PackageId = package_id
				app.Parsing(html)
				app.CategoryRank = category_no

				messages <- *app
			}(package_id, category_no+1)
		}

		// go routine 대기
		for i := 0; i < packageCnt; i++ {
			select {
			case app := <-messages:
				apps = append(apps, app)
			}
		}

		//카테고리 순서대로 정렬
		apps.SortByCategoryRank()

		fmt.Println(apps.ToJson())
	}
}
