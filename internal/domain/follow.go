package domain

import (
	"time"
)

// Follow represents a follow relationship between users
type Follow struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	FollowerID  int       `json:"follower_id" gorm:"not null;index"`
	FollowingID int       `json:"following_id" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Follower  *User `json:"follower,omitempty" gorm:"foreignKey:FollowerID;references:ID"`
	Following *User `json:"following,omitempty" gorm:"foreignKey:FollowingID;references:ID"`
}

// TableName returns the table name for Follow entity
func (Follow) TableName() string {
	return "follows"
}

// NewFollow creates a new follow relationship
func NewFollow(followerID, followingID int) (*Follow, error) {
	follow := &Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := follow.Validate(); err != nil {
		return nil, err
	}

	return follow, nil
}

// Validate validates the follow relationship
func (f *Follow) Validate() error {
	if f.FollowerID <= 0 {
		return NewDomainError("follow.invalid.follower")
	}

	if f.FollowingID <= 0 {
		return NewDomainError("follow.invalid.following")
	}

	if f.FollowerID == f.FollowingID {
		return NewDomainError("follow.cannot.follow.self")
	}

	return nil
}

// IsValid checks if the follow relationship is valid
func (f *Follow) IsValid() bool {
	return f.Validate() == nil
}
