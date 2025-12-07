package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/presentation/middleware"
	"webscraper-v2/internal/presentation/response"
	"webscraper-v2/internal/usecase"
)

type ChatHandler struct {
	chatUseCase     *usecase.ChatUseCase
	scrapingUseCase *usecase.ScrapingUseCase
	scheduleUseCase *usecase.ScheduleUseCase
}

func NewChatHandler(
	chatUseCase *usecase.ChatUseCase,
	scrapingUseCase *usecase.ScrapingUseCase,
	scheduleUseCase *usecase.ScheduleUseCase,
) *ChatHandler {
	return &ChatHandler{
		chatUseCase:     chatUseCase,
		scrapingUseCase: scrapingUseCase,
		scheduleUseCase: scheduleUseCase,
	}
}

func (h *ChatHandler) ParseMessage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}

	var req entity.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}

	if req.Message == "" {
		response.SendErrorResponse(w, "Message is required", http.StatusBadRequest, "")
		return
	}

	log.Printf("Chat message from user %d: %s", user.ID, req.Message)

	intent, err := h.chatUseCase.InterpretMessage(req.Message)
	if err != nil {
		log.Printf("Error interpreting message: %v", err)
		response.SendErrorResponse(w, "Error processing message", http.StatusInternalServerError, err.Error())
		return
	}

	chatResponse := h.chatUseCase.GenerateResponse(intent)

	log.Printf("Chat intent detected: action=%s, url=%s, confidence=%.2f",
		intent.Action, intent.URL, intent.Confidence)

	response.SendSuccessResponse(w, "Message processed", chatResponse)
}

func (h *ChatHandler) ExecuteAction(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}

	var req struct {
		Intent entity.ChatIntent `json:"intent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Executing action: %s for user %d", req.Intent.Action, user.ID)

	switch req.Intent.Action {
	case "scrape_now":
		h.executeScrapeNow(w, user.ID, req.Intent)
	case "create_schedule":
		h.executeCreateSchedule(w, user.ID, req.Intent)
	default:
		response.SendErrorResponse(w, "Unknown action", http.StatusBadRequest, "")
	}
}

func (h *ChatHandler) executeScrapeNow(w http.ResponseWriter, userID int64, intent entity.ChatIntent) {
	if intent.URL == "" {
		response.SendErrorResponse(w, "URL is required", http.StatusBadRequest, "")
		return
	}

	result, err := h.scrapingUseCase.ScrapeURL(intent.URL, userID)
	if err != nil {
		log.Printf("Error scraping URL %s: %v", intent.URL, err)
		response.SendErrorResponse(w, "Failed to scrape URL", http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully scraped URL via chat: %s", intent.URL)
	response.SendSuccessResponse(w, "URL scraped successfully", map[string]interface{}{
		"result":  result,
		"message": "✅ Scraping completado exitosamente",
	})
}

func (h *ChatHandler) executeCreateSchedule(w http.ResponseWriter, userID int64, intent entity.ChatIntent) {
	if intent.URL == "" || intent.CronExpr == "" {
		response.SendErrorResponse(w, "URL and cron expression are required", http.StatusBadRequest, "")
		return
	}

	scheduleReq := &entity.CreateScheduleRequest{
		Name:     "Chat Schedule - " + intent.URL,
		URL:      intent.URL,
		CronExpr: intent.CronExpr,
	}

	schedule, err := h.scheduleUseCase.CreateSchedule(scheduleReq, userID)
	if err != nil {
		log.Printf("Error creating schedule: %v", err)
		response.SendErrorResponse(w, "Failed to create schedule", http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully created schedule via chat: %s (%s)", intent.URL, intent.CronExpr)
	response.SendSuccessResponse(w, "Schedule created successfully", map[string]interface{}{
		"schedule":  schedule,
		"message":   "✅ Schedule creado exitosamente",
		"frequency": intent.Frequency,
	})
}
