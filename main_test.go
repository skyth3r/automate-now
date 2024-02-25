package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}

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

	actual := maxItems(mockItems)
	require.Equal(t, 3, actual)
}

func TestMaxGoFeedItems_LessThanThreeItems(t *testing.T) {
	mockItems := []gofeed.Item{
		{Title: "Title 1"},
		{Title: "Title 2"},
	}

	actual := maxItems(mockItems)
	require.Equal(t, 2, actual)
}

func TestMaxGoFeedItems_ExactlyThreeItems(t *testing.T) {
	mockItems := []gofeed.Item{
		{Title: "Title 1"},
		{Title: "Title 2"},
		{Title: "Title 3"},
	}

	actual := maxItems(mockItems)
	require.Equal(t, 3, actual)
}

func TestLatestFeedItems(t *testing.T) {
	mockItems := []gofeed.Item{
		{Title: "Title 1", Link: "www.test.com"},
		{Title: "Title 2", Link: "www.test.com"},
		{Title: "Title 3", Link: "www.test.com"},
		{Title: "Title 4", Link: "www.test.com"},
		{Title: "Title 5", Link: "www.test.com"},
	}
	expected := []map[string]string{
		{"title": "Title 1", "url": "www.test.com"},
		{"title": "Title 2", "url": "www.test.com"},
	}

	actual := latestGoFeedItems(mockItems, 2)
	assert.Equal(t, expected, actual)
}

func TestRemoveDupes_DupesPresent(t *testing.T) {
	mockTrips := []map[string]string{
		{"name": "America"},
		{"name": "Belgium"},
		{"name": "Croatia"},
		{"name": "America"},
	}
	expected := []map[string]string{
		{"name": "America"},
		{"name": "Croatia"},
		{"name": "Belgium"},
	}

	actual := removeDupes(mockTrips)
	assert.Equal(t, expected, actual)
}

func TestRemoveDupes_NoDupes(t *testing.T) {
	mockTrips := []map[string]string{
		{"name": "America"},
		{"name": "Belgium"},
		{"name": "Croatia"},
	}
	expected := []map[string]string{
		{"name": "Croatia"},
		{"name": "Belgium"},
		{"name": "America"},
	}

	actual := removeDupes(mockTrips)
	assert.Equal(t, expected, actual)
}

func TestFormatMarkdownLink(t *testing.T) {
	mockTitle := "Test Title"
	mockUrl := "https://example.com"
	expected := "* [Test Title](https://example.com)"

	actual := formatMarkdownLink(mockTitle, mockUrl)
	assert.Equal(t, expected, actual)
}

func TestFormatMediaItems(t *testing.T) {
	mockItems := []map[string]string{
		{"title": "Title 1", "url": "https://example.com"},
		{"title": "Title 2", "url": "https://example.com"},
	}
	expected := "* [Title 1](https://example.com)\n* [Title 2](https://example.com)\n\n"

	actual := formatMediaItems(mockItems, "movies")
	assert.Equal(t, expected, actual)
}

func TestFormatCountries(t *testing.T) {
	countries := []map[string]string{
		{"name": "America"},
		{"name": "Belgium"},
		{"name": "Croatia"},
	}
	expected := " America\n\n Belgium\n\n Croatia\n\n"

	actual := formatCountries(countries)
	assert.Equal(t, expected, actual)
}
