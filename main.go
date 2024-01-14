package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/mmcdole/gofeed"
)

func main() {
	const letterboxdRSS = "https://letterboxd.com/akashgoswami/rss/"
	latestMovieItems := []gofeed.Item{}

	feedParser := gofeed.NewParser()
	latestMoviesFeed, err := feedParser.ParseURL(letterboxdRSS)
	if err != nil {
		log.Fatalf("unable to parse letterboxd rss url. Error: %v", err)
	}

	for _, item := range latestMoviesFeed.Items {
		latestMovieItems = append(latestMovieItems, *item)
	}

	max := 3 // Maximum number of movies to retrieve from feed
	if len(latestMovieItems) < max {
		max = len(latestMovieItems)
	}

	// Regex pattern to remove ', YYYY - ★★★★' from movie titles
	// This regex pattern looks for the following in a movie title:
	// - `, 2020` (No rating given)
	// - `, 2020 - ★★★★` (rating given)
	movieTitlePattern := `, (\d{4})(?: - ?[★]{0,5})?$`
	re := regexp.MustCompile(movieTitlePattern)

	for i := 0; i < max; i++ {
		splittedTitle := re.Split(latestMovieItems[i].Title, -1)
		fmt.Printf("Title: %v\n", splittedTitle[0])
	}

}
