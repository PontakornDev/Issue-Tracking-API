package controllers

import (
	"issue-tracking/entities"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IssueController struct {
	db *gorm.DB
}

// NewIssueController creates a new issue controller
func NewIssueController(db *gorm.DB) *IssueController {
	return &IssueController{db: db}
}

// GetAllIssues retrieves all issues with optional status filter
func (ic *IssueController) GetAllIssues(c *gin.Context) {
	var issues []entities.Issue
	query := ic.db

	// Apply status filter if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status_id = (SELECT status_id FROM issue_statuses WHERE status_code = ?)", status)
	}

	if err := query.
		Preload("Reporter").
		Preload("Assignee").
		Preload("Status").
		Preload("Comments").
		Find(&issues).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch issues"})
		return
	}
	c.JSON(200, issues)
}

// GetIssue retrieves a single issue by ID with relations
func (ic *IssueController) GetIssue(c *gin.Context) {
	id := c.Param("id")
	var issue entities.Issue
	if err := ic.db.
		Preload("Reporter").
		Preload("Assignee").
		Preload("Status").
		Preload("StatusHistory").
		Preload("Comments").
		First(&issue, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Issue not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch issue"})
		return
	}
	c.JSON(200, issue)
}

// CreateIssue creates a new issue
func (ic *IssueController) CreateIssue(c *gin.Context) {
	var issue entities.Issue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ic.db.Create(&issue).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create issue"})
		return
	}
	c.JSON(201, issue)
}

// UpdateIssue updates an existing issue
func (ic *IssueController) UpdateIssue(c *gin.Context) {
	id := c.Param("id")
	issueID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid issue ID"})
		return
	}

	var issue entities.Issue
	if err := ic.db.First(&issue, issueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Issue not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch issue"})
		return
	}

	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ic.db.Save(&issue).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update issue"})
		return
	}
	c.JSON(200, issue)
}

// UpdateIssueStatus updates only the status of an issue
func (ic *IssueController) UpdateIssueStatus(c *gin.Context) {
	id := c.Param("id")
	issueID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid issue ID"})
		return
	}

	type StatusUpdate struct {
		NewStatusID uint   `json:"new_status_id" binding:"required"`
		Comment     string `json:"comment"`
	}

	var req StatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get the current issue
	var issue entities.Issue
	if err := ic.db.First(&issue, issueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Issue not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch issue"})
		return
	}

	oldStatusID := issue.StatusID

	// Update the issue status
	if err := ic.db.Model(&issue).Update("status_id", req.NewStatusID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update status"})
		return
	}

	// Record the status history
	history := entities.IssueStatusHistory{
		IssueID:     uint(issueID),
		OldStatusID: &oldStatusID,
		NewStatusID: req.NewStatusID,
		ChangedBy:   1, // Default to officer ID 1, should be from auth context in production
		Comment:     req.Comment,
	}

	if err := ic.db.Create(&history).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to record status history"})
		return
	}

	// Return updated issue with relations
	if err := ic.db.
		Preload("Reporter").
		Preload("Assignee").
		Preload("Status").
		Preload("Comments").
		First(&issue, issueID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated issue"})
		return
	}

	c.JSON(200, issue)
}

// DeleteIssue deletes an issue by ID
func (ic *IssueController) DeleteIssue(c *gin.Context) {
	id := c.Param("id")
	issueID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid issue ID"})
		return
	}

	if err := ic.db.Delete(&entities.Issue{}, issueID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete issue"})
		return
	}
	c.JSON(204, nil)
}
