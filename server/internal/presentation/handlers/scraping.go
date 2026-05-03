package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"webscraper-v2/internal/presentation/middleware"
	"webscraper-v2/internal/presentation/response"
	"webscraper-v2/internal/usecase"

	"github.com/gorilla/mux"
)

type ScrapingHandler struct {
	scrapingUseCase *usecase.ScrapingUseCase
	hub             *SSEHub
}

func NewScrapingHandler(scrapingUseCase *usecase.ScrapingUseCase, hub *SSEHub) *ScrapingHandler {
	return &ScrapingHandler{
		scrapingUseCase: scrapingUseCase,
		hub:             hub,
	}
}

func (h *ScrapingHandler) Scrape(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Scraping URL: %s", req.URL)
	result, err := h.scrapingUseCase.ScrapeURL(r.Context(), req.URL, user.ID)

	if err != nil {
		log.Printf("Error scraping URL %s: %v", req.URL, err)
		response.SendErrorResponse(w, "Failed to scrape URL", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Successfully scraped URL: %s (Status: %d, Words: %d)", req.URL, result.StatusCode, result.WordCount)
	response.SendSuccessResponse(w, "URL scraped successfully", result)
}

func (h *ScrapingHandler) PublicScrape(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Public scraping URL: %s", req.URL)
	result, err := h.scrapingUseCase.ScrapeURL(r.Context(), req.URL, 0)

	if err != nil {
		log.Printf("Error scraping URL %s: %v", req.URL, err)
		response.SendErrorResponse(w, "Failed to scrape URL", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Successfully scraped URL publicly: %s (Status: %d, Words: %d)", req.URL, result.StatusCode, result.WordCount)
	response.SendSuccessResponse(w, "URL scraped successfully", result)
}

func (h *ScrapingHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}

	page, perPage := 1, 50
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if pp, err := strconv.Atoi(r.URL.Query().Get("per_page")); err == nil && pp > 0 {
		perPage = pp
	}

	paginatedResults, err := h.scrapingUseCase.GetAllResultsPaginated(user.ID, page, perPage)
	if err != nil {
		log.Printf("Error getting results: %v", err)
		response.SendErrorResponse(w, "Failed to retrieve results", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Retrieved %d results (page %d/%d) for user %d",
		len(paginatedResults.Data),
		paginatedResults.Pagination.CurrentPage,
		paginatedResults.Pagination.TotalPages,
		user.ID)
	response.SendSuccessResponse(w, "Results retrieved successfully", paginatedResults)
}

func (h *ScrapingHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}

	id, err := parseID(r)
	if err != nil {
		response.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}

	result, err := h.scrapingUseCase.GetResult(id, user.ID)
	if err != nil {
		log.Printf("Error getting result %d by user %s: %v", id, user.Username, err)

		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			response.SendErrorResponse(w, "Result not found", http.StatusNotFound, "")
			return
		}

		response.SendErrorResponse(w, "Failed to retrieve result", http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Retrieved result ID: %d (%s) by user %s", id, result.URL, user.Username)
	response.SendSuccessResponse(w, "Result retrieved successfully", result)
}

func (h *ScrapingHandler) DeleteResult(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}

	id, err := parseID(r)
	if err != nil {
		response.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}

	if err := h.scrapingUseCase.DeleteResult(id, user.ID); err != nil {
		log.Printf("Error deleting result %d by user %s: %v", id, user.Username, err)

		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			response.SendErrorResponse(w, "Result not found", http.StatusNotFound, "")
			return
		}

		response.SendErrorResponse(w, "Failed to delete result", http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Deleted result ID: %d by user %s", id, user.Username)
	response.SendNoContent(w)
}


func (h *ScrapingHandler) StreamResults(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		response.SendErrorResponse(w, "Streaming not supported", http.StatusInternalServerError, "")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch, unsubscribe := h.hub.subscribe(user.ID)
	defer unsubscribe()

	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	keepalive := time.NewTicker(30 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-ch:
			fmt.Fprintf(w, "data: new-result\n\n")
			flusher.Flush()
		case <-keepalive.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
}
