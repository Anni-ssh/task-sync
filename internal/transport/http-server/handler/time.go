package handler

import (
	"TaskSync/internal/entities"
	"TaskSync/pkg/logger"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// Handler methods for Time

// @Summary Start Time Entry
// @Description Start recording time for a task
// @Tags Time
// @Accept json
// @Produce json
// @Param task body entities.Task true "Task to start time entry for"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /time/start [post]
func (h *Handler) timeStartTimeEntry(w http.ResponseWriter, r *http.Request) {
	const op = "handler.timeStartTimeEntry"
	log := h.Logs.With(slog.String("operation", op))

	var task entities.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.services.Time.StartTimeEntry(r.Context(), task); err != nil {
		log.Error("Failed to start time entry", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to start time entry")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error("Failed to write response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to write response")
	}
}

// @Summary End Time Entry
// @Description End recording time for a task
// @Tags Time
// @Accept json
// @Produce json
// @Param task body entities.Task true "Task to end time entry for"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /time/end [post]
func (h *Handler) timeEndTimeEntry(w http.ResponseWriter, r *http.Request) {
	const op = "handler.timeEndTimeEntry"
	log := h.Logs.With(slog.String("operation", op))

	var task entities.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.services.Time.EndTimeEntry(r.Context(), task); err != nil {
		log.Error("Failed to end time entry", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to end time entry")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error("Failed to write response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to write response")
	}
}

// @Summary Get Task Time Spent
// @Description Get time spent on tasks by a person within a specific time range
// @Tags Time
// @Accept json
// @Produce json
// @Param people_id query int true "People ID"
// @Param start_time query string true "Start time in RFC3339 format"
// @Param end_time query string true "End time in RFC3339 format"
// @Success 200 {array} entities.TaskTimeSpent
// @Failure 500 {object} ErrorResponse
// @Router /time/spent [get]
func (h *Handler) timeGetTaskTimeSpent(w http.ResponseWriter, r *http.Request) {
	const op = "handler.timeGetTaskTimeSpent"
	log := h.Logs.With(slog.String("operation", op))

	peopleID, _ := strconv.Atoi(r.URL.Query().Get("people_id"))
	startTime, _ := time.Parse(time.RFC3339, r.URL.Query().Get("start_time"))
	endTime, _ := time.Parse(time.RFC3339, r.URL.Query().Get("end_time"))

	timeSpent, err := h.services.Time.GetTaskTimeSpent(r.Context(), peopleID, startTime, endTime)
	if err != nil {
		log.Error("Failed to get task time spent", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get task time spent")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(timeSpent); err != nil {
		log.Error("Failed to encode response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}
