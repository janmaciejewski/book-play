package config

import (
	"log"

	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase seeds the database with initial data
func SeedDatabase(db *gorm.DB) error {
	// Check if already seeded
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database with initial data...")

	// Create default users
	users := []models.User{
		{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Email:        "admin@bookplay.com",
			PasswordHash: hashPassword("admin123"),
			FirstName:    "Admin",
			LastName:     "User",
			Role:         models.RoleAdmin,
			IsActive:     true,
		},
		{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Email:        "owner@bookplay.com",
			PasswordHash: hashPassword("owner123"),
			FirstName:    "John",
			LastName:     "Owner",
			Phone:        strPtr("+48123456789"),
			Role:         models.RoleFacilityOwner,
			IsActive:     true,
		},
		{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Email:        "player@bookplay.com",
			PasswordHash: hashPassword("player123"),
			FirstName:    "Jane",
			LastName:     "Player",
			Phone:        strPtr("+48987654321"),
			Role:         models.RolePlayer,
			IsActive:     true,
		},
		{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Email:        "captain@bookplay.com",
			PasswordHash: hashPassword("captain123"),
			FirstName:    "Mike",
			LastName:     "Captain",
			Role:         models.RolePlayer,
			IsActive:     true,
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d users", len(users))

	// Create facilities
	facilities := []models.Facility{
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000001"),
			Name:        "Central Football Stadium",
			Description: strPtr("Professional football field with artificial turf"),
			Type:        "football",
			Address:     "ul. Sportowa 1",
			City:        "Warsaw",
			Lat:         floatPtr(52.2297),
			Lng:         floatPtr(21.0122),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(150.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000002"),
			Name:        "City Basketball Arena",
			Description: strPtr("Indoor basketball court with professional equipment"),
			Type:        "basketball",
			Address:     "ul. Koszykowa 2",
			City:        "Warsaw",
			Lat:         floatPtr(52.2185),
			Lng:         floatPtr(21.0138),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(100.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000003"),
			Name:        "Riverside Tennis Club",
			Description: strPtr("Outdoor tennis courts with lighting"),
			Type:        "tennis",
			Address:     "ul. Tenisowa 5",
			City:        "Warsaw",
			Lat:         floatPtr(52.2350),
			Lng:         floatPtr(21.0250),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(80.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000004"),
			Name:        "Beach Volleyball Courts",
			Description: strPtr("Sandy beach volleyball courts near the river"),
			Type:        "volleyball",
			Address:     "ul. Plażowa 10",
			City:        "Warsaw",
			Lat:         floatPtr(52.2400),
			Lng:         floatPtr(21.0300),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(60.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000005"),
			Name:        "Indoor Swimming Pool",
			Description: strPtr("Olympic-size swimming pool with lanes"),
			Type:        "swimming",
			Address:     "ul. Pływacka 8",
			City:        "Krakow",
			Lat:         floatPtr(50.0647),
			Lng:         floatPtr(19.9450),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(45.00),
			IsActive:    true,
		},
	}

	for _, facility := range facilities {
		if err := db.Create(&facility).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d facilities", len(facilities))

	// Create facility slots (opening hours for each facility)
	for _, facility := range facilities {
		slots := []models.FacilitySlot{
			{FacilityID: facility.ID, DayOfWeek: 0, OpenTime: "08:00", CloseTime: "20:00"}, // Sunday
			{FacilityID: facility.ID, DayOfWeek: 1, OpenTime: "06:00", CloseTime: "22:00"}, // Monday
			{FacilityID: facility.ID, DayOfWeek: 2, OpenTime: "06:00", CloseTime: "22:00"}, // Tuesday
			{FacilityID: facility.ID, DayOfWeek: 3, OpenTime: "06:00", CloseTime: "22:00"}, // Wednesday
			{FacilityID: facility.ID, DayOfWeek: 4, OpenTime: "06:00", CloseTime: "22:00"}, // Thursday
			{FacilityID: facility.ID, DayOfWeek: 5, OpenTime: "06:00", CloseTime: "22:00"}, // Friday
			{FacilityID: facility.ID, DayOfWeek: 6, OpenTime: "08:00", CloseTime: "20:00"}, // Saturday
		}
		for _, slot := range slots {
			if err := db.Create(&slot).Error; err != nil {
				return err
			}
		}
	}
	log.Printf("Created facility slots for all facilities")

	// Create teams
	teams := []models.Team{
		{
			ID:          uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			Name:        "Warriors FC",
			Description: strPtr("Amateur football team from Warsaw"),
			CaptainID:   uuid.MustParse("00000000-0000-0000-0000-000000000004"),
		},
		{
			ID:          uuid.MustParse("20000000-0000-0000-0000-000000000002"),
			Name:        "Slam Dunkers",
			Description: strPtr("Basketball enthusiasts team"),
			CaptainID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		},
	}

	for _, team := range teams {
		if err := db.Create(&team).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d teams", len(teams))

	// Create team members
	teamMembers := []models.TeamMember{
		{
			TeamID: uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			UserID: uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Role:   models.TeamRoleCaptain,
		},
		{
			TeamID: uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			UserID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Role:   models.TeamRoleMember,
		},
		{
			TeamID: uuid.MustParse("20000000-0000-0000-0000-000000000002"),
			UserID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Role:   models.TeamRoleCaptain,
		},
		{
			TeamID: uuid.MustParse("20000000-0000-0000-0000-000000000002"),
			UserID: uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Role:   models.TeamRoleMember,
		},
	}

	for _, member := range teamMembers {
		if err := db.Create(&member).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d team members", len(teamMembers))

	log.Println("✅ Database seeding completed successfully!")
	log.Println("")
	log.Println("📋 Test accounts created:")
	log.Println("   admin@bookplay.com   / admin123   (Admin)")
	log.Println("   owner@bookplay.com   / owner123   (Facility Owner)")
	log.Println("   player@bookplay.com / player123   (Player)")
	log.Println("   captain@bookplay.com / captain123 (Team Captain)")
	log.Println("")

	return nil
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashed)
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
