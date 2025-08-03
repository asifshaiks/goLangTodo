// ================== internal/features/todos/handler.go ==================
package todos

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := ValidateCreateTodo(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := &Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Tags:        req.Tags,
		DueDate:     req.DueDate,
	}

	if todo.Priority == "" {
		todo.Priority = "medium"
	}

	if todo.Tags == nil {
		todo.Tags = []string{}
	}

	if err := h.repo.Create(c.Request.Context(), todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"todo": todo})
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	todoID := c.Param("id")

	todo, err := h.repo.GetByID(c.Request.Context(), todoID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": TranslateTodoError(err)})
		return
	}

	if todo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	todoID := c.Param("id")

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := ValidateUpdateTodo(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build update document
	update := bson.M{}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Completed != nil {
		update["completed"] = *req.Completed
	}
	if req.Priority != "" {
		update["priority"] = req.Priority
	}
	if req.Tags != nil {
		update["tags"] = req.Tags
	}
	if req.DueDate != nil {
		update["dueDate"] = req.DueDate
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := h.repo.Update(c.Request.Context(), todoID, userID, update); err != nil {
		if err.Error() == "Todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": TranslateTodoError(err)})
		return
	}

	// Get updated todo
	todo, err := h.repo.GetByID(c.Request.Context(), todoID, userID)
	if err != nil || todo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	todoID := c.Param("id")

	if err := h.repo.Delete(c.Request.Context(), todoID, userID); err != nil {
		if err.Error() == "Todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": TranslateTodoError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("userID")

	// Query params
	completedStr := c.Query("completed")
	limitStr := c.Query("limit")

	var completed *bool
	if completedStr != "" {
		val, err := strconv.ParseBool(completedStr)
		if err == nil {
			completed = &val
		}
	}

	limit := 50 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	todos, err := h.repo.List(c.Request.Context(), userID, completed, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get todos"})
		return
	}

	// Get total count
	total, _ := h.repo.CountByUser(c.Request.Context(), userID)

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"total": total,
		"limit": limit,
	})
}
