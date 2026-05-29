package model

type Url struct {
	ID          int    `json:"id"`
	URLOriginal string `json:"url_original"`
	ShortCode   string `json:"short_code"`
	AccessCount int    `json:"access_count"`
}
