package approvalroom

import (
	"context"
	"time"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalRoomStorage interface {
	Create(ctx context.Context,
		title string,
		dueAt time.Time,
		submitterId string,
	) error
	CreateWithTx(
		title string,
		dueAt time.Time,
		submitterId string,
		approvalRoomId *string,
	) func(ctx context.Context, tx *gorm.DB) error
	GetApprovalRoomByID(ctx context.Context, id string) (*ApprovalRoomDetail, error)
	GetApprovalRoomCountsByID(ctx context.Context, id string) (*ApprovalRoomCounts, error)
	GetApprovalRoomsBySubmitter(
		ctx context.Context,
		submitterId string,
		sortField string,
		orderDir string,
		limit int,
		offset int,
	) ([]model.ApprovalRoom, error)

	GetApprovalRoomsByApprover(
		ctx context.Context,
		approverId string,
		sortField string,
		orderDir string,
		limit int,
		offset int,
	) ([]ApproverInvitedApprovalRoomRequest, error)
	UpdateApprovalDecision(ctx context.Context, approvalRoomId, approvalId, decision string) error
}

type ApprovalRoomRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) ApprovalRoomStorage {
	return &ApprovalRoomRepository{
		db: db,
	}
}

func (ar *ApprovalRoomRepository) Create(
	ctx context.Context,
	title string,
	dueAt time.Time,
	submitterId string,
) error {

	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	newRoom := &model.ApprovalRoom{
		Title:       title,
		DueAt:       dueAt,
		SubmitterId: submitterId,
	}

	return gorm.G[model.ApprovalRoom](ar.db).Create(ctx, newRoom)
}

func (*ApprovalRoomRepository) CreateWithTx(
	title string,
	dueAt time.Time,
	submitterId string,
	approvalRoomId *string,
) func(ctx context.Context, tx *gorm.DB) error {
	roomID := uuid.New().String()
	*approvalRoomId = roomID

	newRoom := &model.ApprovalRoom{
		ID:          roomID,
		Title:       title,
		DueAt:       dueAt,
		SubmitterId: submitterId,
	}

	return func(ctx context.Context, tx *gorm.DB) error {
		ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
		defer cancel()

		return tx.Create(newRoom).Error
	}
}

type ApprovalRoomDetail struct {
	Room      model.ApprovalRoom
	Submitter model.User
	Files     []model.FileMetadata
	Approvers []ApprovalRoomApproverDetail
}

type ApprovalRoomApproverDetail struct {
	Handle   string
	Name     string
	Decision string
}

type ApprovalRoomCounts struct {
	FileUploaded  int64
	TotalApprover int64
	ApprovedCount int64
}

type ApproverInvitedApprovalRoomRequest struct {
	ID        string    `gorm:"column:id" json:"id"`
	Title     string    `gorm:"column:title" json:"title"`
	DueAt     time.Time `gorm:"column:due_at" json:"due_at"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	Decision  string    `gorm:"column:decision" json:"decision"`
}

func (ar *ApprovalRoomRepository) GetApprovalRoomByID(ctx context.Context, id string) (*ApprovalRoomDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var room model.ApprovalRoom
	if err := ar.db.WithContext(ctx).
		Preload("Submitter").
		Preload("Files").
		Preload("ApprovalRoomApprovers.Approval").
		First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}

	approvers := make([]ApprovalRoomApproverDetail, 0, len(room.ApprovalRoomApprovers))
	for _, ap := range room.ApprovalRoomApprovers {
		approvers = append(approvers, ApprovalRoomApproverDetail{
			Handle:   ap.Approval.Handler,
			Name:     ap.Approval.Name,
			Decision: ap.Decision,
		})
	}

	return &ApprovalRoomDetail{
		Room:      room,
		Submitter: room.Submitter,
		Files:     room.Files,
		Approvers: approvers,
	}, nil
}

func (ar *ApprovalRoomRepository) GetApprovalRoomCountsByID(ctx context.Context, id string) (*ApprovalRoomCounts, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var room model.ApprovalRoom
	if err := ar.db.WithContext(ctx).
		Preload("Files").
		Preload("ApprovalRoomApprovers").
		First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}

	counts := &ApprovalRoomCounts{
		FileUploaded:  int64(len(room.Files)),
		TotalApprover: int64(len(room.ApprovalRoomApprovers)),
	}

	var approved int64
	for _, ap := range room.ApprovalRoomApprovers {
		if ap.Decision == "approved" {
			approved++
		}
	}
	counts.ApprovedCount = approved

	return counts, nil
}

func (ar *ApprovalRoomRepository) GetApprovalRoomsBySubmitter(
	ctx context.Context,
	submitterId string,
	sortField string,
	orderDir string,
	limit int,
	offset int,
) ([]model.ApprovalRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	rooms := make([]model.ApprovalRoom, 0)
	orderBy := "created_at desc"
	switch sortField {
	case "due_at":
		orderBy = "due_at " + orderDir
	case "created_at":
		orderBy = "created_at " + orderDir
	}

	q := ar.db.WithContext(ctx).
		Select("id,title,due_at,created_at").
		Where("submitter_id = ?", submitterId).
		Order(orderBy)

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	err := q.Find(&rooms).Error

	return rooms, err
}

func (ar *ApprovalRoomRepository) GetApprovalRoomsByApprover(
	ctx context.Context,
	approverId string,
	sortField string,
	orderDir string,
	limit int,
	offset int,
) ([]ApproverInvitedApprovalRoomRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	orderBy := "rooms.created_at desc"
	switch sortField {
	case "due_at":
		orderBy = "rooms.due_at " + orderDir
	case "created_at":
		orderBy = "rooms.created_at " + orderDir
	}

	if orderDir != "asc" && orderDir != "desc" {
		orderBy = "rooms.created_at desc"
	}

	q := ar.db.WithContext(ctx).
		Table("approval_room_approvers ara").
		Select(
			"ara.approval_room_id as id, rooms.title, rooms.due_at, rooms.created_at, submitter.handler as created_by, ara.decision as decision",
		).
		Joins("JOIN approval_rooms rooms ON rooms.id = ara.approval_room_id").
		Joins("JOIN users submitter ON submitter.id = rooms.submitter_id").
		Where("ara.approval_id = ?", approverId).
		Order(orderBy)

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	rows := make([]ApproverInvitedApprovalRoomRequest, 0)
	err := q.Find(&rows).Error
	return rows, err
}

func (ar *ApprovalRoomRepository) UpdateApprovalDecision(ctx context.Context, approvalRoomId, approvalId, decision string) error {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	return ar.db.Model(&model.ApprovalRoomApprover{}).
		Where("approval_room_id = ? AND approval_id = ?", approvalRoomId, approvalId).
		Update("decision", decision).Error
}
