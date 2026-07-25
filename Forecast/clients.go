package forecast

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func GetForecast(city string) (ForecastResponse, error) {

	apiKey := os.Getenv("WEATHER_API_KEY")

	api_url := "https://api.openweathermap.org/data/2.5/forecast?q=" +
		url.QueryEscape(city) +
		"&appid=" + apiKey +
		"&units=metric"
	fmt.Println("Forecast URL:", api_url)
	resp, err := http.Get(api_url)

	if err != nil {
		return ForecastResponse{}, err
	}
	fmt.Println("Status:", resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ForecastResponse{}, fmt.Errorf("could not fetch forecast")
	}
	var forecast ForecastResponse

	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		return ForecastResponse{}, err
	}
	return forecast, nil
}
