package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var LastUpdated = time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC) // Arbitrary time in the past

// Parse scraped HTML element and return HTML text for the title
func getTitleHTML(element *colly.HTMLElement) string {
	title := element.ChildText("div.timeline-head")
	re := regexp.MustCompile("[^A-Za-z0-9]")
	id := re.ReplaceAll([]byte(title), []byte(""))
	return fmt.Sprintf("\n<h2 id=\"%s\">\n%s\n</h2>\n", id, title)
}

// Parse scraped HTML element and return HTML text for the content
func getContentHTML(element *colly.HTMLElement) string {
	html := ""
	element.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
		html += fmt.Sprintf("\n<p>\n%s\n</p>\n", paragraph.Text)
	})
	return html
}

// Parse scraped HTML element and return Atom XML for the metadata
func getMetaAtom(element *colly.HTMLElement) string {
	title := element.ChildText("div.timeline-head")
	re := regexp.MustCompile("[^A-Za-z0-9]")
	id := re.ReplaceAll([]byte(title), []byte(""))
	timeFormat := "15:04 02/01/2006 -0700"
	t, err := time.Parse(timeFormat, title+" +0700")
	if err != nil {
		fmt.Println("Error happened in parsing string:", err)
	}
	updated := t.Format(time.RFC3339)
	if t.After(LastUpdated) {
		LastUpdated = t
	}
	atom := fmt.Sprintf("<title>\n%s\n</title>\n", title)
	atom += fmt.Sprintf("<id>https://huy-ngo.github.io/cov-news/index.html#%s</id>\n", id)
	atom += fmt.Sprintf("<link>https://huy-ngo.github.io/cov-news/index.html#%s</link>\n", id)
	atom += fmt.Sprintf("<updated>\n%s\n</updated>\n", updated)
	return atom
}

// Parse scraped HTML element and return Atom XML for the content
func getContentAtom(element *colly.HTMLElement) string {
	content := `<content type="xhtml" xml:lang="en" xml:base="http://diveintomark.org/">`
	element.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
		content += fmt.Sprintf("\n<p>\n%s\n</p>\n", paragraph.Text)
	})
	content += "</content>"
	return content
}

// Parse scraped HTML element and return HTML for an entry
func GetEntryHTML(element *colly.HTMLElement) string {
	return getTitleHTML(element) + getContentHTML(element)
}

// Parse scraped HTML element and return Atom XML for an entry
func GetEntryAtom(element *colly.HTMLElement) string {
	return fmt.Sprintf("<entry>\n%s%s</entry>\n", getMetaAtom(element), getContentAtom(element))
}

func main() {
	c := colly.NewCollector()
	html := `<a href="atom.xml">Atom Feed</a>`
	atom_head := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
<title>Diễn biến dịch COVID-19</title>
<link href="https://huy-ngo.github.io/cov-news/index.html"/>
<id>https://huy-ngo.github.io/cov-news</id>`
	atom := ""

	c.OnHTML("div.timeline-detail", func(e *colly.HTMLElement) {
		html += GetEntryHTML(e)
		atom += GetEntryAtom(e)
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

	updated := LastUpdated.Format(time.RFC3339)
	atom_head += fmt.Sprintf("<updated>%s</updated>\n", updated)
	atom = atom_head + atom
	atom += "</feed>"
	err = ioutil.WriteFile("build/atom.xml", []byte(atom), 0644)
	if err != nil {
		panic(err)
	}
}
