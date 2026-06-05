package config

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase wypełnia bazę danych początkowymi danymi testowymi
func SeedDatabase(db *gorm.DB) error {
	// Sprawdza czy dane testowe już istnieją (znacznik: admin z konkretnym ID)
	var existing models.User
	if err := db.First(&existing, "id = ?", uuid.MustParse("00000000-0000-0000-0000-000000000001")).Error; err == nil {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database with initial data...")

	// Tworzy domyślnych użytkowników testowych
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

	// Tworzy obiekty – wszystkie w okolicach Grodziska Wielkopolskiego
	facilities := []models.Facility{
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000001"),
			Name:        "Stadion Miejski Grodzisk Wielkopolski",
			Description: strPtr("Pełnowymiarowe boisko piłkarskie z sztuczną trawą i oświetleniem"),
			Type:        "FOOTBALL",
			Address:     "ul. Sportowa 5",
			City:        "Grodzisk Wielkopolski",
			Lat:         floatPtr(52.22850),
			Lng:         floatPtr(16.36620),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(80.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000002"),
			Name:        "Hala Sportowa OSiR",
			Description: strPtr("Nowoczesna hala sportowa z boiskiem do koszykówki i siatkówki"),
			Type:        "OTHER",
			Address:     "ul. Oświatowa 8",
			City:        "Grodzisk Wielkopolski",
			Lat:         floatPtr(52.22680),
			Lng:         floatPtr(16.36890),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(60.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000003"),
			Name:        "Korty Tenisowe KS Grodzisk",
			Description: strPtr("Dwa korty tenisowe z nawierzchnią ceglaną, dostępne cały rok"),
			Type:        "TENNIS",
			Address:     "ul. Tenisowa 3",
			City:        "Grodzisk Wielkopolski",
			Lat:         floatPtr(52.22910),
			Lng:         floatPtr(16.36250),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(40.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000004"),
			Name:        "Basen Kompleks Sportowy",
			Description: strPtr("Basen rekreacyjny 25m z torami pływackimi i brodzikiem dla dzieci"),
			Type:        "SWIMMING",
			Address:     "ul. Pływania 12",
			City:        "Grodzisk Wielkopolski",
			Lat:         floatPtr(52.22700),
			Lng:         floatPtr(16.37000),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(25.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000005"),
			Name:        "Boisko do Siatkówki Plażowej",
			Description: strPtr("Profesjonalne boisko do siatkówki plażowej przy parku miejskim"),
			Type:        "VOLLEYBALL",
			Address:     "ul. Parkowa 2",
			City:        "Grodzisk Wielkopolski",
			Lat:         floatPtr(52.22650),
			Lng:         floatPtr(16.36380),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(30.00),
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("10000000-0000-0000-0000-000000000007"),
			Name:        "Boisko Orlik",
			Description: strPtr("Boisko wielofunkcyjne do piłki nożnej, koszykówki i siatkówki"),
			Type:        "FOOTBALL",
			Address:     "ul. Szkolna 10",
			City:        "Grodzisk Wielkopolski",
			Lat:         floatPtr(52.23000),
			Lng:         floatPtr(16.36400),
			OwnerID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			HourlyRate:  decimal.NewFromFloat(40.00),
			IsActive:    true,
		},
	}

	for _, facility := range facilities {
		if err := db.Create(&facility).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d facilities", len(facilities))

	// Tworzy godziny otwarcia dla każdego obiektu
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

	// Tworzy drużyny testowe
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

	// Przypisuje członków do drużyn
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

	// Tworzy przykładowe rezerwacje
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	nextWeek := today.AddDate(0, 0, 7)
	nextWeekPlus1 := today.AddDate(0, 0, 8)
	inTwoWeeks := today.AddDate(0, 0, 14)

	teamID1 := uuid.MustParse("20000000-0000-0000-0000-000000000001")
	teamID2 := uuid.MustParse("20000000-0000-0000-0000-000000000002")

	reservations := []models.Reservation{
		// Dzisiejsze rezerwacje – wszystkie należą do kapitana
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000001"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000001"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			TeamID:     &teamID1,
			Date:       today,
			StartTime:  "10:00",
			EndTime:    "12:00",
			Status:     models.StatusConfirmed,
			TotalPrice: decimal.NewFromFloat(60.00),
			Notes:      strPtr("Cotygodniowy trening drużyny"),
		},
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000002"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000002"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Date:       today,
			StartTime:  "14:00",
			EndTime:    "16:00",
			Status:     models.StatusConfirmed,
			TotalPrice: decimal.NewFromFloat(80.00),
			Notes:      strPtr("Mecz koszykówki ze znajomymi"),
		},
		// Rezerwacje na jutro
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000003"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000003"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Date:       tomorrow,
			StartTime:  "09:00",
			EndTime:    "10:00",
			Status:     models.StatusPending,
			TotalPrice: decimal.NewFromFloat(50.00),
			Notes:      strPtr("Lekcja tenisa"),
		},
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000004"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000004"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			TeamID:     &teamID2,
			Date:       tomorrow,
			StartTime:  "16:00",
			EndTime:    "18:00",
			Status:     models.StatusConfirmed,
			TotalPrice: decimal.NewFromFloat(80.00),
			Notes:      strPtr("Turniej siatkówki plażowej"),
		},
		// Rezerwacje na przyszły tydzień
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000005"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000001"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			TeamID:     &teamID1,
			Date:       nextWeek,
			StartTime:  "18:00",
			EndTime:    "20:00",
			Status:     models.StatusConfirmed,
			TotalPrice: decimal.NewFromFloat(60.00),
			Notes:      strPtr("Wieczorny mecz piłki nożnej"),
		},
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000006"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000005"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Date:       nextWeekPlus1,
			StartTime:  "07:00",
			EndTime:    "08:00",
			Status:     models.StatusPending,
			TotalPrice: decimal.NewFromFloat(50.00),
			Notes:      strPtr("Poranny trening pływacki"),
		},
		// Zakończona rezerwacja w przeszłości
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000007"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000002"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Date:       today.AddDate(0, 0, -3),
			StartTime:  "15:00",
			EndTime:    "17:00",
			Status:     models.StatusCompleted,
			TotalPrice: decimal.NewFromFloat(80.00),
			Notes:      strPtr("Mecz koszykówki w zeszłym tygodniu"),
		},
		// Anulowana rezerwacja
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000008"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000003"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Date:       inTwoWeeks,
			StartTime:  "10:00",
			EndTime:    "11:00",
			Status:     models.StatusCancelled,
			TotalPrice: decimal.NewFromFloat(50.00),
			Notes:      strPtr("Odwołane z powodu pogody"),
		},
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000009"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000001"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			TeamID:     &teamID1,
			Date:       nextWeek,
			StartTime:  "10:00",
			EndTime:    "12:00",
			Status:     models.StatusPending,
			TotalPrice: decimal.NewFromFloat(60.00),
			Notes:      strPtr("Sesja treningowa drużyny"),
		},
		{
			ID:         uuid.MustParse("30000000-0000-0000-0000-000000000010"),
			FacilityID: uuid.MustParse("10000000-0000-0000-0000-000000000004"),
			UserID:     uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Date:       inTwoWeeks,
			StartTime:  "14:00",
			EndTime:    "16:00",
			Status:     models.StatusConfirmed,
			TotalPrice: decimal.NewFromFloat(80.00),
			Notes:      strPtr("Weekendowa siatkówka plażowa"),
		},
	}

	for _, reservation := range reservations {
		if err := db.Create(&reservation).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d reservations", len(reservations))

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
