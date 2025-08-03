// ================== internal/features/todos/model.go ==================
package todos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"userId" json:"userId"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	Priority    string             `bson:"priority" json:"priority"` // low, medium, high
	Tags        []string           `bson:"tags" json:"tags"`
	DueDate     *time.Time         `bson:"dueDate,omitempty" json:"dueDate,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Tags        []string   `json:"tags"`
	DueDate     *time.Time `json:"dueDate"`
}

type UpdateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   *bool      `json:"completed"`
	Priority    string     `json:"priority"`
	Tags        []string   `json:"tags"`
	DueDate     *time.Time `json:"dueDate"`
}
