package database

import (
	"letsquiz/server/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

var DB *gorm.DB

func ConnectDatabase(dbType, dsn string) {
	var err error
	log.Println("dbType:", dbType)
	log.Println("dsn:", dsn)
	switch dbType {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported database type: %s", dbType)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Quiz{}, &models.Question{}, &models.Answer{}, &models.UserQuizAttempt{}, &models.UserAnswer{}, &models.Leaderboard{}, &models.Feedback{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate schema: %v", err)
	}

	log.Println("Database connection established and schema migrated")
}
