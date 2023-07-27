package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type FacultyInfo struct {
	Name         string
	Title        string
	Department   string
	EmailAddress string
	ProfileURL   string
}

func main() {
	url := "https://www.stevens.edu/profile"

	facultyInfo := FacultyInfo{}
	facultyInfos := make([]FacultyInfo, 0, 1)

	c := colly.NewCollector(colly.AllowedDomains("www.stevens.edu"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	c.OnHTML(".name", func(h *colly.HTMLElement) {
		facultyInfo.Name = cleanText(h.Text)
	})

	c.OnHTML(".title", func(h *colly.HTMLElement) {
		facultyInfo.Title = cleanText(h.Text)
	})

	c.OnHTML(".department", func(h *colly.HTMLElement) {
		facultyInfo.Department = cleanText(h.Text)
	})

	c.OnHTML(".contact > .email", func(h *colly.HTMLElement) {
		facultyInfo.EmailAddress = cleanText(h.Text)
	})

	c.OnHTML(".name a", func(h *colly.HTMLElement) {
		facultyInfo.ProfileURL = h.Request.AbsoluteURL(h.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		facultyInfos = append(facultyInfos, facultyInfo)
		facultyInfo = FacultyInfo{}
	})

	c.Visit(url)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(facultyInfos)
}

func cleanText(s string) string {
	return strings.TrimSpace(s)
}
