package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	emoji "github.com/jayco/go-emoji-flag"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"

	"github.com/Skyth3r/automate-now/backloggd"
	"github.com/Skyth3r/automate-now/letterboxd"
	"github.com/Skyth3r/automate-now/nomadlist"
	"github.com/Skyth3r/automate-now/serializd"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	// Movies
	latestMovieItems, err := getGoFeedItems(fmt.Sprintf("%s%s/rss/", letterboxd.Url, os.Getenv("LETTERBOXDUSERNAME")))
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}
	itemCount := maxItems(latestMovieItems)
	movies := latestGoFeedItems(latestMovieItems, itemCount)

	// Books
	latestBookItems, err := getGoFeedItems(fmt.Sprintf("%s%s", OkuUrl, os.Getenv("OKUCOLLECTIONID")))
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}
	itemCount = maxItems(latestBookItems)
	books := latestGoFeedItems(latestBookItems, itemCount)

	// TV Shows
	showTitlesAndUrls, err := serializd.GetShows(fmt.Sprintf("%s%s/diary", serializd.Url, os.Getenv("SERIALIZDUSERNAME")))
	if err != nil {
		log.Fatalf("unable to get shows from Serializd. Error: %v", err)
	}
	itemCount = maxItems(showTitlesAndUrls)
	shows := serializd.LatestShows(showTitlesAndUrls, itemCount)

	// Video games
	games, err := backloggd.GetGames(fmt.Sprintf("%s/u/%s/playing/", backloggd.Url, os.Getenv("BACKLOGGDUSERNAME")))
	if err != nil {
		log.Fatalf("unable to get games from Backloggd. Error: %v", err)
	}

	// Travel
	countries, err := nomadlist.GetTravel(fmt.Sprintf("%s%s.json", nomadlist.Url, os.Getenv("NOMADLISTUSERNAME")))
	if err != nil {
		log.Fatalf("unable to get countries from Nomadlist. Error: %v", err)
	}

	var dataString strings.Builder
	// Formatting Travel
	dataString.WriteString("## 🌏 Travel\n\n")
	dataString.WriteString("*Data sourced from [Nomads](https://nomads.com/)*\n\n")

	year := time.Now().Format("2006")
	yearsWithTrips := 0

	for yearsWithTrips < 2 {
		var trips []map[string]string

		// Get trips in current year, if there are no trips in the current year, check previous year
		for {
			trips = nomadlist.TripsInYear(countries, year)
			if len(trips) > 0 {
				break
			}

			year, err = decreaseYear(year)
			if err != nil {
				log.Fatalf("unable to decrease year. Error: %v", err)
			}
		}

		// latest year wth trips
		dataString.WriteString(fmt.Sprintf("### %s\n\n", year))
		trips = removeLondonTrips(trips)
		//check if year is 2023
		if year == "2023" {
			trips = addScotlandTrip2023(trips)
		}
		countriesVisited := removeDupes(trips)
		dataString.WriteString(formatCountries(countriesVisited))

		year, err = decreaseYear(year)
		if err != nil {
			log.Fatalf("unable to decrease year. Error: %v", err)
		}

		trips = nil
		yearsWithTrips++
	}

	// Formatting Books
	dataString.WriteString("## 📚 Books\n\n")
	dataString.WriteString("*Data sourced from [Oku](https://oku.club/)*\n\n")
	dataString.WriteString(formatMediaItems(books, "books"))

	// Formatting Movies and TV Shows
	dataString.WriteString("## 🎬 Movies and TV Shows\n\n")
	// Formatting Movies
	dataString.WriteString("### Recently watched movies\n\n")
	dataString.WriteString("*Data sourced from [Letterboxd](https://letterboxd.com/)*\n\n")
	dataString.WriteString(formatMediaItems(movies, "movies"))

	// Formatting TV Shows
	dataString.WriteString("### Recently watched TV shows\n\n")
	dataString.WriteString("*Data sourced from [Serializd](https://www.serializd.com/)*\n\n")
	dataString.WriteString(formatMediaItems(shows, "tv shows"))

	// Formatting Video games
	dataString.WriteString("## 🎮 Video Games\n\n")
	dataString.WriteString("*Data sourced from [Backloggd](https://backloggd.com/)*\n\n")
	dataString.WriteString(formatMediaItems(games, "video games"))

	dataString.WriteString("---\n\n")
	// Get today's date
	date := time.Now().Format("2 Jan 2006")
	dataString.WriteString("Last updated: ")
	dataString.WriteString(date)

	staticContent, err := os.ReadFile("static.md")
	if err != nil {
		log.Fatalf("unable to read from static.md file. Error: %v", err)
	}

	// Create now.md
	file, err := os.Create("now.md")
	if err != nil {
		log.Fatalf("unable to create now.md file. Error: %v", err)
	}
	defer file.Close()

	data := fmt.Sprintf("%s\n\n%s", staticContent, dataString.String())

	_, err = io.WriteString(file, data)
	if err != nil {
		log.Fatalf("unable to write to now.md file. Error: %v", err)
	}
	err = file.Sync()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//moveFile("now.md", "../content/now.md")

}

func getGoFeedItems(input string) ([]gofeed.Item, error) {
	var feedItems []gofeed.Item

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

func latestGoFeedItems(items []gofeed.Item, count int) []map[string]string {
	var itemSlice []map[string]string

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

func removeDupes(trips []map[string]string) []map[string]string {
	var countries []map[string]string

	// sorts trips from oldest to newest
	slices.Reverse(trips)

	for _, trip := range trips {
		// check if a trip["name"] is present in the slice countries
		if !containsValue(countries, "name", trip["name"]) {
			countries = append(countries, trip)
		}
	}

	return countries
}

func containsValue(slice []map[string]string, key, value string) bool {
	for _, m := range slice {
		if _, ok := m[key]; ok {
			if val, ok := m[key]; ok && val == value {
				return true
			}
		}
	}
	return false
}

func formatMarkdownLink(title string, url string) string {
	title = escapeMarkdown(title)
	return fmt.Sprintf("* [%v](%v)", title, url)
}

func escapeMarkdown(text string) string {
	text = strings.NewReplacer(
		"&", "and",
	).Replace(text)

	text = strings.NewReplacer(
		"'", "\\'",
	).Replace(text)

	return text
}

func formatMediaItems(mediaItems []map[string]string, mediaType string) string {
	var mediaText string

	// check for empty state mediaItems maps
	if len(mediaItems) == 0 {
		switch mediaType {
		case "movies":
			mediaText = NoMovies
		case "books":
			mediaText = NoBooks
		case "tv shows":
			mediaText = NoTvShows
		case "video games":
			mediaText = NoVideoGames
		}
		mediaText += "\n"
		return mediaText
	}

	for i := range mediaItems {
		itemText := formatMarkdownLink(mediaItems[i]["title"], mediaItems[i]["url"])
		mediaText += fmt.Sprintf("%v\n", itemText)
	}
	mediaText += "\n"
	return mediaText

}

func formatCountries(countries []map[string]string) string {
	var formattedText string
	var countryEmoji string

	if len(countries) == 0 {
		formattedText = NoCountries + "\n\n"
		return formattedText
	}

	for i := range countries {
		if countries[i]["code"] == "UK" {
			countries[i]["code"] = "GB"
		}
		if countries[i]["name"] == "Scotland" {
			countryEmoji = "\U0001F3F4\U000E0067\U000E0062\U000E0073\U000E0063\U000E0074\U000E007F"
		} else {
			countryEmoji = emoji.GetFlag(countries[i]["code"])
		}
		countryText := fmt.Sprintf("%s %s\n\n", countryEmoji, countries[i]["name"])
		formattedText += countryText
	}

	return formattedText
}

func maxItems[T gofeed.Item | map[string]string](items []T) int {
	limit := 3
	if len(items) < limit {
		limit = len(items)
	}
	return limit
}

func decreaseYear(y string) (string, error) {
	yearInt, err := strconv.Atoi(y)
	if err != nil {
		return "", fmt.Errorf("unable to convert year to integer. Error: %v", err)
	}
	year := strconv.Itoa(yearInt - 1)
	return year, nil
}
