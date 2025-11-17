package utils

import (
	"issue-tracking/entities"
	"log"

	"gorm.io/gorm"
)

func MockData(db *gorm.DB) error {
	//Create mock data
	mockUser := &entities.User{
		FullName: "John Doe",
	}
	if err := db.Create(mockUser).Error; err != nil {
		log.Fatalf("failed to create mock user: %v", err)
		return err
	}

	mockOfficer := &entities.Officer{
		FullName: "Jane Smith",
	}
	if err := db.Create(mockOfficer).Error; err != nil {
		log.Fatalf("failed to create mock officer: %v", err)
		return err
	}

	issueStatuses := []*entities.IssueStatus{
		{
			StatusCode:   "open",
			DisplayName:  "Open",
			Description:  "Issue is open and pending review",
			Color:        "#FF0000",
			DisplayOrder: 1,
			IsActive:     true,
		},
		{
			StatusCode:   "in-progress",
			DisplayName:  "In Progress",
			Description:  "Issue is being worked on",
			Color:        "#FFA500",
			DisplayOrder: 2,
			IsActive:     true,
		},
		{
			StatusCode:   "closed",
			DisplayName:  "Closed",
			Description:  "Issue is closed",
			Color:        "#00FF00",
			DisplayOrder: 3,
			IsActive:     true,
		},
	}
	for _, status := range issueStatuses {
		if err := db.Create(status).Error; err != nil {
			log.Fatalf("failed to create mock issue status: %v", err)
			return err
		}
	}
	return nil
}
