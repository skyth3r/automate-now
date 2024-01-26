package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"

	"github.com/Skyth3r/automate-now/serializd"
	"github.com/Skyth3r/automate-now/urls"
	"github.com/mmcdole/gofeed"
)

func main() {

	// Movies
	latestMovieItems, err := getFeedItems(urls.LetterboxdRss)
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

	latestBookItems, err := getFeedItems(urls.OkuRss)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}

	itemCount = maxItems(latestBookItems)

	printBookInfo(latestBookItems, itemCount)

	// TV Shows
	rsp, err := http.Get(urls.SerializdDiaryJson)
	if err != nil {
		log.Fatalf("unable to get json from serializd. Error: %v", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %v", rsp.StatusCode)
	}

	var diary serializd.SerializdDiary

	err = json.NewDecoder(rsp.Body).Decode(&diary)
	if err != nil {
		log.Fatalf("unable to decode json. Error: %v", err)
	}

	reviews := diary.Reviews

	var shows []string

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
	fmt.Printf("%v\n", shows)

	// Video games
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

func printBookInfo(items []gofeed.Item, count int) {
	for i := 0; i < count; i++ {
		fmt.Println(items[i].Title)
		fmt.Println(items[i].Link)
		fmt.Println(items[i].Extensions["dc"]["creator"][0].Value) // author
		fmt.Println(items[i].Extensions["oku"]["cover"][0].Value)  // book cover url
	}
}
