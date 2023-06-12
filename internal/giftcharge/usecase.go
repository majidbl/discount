package giftcharge

import (
	"context"

	"github.com/majidbl/discount/internal/models"
)

// UseCase GiftCharge usecase interface
type UseCase interface {
	Create(ctx context.Context, giftCharge *models.GiftChargeCreateReq) error
	GetByID(ctx context.Context, id int64) (*models.GiftCharge, error)
	GetList(ctx context.Context) ([]*models.GiftCharge, error)
	GetValidList(ctx context.Context) ([]*models.GiftCharge, error)
	GetInValidList(ctx context.Context) ([]*models.GiftCharge, error)
	GetByCode(ctx context.Context, code string) (*models.GiftCharge, error)
}
