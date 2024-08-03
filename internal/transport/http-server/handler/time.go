package handler

import (
	"TaskSync/pkg/logger"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// Handler methods for Time

type timeTask struct {
	TaskID int       `json:"task_id"`
	Time   time.Time `json:"time"`
}

// @Summary Start Time Entry
// @Description Start recording time for a task
// @Tags Time
// @Accept json
// @Produce json
// @Param task body timeTask true "Task to start time entry for"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /time/start [post]
func (h *Handler) timeStartTimeEntry(w http.ResponseWriter, r *http.Request) {
	const op = "handler.timeStartTimeEntry"
	log := h.Logs.With(slog.String("operation", op))

	var task timeTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.services.Time.StartTimeEntry(r.Context(), task.TaskID, task.Time); err != nil {
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
// @Description End recording time for a task. FORMAT TIME - RFC 3339 "2024-08-01T08:00:00Z".
// @Tags Time
// @Accept json
// @Produce json
// @Param task body timeTask true "Task to end time entry for"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /time/end [post]
func (h *Handler) timeEndTimeEntry(w http.ResponseWriter, r *http.Request) {
	const op = "handler.timeEndTimeEntry"
	log := h.Logs.With(slog.String("operation", op))

	var task timeTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.services.Time.EndTimeEntry(r.Context(), task.TaskID, task.Time); err != nil {
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

type peopleTimeRange struct {
	PeopleID  int       `json:"people_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// @Summary Task Time Spent
// @Description Get time spent on tasks by a person within a specific time range. FORMAT TIME - RFC 3339 "2024-08-01T08:00:00Z".
// @Tags Time
// @Accept json
// @Produce json
// @Param task body peopleTimeRange true "People id and time range"
// @Success 200 {array} entities.TaskTimeSpent
// @Failure 500 {object} ErrorResponse
// @Router /time/spent [post]
func (h *Handler) TasksTimeSpent(w http.ResponseWriter, r *http.Request) {
	const op = "handler.timeGetTaskTimeSpent"
	log := h.Logs.With(slog.String("operation", op))

	var inputValues peopleTimeRange

	if err := json.NewDecoder(r.Body).Decode(&inputValues); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	timeSpent, err := h.services.Time.TasksTimeSpent(r.Context(), inputValues.PeopleID, inputValues.StartTime, inputValues.EndTime)
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
