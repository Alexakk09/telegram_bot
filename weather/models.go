package weather

type WeatherResponse struct {
	Name string `json:"name"`
	Main Main   `json:"main"`
	Weather []Weather `json:"weather"`
}

type Main struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
}

type Weather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}