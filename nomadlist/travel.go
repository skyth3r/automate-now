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

	rsp, err := http.Get(url)
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
