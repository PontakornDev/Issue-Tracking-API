package controllers

import (
	"issue-tracking/entities"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct {
	db *gorm.DB
}

// NewCommentController creates a new comment controller
func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{db: db}
}

// GetCommentsByIssue retrieves all comments for an issue
func (cc *CommentController) GetCommentsByIssue(c *gin.Context) {
	issueID := c.Param("issue_id")
	var comments []entities.Comment
	if err := cc.db.
		Where("issue_id = ?", issueID).
		Preload("User").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(200, comments)
}

// CreateComment creates a new comment on an issue
func (cc *CommentController) CreateComment(c *gin.Context) {
	issueIDStr := c.Param("id")
	issueID, err := strconv.ParseUint(issueIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid issue ID"})
		return
	}

	type CommentRequest struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	comment := entities.Comment{
		IssueID: uint(issueID),
		UserID:  req.UserID,
		Content: req.Content,
	}

	if err := cc.db.Create(&comment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create comment"})
		return
	}

	// Preload user info
	if err := cc.db.Preload("User").First(&comment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch created comment"})
		return
	}

	c.JSON(201, comment)
}

// GetComment retrieves a single comment by ID
func (cc *CommentController) GetComment(c *gin.Context) {
	id := c.Param("id")
	var comment entities.Comment
	if err := cc.db.Preload("User").Preload("Issue").First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch comment"})
		return
	}
	c.JSON(200, comment)
}

// UpdateComment updates an existing comment
func (cc *CommentController) UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var comment entities.Comment
	if err := cc.db.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch comment"})
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := cc.db.Save(&comment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update comment"})
		return
	}
	c.JSON(200, comment)
}

// DeleteComment deletes a comment
func (cc *CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if err := cc.db.Delete(&entities.Comment{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(204, nil)
}
