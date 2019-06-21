package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
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

	// get the 4 environmental values VLA-ED and VLA-EC
	c.OnHTML("table[class='valores'] tr:not([class='cabecera']) td", func(td *colly.HTMLElement) {
		fmt.Printf("VL :: %s - %d\n", td.Text, td.Index)

		compose.vlas[td.Index] = td.Text
	})

	// get the Notes
	c.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
		if !strings.Contains(td.Request.URL.String(), "vlaallpr.jsp?Bloque=") {
			fmt.Printf("Note :: %s - %d\n", td.Text, td.Index)
		}

	})

	// link has been scrapped
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(" - Finished", r.Request.URL)

		composes = append(composes, compose)
	})

	// Iterate on each page ::
	for i := 1; i <= pages; i++ {
		_ = c.Visit(fmt.Sprintf("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=%d", i))
	}

}

/*
	//composes := make([]Compose, 0, 2000)
	//var a [4]string

	// Create collector
	c := colly.NewCollector(
		colly.AllowedDomains("bdlep.inssbt.es"),
	)

	// Obtain number of pages
	// TODO: implement. Now there are 46, replace after implement 'Obtain number of pages'
	pages := 46

	// Get all compose links
	c.OnHTML("table[class='contents'] a[href*=nombre]", func(a *colly.HTMLElement) {
		fmt.Println(a.Text)

		link := a.Request.AbsoluteURL(a.Attr("href"))

		_ = c.Visit(link)
	})

	// get the 4 environmental values VLA-ED and VLA-EC
	c.OnHTML("table[class='valores'] tr:not([class='cabecera']) td", func(td *colly.HTMLElement) {
		fmt.Printf("VL :: %s - %d\n", td.Text, td.Index)
	})

	// get the hazard advices
	c.OnHTML("a[title='Indicaciones de peligro H']", func(a *colly.HTMLElement) {
		_ = c.Visit(a.Request.AbsoluteURL(a.Attr("href")))
	})

	// get the Notes and
	c.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
		fmt.Printf("Note :: %s - %d\n", td.Text, td.Index)
	})

	// Iterate on each page ::
	for i:= 1; i<= pages; i++ {
		_ = c.Visit(fmt.Sprintf("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=%d", i))
	}
}
*/

//c.OnScraped(func(r *colly.Response) {
//	fmt.Println("Finished", r.Request.URL)
//
//	fmt.Printf("Composes %s", composes)
//})


//c.OnHTML("a[href*=nombre]", func(e *colly.HTMLElement) {
//	link := e.Attr("href")
//
//	c.OnHTML("table[class='valores'] tr:not([class='cabecera']) td", func(td *colly.HTMLElement) {
//		fmt.Printf("VL :: %s - %d\n", td.Text, td.Index)
//		a[td.Index] = td.Text
//	})
//
//	compose := Compose{
//		link: link,
//		name: e.Text,
//		vlas: a,
//	}
//
//	composes = append(composes, compose)
//
//	c.Visit(e.Request.AbsoluteURL(link))
//})
//
///*	c.OnHTML("a[href*=nombre]", func(e *colly.HTMLElement) {
//		link := e.Attr("href")
//
//		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
//
//		compose := Compose{
//			Link: link,
//			Name: e.Text,
//		}
//
//		composes = append(composes, compose)
//
//		fmt.Printf("Compound : %s\n", compose)
//
//		c.OnHTML(".valores td", func(td *colly.HTMLElement) {
//			fmt.Printf("Value found: %s\n", td.Text)
//		})
//
//		c.OnHTML(".contents tr", func(tr *colly.HTMLElement) {
//			fmt.Printf("Notes found: %s\n", tr.Text)
//		})
//
//		c.Visit(e.Request.AbsoluteURL(compose.Link))
//	})
//
//	c.OnRequest(func(r *colly.Request) {
//		fmt.Println("Visiting", r.URL.String())
//	})*/
//
//c.OnScraped(func(r *colly.Response) {
//	fmt.Println("Finished", r.Request.URL)
//
//	fmt.Printf("Composes %s", composes)
//})