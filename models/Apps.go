package models

import (
	"go-market-crawler/lib"
	"encoding/json"
	"fmt"
	"regexp"
//	"sort"
	"strconv"

//	"io/ioutil"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type App struct {
	CoverImage       string  `json:"cover_image"`
	PackageId        string  `json:"package_id"`
	SoftwareTitle    string  `json:"software_title"`
	SoftwareVersion  string  `json:"software_version"`
	DatePublished    string  `json:"date_published"`
	CurrentRating    float64 `json:"current_rating"`
	Reviewers        int     `json:"reviewers"`
	CategoryName     string  `json:"category_name"`
	Genre            string  `json:"genre"`
	OperatingSystems string  `json:"operating_systems"`
	ApkSize          float64 `json:"apk_size"`
	Description      string  `json:"description"`
//	CategoryRank     int     `json:"category_rank"`
//	WholeRank        int     `json:"whole_rank"`
}

type Apps []App

type AppsCategory struct {
	App
}

func (a *AppsCategory) Parsing(html string) []string {
	var arr []string
	pattern := `<a class="title" href="\/store\/apps\/details\?id=(.*?)" title="(?:.*?)" `
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(html, -1)
	if matches != nil && len(matches) > 0 {
		for _, row := range matches {
			arr = append(arr, row[1])
		}
	}

	if false {
		fmt.Println(arr)
	}

	return arr
}

func (app *App) ToJson() string {
	s, err := json.MarshalIndent(app, "\t", "")
	if err != nil {
		return ""
	}

	return string(s)
}

func (appx *App) Parsing(html string) App {
	patterns := map[string]string{
		"cover_image":       `<div class="cover-container">\s*<img class="cover-image" src="(.*?)" alt="Cover art" aria-hidden="true" itemprop="image">\s*</div>`,
		"software_title":    `<div class="document-title" itemprop="name">\s*<div>\s*(.*?)<\/div>\s*`,
		"software_version":  `<div class="content" itemprop="softwareVersion">\s*(\d+\.\d+(?:\.\d+)?)\s*<\/div>`,
		"date_published":    `<div class="content" itemprop="datePublished">\s*(.*?)<\/div>`,
		"current_rating":    `<div class="current-rating" style="width:\s*([0-9.]+)%"><\/div>`,
		"reviewers":         `<span class="rating-count" (?:.*)>\s*([0-9,]+)\s*<\/span>`,
		"category_name":     `<a class="document-subtitle category" href="(?:\/store\/apps\/category\/(.*?))">`,
		"genre":             `<span itemprop="genre">\s*(.*?)\s*<\/span>`,
		"operating_systems": `<div class="content" itemprop="operatingSystems">\s*(.*?)\s*<\/div>`,
		"apk_size":          `<div class="content" itemprop="fileSize">\s*([0-9\.]+)[M|G]\s*<\/div>`,
		"package_id":		 `<link href="https://(?:.*?)?id=(.*?)" rel="canonical">`,
	}

	value := ""
	app := App{}

	for key, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(html)
		if match != nil {
			value = match[1]

			switch key {
			case "cover_image":
				app.CoverImage = value
			case "software_title":
				app.SoftwareTitle = value
			case "software_version":
				app.SoftwareVersion = value
			case "date_published":
				app.DatePublished = value
			case "current_rating":
				app.CurrentRating, _ = strconv.ParseFloat(value, 64)
			case "reviewers":
				app.Reviewers, _ = strconv.Atoi(value)
			case "category_name":
				app.CategoryName = value
			case "genre":
				app.Genre = value
			case "operating_systems":
				app.OperatingSystems = value
			case "apk_size":
				app.ApkSize, _ = strconv.ParseFloat(value, 64)

			case "package_id":
				app.PackageId = value
			}
		}

		if match != nil {
			// fmt.Printf("\n### %s: %q\n\n\n", key, match)
			value = match[1]
			// fmt.Printf("### %s: %s\n", key, match[1])
		} else {
			// fmt.Printf("@@@ %s: %s\n", key, pattern)
		}

		if key == "cover_image" {
			//image_download(value)
		}
	}

	return app
}

func (apps *Apps) ToJson() string {
	s, err := json.MarshalIndent(apps, "\t", "")
	if err != nil {
		return ""
	}

	return string(s)
}

//func (apps *Apps) SortByCategoryRank() {
//	sort.Sort(AppsByCategoryRank(*apps))
//}

// https://golang.org/pkg/sort/
//type AppsByCategoryRank Apps

//func (a AppsByCategoryRank) Len() int           { return len(a) }
//func (a AppsByCategoryRank) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a AppsByCategoryRank) Less(i, j int) bool { return a[i].CategoryRank < a[j].CategoryRank }




func (self *App) Fetch(package_id string) App {
	// httpClient *lib.HttpClient
	httpClient := new(lib.HttpClient)
	httpClient.SetDebugMode(false)

	url := "https://play.google.com/store/apps/details?id=" + package_id
	html := httpClient.Get(url)
//	fmt.Println(html)
	app := self.Parsing(html)
//	fmt.Println(app.ToJson())
	return app
}

func (self *App) FetchAppList() []string {
	db, err := sql.Open("mysql", DSN)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxIdleConns(100)


	querystring := `
		select 	package_id
		from 	rankings
		where 	created_date = DATE_FORMAT(NOW(), '%Y-%m-%d')
		group 	by package_id
	`

	// query
	rows, err := db.Query(querystring)
	var package_ids []string

	for rows.Next() {
		var package_id string
		//err = db.QueryRow(querystring).Scan(&package_id)// .Scan(&new_category.Id, &new_category.Name, &new_category.Text, &new_category.Url)
		rows.Scan(&package_id)
//		fmt.Println(package_id)
		package_ids = append(package_ids, package_id)
	}
	fmt.Println(package_ids)

	return package_ids
}


func (self *App) Save() {
//	fmt.Println(self.PackageId)
//	fmt.Println(self.CategoryName)

// FILE WRITE
//	filename := "/tmp/google_play/" + self.PackageId
//	ioutil.WriteFile(filename, []byte(self.ToJson()), 0x777)
//	return



	db, err := sql.Open("mysql", DSN)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxIdleConns(150)

	// INSERT
	querystring := `
		INSERT INTO apps
		(
			  package_id
			, cover_image
			, software_title
			, software_version
			, date_published
			, current_rating
			, reviewers
			, category_name
			, genre
			, operating_systems
			, apk_size
			, description
			, created_date
			, created_time
		)
		VALUES
		(
			  ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, DATE_FORMAT(NOW(), '%Y-%m-%d')
			, DATE_FORMAT(NOW(), '%H:%i:%s')
		)
	`
	_, err = db.Exec(querystring, self.PackageId, self.CoverImage, self.SoftwareTitle, self.SoftwareVersion, self.DatePublished, self.CurrentRating, self.Reviewers, self.CategoryName, self.Genre, self.OperatingSystems, self.ApkSize, self.Description)
	if err != nil {
		fmt.Println(self.ToJson())
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

