package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

type Note struct {
	Code  string
	Title string
}

type Compose struct {
	Link  string
	Name  string
	VlaEd string
	VlaEc string
	Notes []Note
}

func main() {

	composes := make([]Compose, 0, 2000)

	c := colly.NewCollector(
		colly.AllowedDomains("bdlep.inssbt.es"),
	)

	c.OnHTML("a[href*=nombre]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		compose := Compose{
			Link: link,
			Name: e.Text,
		}

		composes = append(composes, compose)

		fmt.Printf("Compound : %s\n", compose)

		c.OnHTML("td", func(td *colly.HTMLElement) {
			fmt.Printf("Value found: %s\n", td.Text)
		})

		c.Visit(e.Request.AbsoluteURL(compose.Link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)

		fmt.Printf("Composes %s", composes)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=1&submit=Listado+completo+Agentes+Qu%C3%ADmicos")

}
