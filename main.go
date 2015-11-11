package main

import (
	"fmt"
//	"go-market-crawler/lib"
	"go-market-crawler/models"
	"regexp"
	"io/ioutil"
	"time"
	"strings"
	"sync"
)

func main() {
	play()
//	test_regex()
}
func test_regex() {
	filename := "/tmp/google_play_categorySOCIAL.html"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	html := string(content)

	pattern := `<a class="title-link id-track-click" (?:.+?) href="(.+?)">\s*(.+?)\s*</a>\s*<\/h2>`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(html, -1)
	if matches != nil && len(matches) > 0 {
		fmt.Println("Matches: %d", len(matches))
		for _, row := range matches {
			fmt.Println(len(row))
			fmt.Println(row[1], " ### ", row[2])
		}
	}
}

func play() {
	fmt.Printf("##### [%s] ###############################################\n", (time.Now()).Format("2006-01-30"))

	var wg sync.WaitGroup

	go_routine_cnt := 0


	m := new(models.GooglePlay)
	m.Init()
	categories := m.ExtractCategories()

	url := ""
	for _, category := range categories {

		tmp := strings.Split(category.Url, "?")

		for _, selling_type := range []string{"topselling_free", "topselling_paid", "topgrossing"} {
			if len(tmp) == 2 {
				url = tmp[0] + "/collection/" + selling_type + "?" + tmp[1]
			} else {
				url = category.Url + "/collection/" + selling_type
			}
//			fmt.Println(category.Name, category.Text, url)

			wg.Add(1)
			go func(category models.Category, selling_type string, url string) {
				defer wg.Done()

				go_routine_cnt++
				m.ExtractDetailCategories(category, selling_type, models.GOOGLE_PLAY_URL_PREFIX + url)
			}(category, selling_type, url)
//			fmt.Println()
//			break
		}
//		break
	}

	wg.Wait()


//	m.ExplorerMainCategories()
//	m.FetchAppDetail(10)
	fmt.Println("--- END ---------------------------------")
}

//__ func old_main() {
//__ 	http := new(lib.HttpClient)
//__ 	http.SetDebugMode(false)
//__
//__ 	url := ""
//__ 	html := ""
//__
//__ 	if !true {
//__ 		url = "https://play.google.com/store/apps/details?id=com.cashslide"
//__ 		html = http.Get(url)
//__
//__ 		app := new(models.App)
//__ 		app.PackageId = "com.cashslide"
//__ 		app.Parsing(html)
//__ 		fmt.Println("-----------------------------------")
//__ 		fmt.Println(app.ToJson())
//__ 		fmt.Println("-----------------------------------")
//__
//__ 		// apps := make(models.Apps, 1)
//__ 		// apps = append(apps, *app)
//__ 		// apps = append(apps, *app)
//__ 		// fmt.Println(apps.ToJson())
//__ 		// fmt.Println("-----------------------------------")
//__ 		return
//__ 	}
//__
//__ 	// SYNC
//__ 	if !true {
//__ 		//인기 소셜 앱 리스트
//__ 		url = "https://play.google.com/store/apps/category/SOCIAL/collection/topselling_free?authuser=0"
//__ 		params := lib.HttpParams{
//__ 		// "message": "HELLO",
//__ 		// "key": "178",
//__ 		}
//__ 		html = http.Post(url, params)
//__ 		a := new(models.AppsCategory)
//__ 		packageIds := a.Parsing(html)
//__
//__ 		// AppsCategory에서 추출한 PackageID를 순서대로 순회해서 Google에서 Fetch 후 Apps에 추가
//__ 		apps := make(models.Apps, 0)
//__ 		for _, package_id := range packageIds {
//__ 			// fmt.Println(package_id)
//__ 			url = "https://play.google.com/store/apps/details?id=" + package_id
//__ 			html = http.Get(url)
//__
//__ 			app := new(models.App)
//__ 			app.PackageId = package_id
//__ 			app.Parsing(html)
//__ 			apps = append(apps, *app)
//__
//__ 			if len(apps) > 2 {
//__ 				break
//__ 			}
//__ 		}
//__
//__ 		fmt.Println(apps.ToJson())
//__ 	}
//__
//__ 	// ASYNC by Go Routine
//__ 	if !true {
//__ 		//인기 소셜 앱 리스트
//__ 		url = "https://play.google.com/store/apps/category/SOCIAL/collection/topselling_free?authuser=0"
//__ 		params := lib.HttpParams{
//__ 		// "message": "HELLO",
//__ 		// "key": "178",
//__ 		}
//__
//__ 		html = http.Post(url, params)
//__ 		a := new(models.AppsCategory)
//__ 		packageIds := a.Parsing(html)
//__ 		packageCnt := len(packageIds)
//__
//__ 		messages := make(chan models.App, packageCnt)
//__ 		apps := make(models.Apps, 0)
//__
//__ 		// go routine으로 PackageId별 App 정보 수집
//__ 		for category_no, package_id := range packageIds {
//__ 			go func(package_id string, category_no int) {
//__ 				url = "https://play.google.com/store/apps/details?id=" + package_id
//__ 				html = http.Get(url)
//__
//__ 				app := new(models.App)
//__ 				app.PackageId = package_id
//__ 				app.Parsing(html)
//__ 				app.CategoryRank = category_no
//__
//__ 				messages <- *app
//__ 			}(package_id, category_no+1)
//__ 		}
//__
//__ 		// go routine 대기
//__ 		for i := 0; i < packageCnt; i++ {
//__ 			select {
//__ 			case app := <-messages:
//__ 				apps = append(apps, app)
//__ 			}
//__ 		}
//__
//__ 		//카테고리 순서대로 정렬
//__ 		apps.SortByCategoryRank()
//__
//__ 		fmt.Println(apps.ToJson())
//__ 	}
//__
//__ 	if true {
//__ 		// Play! 앱 카테고리ID 추출
//__ 		url = "https://play.google.com/store/apps"
//__ 		html = http.Get(url)
//__
//__ 		categories := make(models.GooglePlayCategories, 0)
//__ 		categories.Parsing(html)
//__ 		// fmt.Println(categories.ToJson())
//__
//__ 		for _, category := range categories {
//__ 			fmt.Println(category.Id, category.Name, category.Url)
//__
//__ 			// 카테고리별 TopSelling LIST 추출
//__ 			// url_pattern := "https://play.google.com/store/apps/category/%s/collection/topselling_free"
//__ 			// url = fmt.Sprintf(url_pattern, category.Id)
//__
//__ 			// params := lib.HttpParams{}
//__ 			// html = http.Post(url, params)
//__ 			// fmt.Println(html)
//__ 			// break
//__ 		}
//__ 	}
//__ }
