package serializd

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const Url = "https://www.serializd.com/api/user/"

func GetShows(url string) ([]map[string]string, error) {
	var shows []map[string]string
	var diary SerializdDiary

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Request headers
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Referer", url)
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="123", "Not:A-Brand";v="8"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Android"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Mobile Safari/537.36")
	req.Header.Set("X-Requested-With", "serializd_vercel")

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", rsp.StatusCode)
	}

	// Check if the response is gzipped
	var reader io.Reader
	if rsp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(rsp.Body)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		reader = gz
	} else {
		reader = rsp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &diary); err != nil {
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
