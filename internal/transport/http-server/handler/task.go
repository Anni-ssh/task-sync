package handler

import (
	"TaskSync/internal/entities"
	"TaskSync/pkg/logger"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Handler methods for Task

// @Summary Create Task
// @Description Create a new task
// @Tags Task
// @Accept json
// @Produce json
// @Param people_id query int true "People ID"
// @Param task body entities.Task true "Task to create"
// @Success 200 {integer} int "Task ID"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task [post]
func (h *Handler) taskCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.taskCreate"
	log := h.Logs.With(slog.String("operation", op))

	var task entities.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	peopleID, err := strconv.Atoi(r.URL.Query().Get("people_id"))
	if err != nil {
		log.Error("Invalid people ID", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid people ID")
		return
	}

	id, err := h.services.Task.Create(r.Context(), peopleID, task)
	if err != nil {
		log.Error("Failed to create task", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
}

// @Summary Get Task by ID
// @Description Get a task by its ID
// @Tags Task
// @Accept json
// @Produce json
// @Param taskID path int true "Task ID"
// @Success 200 {object} entities.Task
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task/{taskID} [get]
func (h *Handler) taskGetByID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.taskGetByID"
	log := h.Logs.With(slog.String("operation", op))

	taskID := chi.URLParam(r, "taskID")
	id, err := strconv.Atoi(taskID)
	if err != nil {
		log.Error("Invalid task ID", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.services.Task.GetByID(r.Context(), id)
	if err != nil {
		log.Error("Failed to get task by ID", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get task by ID")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// @Summary List Tasks
// @Description Get list of all tasks
// @Tags Task
// @Accept json
// @Produce json
// @Success 200 {array} entities.Task
// @Failure 500 {object} ErrorResponse
// @Router /task [get]
func (h *Handler) taskList(w http.ResponseWriter, r *http.Request) {
	const op = "handler.taskList"
	log := h.Logs.With(slog.String("operation", op))

	tasks, err := h.services.Task.List(r.Context())
	if err != nil {
		log.Error("Failed to list tasks", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to list tasks")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Error("Failed to encode response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

// @Summary Update Task
// @Description Update an existing task
// @Tags Task
// @Accept json
// @Produce json
// @Param task body entities.Task true "Task to update"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task [put]
func (h *Handler) taskUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.taskUpdate"
	log := h.Logs.With(slog.String("operation", op))

	var task entities.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.services.Task.Update(r.Context(), task); err != nil {
		log.Error("Failed to update task", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error("Failed to write response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to write response")
	}
}

// @Summary Update People in Task
// @Description Update people associated with a task
// @Tags Task
// @Accept json
// @Produce json
// @Param task_id query int true "Task ID"
// @Param people_id query int true "People ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task/update-people [put]
func (h *Handler) taskUpdatePeople(w http.ResponseWriter, r *http.Request) {
	const op = "handler.taskUpdatePeople"
	log := h.Logs.With(slog.String("operation", op))

	taskID, err := strconv.Atoi(r.URL.Query().Get("task_id"))
	if err != nil {
		log.Error("Invalid task ID", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	peopleID, err := strconv.Atoi(r.URL.Query().Get("people_id"))
	if err != nil {
		log.Error("Invalid people ID", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid people ID")
		return
	}

	if err := h.services.Task.UpdatePeople(r.Context(), peopleID, taskID); err != nil {
		log.Error("Failed to update people in task", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to update people in task")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error("Failed to write response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to write response")
	}
}

// @Summary Delete Task
// @Description Delete a task by its ID
// @Tags Task
// @Accept json
// @Produce json
// @Param taskID path int true "Task ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task/{taskID} [delete]
func (h *Handler) taskDelete(w http.ResponseWriter, r *http.Request) {
	const op = "handler.taskDelete"
	log := h.Logs.With(slog.String("operation", op))

	taskID := chi.URLParam(r, "taskID")
	id, err := strconv.Atoi(taskID)
	if err != nil {
		log.Error("Invalid task ID", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if err := h.services.Task.Delete(r.Context(), id); err != nil {
		log.Error("Failed to delete task", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Error("Failed to write response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to write response")
	}
}
