package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Kitap struct {
	Title string
	Desc  string
	Isbn  string
	Cats  []string
}

//curl -XPOST  -H 'Content-type: application/json' -d '{"deneme" : "hello" }'
func main() {
	kitap := Kitap{}

	for i := 1; i < 1000000; i++ {
		sayfa := strconv.Itoa(i)
		getit(sayfa, &kitap)

	}
}

func getit(sayfa string, kitap *Kitap) {
	url := "http://localhost:9200/kitapyurdu/doc"

	fmt.Printf("requested %s", sayfa)
	c := colly.NewCollector()

	c.OnHTML("h1.product-heading", func(e *colly.HTMLElement) {
		if len(e.Text) < 3 {
			log.Fatalf("sayfa yok %s", sayfa)
		}
		kitap.Title = e.Text
	})
	c.OnHTML("div#description_text", func(e *colly.HTMLElement) {
		tempDesc := e.Text
		desc := strings.Trim(tempDesc, "\t \n \r")
		kitap.Desc = desc
	})
	c.OnHTML("span[itemprop=\"isbn\"]", func(e *colly.HTMLElement) {
		kitap.Isbn = e.Text
	})
	c.OnHTML("div.product-info.grid_9.alpha > div:nth-child(3) > div.book-cover.box-shadow.mg-b-20 > div.grid_6.omega.alpha.book-right > div > div:nth-child(7) > a:nth-child(2)", func(e *colly.HTMLElement) {
		tempCats := e.Text
		cats := strings.Split(tempCats, "»")
		kitap.Cats = cats
	})

	c.OnScraped(func(r *colly.Response) {
		jsonStr, _ := json.Marshal(kitap)

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

	})

	c.Visit("https://www.kitapyurdu.com/kitap/app/" + sayfa + ".html")
}
