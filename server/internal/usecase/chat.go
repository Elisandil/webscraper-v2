package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/infrastructure/config"
	pkgerrors "webscraper-v2/pkg/errors"
)

type ChatUseCase struct {
	config     *config.Config
	httpClient *http.Client
	hfAPIToken string
	hfModelID  string
}

func NewChatUseCase(cfg *config.Config) *ChatUseCase {
	modelID := "google/flan-t5-small"
	if cfg.Chat != nil && cfg.Chat.HFModelID != "" {
		modelID = cfg.Chat.HFModelID
	}

	token := ""
	if cfg.Chat != nil {
		token = cfg.Chat.HFAPIToken
	}

	return &ChatUseCase{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		hfAPIToken: token,
		hfModelID:  modelID,
	}
}

type hfRequest struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type hfResponse []struct {
	GeneratedText string `json:"generated_text"`
}

func (uc *ChatUseCase) InterpretMessage(message string, parentCtx context.Context) (*entity.ChatIntent, error) {
	intent := uc.trySimpleRules(message)
	if intent.Action != "unknown" && intent.Confidence > 0.7 {
		return intent, nil
	}

	if intent.Action != "unknown" && intent.URL != "" {
		return intent, nil
	}

	if uc.hfAPIToken != "" || uc.config.Chat != nil {
		llmIntent, err := uc.interpretWithLLM(message, parentCtx)
		if err == nil {
			return llmIntent, nil
		}
	}

	return intent, nil
}

func (uc *ChatUseCase) trySimpleRules(message string) *entity.ChatIntent {
	originalMessage := message
	message = strings.ToLower(strings.TrimSpace(message))
	intent := &entity.ChatIntent{
		Action:     "unknown",
		Confidence: 0.0,
	}

	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	urlMatches := urlRegex.FindAllString(originalMessage, -1)
	if len(urlMatches) > 0 {
		intent.URL = urlMatches[0]
		intent.Confidence += 0.4
	}

	scheduleKeywords := []string{"programa", "programar", "schedule", "cada", "todos los", "diariamente", "semanalmente", "horario"}
	hasScheduleKeyword := false
	for _, keyword := range scheduleKeywords {
		if strings.Contains(message, keyword) {
			hasScheduleKeyword = true
			break
		}
	}

	if hasScheduleKeyword {
		intent.Action = "create_schedule"
		intent.Confidence += 0.5

		cronExpr, freq := uc.extractCronFromText(message)
		intent.CronExpr = cronExpr
		intent.Frequency = freq
		return intent
	}

	immediateKeywords := []string{"ahora", "inmediatamente", "ya", "scrapea", "escanea", "analiza"}
	for _, keyword := range immediateKeywords {
		if strings.Contains(message, keyword) {
			intent.Action = "scrape_now"
			intent.Confidence += 0.5
			break
		}
	}

	if intent.URL != "" && intent.Action == "unknown" {
		intent.Action = "scrape_now"
		intent.Confidence = 0.6
	}

	return intent
}

func (uc *ChatUseCase) extractCronFromText(message string) (cronExpr, frequency string) {
	message = strings.ToLower(message)

	patterns := map[string]struct {
		regex string
		cron  string
		freq  string
	}{
		"every_n_hours": {
			regex: `cada\s+(\d+)\s+horas?`,
			cron:  "0 0 */%d * * *",
			freq:  "cada %d horas",
		},
		"every_n_minutes": {
			regex: `cada\s+(\d+)\s+minutos?`,
			cron:  "0 */%d * * * *",
			freq:  "cada %d minutos",
		},
		"daily": {
			regex: `(diariamente|cada d[ií]a|todos los d[ií]as)`,
			cron:  "0 0 0 * * *",
			freq:  "diariamente",
		},
		"daily_at_time": {
			regex: `todos los d[ií]as a las? (\d+)(?::(\d+))?`,
			cron:  "0 0 %d * * *",
			freq:  "diariamente a las %d:00",
		},
		"hourly": {
			regex: `cada hora`,
			cron:  "0 0 * * * *",
			freq:  "cada hora",
		},
		"weekly": {
			regex: `semanalmente|cada semana`,
			cron:  "0 0 0 * * 0",
			freq:  "semanalmente",
		},
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.regex)
		matches := re.FindStringSubmatch(message)

		if len(matches) > 0 {
			switch {
			case strings.Contains(pattern.regex, "every_n"):
				n, _ := strconv.Atoi(matches[1])
				return fmt.Sprintf(pattern.cron, n), fmt.Sprintf(pattern.freq, n)
			case strings.Contains(pattern.regex, "daily_at_time"):
				hour, _ := strconv.Atoi(matches[1])
				return fmt.Sprintf("0 0 %d * * *", hour), fmt.Sprintf("diariamente a las %d:00", hour)
			default:
				return pattern.cron, pattern.freq
			}
		}
	}

	return "0 0 0 * * *", "diariamente (por defecto)"
}

func (uc *ChatUseCase) interpretWithLLM(message string, parentCtx context.Context) (*entity.ChatIntent, error) {
	prompt := uc.buildPrompt(message)

	response, err := uc.callHuggingFace(prompt, parentCtx)
	if err != nil {
		return nil, pkgerrors.InternalError("error calling HuggingFace API", err)
	}

	intent, err := uc.parseModelResponse(response)
	if err != nil {
		return nil, pkgerrors.InternalError("error parsing model response", err)
	}

	return intent, nil
}

func (uc *ChatUseCase) buildPrompt(message string) string {
	return fmt.Sprintf(`Extract information from this web scraping request:
Input: "%s"

Extract:
1. Action (scrape_now or create_schedule)
2. URL (if present)
3. Frequency (if scheduling: hourly, daily, every X hours, etc.)

Format your response as: ACTION|URL|FREQUENCY
Example: "scrape_now|https://example.com|"
Example: "create_schedule|https://reddit.com|every 2 hours"`, message)
}

func (uc *ChatUseCase) callHuggingFace(prompt string, parentCtx context.Context) (string, error) {
	apiURL := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", uc.hfModelID)

	reqBody := hfRequest{
		Inputs: prompt,
		Parameters: map[string]interface{}{
			"max_new_tokens": 100,
			"temperature":    0.3,
			"do_sample":      false,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(parentCtx, 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if uc.hfAPIToken != "" {
		req.Header.Set("Authorization", "Bearer "+uc.hfAPIToken)
	}

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HuggingFace API error: %s", string(body))
	}

	var hfResp hfResponse
	if err := json.Unmarshal(body, &hfResp); err != nil {
		return "", err
	}

	if len(hfResp) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	return hfResp[0].GeneratedText, nil
}

func (uc *ChatUseCase) parseModelResponse(response string) (*entity.ChatIntent, error) {
	response = strings.TrimSpace(response)
	lines := strings.Split(response, "\n")
	lastLine := lines[len(lines)-1]

	parts := strings.Split(lastLine, "|")
	if len(parts) < 2 {
		return &entity.ChatIntent{
			Action:     "unknown",
			Confidence: 0.2,
		}, nil
	}

	intent := &entity.ChatIntent{
		Action:     strings.TrimSpace(parts[0]),
		Confidence: 0.8,
	}

	if len(parts) > 1 {
		intent.URL = strings.TrimSpace(parts[1])
	}

	if len(parts) > 2 && parts[2] != "" {
		freq := strings.TrimSpace(parts[2])
		intent.Frequency = freq

		cronExpr, _ := uc.extractCronFromText(freq)
		intent.CronExpr = cronExpr
	}

	return intent, nil
}

func (uc *ChatUseCase) GenerateResponse(intent *entity.ChatIntent) *entity.ChatResponse {
	response := &entity.ChatResponse{
		Intent: intent,
	}

	switch intent.Action {
	case "scrape_now":
		if intent.URL == "" {
			response.Message = "Entendido, pero ¿qué URL deseas scrapear?"
			response.NeedsConfirm = false
			response.Action = "ask_url"
		} else {
			response.Message = fmt.Sprintf("¿Deseas scrapear %s inmediatamente?", intent.URL)
			response.NeedsConfirm = true
			response.Action = "confirm_scrape"
		}

	case "create_schedule":
		if intent.URL == "" {
			response.Message = "¿Qué URL deseas programar?"
			response.NeedsConfirm = false
			response.Action = "ask_url"
		} else if intent.Frequency == "" {
			response.Message = fmt.Sprintf("¿Con qué frecuencia deseas scrapear %s?", intent.URL)
			response.NeedsConfirm = false
			response.Action = "ask_frequency"
		} else {
			response.Message = fmt.Sprintf(
				"Crearé un schedule para scrapear %s %s (expresión cron: %s). ¿Confirmas?",
				intent.URL, intent.Frequency, intent.CronExpr,
			)
			response.NeedsConfirm = true
			response.Action = "confirm_schedule"
		}

	default:
		response.Message = "No entendí tu solicitud. Puedes pedirme:\n• Scrapear una URL ahora\n• Programar scraping automático"
		response.NeedsConfirm = false
		response.Action = "none"
	}

	return response
}
