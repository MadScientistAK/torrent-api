package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type torrent struct {
	tname     string `json:"TorrentName"`
	tlink     string `json:"MagnetLink"`
	tseeders  string `json:"Seeders"`
	tleechers string `json:"Leechers"`
	tsize     string `json:"Size"`
}

var base_link string = "https://1337x.to"

var wg sync.WaitGroup

func getMagnet(tlink *string, wg *sync.WaitGroup) {
	tlink1 := "https://1337x.to" + *tlink
	res, _ := http.Get(tlink1)
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	t, _ := doc.Find(".clearfix").Eq(2).Find("a").Attr("href")
	*tlink = t
	defer wg.Done()
}

func getTorrents(query string, wg *sync.WaitGroup) ([]torrent, string) {
	query = base_link + "/search/" + strings.Replace(query, " ", "+", -1) + "/1/"
	res, err := http.Get(query)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return make([]torrent, 0), "No Results Returned"
	}

	torrentList := make([]torrent, 0)

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			tor := s.Find(".name").Find("a").Eq(1)
			torurl, _ := tor.Attr("href")
			torseeders := s.Find(".seeds").Eq(0).Text()
			torleechers := s.Find(".leeches").Text()
			torsize := strings.Replace(s.Find(".size").Text(), torseeders, "", -1)
			torrentList = append(torrentList, torrent{tname: tor.Text(), tlink: torurl, tseeders: torseeders, tleechers: torleechers, tsize: torsize})
			wg.Add(1)
			go getMagnet(&torrentList[i-1].tlink, wg)
		}
	})
	wg.Wait()
	return torrentList, ""
}

func main() {
	fmt.Println(getTorrents("Jump Force", &wg))
}
