package controllers

import (
	"issue-tracking/entities"
	"issue-tracking/utils"
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
		utils.RespondError(c, 500, "Failed to fetch comments", nil)
		return
	}
	if comments == nil {
		comments = []entities.Comment{}
	}
	utils.RespondSuccess(c, 200, comments)
}

// CreateComment creates a new comment on an issue
func (cc *CommentController) CreateComment(c *gin.Context) {
	issueIDStr := c.Param("id")
	issueID, err := strconv.ParseUint(issueIDStr, 10, 32)
	if err != nil {
		utils.RespondError(c, 400, "Invalid issue ID", "issue_id must be a positive integer")
		return
	}

	type CommentRequest struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate UserID
	var user entities.User
	if err := cc.db.First(&user, req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 400, "User not found", "invalid user_id")
			return
		}
		utils.RespondError(c, 500, "Failed to validate user", nil)
		return
	}

	comment := entities.Comment{
		IssueID: uint(issueID),
		UserID:  req.UserID,
		Content: req.Content,
	}

	// Validate the comment
	if validationErrors := utils.ValidateStruct(comment); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	if err := cc.db.Create(&comment).Error; err != nil {
		utils.RespondError(c, 500, "Failed to create comment", err.Error())
		return
	}

	// Preload user info
	if err := cc.db.Preload("User").First(&comment).Error; err != nil {
		utils.RespondError(c, 500, "Failed to fetch created comment", nil)
		return
	}

	if err := cc.db.Preload("Issue").First(&comment).Error; err != nil {
		utils.RespondError(c, 500, "Failed to fetch created comment", nil)
		return
	}

	utils.RespondSuccess(c, 201, comment)
}

// GetComment retrieves a single comment by ID
func (cc *CommentController) GetComment(c *gin.Context) {
	id := c.Param("id")
	var comment entities.Comment
	if err := cc.db.Preload("User").Preload("Issue").First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 404, "Comment not found", nil)
			return
		}
		utils.RespondError(c, 500, "Failed to fetch comment", nil)
		return
	}
	utils.RespondSuccess(c, 200, comment)
}

// UpdateComment updates an existing comment
func (cc *CommentController) UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var comment entities.Comment
	if err := cc.db.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 404, "Comment not found", nil)
			return
		}
		utils.RespondError(c, 500, "Failed to fetch comment", nil)
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		utils.RespondError(c, 400, "Invalid request body", err.Error())
		return
	}

	if validationErrors := utils.ValidateStruct(comment); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	if err := cc.db.Save(&comment).Error; err != nil {
		utils.RespondError(c, 500, "Failed to update comment", err.Error())
		return
	}
	utils.RespondSuccess(c, 200, comment)
}

// DeleteComment deletes a comment
func (cc *CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if err := cc.db.Delete(&entities.Comment{}, id).Error; err != nil {
		utils.RespondError(c, 500, "Failed to delete comment", err.Error())
		return
	}
	c.JSON(204, nil)
}
