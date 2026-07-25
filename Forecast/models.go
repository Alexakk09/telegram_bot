package forecast

type ForecastResponse struct {
	List []ForecastItem `json:"list"`
	City City           `json:"city"`
}

type ForecastItem struct {
	DateTime string `json:"dt_txt"`
	Main     Main   `json:"main"`
	Weather  []Info `json:"weather"`
}

type Main struct {
	Temp float64 `json:"temp"`
}

type Info struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

type City struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}