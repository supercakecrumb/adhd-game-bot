package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

// CreateDungeonRequest represents the JSON request for creating a dungeon
type CreateDungeonRequest struct {
	Title          string `json:"title"`
	TelegramChatID *int64 `json:"telegram_chat_id,omitempty"`
}

// CreateDungeonResponse represents the JSON response for creating a dungeon
type CreateDungeonResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	AdminUserID    int64  `json:"admin_user_id"`
	TelegramChatID *int64 `json:"telegram_chat_id,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// AddMemberRequest represents the JSON request for adding a member to a dungeon
type AddMemberRequest struct {
	UserID int64 `json:"user_id"`
}

func (s *Server) createDungeonHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateDungeonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get admin user ID from query parameter or request context
	adminUserIDStr := r.URL.Query().Get("admin_user_id")
	if adminUserIDStr == "" {
		http.Error(w, "admin_user_id query parameter is required", http.StatusBadRequest)
		return
	}

	adminUserID, err := strconv.ParseInt(adminUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid admin_user_id", http.StatusBadRequest)
		return
	}

	// Call the use case
	createdDungeon, err := s.DungeonService.CreateDungeon(r.Context(), adminUserID, req.Title, req.TelegramChatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := s.dungeonToResponse(createdDungeon)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) addMemberHandler(w http.ResponseWriter, r *http.Request) {
	// Get dungeon ID from URL parameter
	dungeonID := chi.URLParam(r, "dungeonId")
	if dungeonID == "" {
		http.Error(w, "dungeonId URL parameter is required", http.StatusBadRequest)
		return
	}

	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get admin user ID from query parameter or request context
	adminUserIDStr := r.URL.Query().Get("admin_user_id")
	if adminUserIDStr == "" {
		http.Error(w, "admin_user_id query parameter is required", http.StatusBadRequest)
		return
	}

	adminUserID, err := strconv.ParseInt(adminUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid admin_user_id", http.StatusBadRequest)
		return
	}

	// Call the use case
	err = s.DungeonService.AddMember(r.Context(), adminUserID, dungeonID, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Member added to dungeon"})
}

func (s *Server) listMembersHandler(w http.ResponseWriter, r *http.Request) {
	// Get dungeon ID from URL parameter
	dungeonID := chi.URLParam(r, "dungeonId")
	if dungeonID == "" {
		http.Error(w, "dungeonId URL parameter is required", http.StatusBadRequest)
		return
	}

	// Get admin user ID from query parameter or request context
	adminUserIDStr := r.URL.Query().Get("admin_user_id")
	if adminUserIDStr == "" {
		http.Error(w, "admin_user_id query parameter is required", http.StatusBadRequest)
		return
	}

	adminUserID, err := strconv.ParseInt(adminUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid admin_user_id", http.StatusBadRequest)
		return
	}

	// Call the use case
	members, err := s.DungeonService.ListMembers(r.Context(), adminUserID, dungeonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]int64{"members": members})
}

// Helper method to convert entity.Dungeon to CreateDungeonResponse
func (s *Server) dungeonToResponse(dungeon *entity.Dungeon) CreateDungeonResponse {
	return CreateDungeonResponse{
		ID:             dungeon.ID,
		Title:          dungeon.Title,
		AdminUserID:    dungeon.AdminUserID,
		TelegramChatID: dungeon.TelegramChatID,
		CreatedAt:      dungeon.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
