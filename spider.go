package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"strings"
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
	pages := 1 //46
	var composes []Compose
	advise := Advise{}

	collector := colly.NewCollector(
		colly.AllowedDomains("bdlep.inssbt.es"),
	)

	listRead := collector.Clone()

	// Get all compose links
	listRead.OnHTML("table[class='contents'] a[href*=nombre]", func(a *colly.HTMLElement) {
		fmt.Println(a.Text)

		compose := Compose{Name: a.Text, Url: a.Request.AbsoluteURL(a.Attr("href"))}

		composeRead := listRead.Clone()

		// get nCAS, nCE
		composeRead.OnHTML("span[class='destacado']", func(span *colly.HTMLElement) {
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
		composeRead.OnHTML("table[class='valores'] tr:not([class='cabecera']) td", func(td *colly.HTMLElement) {
			compose.Vlas[td.Index] = td.Text
		})

		// get the Notes and Warnings
		composeRead.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
			parseAdvise(td, &advise)

			if advise.Title != "" {
				fmt.Printf("Note :: %s - %d\n", td.Text, td.Index)
				compose.Notes = append(compose.Notes, advise)
				advise = Advise{}
			}

		})

		// get the hazard advices links
		composeRead.OnHTML("a[title='Indicaciones de peligro H']", func(a *colly.HTMLElement) {
			link := a.Request.AbsoluteURL(a.Attr("href"))

			hazardAdvicesRead := listRead.Clone()

			hazardAdvicesRead.OnHTML("table[class='contents'] td", func(td *colly.HTMLElement) {
				parseAdvise(td, &advise)

				if advise.Title != "" {
					fmt.Printf("Warning :: %s - %d\n", td.Text, td.Index)
					compose.Warns = append(compose.Warns, advise)
					advise = Advise{}
				}
			})

			_ = hazardAdvicesRead.Visit(saneUrl(link))
		})

		// get info for linked component
		composeRead.OnHTML("a[title='Agente Quimico']", func(a *colly.HTMLElement) {
			compose.Parent = a.Text
		})

		composeRead.OnScraped(func(response *colly.Response) {
			fmt.Println("Finished compose: " + compose.Name)
			composes = append(composes, compose)
			sendPost(compose)
		})

		_ = composeRead.Visit(saneUrl(compose.Url))
	})

	listRead.OnScraped(func(response *colly.Response) {
		fmt.Println("Finish Reading> " + response.Request.URL.String())
	})

	listRead.OnRequest(func(request *colly.Request) {
		fmt.Println("Start Reading> " + request.URL.String())
	})

	// Iterate on each page ::
	for i := 1; i <= pages; i++ {
		_ = listRead.Visit(fmt.Sprintf("http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=%d", i))
	}
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
