package timeapi

type TimeResponse struct {
	Datetime string `json:"datetime"`
	Timezone string `json:"timezone"`
	Date     string `json:"date"`
	Day      string `json:"day_of_week"`
	Hour     string `json:"hour"`
	Minute   string `json:"minute"`
	Second   string `json:"second"`
	Year     string `json:"year"`
	Month    string `json:"month"`
	DayNum   string `json:"day"`
}