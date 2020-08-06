package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gofiber/fiber"

	"github.com/PuerkitoBio/goquery"
)

// Torrent is a struct to hold torrent data
type Torrent struct {
	Tname     string `json:"Name"`
	Tlink     string `json:"Magnetlink"`
	TSeeders  string `json:"Seeders"`
	TLeechers string `json:"Leechers"`
	TSize     string `json:"Size"`
}

var Link1337 string = "https://1337x.to"
var LinkNyaa string = "https://nyaa.si"
var wg sync.WaitGroup

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
		torrentList = append(torrentList, Torrent{Tname: tor.Text(), Tlink: torLink, TSeeders: torSeeders, TLeechers: torLeechers, TSize: torSize})
	})
	return torrentList, ""
}

// GetMagnet is a function to get magnet links from 1337x.to
func GetMagnet(Tlink *string) {
	defer wg.Done()
	Tlink1 := "https://1337x.to" + *Tlink
	res, _ := http.Get(Tlink1)
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	t, _ := doc.Find(".col-9.page-content").Find("a").Attr("href")
	*Tlink = t
}

// GetTorrents is a function to get torrents from 1337x.to
func GetTorrents(query string) ([]Torrent, string) {
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
			torLink, _ := tor.Attr("href")
			torSeeders := s.Find(".seeds").Eq(0).Text()
			torLeechers := s.Find(".leeches").Text()
			torSize := strings.Replace(s.Find(".size").Text(), torSeeders, "", -1)
			torrentList = append(torrentList, Torrent{Tname: tor.Text(), Tlink: torLink, TSeeders: torSeeders, TLeechers: torLeechers, TSize: torSize})
		}
	})
	return torrentList, ""
}

func main() {
	app := fiber.New()

	app.Get("/searchNyaa/:query?", func(c *fiber.Ctx) {
		tmp := c.Params("query")
		if tmp != "" {
			tmp = strings.ReplaceAll(tmp, "%20", "+")
			tmp = url.QueryEscape(tmp)
			tmp1, _ := GetAnimeTorrents(tmp)
			tmp2, _ := json.Marshal(tmp1)
			c.Send(tmp2)
		} else {
			c.Send("Where is the query?")
		}
	})

	app.Get("/search1337/:query?", func(c *fiber.Ctx) {
		tmp := c.Params("query")
		if tmp != "" {
			tmp = strings.ReplaceAll(tmp, "%20", "+")
			tmp = url.QueryEscape(tmp)
			tmp1, _ := GetTorrents(tmp)
			wg.Add(len(tmp1))
			for i := range tmp1 {
				go GetMagnet(&tmp1[i].Tlink)
			}
			wg.Wait()
			tmp2, _ := json.Marshal(tmp1)
			c.Send(tmp2)
		} else {
			c.Send("Where is the query?")
		}
	})

	app.Get("/*", func(c *fiber.Ctx) {
		c.Send("This isn't how this is supposed to be used. Usage examples: \n\nSearch for Far Cry on 1337x.to - torscraper.herokuapp.com/search1337/Far Cry \nSearch for Dragon Ball on nyaa.si - torscraper.herokuapp.com/searchNyaa/Dragon Ball")
	})

	port := os.Getenv("PORT")
	defaultPort := "8080"

	if port == "" {
		app.Listen(defaultPort)
	} else {
		app.Listen(port)
	}
}
