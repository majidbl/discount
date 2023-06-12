package usecase

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/discount/internal/giftcharge"
	"github.com/majidbl/discount/internal/models"
	"github.com/majidbl/discount/pkg/logger"
)

type giftChargeUseCase struct {
	log              logger.Logger
	giftChargePGRepo giftcharge.PGRepository
	redisRepo        giftcharge.RedisRepository
}

// NewGiftChargeUseCase giftCharge usecase constructor
func NewGiftChargeUseCase(
	log logger.Logger,
	giftChargePGRepo giftcharge.PGRepository,
	redisRepo giftcharge.RedisRepository,
) *giftChargeUseCase {
	return &giftChargeUseCase{log: log, giftChargePGRepo: giftChargePGRepo, redisRepo: redisRepo}
}

// Create new giftCharge saves in db
func (gu *giftChargeUseCase) Create(ctx context.Context, giftChargeReq *models.GiftChargeCreateReq) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargeUseCase.Create")
	defer span.Finish()

	giftChargeReq.Code = uuid.NewString()

	_, err := gu.giftChargePGRepo.Create(ctx, giftChargeReq)
	if err != nil {
		return errors.Wrap(err, "giftchargePGRepo.Create")
	}

	return nil
}

// GetByCode find giftCharge by code
func (gu *giftChargeUseCase) GetByCode(ctx context.Context, code string) (*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargeUseCase.GetByCode")
	defer span.Finish()

	cached, err := gu.redisRepo.GetGiftCharge(ctx, code)
	if err != nil && err != redis.Nil {
		gu.log.Errorf("redisRepo.GetGiftChargeByID: %v", err)
	}
	if cached != nil {
		return cached, nil
	}

	giftRecharge, err := gu.giftChargePGRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "giftChargePGRepo.GetByCode")
	}

	if err = gu.redisRepo.SetGiftCharge(ctx, giftRecharge); err != nil {
		gu.log.Errorf("redisRepo.SetGiftCharge: %v", err)
	}

	return giftRecharge, nil
}

// GetByID find giftCharge by id
func (gu *giftChargeUseCase) GetByID(ctx context.Context, id int64) (*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargeUseCase.GetByID")
	defer span.Finish()

	cached, err := gu.redisRepo.GetGiftCharge(ctx, strconv.FormatInt(id, 10))
	if err != nil && err != redis.Nil {
		gu.log.Errorf("redisRepo.GetGiftChargeByID: %v", err)
	}
	if cached != nil {
		return cached, nil
	}

	giftRecharge, err := gu.giftChargePGRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "giftChargePGRepo.GetByID")
	}

	if err = gu.redisRepo.SetGiftCharge(ctx, giftRecharge); err != nil {
		gu.log.Errorf("redisRepo.SetGiftCharge: %v", err)
	}

	return giftRecharge, nil
}

// GetList find giftCharge by id
func (gu *giftChargeUseCase) GetList(ctx context.Context) ([]*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargeUseCase.GetList")
	defer span.Finish()

	giftRecharge, err := gu.giftChargePGRepo.GetList(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "giftChargePGRepo.GetList")
	}

	return giftRecharge, nil
}

// GetValidList find list valid giftCharge
func (gu *giftChargeUseCase) GetValidList(ctx context.Context) ([]*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargeUseCase.GetValidList")
	defer span.Finish()

	giftRecharge, err := gu.giftChargePGRepo.GetValidList(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "giftChargePGRepo.GetValidList")
	}

	return giftRecharge, nil
}

// GetInValidList find list invalid giftCharge
func (gu *giftChargeUseCase) GetInValidList(ctx context.Context) ([]*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargeUseCase.GetInValidList")
	defer span.Finish()

	giftRecharge, err := gu.giftChargePGRepo.GetInValidList(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "giftChargePGRepo.GetInValidList")
	}

	return giftRecharge, nil
}
