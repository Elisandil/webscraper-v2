package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	if r.URL.Query().Get("page") != "" || r.URL.Query().Get("per_page") != "" {
		h.GetResultsPaginated(w, r)
		return
	}
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
	if result.UserID != user.ID {
		response.SendErrorResponse(w, "Result not found", http.StatusNotFound, fmt.Sprintf("No result found with ID %d", id))
		return
	}
	log.Printf("Retrieved result ID: %d (%s)", id, result.URL)

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
	if result.UserID != user.ID {
		response.SendErrorResponse(w, "Result not found", http.StatusNotFound, fmt.Sprintf("No result found with ID %d", id))
		return
	}
	if err := h.scrapingUseCase.DeleteResult(id); err != nil {
		log.Printf("Error deleting result %d: %v", id, err)
		response.SendErrorResponse(w, "Failed to delete result", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Deleted result ID: %d (%s) by user %s", id, result.URL, user.Username)

	response.SendNoContent(w)
}

// Handler for paginated results
//-------------------------------------------------------------------------------------------------------

func (h *ScrapingHandler) GetResultsPaginated(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")
	page := 1
	perPage := 10

	if pageStr != "" {

		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if perPageStr != "" {

		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
			perPage = pp
		}
	}
	paginatedResults, err := h.scrapingUseCase.GetAllResultsPaginated(user.ID, page, perPage)

	if err != nil {
		log.Printf("Error getting paginated results: %v", err)
		response.SendErrorResponse(w, "Failed to retrieve results", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Retrieved %d results (page %d of %d) for user %d",
		len(paginatedResults.Data),
		paginatedResults.Pagination.CurrentPage,
		paginatedResults.Pagination.TotalPages,
		user.ID)

	response.SendSuccessResponse(w, "Results retrieved successfully", paginatedResults)
}

//--------------------------------------------------------------------------------------------------------

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
}
