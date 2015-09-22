package models

import (
	"encoding/json"
	"regexp"
)

type GooglePlayCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type GooglePlayCategories []GooglePlayCategory

func (self *GooglePlayCategories) Parsing(html string) {
	pattern := `<a class="child-submenu-link" href="((?:\S+)\/(\S+))" title="(.+?)" `

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(html, -1)
	if matches != nil && len(matches) > 0 {
		for _, row := range matches {
			category := GooglePlayCategory{Id: row[2], Name: row[3], Url: row[1]}
			*self = append(*self, category)
			// // fmt.Println(row[2])//, row[3], row[1])
			// fmt.Println(category)
			// break
		}
	}
}

func (self *GooglePlayCategories) ToJson() string {
	s, err := json.MarshalIndent(self, "\t", "")
	if err != nil {
		return ""
	}

	return string(s)
}
