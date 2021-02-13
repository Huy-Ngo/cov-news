package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
)

func getTitle(html *string, element *colly.HTMLElement) {
	*html += "\n<h2>\n"
	*html += element.ChildText("div.timeline-head")
	*html += "\n</h2>\n"
}

func getContent(html *string, element *colly.HTMLElement) {
	element.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
		*html += "\n<p>\n"
		*html += paragraph.Text
		*html += "\n</p>\n"
	})
}

func main() {
	c := colly.NewCollector()
	html := ""

	c.OnHTML("div.timeline-detail", func(e *colly.HTMLElement) {
		getTitle(&html, e)
		getContent(&html, e)
	})

	// Ignore unknown certificate
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.Visit("https://ncov.moh.gov.vn/dong-thoi-gian")

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	err := ioutil.WriteFile("build/index.html", []byte(html), 0644)
	if err != nil {
		panic(err)
	}
}
