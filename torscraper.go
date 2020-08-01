package torscraper

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Torrent is a struct to hold torrent data
type Torrent struct {
	tname     string `json:"TorrentName"`
	tlink     string `json:"MagnetLink"`
	tseeders  string `json:"Seeders"`
	tleechers string `json:"Leechers"`
	tsize     string `json:"Size"`
}

var Link1337 string = "https://1337x.to"
var LinkNyaa string = "https://nyaa.si"

// GetAnimeTorrents is a function to get English translated Anime from Nyaa.si
func GetAnimeTorrents(query string) ([]Torrent, string) {
	query = LinkNyaa + "/?f=0&c=1_2&q=" + strings.ReplaceAll(query, " ", "+")
	res, err := http.Get(query)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	torrentList := make([]Torrent, 0)

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	doc.Find(".danger, .success, .default").Each(func(i int, s *goquery.Selection) {
		tmp := s.Find("td")
		tor := tmp.Eq(1).Find("a").Eq(0)

		if tor.HasClass("comments") {
			tor = tor.NextFiltered("a")
		}
		torLink, _ := tmp.Eq(2).Find("a").Last().Attr("href")
		torSize := tmp.Eq(3).Text()
		torSeeders := tmp.Eq(5).Text()
		torLeechers := tmp.Eq(6).Text()
		torrentList = append(torrentList, torrent{tname: tor.Text(), tlink: torUrl, tseeders: torSeeders, tleechers: torLeechers, tsize: torSize})
	})
	return torrentList, ""
}

var wg sync.WaitGroup

// GetMagnet is a function to get magnet links from 1337x.to
func GetMagnet(tlink *string, wg *sync.WaitGroup) {
	tlink1 := "https://1337x.to" + *tlink
	res, _ := http.Get(tlink1)
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	t, _ := doc.Find(".clearfix").Eq(2).Find("a").Attr("href")
	*tlink = t
	defer wg.Done()
}

// GetTorrents is a function to get torrents from 1337x.to
func GetTorrents(query string, wg *sync.WaitGroup) ([]Torrent, string) {
	query = Link1337 + "/search/" + strings.ReplaceAll(query, " ", "+") + "/1/"
	res, err := http.Get(query)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return make([]Torrent, 0), "No Results Returned"
	}

	torrentList := make([]Torrent, 0)

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			tor := s.Find(".name").Find("a").Eq(1)
			torUrl, _ := tor.Attr("href")
			torSeeders := s.Find(".seeds").Eq(0).Text()
			torLeechers := s.Find(".leeches").Text()
			torSize := strings.Replace(s.Find(".size").Text(), torSeeders, "", -1)
			torrentList = append(torrentList, torrent{tname: tor.Text(), tlink: torUrl, tseeders: torSeeders, tleechers: torLeechers, tsize: torSize})
			wg.Add(1)
			go GetMagnet(&torrentList[i-1].tlink, wg)
		}
	})
	wg.Wait()
	return torrentList, ""
}
