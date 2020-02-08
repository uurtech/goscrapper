package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

func main() {
	var m map[string]interface{}
	data := []byte("")
	err := json.Unmarshal(data, &m)
	if err != nil {

	}
	for i := 1; i < 200000; i++ {
		sayfa := strconv.Itoa(i)
		geziyor.NewGeziyor(&geziyor.Options{
			StartURLs:                   []string{"https://www.kitapyurdu.com/kitap/app/" + sayfa + ".html"},
			ConcurrentRequests:          3,
			ConcurrentRequestsPerDomain: 3,
			ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
				r.HTMLDoc.Find("#productDetail-" + sayfa + " > div.product-info.grid_9.alpha > div:nth-child(3) > div.book-cover.box-shadow.mg-b-20").Each(func(_ int, s *goquery.Selection) {
					title := s.Find("h1.product-heading").Text()
					if len(title) < 3 {
						fmt.Printf("kalınan sayfa i : %d", i)
						os.Exit(1)
					}
					tempDesc := s.Find("div#description_text").Text()
					desc := strings.Trim(tempDesc, "\t \n \r")
					isbn := s.Find("span[itemprop=\"isbn\"]").Text()
					tempCats := s.Find("div.product-info.grid_9.alpha > div:nth-child(3) > div.book-cover.box-shadow.mg-b-20 > div.grid_6.omega.alpha.book-right > div > div:nth-child(7) > a:nth-child(2)").Text()
					cats := strings.Split(tempCats, "»")
					g.Exports <- map[string]interface{}{
						"title":       title,
						"description": desc,
						"isbn":        isbn,
						"cats":        cats,
					}
				})
			},
			Exporters: []export.Exporter{&export.JSON{}},
		}).Start()
	}

}
