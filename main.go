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

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&entities.User{}); err != nil {
		log.Fatalf("failed to migrate User: %v", err)
	}

	if err := db.AutoMigrate(&entities.Officer{}); err != nil {
		log.Fatalf("failed to migrate Officer: %v", err)
	}

	if err := db.AutoMigrate(&entities.IssueStatus{}); err != nil {
		log.Fatalf("failed to migrate IssueStatus: %v", err)
	}

	if err := db.AutoMigrate(&entities.Issue{}); err != nil {
		log.Fatalf("failed to migrate Issue: %v", err)
	}

	if err := db.AutoMigrate(&entities.IssueStatusHistory{}); err != nil {
		log.Fatalf("failed to migrate IssueStatusHistory: %v", err)
	}

	if err := db.AutoMigrate(&entities.Comment{}); err != nil {
		log.Fatalf("failed to migrate Comment: %v", err)
	}

	//! create mock data
	// utils.MockData(db)

	fmt.Println("Database initialized successfully")
}

func main() {
	router := gin.Default()

	// CORS middleware to allow all origins
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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
