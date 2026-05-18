package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/config"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/mail"
)

func main() {
	fmt.Println("=== BookPlay Email Reminder Test ===")
	fmt.Println()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to DB
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	// Run auto-migrate to ensure reminder_sent column exists
	if err := db.AutoMigrate(&models.Reservation{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	mailService := mail.NewService(db)

	// First, check if test user and reservation exist
	testUserID := uuid.MustParse("40000000-0000-0000-0000-000000000001")
	testReservationID := uuid.MustParse("40000000-0000-0000-0000-000000000001")

	var testUser models.User
	if err := db.First(&testUser, "id = ?", testUserID).Error; err != nil {
		fmt.Printf("❌ Test user not found: %v\n", err)
		fmt.Println("   Make sure the server has been started at least once after the changes.")
		os.Exit(1)
	}
	fmt.Printf("✅ Found test user: %s <%s>\n", testUser.FirstName, testUser.Email)

	var testReservation models.Reservation
	if err := db.Preload("User").Preload("Facility").First(&testReservation, "id = ?", testReservationID).Error; err != nil {
		fmt.Printf("❌ Test reservation not found: %v\n", err)
		fmt.Println("   Make sure the server has been started at least once after the changes.")
		os.Exit(1)
	}
	fmt.Printf("✅ Found test reservation:\n")
	fmt.Printf("   Facility: %s\n", testReservation.Facility.Name)
	fmt.Printf("   Date: %s\n", testReservation.Date.Format("2006-01-02"))
	fmt.Printf("   Time: %s - %s\n", testReservation.StartTime, testReservation.EndTime)
	fmt.Printf("   Price: %s PLN\n", testReservation.TotalPrice.String())
	fmt.Printf("   Status: %s\n", testReservation.Status)
	fmt.Printf("   ReminderSent: %v\n", testReservation.ReminderSent)
	fmt.Println()

	// Reset reminder_sent so we can test again
	if testReservation.ReminderSent {
		fmt.Println("🔄 Resetting reminder_sent flag...")
		db.Model(&testReservation).Update("reminder_sent", false)
	}

	// Re-fetch with preloaded relations after reset
	if err := db.Preload("User").Preload("Facility").First(&testReservation, "id = ?", testReservationID).Error; err != nil {
		log.Fatalf("Failed to refetch reservation: %v", err)
	}

	// Send the reminder email
	fmt.Println("📧 Sending reminder email...")
	if err := mailService.SendReservationReminder(&testReservation); err != nil {
		fmt.Printf("❌ Failed to send email: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Reminder email sent successfully!")
	fmt.Println()
	fmt.Println("Check your inbox at mufipl@gmail.com (and spam folder).")
}
