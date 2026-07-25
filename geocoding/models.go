package geocoding

type GeoResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Name     string `json:"name"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone string `json:"timezone"`
	Country string `json:"country"`
}