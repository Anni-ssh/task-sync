package handler

import (
	_ "TaskSync/docs"
	"TaskSync/internal/service"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services *service.Service
	Logs     *slog.Logger
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitLogger(l *slog.Logger) {
	h.Logs = l
}

func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer) // Recovery из panic
	r.Use(middleware.CleanPath) // Исправление путей

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"}, // Разрешаем только запросы с этого домена
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "Content-Length", "Cache-Control",
			"Connection", "Host", "Origin"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(c.Handler)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// API people
	r.Route("/people", func(r chi.Router) {
		r.Get("/", h.peopleList)
		r.Post("/", h.peopleCreate)
		r.Get("/{peopleID}", h.peopleGetByID)
		r.Get("/filter", h.peopleGetByFilter)
		r.Put("/", h.peopleUpdate)
		r.Delete("/{peopleID}", h.peopleDelete)
	})

	// API task
	r.Route("/task", func(r chi.Router) {
		r.Get("/", h.taskList)
		r.Post("/", h.taskCreate)
		r.Get("/{taskID}", h.taskGetByID)
		r.Put("/", h.taskUpdate)
		r.Put("/update-people", h.taskUpdatePeople)
		r.Delete("/{taskID}", h.taskDelete)
	})

	// API time
	r.Route("/time", func(r chi.Router) {
		r.Post("/start", h.timeStartTimeEntry)
		r.Post("/end", h.timeEndTimeEntry)
		r.Get("/spent", h.timeGetTaskTimeSpent)
	})

	return r
}
