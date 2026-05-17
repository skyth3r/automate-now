package backloggd

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

const Url = "https://backloggd.com"

func GetGames(url string) ([]map[string]string, error) {
	var games []map[string]string

	wsURL, err := launcher.New().Headless(true).NoSandbox(true).Launch()
	if err != nil {
		return nil, fmt.Errorf("browser launch failed: %w", err)
	}

	browser := rod.New().ControlURL(wsURL)
	if err = browser.Connect(); err != nil {
		return nil, fmt.Errorf("browser connect failed: %w", err)
	}
	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		return nil, fmt.Errorf("page open failed: %w", err)
	}

	err = rod.Try(func() {
		page.Timeout(30 * time.Second).MustElement("div.rating-hover")
	})
	if err != nil {
		return games, nil
	}

	elements, err := page.Elements("div.rating-hover")
	if err != nil {
		return nil, err
	}

	for _, el := range elements {
		titleEl, err := el.Element("div.game-text-centered")
		if err != nil {
			fmt.Println("no game title element found, skipping")
			continue
		}
		title, err := titleEl.Text()
		if err != nil {
			fmt.Println("no game title found, skipping")
			continue
		}

		linkEl, err := el.Element("a")
		if err != nil {
			fmt.Println("no game url element found, skipping")
			continue
		}
		href, err := linkEl.Attribute("href")
		if err != nil || href == nil {
			fmt.Println("no game url found, skipping")
			continue
		}

		games = append(games, map[string]string{
			"title": title,
			"url":   Url + *href,
		})
	}

	return games, nil
}
