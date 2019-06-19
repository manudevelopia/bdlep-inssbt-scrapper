package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main(){

	c := colly.NewCollector(
		colly.AllowedDomains("bdlep.inssbt.es"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		c.OnHTML("td", func(td *colly.HTMLElement) {
			fmt.Printf("Value found: %s\n", td.Text)
		})

		c.Visit(e.Request.AbsoluteURL(link))
	})


	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=1&submit=Listado+completo+Agentes+Qu%C3%ADmicos")
}
