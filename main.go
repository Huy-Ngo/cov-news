package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
)

func getTitleHTML(element *colly.HTMLElement) string {
	html := "\n<h2>\n"
	html += element.ChildText("div.timeline-head")
	html += "\n</h2>\n"
	return html
}

func getContentHTML(element *colly.HTMLElement) string {
	html := ""
	element.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
		html += "\n<p>\n"
		html += paragraph.Text
		html += "\n</p>\n"
	})
	return html
}

func main() {
	c := colly.NewCollector()
	html := "<a href=\"#\">Atom Feed</a>"
	// atom := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<feed xmlns=\"http://www.w3.org/2005/Atom\">"

	c.OnHTML("div.timeline-detail", func(e *colly.HTMLElement) {
		html += getTitleHTML(e)
		html += getContentHTML(e)
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
