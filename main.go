package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/mmcdole/gofeed"
)

func main() {
	const letterboxdRSS = "https://letterboxd.com/akashgoswami/rss/"

	latestMovieItems, err := getFeedItems(letterboxdRSS)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}

	itemCount := maxItems(latestMovieItems)

	// Regex pattern to remove ', YYYY - ★★★★' from movie titles
	// This regex pattern looks for the following in a movie title:
	// - `, 2020` (No rating given)
	// - `, 2020 - ★★★★` (rating given)
	const movieTitlePattern = `, (\d{4})(?: - ?[★]{0,5})?$`
	re := regexp.MustCompile(movieTitlePattern)

	printMovieTitles(latestMovieItems, itemCount, re)
}

func getFeedItems(input string) ([]gofeed.Item, error) {
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

func maxItems(items []gofeed.Item) int {
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
