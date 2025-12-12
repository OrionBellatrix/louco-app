package domain

import (
	"time"
)

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusApproved InvitationStatus = "approved"
	InvitationStatusRejected InvitationStatus = "rejected"
)

type Invitation struct {
	ID            int              `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID       int              `json:"event_id" gorm:"not null;index"`
	InvitedUserID *int             `json:"invited_user_id" gorm:"index"` // nullable for non-members
	InvitedEmail  string           `json:"invited_email" gorm:"type:varchar(255);not null;index"`
	Status        InvitationStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	InvitedAt     time.Time        `json:"invited_at" gorm:"autoCreateTime"`
	RespondedAt   *time.Time       `json:"responded_at"`
	CreatedAt     time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Event       Event `json:"event" gorm:"foreignKey:EventID;references:ID"`
	InvitedUser *User `json:"invited_user,omitempty" gorm:"foreignKey:InvitedUserID;references:ID"`
}

func NewInvitation(eventID int, invitedEmail string, invitedUserID *int) *Invitation {
	return &Invitation{
		EventID:       eventID,
		InvitedEmail:  invitedEmail,
		InvitedUserID: invitedUserID,
		Status:        InvitationStatusPending,
		InvitedAt:     time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (i *Invitation) Approve() error {
	if i.Status != InvitationStatusPending {
		return ErrInvitationInvalidStatusTransition
	}

	i.Status = InvitationStatusApproved
	now := time.Now()
	i.RespondedAt = &now
	i.UpdatedAt = now
	return nil
}

func (i *Invitation) Reject() error {
	if i.Status != InvitationStatusPending {
		return ErrInvitationInvalidStatusTransition
	}

	i.Status = InvitationStatusRejected
	now := time.Now()
	i.RespondedAt = &now
	i.UpdatedAt = now
	return nil
}

func (i *Invitation) ResetToPending() error {
	if i.Status == InvitationStatusPending {
		return ErrInvitationAlreadyPending
	}

	i.Status = InvitationStatusPending
	i.RespondedAt = nil
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Invitation) UpdateInvitedUser(userID int) {
	i.InvitedUserID = &userID
	i.UpdatedAt = time.Now()
}

func (i *Invitation) RemoveInvitedUser() {
	i.InvitedUserID = nil
	i.UpdatedAt = time.Now()
}

func (i *Invitation) UpdateEmail(email string) {
	i.InvitedEmail = email
	i.UpdatedAt = time.Now()
}

// Status check methods
func (i *Invitation) IsPending() bool {
	return i.Status == InvitationStatusPending
}

func (i *Invitation) IsApproved() bool {
	return i.Status == InvitationStatusApproved
}

func (i *Invitation) IsRejected() bool {
	return i.Status == InvitationStatusRejected
}

func (i *Invitation) HasResponded() bool {
	return i.RespondedAt != nil
}

func (i *Invitation) IsForSystemUser() bool {
	return i.InvitedUserID != nil
}

func (i *Invitation) IsForExternalUser() bool {
	return i.InvitedUserID == nil
}

func (i *Invitation) CanBeDeleted() bool {
	// Invitations can be deleted if they haven't been approved yet
	return i.Status != InvitationStatusApproved
}

func (i *Invitation) CanBeResent() bool {
	// Invitations can be resent if they are pending or rejected
	return i.Status == InvitationStatusPending || i.Status == InvitationStatusRejected
}

func (i *Invitation) GetResponseTime() *time.Duration {
	if i.RespondedAt == nil {
		return nil
	}

	duration := i.RespondedAt.Sub(i.InvitedAt)
	return &duration
}

func (i *Invitation) ValidateRequiredFields() error {
	if i.InvitedEmail == "" {
		return ErrInvitationEmailRequired
	}
	if i.EventID <= 0 {
		return ErrInvitationEventRequired
	}
	return nil
}

func (i *Invitation) ValidateEmail() error {
	// Basic email validation - in real implementation, use proper email validation
	if i.InvitedEmail == "" {
		return ErrInvitationEmailRequired
	}
	// Add more sophisticated email validation here if needed
	return nil
}

func (i *Invitation) IsExpired(expirationHours int) bool {
	if expirationHours <= 0 {
		return false // No expiration
	}

	expirationTime := i.InvitedAt.Add(time.Duration(expirationHours) * time.Hour)
	return time.Now().After(expirationTime) && i.IsPending()
}

func (i *Invitation) GetDaysUntilExpiration(expirationHours int) int {
	if expirationHours <= 0 {
		return -1 // No expiration
	}

	expirationTime := i.InvitedAt.Add(time.Duration(expirationHours) * time.Hour)
	if time.Now().After(expirationTime) {
		return 0 // Already expired
	}

	duration := expirationTime.Sub(time.Now())
	return int(duration.Hours() / 24)
}

// Helper method to check if invitation belongs to a specific user
func (i *Invitation) BelongsToUser(userID int) bool {
	return i.InvitedUserID != nil && *i.InvitedUserID == userID
}

// Helper method to check if invitation is for a specific email
func (i *Invitation) IsForEmail(email string) bool {
	return i.InvitedEmail == email
}

// Invitation domain errors
var (
	ErrInvitationEmailRequired           = NewDomainError("invited email is required")
	ErrInvitationEventRequired           = NewDomainError("event ID is required")
	ErrInvitationInvalidStatusTransition = NewDomainError("invalid invitation status transition")
	ErrInvitationAlreadyPending          = NewDomainError("invitation is already pending")
	ErrInvitationNotFound                = NewDomainError("invitation not found")
	ErrInvitationAlreadyResponded        = NewDomainError("invitation has already been responded to")
	ErrInvitationExpired                 = NewDomainError("invitation has expired")
	ErrInvitationCannotBeDeleted         = NewDomainError("invitation cannot be deleted")
	ErrInvitationCannotBeResent          = NewDomainError("invitation cannot be resent")
	ErrInvitationUnauthorized            = NewDomainError("unauthorized to access this invitation")
	ErrInvitationDuplicateEmail          = NewDomainError("invitation already exists for this email")
)
