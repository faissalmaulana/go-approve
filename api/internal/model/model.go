package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	Handler  string `gorm:"type:varchar(50);unique;not null"`
	Email    string `gorm:"type:varchar(255);unique;not null"`
	Password string `gorm:"type:varchar(60);not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}

type BlocklistToken struct {
	Token string `gorm:"primaryKey;type:varchar(255);not null"`
}

type ApprovalRoom struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Filepaths   string     `json:"filepaths"` // this is the raw form which is each filepath is seperated with ";"
	DueAt       time.Time  `json:"due_at"`
	SubmitterId string     `json:"submitter_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (a *ApprovalRoom) BeforeCreate(tx *gorm.DB) error {
	// somehow in the repository they create a new one manually
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

type ApprovalRoomApprover struct {
	ApprovalId     string `json:"approval_id"`
	ApprovalRoomId string `json:"approval_room_id"`
	Decision       string `json:"decision"`
}

type ReviewRequest struct {
	ID             string    `json:"id"`
	IsRead         bool      `json:"is_read"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	ApprovalRoomId string    `json:"approval_room_id"`
	InviteeId      string    `json:"invitee_id"`
	RequesterId    string    `json:"requester_id"`
}

func (r *ReviewRequest) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New().String()
	return nil
}
