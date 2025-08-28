package main

import (
	"backend/app"
	"backend/internal/model"
	"log"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	app.Connect()

	// Create admin user if not exists
	createAdminUser()

	log.Println("✅ Database setup completed!")
}

func createAdminUser() {
	var user model.User
	result := app.DB.Where("username = ?", "admin").First(&user)
	
	if result.Error != nil {
		// Admin user doesn't exist, create one
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("❌ Failed to hash password:", err)
		}

		admin := model.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			FullName: "System Administrator",
			Role:     "admin",
			IsActive: true,
		}

		if err := app.DB.Create(&admin).Error; err != nil {
			log.Fatal("❌ Failed to create admin user:", err)
		}

		log.Println("✅ Admin user created successfully!")
		log.Println("   Username: admin")
		log.Println("   Password: admin123")
		log.Println("   Please change the password after first login!")
	} else {
		log.Println("ℹ️ Admin user already exists")
	}
}
