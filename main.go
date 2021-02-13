package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
)

func main() {
	c := colly.NewCollector()
	html := ""

	c.OnHTML("div.timeline-detail", func(e *colly.HTMLElement) {
		html += "\n<h2>\n"
		html += e.ChildText("div.timeline-head")
		html += "\n</h2>\n"
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			html += "\n<p>\n"
			html += paragraph.Text
			html += "\n</p>\n"
		})
	})

	// Ignore unknown certificate
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://ncov.moh.gov.vn/dong-thoi-gian")

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	err := ioutil.WriteFile("output.html", []byte(html), 0644)
	if err != nil {
		panic(err)
	}
}
