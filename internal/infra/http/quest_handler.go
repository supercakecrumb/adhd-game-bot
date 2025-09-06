package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

// QuestResponse represents the JSON response for a quest
type QuestResponse struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Category         string  `json:"category"`
	Difficulty       string  `json:"difficulty"`
	Mode             string  `json:"mode"`
	PointsAward      string  `json:"points_award"`
	RatePointsPerMin *string `json:"rate_points_per_min,omitempty"`
	MinMinutes       *int    `json:"min_minutes,omitempty"`
	MaxMinutes       *int    `json:"max_minutes,omitempty"`
	DailyPointsCap   *string `json:"daily_points_cap,omitempty"`
	CooldownSec      int     `json:"cooldown_sec"`
	StreakEnabled    bool    `json:"streak_enabled"`
	Status           string  `json:"status"`
}

// CreateQuestRequest represents the JSON request for creating a quest
type CreateQuestRequest struct {
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Category         string  `json:"category"`
	Difficulty       string  `json:"difficulty"`
	Mode             string  `json:"mode"`
	PointsAward      string  `json:"points_award"`
	RatePointsPerMin *string `json:"rate_points_per_min,omitempty"`
	MinMinutes       *int    `json:"min_minutes,omitempty"`
	MaxMinutes       *int    `json:"max_minutes,omitempty"`
	DailyPointsCap   *string `json:"daily_points_cap,omitempty"`
	CooldownSec      *int    `json:"cooldown_sec,omitempty"`
	StreakEnabled    *bool   `json:"streak_enabled,omitempty"`
	Status           *string `json:"status,omitempty"`
}

// CompleteQuestRequest represents the JSON request for completing a quest
type CompleteQuestRequest struct {
	IdempotencyKey  string   `json:"idempotency_key"`
	CompletionRatio *float64 `json:"completion_ratio,omitempty"`
	Minutes         *int     `json:"minutes,omitempty"`
}

// CompleteQuestResponse represents the JSON response for completing a quest
type CompleteQuestResponse struct {
	AwardedPoints string `json:"awarded_points"`
	SubmittedAt   string `json:"submitted_at"`
	StreakCount   *int   `json:"streak_count,omitempty"`
}

func (s *Server) createQuestHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateQuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get dungeon ID from URL parameter
	dungeonID := chi.URLParam(r, "dungeonId")
	if dungeonID == "" {
		http.Error(w, "dungeonId URL parameter is required", http.StatusBadRequest)
		return
	}

	// Get user ID from query parameter or request context
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Convert string values to Decimal
	pointsAward := valueobject.NewDecimal(req.PointsAward)

	var ratePointsPerMin *valueobject.Decimal
	if req.RatePointsPerMin != nil {
		decimal := valueobject.NewDecimal(*req.RatePointsPerMin)
		ratePointsPerMin = &decimal
	}

	var dailyPointsCap *valueobject.Decimal
	if req.DailyPointsCap != nil {
		decimal := valueobject.NewDecimal(*req.DailyPointsCap)
		dailyPointsCap = &decimal
	}

	// Create input for service
	input := usecase.CreateQuestInput{
		Title:            req.Title,
		Description:      req.Description,
		Category:         req.Category,
		Difficulty:       req.Difficulty,
		Mode:             req.Mode,
		PointsAward:      pointsAward,
		RatePointsPerMin: ratePointsPerMin,
		MinMinutes:       req.MinMinutes,
		MaxMinutes:       req.MaxMinutes,
		DailyPointsCap:   dailyPointsCap,
		CooldownSec:      0,
		StreakEnabled:    true,
		Status:           "active",
		TimeZone:         "UTC",
	}

	if req.CooldownSec != nil {
		input.CooldownSec = *req.CooldownSec
	}

	if req.StreakEnabled != nil {
		input.StreakEnabled = *req.StreakEnabled
	}

	if req.Status != nil {
		input.Status = *req.Status
	}

	// Call the use case
	createdQuest, err := s.QuestService.CreateQuest(r.Context(), userID, dungeonID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := s.questToResponse(createdQuest)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) listQuestsHandler(w http.ResponseWriter, r *http.Request) {
	// Get dungeon ID from URL parameter
	dungeonID := chi.URLParam(r, "dungeonId")
	if dungeonID == "" {
		http.Error(w, "dungeonId URL parameter is required", http.StatusBadRequest)
		return
	}

	// Get user ID from query parameter or request context
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	quests, err := s.QuestService.ListQuests(r.Context(), userID, dungeonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]QuestResponse, len(quests))
	for i, quest := range quests {
		response[i] = s.questToResponse(quest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) completeQuestHandler(w http.ResponseWriter, r *http.Request) {
	questID := chi.URLParam(r, "questId")
	if questID == "" {
		http.Error(w, "questId is required", http.StatusBadRequest)
		return
	}

	var req CompleteQuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get user ID from query parameter or request context
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Call the use case
	input := usecase.CompleteQuestInput{
		IdempotencyKey:  req.IdempotencyKey,
		CompletionRatio: req.CompletionRatio,
		Minutes:         req.Minutes,
	}

	err = s.QuestService.CompleteQuest(r.Context(), userID, questID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Quest completed"})
}

// Helper method to convert entity.Quest to QuestResponse
func (s *Server) questToResponse(quest *entity.Quest) QuestResponse {
	var ratePointsPerMin, dailyPointsCap *string

	if quest.RatePointsPerMin != nil {
		str := quest.RatePointsPerMin.String()
		ratePointsPerMin = &str
	}

	if quest.DailyPointsCap != nil {
		str := quest.DailyPointsCap.String()
		dailyPointsCap = &str
	}

	return QuestResponse{
		ID:               quest.ID,
		Title:            quest.Title,
		Description:      quest.Description,
		Category:         quest.Category,
		Difficulty:       quest.Difficulty,
		Mode:             quest.Mode,
		PointsAward:      quest.PointsAward.String(),
		RatePointsPerMin: ratePointsPerMin,
		MinMinutes:       quest.MinMinutes,
		MaxMinutes:       quest.MaxMinutes,
		DailyPointsCap:   dailyPointsCap,
		CooldownSec:      quest.CooldownSec,
		StreakEnabled:    quest.StreakEnabled,
		Status:           quest.Status,
	}
}
