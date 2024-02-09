package backloggd

import (
	"errors"

	"github.com/gocolly/colly"
)

const Url = "https://backloggd.com"

func GetGames(url string) ([]map[string]string, error) {
	var games []map[string]string

	c := colly.NewCollector()

	c.OnHTML("div.rating-hover", func(e *colly.HTMLElement) {
		game := make(map[string]string)
		game["title"] = e.ChildText("div.game-text-centered")
		game["url"] = Url + e.ChildAttr("a", "href")
		games = append(games, game)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		err := errors.New("no games found")
		return nil, err
	}

	return games, nil
}
