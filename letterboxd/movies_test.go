package letterboxd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetMovieTitle(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"Movie Title, 2024", "Movie Title"},
		{"Movie Title, the sequel, 2023 - ★★★★★", "Movie Title, the sequel"},
		{"Movie - Title, 2022 - ★★★★", "Movie - Title"},
		{"Movie Title and the movie title, 2021 - ★★★", "Movie Title and the movie title"},
		{"Movie, Title, 2020 - ★★", "Movie, Title"},
		{"Movie, - Title, 2019 - ★", "Movie, - Title"},
		{"Movie Title, 2018 - ", "Movie Title"},
		{"Movie Title", "Movie Title"},                 // Edge case: No year or rating
		{"2024, Movie Title", "2024, Movie Title"},     // Edge case: Year at the start
		{"Movie Title - ★★★★★", "Movie Title - ★★★★★"}, // Edge case: Rating but no year
	}

	for i := range tests {
		title := tests[i].title
		expected := tests[i].expected

		actual := GetMovieTitle(title)
		require.Equal(t, expected, actual)
	}
}

func TestGetMovieUrl(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://letterboxd.com/USERNAME_HERE/film/Movie/", "https://letterboxd.com/film/Movie/"},
		{"https://letterboxd.com/USERNAME_HERE/film/Movie-Title", "https://letterboxd.com/film/Movie-Title"},
		{"https://letterboxd.com/USERNAME_HERE/film/Movie-Title-and-the-movie-title", "https://letterboxd.com/film/Movie-Title-and-the-movie-title"},
	}

	for i := range tests {
		url := tests[i].url
		expected := tests[i].expected

		actual := GetMovieUrl(url)
		require.Equal(t, expected, actual)
	}
}
