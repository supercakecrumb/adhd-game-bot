package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

// TaskResponse represents the JSON response for a task
type TaskResponse struct {
	ID              string         `json:"id"`
	ChatID          int64          `json:"chat_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Category        string         `json:"category"`
	Difficulty      string         `json:"difficulty"`
	ScheduleJSON    string         `json:"schedule_json"`
	BaseDuration    int            `json:"base_duration"`
	GracePeriod     int            `json:"grace_period"`
	Cooldown        int            `json:"cooldown"`
	RewardCurveJSON string         `json:"reward_curve_json"`
	PartialCredit   *entity.Reward `json:"partial_credit,omitempty"`
	StreakEnabled   bool           `json:"streak_enabled"`
	Status          string         `json:"status"`
	LastCompletedAt *string        `json:"last_completed_at,omitempty"`
	StreakCount     int            `json:"streak_count"`
	TimeZone        string         `json:"time_zone"`
}

// TaskCreateRequest represents the JSON request for creating a task
type TaskCreateRequest struct {
	ChatID          int64          `json:"chat_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Category        string         `json:"category"`
	Difficulty      string         `json:"difficulty"`
	ScheduleJSON    string         `json:"schedule_json"`
	BaseDuration    int            `json:"base_duration"`
	GracePeriod     int            `json:"grace_period"`
	Cooldown        int            `json:"cooldown"`
	RewardCurveJSON string         `json:"reward_curve_json"`
	PartialCredit   *entity.Reward `json:"partial_credit,omitempty"`
	StreakEnabled   bool           `json:"streak_enabled"`
	Status          string         `json:"status"`
	TimeZone        string         `json:"time_zone"`
}

// TaskUpdateRequest represents the JSON request for updating a task
type TaskUpdateRequest struct {
	Title           string         `json:"title,omitempty"`
	Description     string         `json:"description,omitempty"`
	Category        string         `json:"category,omitempty"`
	Difficulty      string         `json:"difficulty,omitempty"`
	ScheduleJSON    string         `json:"schedule_json,omitempty"`
	BaseDuration    *int           `json:"base_duration,omitempty"`
	GracePeriod     *int           `json:"grace_period,omitempty"`
	Cooldown        *int           `json:"cooldown,omitempty"`
	RewardCurveJSON string         `json:"reward_curve_json,omitempty"`
	PartialCredit   *entity.Reward `json:"partial_credit,omitempty"`
	StreakEnabled   *bool          `json:"streak_enabled,omitempty"`
	Status          string         `json:"status,omitempty"`
	TimeZone        string         `json:"time_zone,omitempty"`
}

func (s *Server) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req TaskCreateRequest
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

	// Create task entity
	task := &entity.Task{
		ChatID:          req.ChatID,
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		Difficulty:      req.Difficulty,
		ScheduleJSON:    req.ScheduleJSON,
		BaseDuration:    req.BaseDuration,
		GracePeriod:     req.GracePeriod,
		Cooldown:        req.Cooldown,
		RewardCurveJSON: req.RewardCurveJSON,
		PartialCredit:   req.PartialCredit,
		StreakEnabled:   req.StreakEnabled,
		Status:          req.Status,
		TimeZone:        req.TimeZone,
	}

	// Call the use case
	createdTask, err := s.TaskService.CreateTask(r.Context(), userID, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := s.taskToResponse(createdTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		http.Error(w, "taskID is required", http.StatusBadRequest)
		return
	}

	task, err := s.TaskService.GetTask(r.Context(), taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := s.taskToResponse(task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		http.Error(w, "taskID is required", http.StatusBadRequest)
		return
	}

	var req TaskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get existing task
	task, err := s.TaskService.GetTask(r.Context(), taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Category != "" {
		task.Category = req.Category
	}
	if req.Difficulty != "" {
		task.Difficulty = req.Difficulty
	}
	if req.ScheduleJSON != "" {
		task.ScheduleJSON = req.ScheduleJSON
	}
	if req.BaseDuration != nil {
		task.BaseDuration = *req.BaseDuration
	}
	if req.GracePeriod != nil {
		task.GracePeriod = *req.GracePeriod
	}
	if req.Cooldown != nil {
		task.Cooldown = *req.Cooldown
	}
	if req.RewardCurveJSON != "" {
		task.RewardCurveJSON = req.RewardCurveJSON
	}
	if req.PartialCredit != nil {
		task.PartialCredit = req.PartialCredit
	}
	if req.StreakEnabled != nil {
		task.StreakEnabled = *req.StreakEnabled
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.TimeZone != "" {
		task.TimeZone = req.TimeZone
	}

	// Call the use case
	err = s.TaskService.UpdateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := s.taskToResponse(task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		http.Error(w, "taskID is required", http.StatusBadRequest)
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
	err = s.TaskService.CompleteTask(r.Context(), userID, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Task completed"})
}

func (s *Server) listTasksByUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	if userIDStr == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	tasks, err := s.TaskService.ListTasksByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = s.taskToResponse(task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper method to convert entity.Task to TaskResponse
func (s *Server) taskToResponse(task *entity.Task) TaskResponse {
	var lastCompletedAt *string
	if task.LastCompletedAt != nil {
		t := task.LastCompletedAt.Format("2006-01-02T15:04:05Z07:00")
		lastCompletedAt = &t
	}

	return TaskResponse{
		ID:              task.ID,
		ChatID:          task.ChatID,
		Title:           task.Title,
		Description:     task.Description,
		Category:        task.Category,
		Difficulty:      task.Difficulty,
		ScheduleJSON:    task.ScheduleJSON,
		BaseDuration:    task.BaseDuration,
		GracePeriod:     task.GracePeriod,
		Cooldown:        task.Cooldown,
		RewardCurveJSON: task.RewardCurveJSON,
		PartialCredit:   task.PartialCredit,
		StreakEnabled:   task.StreakEnabled,
		Status:          task.Status,
		LastCompletedAt: lastCompletedAt,
		StreakCount:     task.StreakCount,
		TimeZone:        task.TimeZone,
	}
}
