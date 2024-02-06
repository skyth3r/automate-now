package letterboxd

import (
	"regexp"

	"github.com/mmcdole/gofeed"
)

func MovieTitles(items []gofeed.Item, count int) []map[string]string {
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
