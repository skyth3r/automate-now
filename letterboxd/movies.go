package letterboxd

import (
	"regexp"
)

const Url = "https://letterboxd.com/"

func GetMovieTitle(input string) string {
	// Regex pattern to remove ', YYYY - ★★★★' from movie titles
	// This regex pattern looks for the following in a movie title:
	// - `, 2020` (No rating given)
	// - `, 2020 - ★★★★` (rating given)
	const movieTitlePattern = `, (\d{4})(?: - ?[★]{0,5})?$`
	re := regexp.MustCompile(movieTitlePattern)

	title := re.Split(input, -1)

	return title[0]
}

func GetMovieUrl(url string) string {
	// Get Letterboxd item link without the username
	// Replaces "https://letterboxd.com/USERNAME_HERE/film/MOVIE_TITLE/" with "https://letterboxd.com/film/MOVIE_TITLE/"
	regexPattern := regexp.MustCompile(`https:\/\/letterboxd\.com\/([^\/]+)\/`)
	match := regexPattern.ReplaceAllString(url, "https://letterboxd.com/")
	//fmt.Printf("Movie URL: %v\n", match)
	return match
}
