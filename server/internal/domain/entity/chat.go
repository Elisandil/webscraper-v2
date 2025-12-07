package entity

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatIntent struct {
	Action     string  `json:"action"`
	URL        string  `json:"url"`
	Frequency  string  `json:"frequency"`
	CronExpr   string  `json:"cron_expr"`
	Confidence float64 `json:"confidence"` // 0.0 to 1.0
}

type ChatResponse struct {
	Message      string      `json:"message"`
	Intent       *ChatIntent `json:"intent,omitempty"`
	NeedsConfirm bool        `json:"needs_confirm"`
	Action       string      `json:"action"`
}

type ChatConfirmation struct {
	Confirmed bool   `json:"confirmed"`
	IntentID  string `json:"intent_id"`
}
