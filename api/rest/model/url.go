package model

type (
	RequestURL struct {
		URL string `json:"url"`
	}

	ResponseURL struct {
		ShortURL string `json:"result"`
	}

	ResponseFullURL struct {
		URL      string `json:"original_url"`
		ShortURL string `json:"short_url"`
	}
)
