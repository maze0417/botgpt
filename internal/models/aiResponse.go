package models

type AiResponse struct {
	IsImage     bool   `json:"is_image"`
	IsText      bool   `json:"is_text"`
	Text        string `json:"text,omitempty"`
	TgParseMode string `json:"tg_parse_mode,omitempty"`
	CommandMode string `json:"command_mode,omitempty"`
}
