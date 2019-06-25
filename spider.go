package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
	"time"
)

type Note struct {
	code  string
	title string
}

type Compose struct {
	link  string
	name  string
	vlas  [4]string
	notes []Note
}

func main() {

	var composes []Compose
	compose := Compose{}
	pages := 1 //46

	// Create collector
	c := colly.NewCollector(
		colly.AllowedDomains("bdlep.inssbt.es"),
	)

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting: ", request.URL)
	})

	// Get all compose links
	c.OnHTML("table[class='contents'] a[href*=nombre]", func(a *colly.HTMLElement) {
		fmt.Println(a.Text)

		compose.link = a.Request.AbsoluteURL(a.Attr("href"))
		compose.name = a.Text

		_ = c.Visit(compose.link)
	})

	// get nCAS, nCE and Compose name
	c.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
		if strings.Contains(td.Request.URL.String(), "vlaallpr.jsp?Bloque=") {
			fmt.Println(td.Text, td.Index)
		}
	})

	// get the 4 environmental values VLA-ED and VLA-EC
	c.OnHTML("table[class='valores'] tr:not([class='cabecera']) td", func(td *colly.HTMLElement) {
		fmt.Printf("VL :: %s - %d\n", td.Text, td.Index)

		compose.vlas[td.Index] = td.Text
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

	})

	// link has been scrapped
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(" - Finished", r.Request.URL)

		// add compose to collection only when finish compose information page
		if strings.Contains(r.Request.URL.String(), "vlapr.jsp?") {
			composes = append(composes, compose)
		}
	})

	// get the hazard advices links
	c.OnHTML("a[title='Indicaciones de peligro H']", func(a *colly.HTMLElement) {
		_ = c.Visit(a.Request.AbsoluteURL(a.Attr("href")))
	})

	// Iterate on each page ::
	for i := 1; i <= pages; i++ {
		_ = c.Visit(fmt.Sprintf("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=%d", i))
	}

}
