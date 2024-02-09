package serializd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const Url = "https://www.serializd.com/api/user/"

func GetShows(url string) ([]map[string]string, error) {
	var shows []map[string]string
	var diary SerializdDiary

	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", rsp.StatusCode)
	}

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

		// Loop through review.showSeasons to find season name using review.SeasonID
		for s := range review.ShowSeasons {
			season := review.ShowSeasons[s]
			if reviewSeasonID == season.ID {
				review.SeasonName = season.Name
			}
		}

		// format showName with SeasonName and store in output
		showAndSeason = fmt.Sprintf("%v, %v", review.ShowName, review.SeasonName)
		show["title"] = showAndSeason

		// get show url
		const showBaseUrl = "https://www.serializd.com/show/"
		showUrl := showBaseUrl + fmt.Sprint(review.ShowID)
		show["url"] = showUrl

		// Append show to shows if shows["title"] is not present in the map
		if !containsValue(shows, "title", show["title"]) {
			shows = append(shows, show)
		}
	}

	return shows, nil
}

func LatestShows(items []map[string]string, count int) []map[string]string {
	var shows []map[string]string
	for i := 0; i < count; i++ {
		shows = append(shows, items[i])
	}
	return shows
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
