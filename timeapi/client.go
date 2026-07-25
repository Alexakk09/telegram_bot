package timeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func GetCurrentTime(timezone string) (TimeResponse, error) {

	apiKey := os.Getenv("TIME_API_KEY")
	fmt.Println("API Key:", apiKey)

	apiURL := "https://api.api-ninjas.com/v1/worldtime?timezone=" +
		url.QueryEscape(timezone)

	fmt.Println("URL:", apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return TimeResponse{}, err
	}

	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return TimeResponse{}, err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TimeResponse{}, err
	}

	fmt.Println("Response:", string(body))

	if resp.StatusCode != http.StatusOK {
		return TimeResponse{}, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var timeData TimeResponse

	err = json.Unmarshal(body, &timeData)
	if err != nil {
		return TimeResponse{}, err
	}

	return timeData, nil
}