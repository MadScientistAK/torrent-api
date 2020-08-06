package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var Link1337 string = "https://1337x.to"
var wg sync.WaitGroup

func GetTorrents(query string) {
	query = Link1337 + "/search/" + strings.ReplaceAll(query, " ", "+") + "/1/"
	res, err := http.Get(query)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
		// return make([]Torrent, 0), "No Results Returned"
	}

	// torrentList := make([]Torrent, 0)

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	a := make([]string, 0)
	b := make([]string, 0)

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			tor := s.Find(".name").Find("a").Eq(1)
			torLink, _ := tor.Attr("href")
			// torSeeders := s.Find(".seeds").Eq(0).Text()
			// torLeechers := s.Find(".leeches").Text()
			// torSize := strings.Replace(s.Find(".size").Text(), torSeeders, "", -1)
			// torrentList = append(torrentList, Torrent{Tname: tor.Text(), Tlink: torLink, TSeeders: torSeeders, TLeechers: torLeechers, TSize: torSize})
			a = append(a, tor.Text())
			wg.Add(1)
			go GetMagnet(&torLink, b, &wg)
		}
	})
	wg.Wait()
	fmt.Println(a)
	fmt.Println(b)
}

func GetMagnet(Tlink *string, sli []string, wg *sync.WaitGroup) {
	Tlink1 := "https://1337x.to" + *Tlink
	res, _ := http.Get(Tlink1)
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	t, _ := doc.Find(".clearfix").Eq(2).Find("a").Attr("href")
	sli = append(sli, t)
	*Tlink = t
	wg.Done()
}

func main() {
	GetTorrents("Far Cry")
	wg.Wait()
}
