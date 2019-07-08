package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Advise struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

type Compose struct {
	Url    string    `json:"url"`
	Name   string    `json:"name"`
	Parent string    `json:"parent"`
	Ncas   string    `json:"ncas"`
	Nce    string    `json:"nce"`
	Vlas   [4]string `json:"vlas"`
	Notes  []Advise  `json:"notes"`
	Warns  []Advise  `json:"warns"`
}

func main() {

	var composes []Compose
	compose := Compose{}
	advise := Advise{}
	pages := 46

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
		log.Println("Something went wrong:", err)
	})

	// Get all compose links
	c.OnHTML("table[class='contents'] a[href*=nombre]", func(a *colly.HTMLElement) {
		fmt.Println(a.Text)

		compose.Url = a.Request.AbsoluteURL(a.Attr("href"))
		compose.Name = a.Text

		_ = c.Visit(saneUrl(compose.Url))
	})

	// get nCAS, nCE and Compose Name
	c.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
		if strings.Contains(td.Request.URL.String(), "vlaallpr.jsp?Bloque=") {
			fmt.Println(td.Text, td.Index)
		}
	})

	// get nCAS, nCE
	c.OnHTML("span[class='destacado']", func(span *colly.HTMLElement) {
		if span.Text == "Indicaciones de peligro H" || len(span.Text) <= 4 {
			return
		}

		if span.Index == 0 {
			compose.Ncas = span.Text
		} else if span.Index == 1 {
			compose.Nce = span.Text
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

		parseAdvise(td, &advise)

		if advise.Title != "" {
			if strings.Contains(td.Request.URL.String(), "&FH=") {
				fmt.Printf("Warning :: %s - %d\n", td.Text, td.Index)
				compose.Warns = append(compose.Warns, advise)
				advise = Advise{}
			} else if strings.Contains(td.Request.URL.String(), "&nombre=") {
				fmt.Printf("Note :: %s - %d\n", td.Text, td.Index)
				compose.Notes = append(compose.Notes, advise)
				advise = Advise{}
			}
		}

	})

	// Url has been scrapped
	c.OnScraped(func(r *colly.Response) {
		url := r.Request.URL.String()
		fmt.Println(" - Finished", url)

		// add compose to collection only when finish compose information page
		if strings.Contains(url, "vlapr.jsp?") {
			composes = append(composes, compose)
			sendPost(compose)
			compose = Compose{}
			advise = Advise{}
		}
	})

	// get the hazard advices links
	c.OnHTML("a[title='Indicaciones de peligro H']", func(a *colly.HTMLElement) {
		link := a.Request.AbsoluteURL(a.Attr("href"))
		_ = c.Visit(saneUrl(link))
	})

	// Iterate on each page ::
	for i := 1; i <= pages; i++ {
		//		_ = c.Visit(fmt.Sprintf("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=%d", i))
	}

	_ = c.Visit("http://bdlep.inssbt.es/LEP/vlapr.jsp?ID=258&nombre=beta-Cloropreno")

}

func parseAdvise(tdAdvise *colly.HTMLElement, advice *Advise) {
	if tdAdvise.Index%2 == 0 {
		advice.Code = tdAdvise.Text
	} else {
		advice.Title = tdAdvise.Text
	}
}

func saneUrl(url string) string {
	return strings.ReplaceAll(url, " ", "%20")
}

func sendPost(compose Compose) {
	url := "http://localhost:8080/api/scrapper/compound"
	fmt.Println("URL:>", url)

	var jsonData []byte
	jsonData, _ = json.Marshal(compose)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
