package main

import (
	"fmt"
	"log"
	"os"

	"issue-tracking/entities"
	"issue-tracking/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	// Read PostgreSQL connection string from environment or use default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=issue_tracking port=5432 sslmode=disable"
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(
		&entities.User{},
		&entities.Officer{},
		&entities.IssueStatus{},
		&entities.Issue{},
		&entities.IssueStatusHistory{},
		&entities.Comment{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database initialized successfully")
}

func main() {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register all routes
	routes.RegisterRoutes(router, db)

	fmt.Println("Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
