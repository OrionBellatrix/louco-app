package domain

import (
	"time"
)

type UserType string

const (
	UserTypeUser    UserType = "user"
	UserTypeCreator UserType = "creator"
)

type User struct {
	ID              int        `json:"id" db:"id"`
	FullName        string     `json:"full_name" db:"full_name"`
	Username        *string    `json:"username" db:"username"`
	Email           *string    `json:"email" db:"email"`
	Phone           *string    `json:"phone" db:"phone"`
	Password        string     `json:"-" db:"password"`
	UserType        UserType   `json:"user_type" db:"user_type"`
	AppleID         *string    `json:"apple_id" db:"apple_id"`
	GoogleID        *string    `json:"google_id" db:"google_id"`
	Biography       *string    `json:"biography" db:"biography"`
	BirthDate       *time.Time `json:"birth_date" db:"birth_date"`
	ProfilePicID    *int       `json:"profile_pic_id" db:"profile_pic_id"`
	CoverPicID      *int       `json:"cover_pic_id" db:"cover_pic_id"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" db:"email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at" db:"phone_verified_at"`
	FollowersCount  int        `json:"followers_count" db:"followers_count" gorm:"default:0"`
	FollowingCount  int        `json:"following_count" db:"following_count" gorm:"default:0"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`

	// Relations
	Creator        *Creator `json:"creator,omitempty" gorm:"foreignKey:UserID;references:ID"`
	ProfilePicture *Media   `json:"profile_picture,omitempty" gorm:"foreignKey:ProfilePicID;references:ID"`
	CoverPicture   *Media   `json:"cover_picture,omitempty" gorm:"foreignKey:CoverPicID;references:ID"`
}

func NewUser(fullName, password string, userType UserType) *User {
	return &User{
		FullName:  fullName,
		Password:  password,
		UserType:  userType,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) IsValidForLogin() bool {
	hasLoginCredential := (u.Email != nil && *u.Email != "") || (u.Username != nil && *u.Username != "") || (u.Phone != nil && *u.Phone != "")
	return u.IsActive && hasLoginCredential && u.Password != ""
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

func (u *User) SetEmail(email string) {
	u.Email = &email
	u.UpdatedAt = time.Now()
}

func (u *User) SetPhone(phone string) {
	u.Phone = &phone
	u.UpdatedAt = time.Now()
}

func (u *User) SetUsername(username string) {
	u.Username = &username
	u.UpdatedAt = time.Now()
}

func (u *User) SetProfilePicID(profilePicID int) {
	u.ProfilePicID = &profilePicID
	u.UpdatedAt = time.Now()
}

func (u *User) SetCoverPicID(coverPicID int) {
	u.CoverPicID = &coverPicID
	u.UpdatedAt = time.Now()
}

func (u *User) SetEmailVerified() {
	now := time.Now()
	u.EmailVerifiedAt = &now
	u.UpdatedAt = now
}

func (u *User) SetPhoneVerified() {
	now := time.Now()
	u.PhoneVerifiedAt = &now
	u.UpdatedAt = now
}

func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) IsPhoneVerified() bool {
	return u.PhoneVerifiedAt != nil
}

func (u *User) IsVerified() bool {
	// User is verified if they have verified either email or phone
	return u.IsEmailVerified() || u.IsPhoneVerified()
}

func (u *User) RequiresVerification() bool {
	// If user has email, email must be verified
	// If user has phone, phone must be verified
	if u.Email != nil && *u.Email != "" && !u.IsEmailVerified() {
		return true
	}
	if u.Phone != nil && *u.Phone != "" && !u.IsPhoneVerified() {
		return true
	}
	return false
}

func (u *User) UpdateProfile(fullName string, biography *string, birthDate *time.Time) {
	if fullName != "" {
		u.FullName = fullName
	}
	if biography != nil {
		u.Biography = biography
	}
	if birthDate != nil {
		u.BirthDate = birthDate
	}
	u.UpdatedAt = time.Now()
}

func (u *User) SetSocialID(provider string, socialID string) {
	switch provider {
	case "apple":
		u.AppleID = &socialID
	case "google":
		u.GoogleID = &socialID
	}
	u.UpdatedAt = time.Now()
}

func (u *User) HasSocialLogin() bool {
	return u.AppleID != nil || u.GoogleID != nil
}

func (u *User) IsCreator() bool {
	return u.UserType == UserTypeCreator
}

func (u *User) ValidateRequiredFields() error {
	// Creator-specific validations are now handled in Creator entity
	return nil
}

func (u *User) HasCreatorProfile() bool {
	return u.Creator != nil
}

// Follow count management methods
func (u *User) IncrementFollowersCount() {
	u.FollowersCount++
	u.UpdatedAt = time.Now()
}

func (u *User) DecrementFollowersCount() {
	if u.FollowersCount > 0 {
		u.FollowersCount--
	}
	u.UpdatedAt = time.Now()
}

func (u *User) IncrementFollowingCount() {
	u.FollowingCount++
	u.UpdatedAt = time.Now()
}

func (u *User) DecrementFollowingCount() {
	if u.FollowingCount > 0 {
		u.FollowingCount--
	}
	u.UpdatedAt = time.Now()
}

// Domain errors - removed creator-specific errors as they're now in Creator entity
var ()

type DomainError struct {
	Message string
}

func NewDomainError(message string) *DomainError {
	return &DomainError{Message: message}
}

func (e *DomainError) Error() string {
	return e.Message
}
