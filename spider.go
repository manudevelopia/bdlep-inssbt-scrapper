package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
	"time"
)

type Advise struct {
	code  string
	title string
}

type Compose struct {
	Link  string
	Name  string
	Vlas  [4]string
	Notes []Advise
	Warns []Advise
}

func main() {

	var composes []Compose
	compose := Compose{}
	advise := Advise{}
	pages := 1 //46

	// Create collector
	c := colly.NewCollector(
		colly.AllowedDomains("bdlep.inssbt.es"),
	)

	_ = c.Limit(&colly.LimitRule{
		Delay: 10 * time.Second,
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting: ", request.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Somethiwent wrong:", err)
	})

	// Get all compose links
	c.OnHTML("table[class='contents'] a[href*=nombre]", func(a *colly.HTMLElement) {
		fmt.Println(a.Text)

		compose.Link = a.Request.AbsoluteURL(a.Attr("href"))
		compose.Name = a.Text

		_ = c.Visit(saneUrl(compose.Link))
	})

	// get nCAS, nCE and Compose Name
	c.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
		if strings.Contains(td.Request.URL.String(), "vlaallpr.jsp?Bloque=") {
			fmt.Println(td.Text, td.Index)
		}
	})

	// get the 4 environmental values VLA-ED and VLA-EC
	c.OnHTML("table[class='valores'] tr:not([class='cabecera']) td", func(td *colly.HTMLElement) {
		fmt.Printf("VL :: %s - %d\n", td.Text, td.Index)

		compose.Vlas[td.Index] = td.Text
	})

	// get the Notes and Warnings
	c.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
		if strings.Contains(td.Request.URL.String(), "vlaallpr.jsp?Bloque=") {
			return
		}

		if strings.Contains(td.Request.URL.String(), "&FH=") {
			fmt.Printf("Warning :: %s - %d\n", td.Text, td.Index)
		} else if strings.Contains(td.Request.URL.String(), "&nombre=") {
			fmt.Printf("Note :: %s - %d\n", td.Text, td.Index)
		}

		parseAdvise(td, &advise)
	})

	// Link has been scrapped
	c.OnScraped(func(r *colly.Response) {
		url := r.Request.URL.String()
		fmt.Println(" - Finished", url)

		// add compose to collection only when finish compose information page
		if strings.Contains(url, "vlapr.jsp?") {
			composes = append(composes, compose)
		} else if strings.Contains(url, "&FH=") {
			compose.Warns = append(compose.Warns, advise)
		} else if strings.Contains(url, "&nombre=") {
			compose.Notes = append(compose.Notes, advise)
		}
	})

	// get the hazard advices links
	c.OnHTML("a[title='Indicaciones de peligro H']", func(a *colly.HTMLElement) {
		link := a.Request.AbsoluteURL(a.Attr("href"))
		_ = c.Visit(saneUrl(link))
	})

	// Iterate on each page ::
	for i := 1; i <= pages; i++ {
		_ = c.Visit(fmt.Sprintf("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=%d", i))
	}

}

func parseAdvise(tdAdvise *colly.HTMLElement, advice *Advise) {
	if tdAdvise.Index%2 == 0 {
		advice.code = tdAdvise.Text
	} else {
		advice.title = tdAdvise.Text
	}
}

func saneUrl(url string) string {
	return strings.ReplaceAll(url, " ", "%20")
}
