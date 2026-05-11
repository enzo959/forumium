package migrations

import (
	"log"
	"github.com/enzo959/forumium/models"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)

	if count > 0 {
		log.Println("Categories already seeded, skipping...")
		return
	}

	categories := []models.Category{
		{Name: "Général"},
		{Name: "Critiques"},
		{Name: "Recommandations"},
	}

	result := db.Create(&categories)
	if result.Error != nil {
		log.Fatal("Failed to seed categories:", result.Error)
	}

	log.Println("Categories seeded successfully")
}