package handlers

import (
	"fmt"
	"net/http"
	"time"
	"webscraper-v2/internal/infrastructure/config"
	"webscraper-v2/internal/presentation/response"
)

type CommonHandler struct {
	config *config.Config
}

func NewCommonHandler(config *config.Config) *CommonHandler {
	return &CommonHandler{
		config: config,
	}
}

func (h *CommonHandler) Health(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "webscraper",
		"version":   "2.0",
	}
	response.SendSuccessResponse(w, "Service is healthy", health)
}

func (h *CommonHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, "Endpoint not found", http.StatusNotFound,
		fmt.Sprintf("The requested endpoint %s %s does not exist", r.Method, r.URL.Path))
}
