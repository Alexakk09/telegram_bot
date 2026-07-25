package geocoding

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
)

func GetCity(city string) (Result, error) {

	escapedCity := url.QueryEscape(city)

	apiURL := "https://geocoding-api.open-meteo.com/v1/search?name=" +
		escapedCity +
		"&count=1"

	resp, err := http.Get(apiURL)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	var geo GeoResponse

	err = json.NewDecoder(resp.Body).Decode(&geo)
	if err != nil {
		return Result{}, err
	}

	if len(geo.Results) == 0 {
	return Result{}, fmt.Errorf("city not found")
}

return geo.Results[0], nil

}
