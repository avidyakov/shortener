package models

type ShortURLRequest struct {
	URL string `json:"url"`
}

type RequestLinkBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginURL     string `json:"original_url"`
}

type ResponseLinkBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
