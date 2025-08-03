// ================== internal/features/todos/validator.go ==================
package todos

import (
	"errors"
	"strings"
	"time"
)

func ValidateCreateTodo(req *CreateTodoRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		return errors.New("Title is required")
	}

	if len(req.Title) < 3 {
		return errors.New("Title must be at least 3 characters")
	}

	if len(req.Title) > 200 {
		return errors.New("Title cannot exceed 200 characters")
	}

	if len(req.Description) > 1000 {
		return errors.New("Description cannot exceed 1000 characters")
	}

	if req.Priority != "" && req.Priority != "low" && req.Priority != "medium" && req.Priority != "high" {
		return errors.New("Priority must be low, medium, or high")
	}

	if req.DueDate != nil && req.DueDate.Before(time.Now().Add(-24*time.Hour)) {
		return errors.New("Due date cannot be in the past")
	}

	if len(req.Tags) > 10 {
		return errors.New("Cannot have more than 10 tags")
	}

	for i, tag := range req.Tags {
		req.Tags[i] = strings.TrimSpace(tag)
		if len(req.Tags[i]) > 20 {
			return errors.New("Tag cannot exceed 20 characters")
		}
	}

	return nil
}

func ValidateUpdateTodo(req *UpdateTodoRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title != "" && len(req.Title) < 3 {
		return errors.New("Title must be at least 3 characters")
	}

	if len(req.Title) > 200 {
		return errors.New("Title cannot exceed 200 characters")
	}

	if len(req.Description) > 1000 {
		return errors.New("Description cannot exceed 1000 characters")
	}

	if req.Priority != "" && req.Priority != "low" && req.Priority != "medium" && req.Priority != "high" {
		return errors.New("Priority must be low, medium, or high")
	}

	if len(req.Tags) > 10 {
		return errors.New("Cannot have more than 10 tags")
	}

	for i, tag := range req.Tags {
		req.Tags[i] = strings.TrimSpace(tag)
		if len(req.Tags[i]) > 20 {
			return errors.New("Tag cannot exceed 20 characters")
		}
	}

	return nil
}

// TranslateTodoError converts database errors to user-friendly messages
func TranslateTodoError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	if strings.Contains(errStr, "not found") {
		return "Todo not found"
	}

	if strings.Contains(errStr, "invalid") && strings.Contains(errStr, "ID") {
		return "Invalid todo ID format"
	}

	return "Something went wrong. Please try again"
}
