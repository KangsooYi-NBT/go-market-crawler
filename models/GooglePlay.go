package models

import (
	"fmt"
	"go-market-crawler/lib"
	"io/ioutil"
	//	"encoding/json"
	"regexp"
	"strconv"
	"strings"
//	"time"


	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const GOOGLE_PLAY_URL_PREFIX string = "https://play.google.com"
const GOOGLE_CATEGORY_URL string = "https://play.google.com/store/apps"

const CATEGORY_RETRIEVE_OFFSET = 0   // 카테고리에서 몇번째 앱 부터 가져올지 지정
const CATEGORY_RETRIEVE_LIMIT = 120  // 카테고리 한 페이지에서 최대로 읽어올 앱 개수 (MAX. 120)
const CATEGORY_RETRIEVE_MAX_PAGE = 20 // 카테고리 페이지로 부터 앱 리스트 추출을 최대 몇 페이지까지 할지 지정

//const CATEGORY_RETRIEVE_LIMIT = 10  // 카테고리 한 페이지에서 최대로 읽어올 앱 개수 (MAX. 120)
//const CATEGORY_RETRIEVE_MAX_PAGE = 1 // 카테고리 페이지로 부터 앱 리스트 추출을 최대 몇 페이지까지 할지 지정



var db *sql.DB
const DSN = "{USER_ID}:{PASSWORD}@tcp({SERVER_IP}:{SERVER_PORT})/{DATABASE_NAME}"






type Category struct {
	Id   int    `json:"id"`   //ID
	Name string `json:"name"` //영문 이름
	Text string `json:"text"` //한글 이름/설명
	Url  string `json:"url"`  //URL
}

type Rank struct {
	Id         int64  `json:"id"`          //ID
	CategoryId int    `json:"category_id"` //카테고리ID
	Type       string `json:"type"`        //[ topselling_free | topselling_paid | topgrossing ]
	Ranking    int    `json:"ranking"`     //순위
	PackageId  string `json:"package_id"`  //안드로이드 PackageID
	CreatedDate  string `json:"created_date"`  //생성일자
	CreatedTime  string    `json:"created_time"`  //생성일시
}

type GooglePlay struct {
	httpClient *lib.HttpClient
	categories []Category
}

func (self *GooglePlay) Init() {
	self.httpClient = new(lib.HttpClient)
	self.httpClient.SetDebugMode(false)
	//	fmt.Println("Init()")
}

// 앱 카테고리 정보 조회
func (self *GooglePlay) ExtractCategories() []Category {
	html := ""
	filename := "/tmp/google_play_category.html"
	content, err := ioutil.ReadFile(filename)
	if true || err != nil {
		//fmt.Println("---------------------------")
		html = self.httpClient.Get(GOOGLE_CATEGORY_URL)
		ioutil.WriteFile(filename, []byte(html), 0x777)
	} else {
		html = string(content)
	}

	// (전체) 인기차트
	if true {
		// 인기 앱
		category := Category{Name: "APPS", Text: "전체 인기 앱", Url: "/store/apps"}
		category.save()
		self.categories = append(self.categories, category)

		// 게임
		category = Category{Name: "GAMES", Text: "전체 게임", Url: "/store/apps/category/GAME"}
		category.save()
		self.categories = append(self.categories, category)
		//result := []string{category.Id, category.Name, category.Url}; fmt.Println("[" + strings.Join(result, "]\t[") + "]")
	}

	self.ParseCategories(html)
	return self.categories
}
func (self *GooglePlay) ParseCategories(html string) {
	pattern := `<a class="child-submenu-link" href="((?:\S+)\/(\S+))" title="(.+?)" `

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(html, -1)
	if matches != nil && len(matches) > 0 {
		for _, row := range matches {
			// Url이 /store/apps/로 시작하는 것만 추출
			if strings.HasPrefix(row[1], "/store/apps/") {
				category := Category{Name: row[2], Text: row[3], Url: row[1]}
				category.save()
				self.categories = append(self.categories, category)
				//result := []string{category.Id, category.Name, category.Url}; fmt.Println("[" + strings.Join(result, "]\t[") + "]")
			}
		}
	}
	//	fmt.Println(self.categories)
}

func (self *GooglePlay) ExplorerMainCategories() {
	for _, category := range self.categories {
		self.ExtractSubCategories(category)
	}
}

func (self *GooglePlay) ExtractSubCategories(category Category) {
	//	https://play.google.com/store/apps/category/HEALTH_AND_FITNESS
	url := GOOGLE_PLAY_URL_PREFIX + category.Url
	html := ""
	filename := "/tmp/google_play_category" + category.Name + ".html"
	content, err := ioutil.ReadFile(filename)
	if true || err != nil {
		fmt.Println("---------------------------")
		html = self.httpClient.Get(url)

		pattern := `<div class="cluster-heading">`
		re := regexp.MustCompile(pattern)
		html = re.ReplaceAllString(html, "\n\n\n"+pattern)
		ioutil.WriteFile(filename, []byte(html), 0x777)
	} else {
		html = string(content)
	}

	//	fmt.Println(html)
	self.ParseSubCategories(category, html)

}

func (self *GooglePlay) ParseSubCategories(_category Category, html string) {
	pattern := `<h2>\s*<a class="title-link id-track-click" (?:.+?) href="(.+?)">\s*(.+?)\s*</a>\s*<\/h2>`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(html, -1)
	if matches != nil && len(matches) > 0 {
		fmt.Println("Matches: ", len(matches))
		for _, row := range matches {
			fmt.Println("- Cols: ", len(row))

			result := []string{row[1], row[2]}
			fmt.Println("[" + strings.Join(result, "]\t[") + "]")
			//			fmt.Println(Category{Id: row[1], Name: row[2], Url: row[1]})
		}
	}

	fmt.Println("END")
}

func (self *GooglePlay) ExtractDetailCategories(category Category, selling_type string, url string) {
	html := ""

	num := CATEGORY_RETRIEVE_LIMIT
	start := CATEGORY_RETRIEVE_OFFSET
	last_no := 0
	loop_cnt := 0

	for {
		start = num * loop_cnt

		filename := "/tmp/google_play_category_tmp.html"
		content, err := ioutil.ReadFile(filename)
		if true || err != nil {
			params := lib.HttpParams{
				"start":       strconv.Itoa(start), //카테고리 상위 몇번째 부터 가져올지
				"num":         strconv.Itoa(num),   //한번에 가져올 App 갯수 최대 120으로 확인
				"numChildren": "0",
				"ipf":         "1",
				"xhr":         "1",
			}
			html = self.httpClient.Post(url, params)
			ioutil.WriteFile(filename, []byte(html), 0x777)
		} else {
			html = string(content)
		}

		last_no = self.ParseDetailCategories(category, selling_type, html, last_no)

		loop_cnt++
		if loop_cnt >= CATEGORY_RETRIEVE_MAX_PAGE || last_no == -1 {
			break
		}
	}
}

func (self *GooglePlay) ParseDetailCategories(category Category, selling_type string, html string, prev_last_no int) int {
	last_no := 0

	var arr []string
	pattern := `<a class="title" href="\/store\/apps\/details\?id=(.*?)" title=(?:"|')(?:.*?)(?:"|') (?:.*?)>\s*(\d+)\.\s*(.*?)\s+<`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(html, -1)
	if matches != nil && len(matches) > 0 {
		for _, row := range matches {
			arr = append(arr, row[1])
			curr_no, _ := strconv.Atoi(row[2])

			if prev_last_no > 0 && prev_last_no > curr_no {
				fmt.Println("### Category에 더 이상 App이 없음")
				return -1
			}

			ranking, _ := strconv.Atoi(row[2])
			rank := Rank{
				CategoryId: category.Id,
				Type: selling_type,
				Ranking: ranking,
				PackageId: row[1],
			}
			rank.save()

			result := []string{category.Name, selling_type, row[2], row[3], row[1]}; fmt.Println("[" + strings.Join(result, "]\t[") + "]")
			last_no = curr_no
		}
	}

	//	fmt.Println(arr)
	return last_no
}

// 앱 상세 정보 조회
func (self *GooglePlay) FetchAppDetail(app_id int) {
	fmt.Println(app_id)
}


func (self *Category) save() {
	var err error
	if db == nil {
		db, err = sql.Open("mysql", DSN)
//		defer db.Close()
		if err != nil {
			panic(err.Error())
		}

		db.SetMaxIdleConns(100)
	}

	new_category := Category{}
	err = db.QueryRow("SELECT `id`, `name`, `text`, `url` FROM `categories` WHERE url = ?", self.Url).Scan(&new_category.Id, &new_category.Name, &new_category.Text, &new_category.Url)
	if err != nil {
		// INSERT
		_, err := db.Exec("INSERT INTO `categories` (`name`, `text`, `url`) VALUE (?, ?, ?)", self.Name, self.Text, self.Url)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// SELECT
		err = db.QueryRow("SELECT `id`, `name`, `text`, `url` FROM `categories` WHERE url = ?", self.Url).Scan(&new_category.Id, &new_category.Name, &new_category.Text, &new_category.Url)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	//	rows, err := db.Query("SELECT `id`, `name`, `text`, `url` FROM `categories` WHERE url = ?", self.Url)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer rows.Close()
//
//	var old_category *Category
//	rows.Next()
////	for rows.Next() {
//		err = rows.Scan(&old_category)
////		if err != nil {
////			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
////		}
//		fmt.Println(new_category)
////	}

	*self = new_category
}

func (self *Rank) save() {
	var err error
	if db == nil {
		db, err = sql.Open("mysql", DSN)
//		defer db.Close()
		if err != nil {
			panic(err.Error())
		}

		db.SetMaxIdleConns(100)
	}

	new_rank := Rank{}
	err = db.QueryRow("SELECT `id`, `category_id`, `type`, `ranking`, `package_id`, `created_date`, `created_time` FROM `rankings` WHERE category_id = ? AND type = ? AND ranking = ? AND created_date = DATE_FORMAT(NOW(), '%Y-%m-%d')", self.CategoryId, self.Type, self.Ranking).Scan(&new_rank.Id, &new_rank.CategoryId, &new_rank.Type, &new_rank.Ranking, &new_rank.PackageId, &new_rank.CreatedDate, &new_rank.CreatedTime)
	if err != nil {
		// INSERT
		_, err := db.Exec("INSERT INTO `rankings` (`category_id`, `type`, `ranking`, `package_id`, `created_date`, `created_time`) VALUE (?, ?, ?, ?, DATE_FORMAT(NOW(), '%Y-%m-%d'), DATE_FORMAT(NOW(), '%H:%m:%d') )", self.CategoryId, self.Type, self.Ranking, self.PackageId)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// SELECT
		err = db.QueryRow("SELECT `id`, `category_id`, `type`, `ranking`, `package_id`, `created_date`, `created_time` FROM `rankings` WHERE category_id = ? AND type = ? AND ranking = ? AND created_date = DATE_FORMAT(NOW(), '%Y-%m-%d')", self.CategoryId, self.Type, self.Ranking).Scan(&new_rank.Id, &new_rank.CategoryId, &new_rank.Type, &new_rank.Ranking, &new_rank.PackageId, &new_rank.CreatedDate, &new_rank.CreatedTime)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	*self = new_rank
}
