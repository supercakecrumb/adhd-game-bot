package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

type Server struct {
	Router         *chi.Mux
	QuestService   *usecase.QuestService
	DungeonService *usecase.DungeonService
}

func NewServer(questService *usecase.QuestService, dungeonService *usecase.DungeonService) *Server {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	server := &Server{
		Router:         r,
		QuestService:   questService,
		DungeonService: dungeonService,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.Router.Route("/api/v1", func(r chi.Router) {
		// Quest routes
		r.Route("/dungeons/{dungeonId}/quests", func(r chi.Router) {
			r.Get("/", s.listQuestsHandler)
		})

		r.Route("/quests/{questId}", func(r chi.Router) {
			r.Post("/complete", s.completeQuestHandler)
		})

		// Dungeon routes
		r.Route("/dungeons", func(r chi.Router) {
			r.Post("/", s.createDungeonHandler)
			r.Route("/{dungeonId}", func(r chi.Router) {
				r.Post("/quests", s.createQuestHandler)
				r.Post("/members", s.addMemberHandler)
				r.Get("/members", s.listMembersHandler)
			})
		})
	})
}
