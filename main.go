package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Skyth3r/automate-now/backloggd"
	"github.com/Skyth3r/automate-now/letterboxd"
	"github.com/Skyth3r/automate-now/serializd"
	"github.com/mmcdole/gofeed"
)

func main() {

	// Movies
	latestMovieItems, err := getGoFeedItems(letterboxdRss)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}
	itemCount := maxGoFeedItems(latestMovieItems)
	movies := latestFeedItems(latestMovieItems, itemCount)

	// Books
	latestBookItems, err := getGoFeedItems(okuRss)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}
	itemCount = maxGoFeedItems(latestBookItems)
	books := latestFeedItems(latestBookItems, itemCount)

	// TV Shows
	showTitlesAndUrls, err := serializd.GetShows(serializdDiaryJson)
	if err != nil {
		log.Fatalf("unable to get shows from Serializd. Error: %v", err)
	}
	itemCount = maxItems(showTitlesAndUrls)
	shows := serializd.LatestShows(showTitlesAndUrls, itemCount)

	// Video games
	backloggdUrl := backloggd.Url + "/u/" + backloggdUsername + "/playing/"
	games, err := backloggd.GetGames(backloggdUrl)
	if err != nil {
		log.Fatalf("unable to get games from Backloggd. Error: %v", err)
	}

	// formatting Books
	booksHeader := "## 📚 Books\n"
	booksBody := formatMediaItems(books)

	moviesAndTvShowsHeader := "## 🎬 Movies and TV Shows\n"
	// formatting Movies
	moviesSubHeader := "### Recently watched movies\n"
	moviesBody := formatMediaItems(movies)

	// formatting TV Shows
	showsSubHeader := "### Recently watched TV shows\n"
	showsBody := formatMediaItems(shows)

	// formatting Video games
	gamesHeader := "## 🎮 Video Games\n"
	gamesBody := formatMediaItems(games)

	// get date
	date := time.Now().Format("2 Jan 2006")
	updated := fmt.Sprintf("\nLast updated: %v", date)

	staticContent, err := os.ReadFile("static.md")
	if err != nil {
		log.Fatalf("unable to read from static.md file. Error: %v", err)
	}

	// create now.md
	file, err := os.Create("now.md")
	if err != nil {
		log.Fatalf("unable to create now.md file. Error: %v", err)
	}
	defer file.Close()

	data := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n---\n%s", booksHeader, booksBody, moviesAndTvShowsHeader, moviesSubHeader, moviesBody, showsSubHeader, showsBody, gamesHeader, gamesBody, updated)
	data = fmt.Sprintf("%s\n\n%s", staticContent, data)

	_, err = io.WriteString(file, data)
	if err != nil {
		log.Fatalf("unable to write to now.md file. Error: %v", err)
	}
	file.Sync()

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

func latestFeedItems(items []gofeed.Item, count int) []map[string]string {
	var itemSlice = []map[string]string{}

	for i := 0; i < count; i++ {
		item := make(map[string]string)
		if strings.HasPrefix(items[i].Link, "https://letterboxd.com") {
			item["title"] = letterboxd.GetMovieTitle(items[i].Title)
			item["url"] = letterboxd.GetMovieUrl(items[i].Link)
		} else {
			item["title"] = items[i].Title
			item["url"] = items[i].Link
		}
		itemSlice = append(itemSlice, item)
	}
	return itemSlice
}

func maxItems(items []map[string]string) int {
	max := 3
	if len(items) < max {
		max = len(items)
	}
	return max
}

func formatMarkdownLink(title string, url string) string {
	return fmt.Sprintf("* [%v](%v)", title, url)
}

func formatMediaItems(mediaItems []map[string]string) string {
	var mediaText string
	for i := range mediaItems {
		itemText := formatMarkdownLink(mediaItems[i]["title"], mediaItems[i]["url"])
		mediaText += fmt.Sprintf("%v\n", itemText)
	}
	return mediaText
}
