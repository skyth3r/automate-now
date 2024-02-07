package nomadlist

type Profile struct {
	Stats Stats  `json:"stats"`
	Trips []Trip `json:"trips"`
}

type Stats struct {
	Cities                     int     `json:"cities"`
	Countries                  int     `json:"countries"`
	DistanceTraveledKM         int     `json:"distance_traveled_km"`
	DistanceTraveledMiles      int     `json:"distance_traveled_miles"`
	CountriesVisitedPercentage float64 `json:"countries_visited_percentage"`
	CitiesVisitedPercentage    float64 `json:"cities_visited_percentage"`
}

type Trip struct {
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Length      string `json:"length"`
	Place       string `json:"place"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}
