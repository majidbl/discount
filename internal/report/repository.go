package report

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/majidbl/discount/internal/models"
)

// PGRepository defines the methods for interacting with the report data.
type PGRepository interface {
	Create(ctx context.Context, createReq *models.CreateReportReq) (*models.Report, error)
	CreateX(ctx context.Context, createReq *models.CreateReportReq, tx pgx.Tx) (*models.Report, pgx.Tx, error)
	GetByGiftCode(ctx context.Context, giftCode string) ([]*models.Report, error)
	GetByMobile(ctx context.Context, giftCode string) ([]*models.Report, error)
	GetCountByGiftCode(ctx context.Context, mobile string) (int64, error)
	GetCountByGiftCodeX(ctx context.Context, mobile string, tx pgx.Tx) (int64, pgx.Tx, error)
	CheckGiftUsage(ctx context.Context, usage *models.CheckGiftUsage) (*models.Report, error)
	CheckGiftUsageX(ctx context.Context, usage *models.CheckGiftUsage, tx pgx.Tx) (*models.Report, pgx.Tx, error)
	RollBack(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}

// RedisRepository redis report discount_repository.go interface
type RedisRepository interface {
	SetReport(ctx context.Context, wallet []*models.Report, giftCode string) error
	GetReport(ctx context.Context, giftCode string) ([]*models.Report, error)
	SetCountUsage(ctx context.Context, usage *models.CountUsage) error
	GetCountUsage(ctx context.Context, giftCode string) (*models.CountUsage, error)
}
