package handler

import (
	"TaskSync/internal/entities"
	"TaskSync/pkg/logger"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Handler methods for People

// @Summary Create People
// @Description Create a new people
// @Tags People
// @Accept json
// @Produce json
// @Param people body entities.People true "People to create"
// @Success 201 {integer} int
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /people [post]
func (h *Handler) peopleCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.peopleCreate"
	log := h.Logs.With(slog.String("operation", op))

	var people entities.People

	if err := json.NewDecoder(r.Body).Decode(&people); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	id, err := h.services.People.Create(r.Context(), people)
	if err != nil {
		log.Error("Failed to create person", logger.Err(err))
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Failed to create person")
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		log.Error("Failed to encode response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

// @Summary List People
// @Description Get all people
// @Tags People
// @Accept json
// @Produce json
// @Success 200 {array} entities.People
// @Failure 500 {object} ErrorResponse
// @Router /people [get]
func (h *Handler) peopleList(w http.ResponseWriter, r *http.Request) {
	const op = "handler.peopleList"
	log := h.Logs.With(slog.String("operation", op))

	people, err := h.services.People.List(r.Context())
	if err != nil {
		log.Error("Failed to fetch people list", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch people list")
		return
	}

	if err := json.NewEncoder(w).Encode(people); err != nil {
		log.Error("Failed to encode response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

// @Summary Get People by ID
// @Description Get details of a people by ID
// @Tags People
// @Accept json
// @Produce json
// @Param peopleID path int true "People ID"
// @Success 200 {object} entities.People
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /people/{peopleID} [get]
func (h *Handler) peopleGetByID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.peopleGetByID"
	log := h.Logs.With(slog.String("operation", op))

	peopleID := chi.URLParam(r, "peopleID")
	id, err := strconv.Atoi(peopleID)
	if err != nil {
		log.Error("Invalid people ID", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid people ID")
		return
	}

	people, err := h.services.People.GetByID(r.Context(), id)
	if err != nil {
		log.Error("Failed to fetch person by ID", logger.Err(err))
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Failed to fetch person by ID")
		return
	}

	if err := json.NewEncoder(w).Encode(people); err != nil {
		log.Error("Failed to encode response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

// @Summary Get People by Filter
// @Description Get people based on filters
// @Tags People
// @Accept json
// @Produce json
// @Param id query int false "Person ID"
// @Param passport_series query int false "Passport Series"
// @Param passport_number query int false "Passport Number"
// @Param surname query string false "Surname"
// @Param name query string false "Name"
// @Param patronymic query string false "Patronymic"
// @Param address query string false "Address"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} entities.People
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /people/filter [get]
func (h *Handler) peopleGetByFilter(w http.ResponseWriter, r *http.Request) {
	const op = "handler.peopleGetByFilter"
	log := h.Logs.With(slog.String("operation", op))

	filter := entities.People{
		ID:             parseQueryInt(r.URL.Query().Get("id")),
		PassportSeries: parseQueryInt(r.URL.Query().Get("passport_series")),
		PassportNumber: parseQueryInt(r.URL.Query().Get("passport_number")),
		Surname:        r.URL.Query().Get("surname"),
		Name:           r.URL.Query().Get("name"),
		Patronymic:     r.URL.Query().Get("patronymic"),
		Address:        r.URL.Query().Get("address"),
	}

	limit := parseQueryInt(r.URL.Query().Get("limit"))
	offset := parseQueryInt(r.URL.Query().Get("offset"))

	fmt.Println(filter)
	fmt.Println(limit)
	fmt.Println(offset)

	people, err := h.services.People.GetByFilter(r.Context(), filter, limit, offset)
	if err != nil {
		log.Error("Failed to fetch people by filter", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch people by filter")
		return
	}

	if err := json.NewEncoder(w).Encode(people); err != nil {
		log.Error("Failed to encode response", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

// @Summary Update People
// @Description Update an existing person
// @Tags People
// @Accept json
// @Produce json
// @Param person body entities.People true "Person to update"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /people [put]
func (h *Handler) peopleUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.peopleUpdate"
	log := h.Logs.With(slog.String("operation", op))

	var people entities.People
	if err := json.NewDecoder(r.Body).Decode(&people); err != nil {
		log.Error("Failed to decode request body", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.services.People.Update(r.Context(), people); err != nil {
		log.Error("Failed to update person", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to update person")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete people
// @Description Delete a people by ID
// @Tags People
// @Accept json
// @Produce json
// @Param peopleID path int true "People ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /people/{peopleID} [delete]
func (h *Handler) peopleDelete(w http.ResponseWriter, r *http.Request) {
	const op = "handler.peopleDelete"
	log := h.Logs.With(slog.String("operation", op))

	peopleID := chi.URLParam(r, "peopleID")
	id, err := strconv.Atoi(peopleID)
	if err != nil {
		log.Error("Invalid people ID", logger.Err(err))
		writeErrorResponse(w, http.StatusBadRequest, "Invalid people ID")
		return
	}

	if err := h.services.People.Delete(r.Context(), id); err != nil {
		log.Error("Failed to delete person", logger.Err(err))
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete person")
		return
	}

	w.WriteHeader(http.StatusOK)
}
