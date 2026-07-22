package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetCurrent(city string) (WeatherResponse, error) {

	apiKey := os.Getenv("WEATHER_API_KEY")

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		city,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return WeatherResponse{},fmt.Errorf("city not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weather WeatherResponse

	err = json.Unmarshal(body, &weather)
	if err != nil {
		return WeatherResponse{}, err
	}

	return weather, nil
}