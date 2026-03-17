package filemetadata

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"gorm.io/gorm"
)

type FileMetadataStorage interface {
	CreateBatchWithTx(
		fileMetadataMap map[string]model.FileMetadata,
		approvalRoomId string,
	) func(ctx context.Context, tx *gorm.DB) error
}

type FileMetadataRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) FileMetadataStorage {
	return &FileMetadataRepository{
		db: db,
	}
}

func (f *FileMetadataRepository) CreateBatchWithTx(
	fileMetadataMap map[string]model.FileMetadata,
	approvalRoomId string,
) func(ctx context.Context, tx *gorm.DB) error {
	fileMetadatas := make([]*model.FileMetadata, 0, len(fileMetadataMap))
	for _, fm := range fileMetadataMap {
		fileMetadatas = append(fileMetadatas, &model.FileMetadata{
			Link:           fm.Link,
			Filename:       fm.Filename,
			ApprovalRoomId: approvalRoomId,
			Size:           fm.Size,
		})
	}

	return func(ctx context.Context, tx *gorm.DB) error {
		ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
		defer cancel()

		return tx.Create(fileMetadatas).Error
	}
}
