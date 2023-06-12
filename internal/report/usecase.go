package report

import (
	"context"

	"github.com/majidbl/discount/internal/models"
)

// UseCase Report usecase interface
type UseCase interface {
	Create(ctx context.Context, wallet *models.Report) error
	GetByGiftCode(ctx context.Context, giftCode string) ([]*models.Report, error)
	GetByMobile(ctx context.Context, mobile string) ([]*models.Report, error)
	GetGiftCodeCountUsage(ctx context.Context, usage string) (int64, error)
}
