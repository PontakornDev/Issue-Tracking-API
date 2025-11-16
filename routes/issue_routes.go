package routes

import (
	"issue-tracking/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize controllers
	issueController := controllers.NewIssueController(db)
	commentController := controllers.NewCommentController(db)

	// Issues routes
	issues := router.Group("/api/issues")
	{
		issues.POST("", issueController.CreateIssue)
		issues.GET("", issueController.GetAllIssues)
		issues.GET("/:id", issueController.GetIssue)
		issues.PATCH("/:id/status", issueController.UpdateIssueStatus)
		issues.POST("/:id/comment", commentController.CreateComment)
	}
}
