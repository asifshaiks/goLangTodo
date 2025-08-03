// ================== internal/features/todos/routes.go ==================
package todos

import (
	"gotodo/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *gin.RouterGroup, db *mongo.Database) {
	repo := NewRepository(db)
	handler := NewHandler(repo)

	todos := router.Group("/todos")
	todos.Use(middleware.Auth()) // All todo routes require authentication
	{
		todos.POST("/", handler.Create)
		todos.GET("/", handler.List)
		todos.GET("/:id", handler.Get)
		todos.PUT("/:id", handler.Update)
		todos.DELETE("/:id", handler.Delete)
	}
}
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("User not found")
	}

	return nil
}
