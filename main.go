package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

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
	movies := movieTitles(latestMovieItems, itemCount)
	fmt.Println(movies)

	// Books
	latestBookItems, err := getGoFeedItems(urls.OkuRss)
	if err != nil {
		log.Fatalf("unable to parse rss url. Error: %v", err)
	}
	itemCount = maxGoFeedItems(latestBookItems)
	books := booksInfo(latestBookItems, itemCount)
	fmt.Println(books)

	// TV Shows
	showTitlesAndUrls, err := getShowDetails(urls.SerializdDiaryJson)
	if err != nil {
		log.Fatalf("unable to get shows from Serializd. Error: %v", err)
	}
	itemCount = maxItems(showTitlesAndUrls)
	shows := showDetails(showTitlesAndUrls, itemCount)
	fmt.Println(shows)

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

func movieTitles(items []gofeed.Item, count int) []map[string]string {
	var movies = []map[string]string{}

	for i := 0; i < count; i++ {
		movie := make(map[string]string)
		movie["title"] = getMovieTitle(items[i].Title)
		movie["url"] = getMovieUrl(items[i].Link)
		movies = append(movies, movie)
	}

	return movies
}

func getMovieTitle(input string) string {
	// Regex pattern to remove ', YYYY - ★★★★' from movie titles
	// This regex pattern looks for the following in a movie title:
	// - `, 2020` (No rating given)
	// - `, 2020 - ★★★★` (rating given)
	const movieTitlePattern = `, (\d{4})(?: - ?[★]{0,5})?$`
	re := regexp.MustCompile(movieTitlePattern)

	title := re.Split(input, -1)

	return title[0]
}

func getMovieUrl(url string) string {
	// Get Letterboxd item link without the username
	// Replaces "https://letterboxd.com/USERNAME_HERE/film/MOVIE_TITLE/" with "https://letterboxd.com/film/MOVIE_TITLE/"
	regexPattern := regexp.MustCompile(`https:\/\/letterboxd\.com\/([^\/]+)\/`)
	match := regexPattern.ReplaceAllString(url, "https://letterboxd.com/")
	//fmt.Printf("Movie URL: %v\n", match)
	return match
}

func booksInfo(items []gofeed.Item, count int) []map[string]string {
	var books = []map[string]string{}

	for i := 0; i < count; i++ {
		book := make(map[string]string)
		book["title"] = items[i].Title
		book["url"] = items[i].Link
		books = append(books, book)
	}

	return books
}

func getShowDetails(url string) ([]map[string]string, error) {
	var shows = []map[string]string{}

	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", rsp.StatusCode)
	}

	var diary serializd.SerializdDiary

	err = json.NewDecoder(rsp.Body).Decode(&diary)
	if err != nil {
		return nil, err
	}

	reviews := diary.Reviews

	for r := range reviews {
		show := make(map[string]string)
		var showAndSeason string
		review := reviews[r]
		reviewSeasonID := review.SeasonID

		// Loop through review.showSeasons to find season name using reviewSesonID
		for s := range review.ShowSeasons {
			season := review.ShowSeasons[s]
			if reviewSeasonID == season.ID {
				review.SeasonName = season.Name
			}
		}

		// format showName with SeasonName and store in output
		showAndSeason = fmt.Sprintf("%v, %v", review.ShowName, review.SeasonName)
		show["title"] = showAndSeason
		//fmt.Printf("Show Title: %v\n", showAndSeason)

		// get show url
		const showBaseUrl = "https://www.serializd.com/show/"
		showUrl := showBaseUrl + fmt.Sprint(review.ShowID)
		show["url"] = showUrl

		// Append show to shows if shows["title'"] is not present in the map
		if !containsValue(shows, "title", show["title"]) {
			shows = append(shows, show)
		}
	}

	return shows, nil
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

func maxItems(items []map[string]string) int {
	max := 3
	if len(items) < max {
		max = len(items)
	}
	return max
}
func showDetails(items []map[string]string, count int) []map[string]string {
	var cappedShows = []map[string]string{}
	for i := 0; i < count; i++ {
		cappedShows = append(cappedShows, items[i])
		//fmt.Printf("%v\n", items[i])
	}
	return cappedShows
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
