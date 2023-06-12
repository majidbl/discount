package giftcharge

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/majidbl/discount/internal/models"
)

// PGRepository GiftCharge postgresql discount_repository.go interface
type PGRepository interface {
	Create(ctx context.Context, giftCharge *models.GiftChargeCreateReq) (*models.GiftCharge, error)
	GetByID(ctx context.Context, id int64) (*models.GiftCharge, error)
	GetByCode(ctx context.Context, code string) (*models.GiftCharge, error)
	GetByCodeX(ctx context.Context, code string, tx pgx.Tx) (*models.GiftCharge, pgx.Tx, error)
	GetList(ctx context.Context) ([]*models.GiftCharge, error)
	GetValidList(ctx context.Context) ([]*models.GiftCharge, error)
	GetInValidList(ctx context.Context) ([]*models.GiftCharge, error)
}

// RedisRepository redis giftCharge discount_repository.go interface
type RedisRepository interface {
	SetGiftCharge(ctx context.Context, giftCharge *models.GiftCharge) error
	GetGiftCharge(ctx context.Context, key string) (*models.GiftCharge, error)
	DeleteGiftCharge(ctx context.Context, id string) error
}
