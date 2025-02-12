package models

// URLMapping структура для хранения URL и его сокращённого идентификатора.
type URLModel struct {
	ID  string
	URL string
}
type URLBatchModel struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponseModel struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
