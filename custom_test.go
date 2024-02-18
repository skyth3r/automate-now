package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveLondonTrips(t *testing.T) {
	countries := []map[string]string{
		{
			"place": "London",
		},
		{
			"place": "New York City",
		},
		{
			"place": "Paris",
		},
		{
			"place": "London",
		},
		{
			"place": "Berlin",
		},
	}

	expected := []map[string]string{
		{
			"place": "New York City",
		},
		{
			"place": "Paris",
		},
		{
			"place": "Berlin",
		},
	}

	actual := removeLondonTrips(countries)

	assert.Equal(t, expected, actual)

}

func TestAddScotlandTrip2023(t *testing.T) {
	countries := []map[string]string{
		{
			"name": "United Kingdom",
		},
		{
			"name": "United States",
		},
		{
			"name": "Germany",
		},
	}

	expected := []map[string]string{
		{
			"name": "Scotland",
		},
		{
			"name": "United States",
		},
		{
			"name": "Germany",
		},
	}

	actual := addScotlandTrip2023(countries)

	assert.Equal(t, expected, actual)
}
