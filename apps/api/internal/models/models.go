package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Role enum
type Role string

const (
	RolePlayer        Role = "PLAYER"
	RoleTeamAdmin     Role = "TEAM_ADMIN"
	RoleFacilityOwner Role = "FACILITY_OWNER"
	RoleAdmin         Role = "ADMIN"
)

// ReservationStatus enum
type ReservationStatus string

const (
	StatusPending   ReservationStatus = "PENDING"
	StatusConfirmed ReservationStatus = "CONFIRMED"
	StatusCancelled ReservationStatus = "CANCELLED"
	StatusCompleted ReservationStatus = "COMPLETED"
)

// TeamRole enum
type TeamRole string

const (
	TeamRoleCaptain TeamRole = "CAPTAIN"
	TeamRoleMember  TeamRole = "MEMBER"
)

// User model
type User struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email             string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash      string    `gorm:"not null;size:255" json:"-"`
	FirstName         string    `gorm:"not null;size:100" json:"first_name"`
	LastName          string    `gorm:"not null;size:100" json:"last_name"`
	Phone             *string   `gorm:"size:20" json:"phone,omitempty"`
	Role              Role      `gorm:"type:varchar(20);default:PLAYER;not null" json:"role"`
	Avatar            *string   `gorm:"size:500" json:"avatar,omitempty"`
	Bio               *string   `gorm:"type:text" json:"bio,omitempty"`
	City              *string   `gorm:"size:100" json:"city,omitempty"`
	Country           *string   `gorm:"size:100" json:"country,omitempty"`
	Position          *string   `gorm:"size:100" json:"position,omitempty"`
	PreferredPosition *string   `gorm:"size:100" json:"preferred_position,omitempty"`
	Age               *int      `json:"age,omitempty"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Facilities    []Facility     `gorm:"foreignKey:OwnerID" json:"facilities,omitempty"`
	Reservations  []Reservation  `gorm:"foreignKey:UserID" json:"reservations,omitempty"`
	Teams         []TeamMember   `gorm:"foreignKey:UserID" json:"teams,omitempty"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string { return "users" }

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// RefreshToken model
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Token     string    `gorm:"uniqueIndex;not null;size:500" json:"token"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Facility model
type Facility struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"not null;size:200" json:"name"`
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	Type        string          `gorm:"not null;size:50" json:"type"`
	Address     string          `gorm:"not null;size:500" json:"address"`
	City        string          `gorm:"not null;size:100;index" json:"city"`
	Lat         *float64        `gorm:"type:decimal(10,8)" json:"lat,omitempty"`
	Lng         *float64        `gorm:"type:decimal(11,8)" json:"lng,omitempty"`
	OwnerID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"ownerId"`
	Owner       User            `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	HourlyRate  decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"hourlyRate"`
	Amenities   *JSONB          `gorm:"type:jsonb" json:"amenities,omitempty"`
	Images      StringArray     `gorm:"type:text[]" json:"images,omitempty"`
	IsActive    bool            `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`

	Slots        []FacilitySlot `gorm:"foreignKey:FacilityID" json:"slots,omitempty"`
	Reservations []Reservation  `gorm:"foreignKey:FacilityID" json:"reservations,omitempty"`
}

func (Facility) TableName() string { return "facilities" }

func (f *Facility) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type JSONB map[string]interface{}

// FacilitySlot model
type FacilitySlot struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FacilityID uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:facility_day" json:"facility_id"`
	Facility   Facility         `gorm:"foreignKey:FacilityID" json:"facility,omitempty"`
	DayOfWeek  int              `gorm:"not null;uniqueIndex:facility_day" json:"day_of_week"`
	OpenTime   string           `gorm:"not null;size:5" json:"open_time"`
	CloseTime  string           `gorm:"not null;size:5" json:"close_time"`
	Price      *decimal.Decimal `gorm:"type:decimal(10,2)" json:"price,omitempty"`
}

func (FacilitySlot) TableName() string { return "facility_slots" }

func (f *FacilitySlot) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// Reservation model
type Reservation struct {
	ID           uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FacilityID   uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:facility_date_time;index" json:"facilityId"`
	Facility     Facility          `gorm:"foreignKey:FacilityID" json:"facility,omitempty"`
	UserID       uuid.UUID         `gorm:"type:uuid;not null;index" json:"userId"`
	User         User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TeamID       *uuid.UUID        `gorm:"type:uuid;index" json:"teamId,omitempty"`
	Team         *Team             `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Date         time.Time         `gorm:"type:date;not null;uniqueIndex:facility_date_time" json:"date"`
	StartTime    string            `gorm:"not null;size:5;uniqueIndex:facility_date_time" json:"startTime"`
	EndTime      string            `gorm:"not null;size:5;uniqueIndex:facility_date_time" json:"endTime"`
	Status       ReservationStatus `gorm:"type:varchar(20);default:PENDING;not null" json:"status"`
	TotalPrice   decimal.Decimal   `gorm:"type:decimal(10,2);not null" json:"totalPrice"`
	Notes        *string           `gorm:"type:text" json:"notes,omitempty"`
	Version      int               `gorm:"default:0" json:"version"`
	ReminderSent bool              `gorm:"default:false" json:"-"`
	CreatedAt    time.Time         `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time         `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Reservation) TableName() string { return "reservations" }

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Team model
type Team struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string    `gorm:"not null;size:200" json:"name"`
	Description     *string   `gorm:"type:text" json:"description,omitempty"`
	CaptainID       uuid.UUID `gorm:"type:uuid;not null" json:"captain_id"`
	Logo            *string   `gorm:"size:500" json:"logo,omitempty"`
	RecruitmentOpen bool      `gorm:"default:false" json:"recruitment_open"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Members      []TeamMember  `gorm:"foreignKey:TeamID" json:"members,omitempty"`
	Reservations []Reservation `gorm:"foreignKey:TeamID" json:"reservations,omitempty"`
}

func (Team) TableName() string { return "teams" }

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TeamMember model
type TeamMember struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TeamID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:team_user;index" json:"team_id"`
	Team     Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:team_user;index" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role     TeamRole  `gorm:"type:varchar(20);default:MEMBER;not null" json:"role"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

func (TeamMember) TableName() string { return "team_members" }

func (t *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TeamRecruitmentApplication model
type TeamRecruitmentApplication struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null;index" json:"team_id"`
	Team      Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Status    string    `gorm:"type:varchar(20);default:PENDING;not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (TeamRecruitmentApplication) TableName() string { return "team_recruitment_applications" }

func (a *TeamRecruitmentApplication) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// EmailVerificationOTP model
type EmailVerificationOTP struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"not null;size:255;index" json:"email"`
	Code      string    `gorm:"not null;size:6" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (EmailVerificationOTP) TableName() string { return "email_verification_otps" }

func (o *EmailVerificationOTP) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type StringArray []string
