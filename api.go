package main

import (
	"net/url"
	"os"
	"strings"

	"github.com/MadScientistAK/torscraper/torscraper"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	app.Get("/searchNyaa/:query?", func(c *fiber.Ctx) {
		tmp := c.Params("query")
		if tmp != "" {
			tmp = strings.ReplaceAll(tmp, "%20", "+")
			tmp = url.QueryEscape(tmp)
			c.Send(torscraper.GetAnimeTorrents(tmp))
		} else {
			c.Send("Where is the query?")
		}
	})

	app.Get("/search1337/:query?", func(c *fiber.Ctx) {
		tmp := c.Params("query")
		if tmp != "" {
			tmp = strings.ReplaceAll(tmp, "%20", "+")
			tmp = url.QueryEscape(tmp)
			c.Send(torscraper.GetTorrents(tmp))
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
