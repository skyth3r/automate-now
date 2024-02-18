package main

import (
	"log"
	"os"
)

func removeLondonTrips(countries []map[string]string) []map[string]string {
	var filteredCountries []map[string]string

	for _, trip := range countries {
		if trip["place"] == "London" {
			continue
		}
		filteredCountries = append(filteredCountries, trip)
	}

	return filteredCountries
}

func addScotlandTrip2023(countries []map[string]string) []map[string]string {
	var filteredCountries []map[string]string

	for _, trip := range countries {
		if trip["name"] == "United Kingdom" {
			trip["name"] = "Scotland"
		}
		filteredCountries = append(filteredCountries, trip)
	}

	return filteredCountries
}

func moveFile(fileName, filePath string) {
	if err := os.Rename(fileName, filePath); err != nil {
		log.Fatalf("unable to move %s to '%s'. Error: %v", fileName, filePath, err)
	}
}
