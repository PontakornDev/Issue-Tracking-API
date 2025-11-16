package controllers

import (
	"issue-tracking/entities"
	"issue-tracking/utils"
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
		utils.RespondError(c, 500, "Failed to fetch issues", nil)
		return
	}

	if issues == nil {
		issues = []entities.Issue{}
	}
	utils.RespondSuccess(c, 200, issues)
}

// GetIssue retrieves a single issue by ID with relations
func (ic *IssueController) GetIssue(c *gin.Context) {
	id := c.Param("id")
	issueID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.RespondError(c, 400, "Invalid issue ID", "issue_id must be a positive integer")
		return
	}

	var issue entities.Issue
	if err := ic.db.
		Preload("Reporter").
		Preload("Assignee").
		Preload("Status").
		Preload("StatusHistory").
		Preload("Comments").
		First(&issue, issueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 404, "Issue not found", nil)
			return
		}
		utils.RespondError(c, 500, "Failed to fetch issue", nil)
		return
	}
	utils.RespondSuccess(c, 200, issue)
}

// CreateIssue creates a new issue
func (ic *IssueController) CreateIssue(c *gin.Context) {
	var issue entities.Issue
	if err := c.ShouldBindJSON(&issue); err != nil {
		utils.RespondError(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate the issue
	if validationErrors := utils.ValidateStruct(issue); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	// Validate reporter exists
	var reporter entities.User
	if err := ic.db.First(&reporter, issue.ReporterID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 400, "Reporter not found", "invalid reporter_id")
			return
		}
		utils.RespondError(c, 500, "Failed to validate reporter", nil)
		return
	}

	// Validate status exists
	var status entities.IssueStatus
	if err := ic.db.First(&status, issue.StatusID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 400, "Status not found", "invalid status_id")
			return
		}
		utils.RespondError(c, 500, "Failed to validate status", nil)
		return
	}

	// Validate assignee if provided
	if issue.AssigneeID != nil {
		var assignee entities.Officer
		if err := ic.db.First(&assignee, *issue.AssigneeID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.RespondError(c, 400, "Assignee not found", "invalid assignee_id")
				return
			}
			utils.RespondError(c, 500, "Failed to validate assignee", nil)
			return
		}
	}

	if err := ic.db.Create(&issue).Error; err != nil {
		utils.RespondError(c, 500, "Failed to create issue", err.Error())
		return
	}

	// Reload with relations
	ic.db.
		Preload("Reporter").
		Preload("Assignee").
		Preload("Status").
		Preload("Comments").
		First(&issue, issue.IssueID)

	utils.RespondSuccess(c, 201, issue)
}

// UpdateIssue updates an existing issue
func (ic *IssueController) UpdateIssue(c *gin.Context) {
	id := c.Param("id")
	issueID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.RespondError(c, 400, "Invalid issue ID", "issue_id must be a positive integer")
		return
	}

	var issue entities.Issue
	if err := ic.db.First(&issue, issueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 404, "Issue not found", nil)
			return
		}
		utils.RespondError(c, 500, "Failed to fetch issue", nil)
		return
	}

	if err := c.ShouldBindJSON(&issue); err != nil {
		utils.RespondError(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate the updated issue
	if validationErrors := utils.ValidateStruct(issue); len(validationErrors) > 0 {
		utils.RespondValidationError(c, validationErrors)
		return
	}

	if err := ic.db.Save(&issue).Error; err != nil {
		utils.RespondError(c, 500, "Failed to update issue", err.Error())
		return
	}

	// Reload with relations
	ic.db.
		Preload("Reporter").
		Preload("Assignee").
		Preload("Status").
		Preload("Comments").
		First(&issue, issueID)

	utils.RespondSuccess(c, 200, issue)
}

// UpdateIssueStatus updates only the status of an issue
func (ic *IssueController) UpdateIssueStatus(c *gin.Context) {
	id := c.Param("id")
	issueID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.RespondError(c, 400, "Invalid issue ID", "issue_id must be a positive integer")
		return
	}

	type StatusUpdate struct {
		NewStatusID uint   `json:"new_status_id" binding:"required"`
		Comment     string `json:"comment"`
	}

	var req StatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, 400, "Invalid request body", err.Error())
		return
	}

	// Get the current issue
	var issue entities.Issue
	if err := ic.db.First(&issue, issueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 404, "Issue not found", nil)
			return
		}
		utils.RespondError(c, 500, "Failed to fetch issue", nil)
		return
	}

	// Validate new status exists
	var status entities.IssueStatus
	if err := ic.db.First(&status, req.NewStatusID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, 400, "Invalid status", "status_id does not exist")
			return
		}
		utils.RespondError(c, 500, "Failed to validate status", nil)
		return
	}

	oldStatusID := issue.StatusID

	// Update the issue status
	if err := ic.db.Model(&issue).Update("status_id", req.NewStatusID).Error; err != nil {
		utils.RespondError(c, 500, "Failed to update status", err.Error())
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
		utils.RespondError(c, 500, "Failed to record status history", err.Error())
		return
	}

	// Return updated issue with relations
	if err := ic.db.
		Preload("Reporter").
		Preload("Assignee").
		Preload("Status").
		Preload("Comments").
		First(&issue, issueID).Error; err != nil {
		utils.RespondError(c, 500, "Failed to fetch updated issue", nil)
		return
	}

	utils.RespondSuccess(c, 200, issue)
}

// DeleteIssue deletes an issue by ID
func (ic *IssueController) DeleteIssue(c *gin.Context) {
	id := c.Param("id")
	issueID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.RespondError(c, 400, "Invalid issue ID", "issue_id must be a positive integer")
		return
	}

	if err := ic.db.Delete(&entities.Issue{}, issueID).Error; err != nil {
		utils.RespondError(c, 500, "Failed to delete issue", err.Error())
		return
	}
	c.JSON(204, nil)
}
