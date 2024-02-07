package backloggd

import (
	"errors"

	"github.com/Skyth3r/automate-now/urls"
	"github.com/gocolly/colly"
)

func GetGames(url string) ([]map[string]string, error) {
	var games = []map[string]string{}

	c := colly.NewCollector()

	c.OnHTML("div.rating-hover", func(e *colly.HTMLElement) {
		game := make(map[string]string)
		game["title"] = e.ChildText("div.game-text-centered")
		game["url"] = urls.BackloggdBase + e.ChildAttr("a", "href")
		games = append(games, game)
	})

	c.Visit(url)

	if len(games) == 0 {
		err := errors.New("no games found")
		return nil, err
	}

	return games, nil
}
