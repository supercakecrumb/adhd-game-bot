package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

type Server struct {
	Router      *chi.Mux
	TaskService *usecase.TaskService
}

func NewServer(taskService *usecase.TaskService) *Server {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	server := &Server{
		Router:      r,
		TaskService: taskService,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.Router.Route("/api", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Post("/", s.createTaskHandler)
			r.Get("/{taskID}", s.getTaskHandler)
			r.Put("/{taskID}", s.updateTaskHandler)
			r.Post("/{taskID}/complete", s.completeTaskHandler)
		})

		r.Route("/users/{userID}/tasks", func(r chi.Router) {
			r.Get("/", s.listTasksByUserHandler)
		})
	})
}
