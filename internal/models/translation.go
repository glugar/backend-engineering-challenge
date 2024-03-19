package models

type Translation struct {
	Timestamp      string  `json:"timestamp"`
	TranslationID  string  `json:"translation_id"`
	SourceLanguage string  `json:"source_language"`
	TargetLanguage string  `json:"target_language"`
	ClientName     string  `json:"client_name"`
	EventName      string  `json:"event_name"`
	NrWords        int     `json:"nr_words"`
	Duration       float64 `json:"duration"`
}
