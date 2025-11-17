package controllers

import (
	"issue-tracking/entities"
	"issue-tracking/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OfficerController struct {
	db *gorm.DB
}

// NewOfficerController creates a new officer controller
func NewOfficerController(db *gorm.DB) *OfficerController {
	return &OfficerController{db: db}
}

// GetAllOfficers retrieves all officers
func (oc *OfficerController) GetAllOfficers(c *gin.Context) {
	var officers []entities.Officer

	if err := oc.db.Find(&officers).Error; err != nil {
		utils.RespondError(c, 500, "Failed to fetch officers", nil)
		return
	}

	if officers == nil {
		officers = []entities.Officer{}
	}

	utils.RespondSuccess(c, 200, officers)
}
