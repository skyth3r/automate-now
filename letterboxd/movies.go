package letterboxd

import (
	"regexp"
)

const (
	Url                  = "https://letterboxd.com/"
	movieTitlePattern    = `, (\d{4})(?: - ?[★]{0,5}(½)?)?$`
	movieUrlWithUsername = `https:\/\/letterboxd\.com\/([^\/]+)\/`
)

func GetMovieTitle(input string) string {
	// Removes ', YYYY - ★★★★' from movie titles
	// The regex pattern looks for the following in a movie title:
	// - `, 2020` (No rating given)
	// - `, 2020 - ★★★★` (rating given)
	re := regexp.MustCompile(movieTitlePattern)
	title := re.Split(input, -1)
	return title[0]
}

func GetMovieUrl(movieUrl string) string {
	// Get Letterboxd item link without the username
	// Replaces "https://letterboxd.com/USERNAME_HERE/film/MOVIE_TITLE/" with "https://letterboxd.com/film/MOVIE_TITLE/"
	usernamePattern := regexp.MustCompile(movieUrlWithUsername)
	formattedUrl := usernamePattern.ReplaceAllString(movieUrl, Url)
	return formattedUrl
}
