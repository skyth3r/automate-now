package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/mmcdole/gofeed"
)

type SerializdDiary struct {
	Reviews      []SerializdDiaryReview `json:"reviews"`
	TotalPages   int                    `json:"total_pages"`
	TotalReviews int                    `json:"total_reviews"`
}

type SerializdDiaryReview struct {
	ID              int          `json:"id"`
	ShowID          string       `json:"show_id"`
	SeasonID        string       `json:"season_id"`
	DateAdded       string       `json:"date_added"`
	Rating          int          `json:"rating"`
	ReviewText      string       `json:"review_text"`
	Author          string       `json:"author"`
	AuthorImageUrl  string       `json:"author_image_url"`
	ContainsSpoiler bool         `json:"containsSpoilers"`
	BackDate        string       `json:"backdate"`
	ShowName        string       `json:"showName"`
	ShowBannerImage string       `json:"showBannerImage"`
	ShowSeasons     []ShowSeason `json:"showSeasons"`
}

type ShowSeason struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SeasonNumber int    `json:"seasonNumber"`
	PosterPath   string `json:"posterPath"`
}

func main() {
	const letterboxdRSS = "https://letterboxd.com/akashgoswami/rss/"
	const okuRSS = "https://oku.club/rss/collection/T8k9M"

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

	latestBookItems, err := getFeedItems(okuRSS)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}

	itemCount = maxItems(latestBookItems)

	printBookInfo(latestBookItems, itemCount)

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
