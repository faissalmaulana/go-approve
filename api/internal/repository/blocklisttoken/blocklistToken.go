package blocklisttoken

import (
	"context"

	"github.com/faissalmaulana/go-approve/internal/constant"
	"github.com/faissalmaulana/go-approve/internal/model"
	"gorm.io/gorm"
)

type BlocklistStorage interface {
	Create(ctx context.Context, token string) error
	FindByToken(ctx context.Context, token string) (*model.BloclistToken, error)
}

type BlocklistRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) BlocklistStorage {
	return &BlocklistRepository{
		db: db,
	}
}

func (r *BlocklistRepository) Create(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	blocklist := &model.BloclistToken{
		Token: token,
	}

	if err := gorm.G[model.BloclistToken](r.db).Create(ctx, blocklist); err != nil {
		return err
	}

	return nil
}

func (r *BlocklistRepository) FindByToken(ctx context.Context, token string) (*model.BloclistToken, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	blocklist, err := gorm.G[model.BloclistToken](r.db).Where("token = ?", token).First(ctx)
	if err != nil {
		return nil, err
	}

	return &blocklist, nil
}
