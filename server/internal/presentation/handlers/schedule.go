package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/presentation/middleware"
	"webscraper-v2/internal/presentation/response"
	"webscraper-v2/internal/usecase"
)

type ScheduleHandler struct {
	scheduleUseCase *usecase.ScheduleUseCase
}

func NewScheduleHandler(scheduleUseCase *usecase.ScheduleUseCase) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleUseCase: scheduleUseCase,
	}
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	var req entity.CreateScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}

	schedule, err := h.scheduleUseCase.CreateSchedule(&req, user.ID)

	if err != nil {
		log.Printf("Error creating schedule: %v", err)
		response.SendErrorResponse(w, "Failed to create schedule", http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("Schedule created: %s (ID: %d) by user %s", schedule.Name, schedule.ID, user.Username)
	response.SendSuccessResponse(w, "Schedule created successfully", schedule)
}

func (h *ScheduleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	schedules, err := h.scheduleUseCase.GetSchedulesByUser(user.ID)

	if err != nil {
		log.Printf("Error getting schedules: %v", err)
		response.SendErrorResponse(w, "Failed to retrieve schedules", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Retrieved %d schedules for user %s", len(schedules), user.Username)
	response.SendSuccessResponse(w, fmt.Sprintf("Retrieved %d schedules", len(schedules)), schedules)
}

func (h *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	id, err := parseIDFromRequest(r)

	if err != nil {
		response.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}
	schedule, err := h.scheduleUseCase.GetSchedule(id)

	if err != nil {
		log.Printf("Error getting schedule %d: %v", id, err)
		response.SendErrorResponse(w, "Failed to retrieve schedule", http.StatusInternalServerError, err.Error())
		return
	}
	if schedule == nil {
		response.SendErrorResponse(w, "Schedule not found", http.StatusNotFound, fmt.Sprintf("No schedule found with ID %d", id))
		return
	}
	if schedule.UserID != user.ID {
		response.SendErrorResponse(w, "Schedule not found", http.StatusNotFound, fmt.Sprintf("No schedule found with ID %d", id))
		return
	}
	log.Printf("Retrieved schedule ID: %d (%s)", id, schedule.Name)
	response.SendSuccessResponse(w, "Schedule retrieved successfully", schedule)
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	id, err := parseIDFromRequest(r)

	if err != nil {
		response.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}
	var req entity.UpdateScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}
	schedule, err := h.scheduleUseCase.UpdateSchedule(id, &req, user.ID)

	if err != nil {
		log.Printf("Error updating schedule %d: %v", id, err)

		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			response.SendErrorResponse(w, "Schedule not found", http.StatusNotFound, "")
		} else {
			response.SendErrorResponse(w, "Failed to update schedule", http.StatusBadRequest, err.Error())
		}
		return
	}
	log.Printf("Updated schedule ID: %d (%s)", id, schedule.Name)
	response.SendSuccessResponse(w, "Schedule updated successfully", schedule)
}

func (h *ScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	id, err := parseIDFromRequest(r)

	if err != nil {
		response.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}
	err = h.scheduleUseCase.DeleteSchedule(id, user.ID)

	if err != nil {
		log.Printf("Error deleting schedule %d: %v", id, err)

		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			response.SendErrorResponse(w, "Schedule not found", http.StatusNotFound, "")
		} else {
			response.SendErrorResponse(w, "Failed to delete schedule", http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("Deleted schedule ID: %d by user %s", id, user.Username)
	response.SendNoContent(w)
}

func parseIDFromRequest(r *http.Request) (int64, error) {
	return parseID(r)
}
