package main

import (
	"regexp"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
)

func TestGetGoFeedItems_Success(t *testing.T) {
	const mockRSS = "https://akashgoswami.com/index.xml"

	_, err := getGoFeedItems(mockRSS)

	require.Nil(t, err)
}

func TestGetGoFeedItems_Error(t *testing.T) {
	const mockRSS = "https://akashgoswami.com/invalid_rss"

	_, err := getGoFeedItems(mockRSS)

	require.NotNil(t, err)
}

func TestMaxGoFeedItems_MoreThanThreeItems(t *testing.T) {
	mockItems := []gofeed.Item{
		{Title: "Title 1"},
		{Title: "Title 2"},
		{Title: "Title 3"},
		{Title: "Title 4"},
	}

	const expectedMax = 3

	got := maxGoFeedItems(mockItems)

	require.Equal(t, expectedMax, got)
}

func TestMaxGoFeedItems_LessThanThreeItems(t *testing.T) {
	mockItems := []gofeed.Item{
		{Title: "Title 1"},
		{Title: "Title 2"},
	}

	const expectedMax = 2

	got := maxGoFeedItems(mockItems)

	require.Equal(t, expectedMax, got)
}

func TestMaxGoFeedItems_ExactlyThreeItems(t *testing.T) {
	mockItems := []gofeed.Item{
		{Title: "Title 1"},
		{Title: "Title 2"},
		{Title: "Title 3"},
	}

	const expectedMax = 3

	got := maxGoFeedItems(mockItems)

	require.Equal(t, expectedMax, got)
}

func TestMovieTitlePattern(t *testing.T) {
	const movieTitlePattern = `, (\d{4})(?: - ?[★]{0,5})?$`

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

	re := regexp.MustCompile(movieTitlePattern)

	for _, tc := range tests {
		got := re.Split(tc.title, -1)[0]

		require.Equal(t, tc.expected, got)
	}
}
