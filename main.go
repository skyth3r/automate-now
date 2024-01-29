package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"

	"github.com/Skyth3r/automate-now/backloggd"
	"github.com/Skyth3r/automate-now/serializd"
	"github.com/Skyth3r/automate-now/urls"
	"github.com/gocolly/colly"
	"github.com/mmcdole/gofeed"
)

func main() {

	// Movies
	latestMovieItems, err := getGoFeedItems(urls.LetterboxdRss)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}

	itemCount := maxGoFeedItems(latestMovieItems)

	// Regex pattern to remove ', YYYY - ★★★★' from movie titles
	// This regex pattern looks for the following in a movie title:
	// - `, 2020` (No rating given)
	// - `, 2020 - ★★★★` (rating given)
	const movieTitlePattern = `, (\d{4})(?: - ?[★]{0,5})?$`
	re := regexp.MustCompile(movieTitlePattern)

	printMovieTitles(latestMovieItems, itemCount, re)

	latestBookItems, err := getGoFeedItems(urls.OkuRss)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}

	itemCount = maxGoFeedItems(latestBookItems)

	printBookInfo(latestBookItems, itemCount)

	// TV Shows
	shows, err := getShowNames(urls.SerializdDiaryJson)
	if err != nil {
		log.Fatalf("unable to get show names. Error: %v", err)
	}

	itemCount = maxItems(shows)

	printShows(shows, itemCount)

	// Video games
	backloggdUrl := urls.BackloggdBase + "/u/" + urls.BackloggdUsername + "/playing/"

	currentGames := getBackloggdGames(backloggdUrl)

	for i := range currentGames {
		fmt.Printf("%v\n", currentGames[i].Name)
		fmt.Printf("%v\n", currentGames[i].Url)
	}

}

func getGoFeedItems(input string) ([]gofeed.Item, error) {
	feedItems := []gofeed.Item{}

	feedParser := gofeed.NewParser()
	feed, err := feedParser.ParseURL(input)

	if err != nil {
		return nil, err
	}

	for _, item := range feed.Items {
		feedItems = append(feedItems, *item)
	}

	return feedItems, nil
}

func maxGoFeedItems(items []gofeed.Item) int {
	max := 3 // Maximum number of movies to retrieve from feed
	if len(items) < max {
		max = len(items)
	}
	return max
}

func printMovieTitles(items []gofeed.Item, count int, re *regexp.Regexp) {
	for i := 0; i < count; i++ {
		title := re.Split(items[i].Title, -1)
		fmt.Printf("Title: %v\n", title[0])
	}
}

func printBookInfo(items []gofeed.Item, count int) {
	for i := 0; i < count; i++ {
		fmt.Println(items[i].Title)
		fmt.Println(items[i].Link)
		fmt.Println(items[i].Extensions["dc"]["creator"][0].Value) // author
		fmt.Println(items[i].Extensions["oku"]["cover"][0].Value)  // book cover url
	}
}

func getShowNames(url string) ([]string, error) {
	var shows []string

	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
		//log.Fatalf("unable to get json from serializd. Error: %v", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", rsp.StatusCode)
		//log.Fatalf("unexpected status code: %v", rsp.StatusCode)
	}

	var diary serializd.SerializdDiary

	err = json.NewDecoder(rsp.Body).Decode(&diary)
	if err != nil {
		return nil, err
		//log.Fatalf("unable to decode json. Error: %v", err)
	}

	reviews := diary.Reviews

	for r := range reviews {
		var showAndSeason string
		review := reviews[r]
		reviewSesonID := review.SeasonID

		// Loop through review.showSeasons to find season name using reviewSesonID
		for s := range review.ShowSeasons {
			season := review.ShowSeasons[s]
			if reviewSesonID == season.ID {
				review.SeasonName = season.Name
			}
		}

		// format showName with SeasonName and store in output
		showAndSeason = fmt.Sprintf("%v, %v", review.ShowName, review.SeasonName)

		// Add show name to showNames array
		if !slices.Contains(shows, showAndSeason) {
			shows = append(shows, showAndSeason)
		}
	}

	return shows, nil
}

func maxItems(items []string) int {
	max := 3
	if len(items) < max {
		max = len(items)
	}
	return max
}

func printShows(items []string, count int) {
	for i := 0; i < count; i++ {
		fmt.Printf("%v\n", items[i])
	}
}

func getBackloggdGames(url string) []backloggd.CurrentGame {
	var currentGames []backloggd.CurrentGame

	c := colly.NewCollector()

	c.OnHTML("div.rating-hover", func(e *colly.HTMLElement) {
		game := backloggd.CurrentGame{}

		game.Name = e.ChildText("div.game-text-centered")
		partialUrl := e.ChildAttr("a", "href")
		game.Url = urls.BackloggdBase + partialUrl

		currentGames = append(currentGames, game)
	})

	c.Visit(url)

	return currentGames
}
