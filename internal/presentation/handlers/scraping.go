package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"webscraper-v2/internal/presentation/middleware"
	"webscraper-v2/internal/presentation/response"
	"webscraper-v2/internal/usecase"

	"github.com/gorilla/mux"
)

type ScrapingHandler struct {
	scrapingUseCase *usecase.ScrapingUseCase
}

func NewScrapingHandler(scrapingUseCase *usecase.ScrapingUseCase) *ScrapingHandler {
	return &ScrapingHandler{
		scrapingUseCase: scrapingUseCase,
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
	if req.URL = strings.TrimSpace(req.URL); req.URL == "" {
		response.SendErrorResponse(w, "URL is required", http.StatusBadRequest, "")
		return
	}
	parsedURL, err := url.ParseRequestURI(req.URL)

	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		response.SendErrorResponse(w, "Invalid URL format", http.StatusBadRequest, "URL must include protocol (http:// or https://) and valid domain")
		return
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		response.SendErrorResponse(w, "Invalid URL scheme", http.StatusBadRequest, "Only HTTP and HTTPS protocols are supported")
		return
	}
	log.Printf("Scraping URL: %s", req.URL)
	result, err := h.scrapingUseCase.ScrapeURL(req.URL, user.ID)

	if err != nil {
		log.Printf("Error scraping URL %s: %v", req.URL, err)
		response.SendErrorResponse(w, "Failed to scrape URL", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Successfully scraped URL: %s (Status: %d, Words: %d)", req.URL, result.StatusCode, result.WordCount)
	response.SendSuccessResponse(w, "URL scraped successfully", result)
}

func (h *ScrapingHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	results, err := h.scrapingUseCase.GetAllResults(user.ID)

	if err != nil {
		log.Printf("Error getting results: %v", err)
		response.SendErrorResponse(w, "Failed to retrieve results", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Retrieved %d scraping results", len(results))
	response.SendSuccessResponse(w, fmt.Sprintf("Retrieved %d results", len(results)), results)
}

func (h *ScrapingHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)

	if err != nil {
		response.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}
	result, err := h.scrapingUseCase.GetResult(id)

	if err != nil {
		log.Printf("Error getting result %d: %v", id, err)
		response.SendErrorResponse(w, "Failed to retrieve result", http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		response.SendErrorResponse(w, "Result not found", http.StatusNotFound, fmt.Sprintf("No result found with ID %d", id))
		return
	}
	log.Printf("Retrieved result ID: %d (%s)", id, result.URL)
	response.SendSuccessResponse(w, "Result retrieved successfully", result)
}

func (h *ScrapingHandler) DeleteResult(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)

	if err != nil {
		response.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}
	result, err := h.scrapingUseCase.GetResult(id)

	if err != nil {
		log.Printf("Error checking result %d: %v", id, err)
		response.SendErrorResponse(w, "Failed to check result", http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		response.SendErrorResponse(w, "Result not found", http.StatusNotFound, fmt.Sprintf("No result found with ID %d", id))
		return
	}
	if err := h.scrapingUseCase.DeleteResult(id); err != nil {
		log.Printf("Error deleting result %d: %v", id, err)
		response.SendErrorResponse(w, "Failed to delete result", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Deleted result ID: %d (%s)", id, result.URL)
	response.SendNoContent(w)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
}
