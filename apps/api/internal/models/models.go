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
	TeamRoleAdmin   TeamRole = "ADMIN"
	TeamRoleMember  TeamRole = "MEMBER"
)

// User model
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string    `gorm:"not null;size:255" json:"-"`
	FirstName    string    `gorm:"not null;size:100" json:"first_name"`
	LastName     string    `gorm:"not null;size:100" json:"last_name"`
	Phone        *string   `gorm:"size:20" json:"phone,omitempty"`
	Role         Role      `gorm:"type:varchar(20);default:PLAYER;not null" json:"role"`
	Avatar       *string   `gorm:"size:500" json:"avatar,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Facilities    []Facility     `gorm:"foreignKey:OwnerID" json:"facilities,omitempty"`
	Reservations  []Reservation  `gorm:"foreignKey:UserID" json:"reservations,omitempty"`
	Teams         []TeamMember   `gorm:"foreignKey:UserID" json:"teams,omitempty"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID" json:"-"`
}

// TableName for User
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook for User
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

// TableName for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// BeforeCreate hook for RefreshToken
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
	Type        string          `gorm:"not null;size:50" json:"type"` // football, basketball, tennis
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

	// Relations
	Slots        []FacilitySlot `gorm:"foreignKey:FacilityID" json:"slots,omitempty"`
	Reservations []Reservation  `gorm:"foreignKey:FacilityID" json:"reservations,omitempty"`
}

// TableName for Facility
func (Facility) TableName() string {
	return "facilities"
}

// BeforeCreate hook for Facility
func (f *Facility) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// JSONB type for PostgreSQL JSONB
type JSONB map[string]interface{}

// FacilitySlot model
type FacilitySlot struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FacilityID uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:facility_day" json:"facility_id"`
	Facility   Facility         `gorm:"foreignKey:FacilityID" json:"facility,omitempty"`
	DayOfWeek  int              `gorm:"not null;uniqueIndex:facility_day" json:"day_of_week"` // 0-6 (Sunday-Saturday)
	OpenTime   string           `gorm:"not null;size:5" json:"open_time"`                     // "08:00"
	CloseTime  string           `gorm:"not null;size:5" json:"close_time"`                    // "22:00"
	Price      *decimal.Decimal `gorm:"type:decimal(10,2)" json:"price,omitempty"`            // Optional override
}

// TableName for FacilitySlot
func (FacilitySlot) TableName() string {
	return "facility_slots"
}

// BeforeCreate hook for FacilitySlot
func (f *FacilitySlot) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// Reservation model
type Reservation struct {
	ID         uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FacilityID uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:facility_date_time;index" json:"facilityId"`
	Facility   Facility          `gorm:"foreignKey:FacilityID" json:"facility,omitempty"`
	UserID     uuid.UUID         `gorm:"type:uuid;not null;index" json:"userId"`
	User       User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TeamID     *uuid.UUID        `gorm:"type:uuid;index" json:"teamId,omitempty"`
	Team       *Team             `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Date       time.Time         `gorm:"type:date;not null;uniqueIndex:facility_date_time" json:"date"`
	StartTime  string            `gorm:"not null;size:5;uniqueIndex:facility_date_time" json:"startTime"` // "10:00"
	EndTime    string            `gorm:"not null;size:5;uniqueIndex:facility_date_time" json:"endTime"`   // "11:00"
	Status     ReservationStatus `gorm:"type:varchar(20);default:PENDING;not null" json:"status"`
	TotalPrice decimal.Decimal   `gorm:"type:decimal(10,2);not null" json:"totalPrice"`
	Notes      *string           `gorm:"type:text" json:"notes,omitempty"`
	Version    int               `gorm:"default:0" json:"version"` // Optimistic locking
	CreatedAt  time.Time         `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time         `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName for Reservation
func (Reservation) TableName() string {
	return "reservations"
}

// BeforeCreate hook for Reservation
func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Team model
type Team struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null;size:200" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CaptainID   uuid.UUID `gorm:"type:uuid;not null" json:"captain_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Members      []TeamMember  `gorm:"foreignKey:TeamID" json:"members,omitempty"`
	Reservations []Reservation `gorm:"foreignKey:TeamID" json:"reservations,omitempty"`
}

// TableName for Team
func (Team) TableName() string {
	return "teams"
}

// BeforeCreate hook for Team
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
	Role     TeamRole  `gorm:"type:varchar(20);default:MEMBER;not null" json:"role"` // CAPTAIN, ADMIN, MEMBER
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

// TableName for TeamMember
func (TeamMember) TableName() string {
	return "team_members"
}

// BeforeCreate hook for TeamMember
func (t *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// StringArray type for PostgreSQL text[]
type StringArray []string
