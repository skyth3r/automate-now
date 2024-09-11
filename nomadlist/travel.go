package nomadlist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const Url = "https://nomadlist.com/@"

func GetTravel(url string) ([]map[string]string, error) {
	var countries []map[string]string
	var nomadListProfile Profile

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Request headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Mobile Safari/537.36")

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", rsp.StatusCode)
	}

	err = json.NewDecoder(rsp.Body).Decode(&nomadListProfile)
	if err != nil {
		return nil, err
	}

	trips := nomadListProfile.Trips

	for t := range trips {
		country := make(map[string]string)
		country["name"] = trips[t].Country
		country["place"] = trips[t].Place
		country["code"] = trips[t].CountryCode
		country["start_date"] = trips[t].DateStart
		countries = append(countries, country)
	}

	return countries, nil
}

func TripsInYear(tripsInput []map[string]string, year string) []map[string]string {
	var tripsOutput []map[string]string

	for _, trip := range tripsInput {
		if trip["start_date"][0:4] == year {
			tripsOutput = append(tripsOutput, trip)
		}
	}

	return tripsOutput
}
