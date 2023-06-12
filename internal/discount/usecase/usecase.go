package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/discount/internal/discount"
	"github.com/majidbl/discount/internal/discount/delivery/nats"
	"github.com/majidbl/discount/internal/giftcharge"
	"github.com/majidbl/discount/internal/models"
	"github.com/majidbl/discount/internal/report"
	"github.com/majidbl/discount/pkg/logger"
	"github.com/majidbl/discount/pkg/sql_errors"
	walletService "github.com/majidbl/discount/proto/wallet"
)

const (
	createDiscountSubject = "discount:create"
)

type discountUseCase struct {
	log                 logger.Logger
	reportPGRepo        report.PGRepository
	giftChargePGRepo    giftcharge.PGRepository
	giftChargeRedisRepo giftcharge.RedisRepository
	walletGrpcClient    walletService.WalletServiceClient
	publisher           nats.Publisher
	redisRepo           discount.RedisRepository
}

// NewDiscountUseCase report usecase constructor
func NewDiscountUseCase(
	log logger.Logger,
	reportPGRepo report.PGRepository,
	walletGrpcClient walletService.WalletServiceClient,
	publisher nats.Publisher,
	redisRepo discount.RedisRepository,
	giftChargePGRepo giftcharge.PGRepository,
	giftChargeRedisRepo giftcharge.RedisRepository,
) *discountUseCase {
	return &discountUseCase{
		log:                 log,
		reportPGRepo:        reportPGRepo,
		walletGrpcClient:    walletGrpcClient,
		giftChargePGRepo:    giftChargePGRepo,
		giftChargeRedisRepo: giftChargeRedisRepo,
		publisher:           publisher,
		redisRepo:           redisRepo,
	}
}

var _ discount.UseCase = &discountUseCase{}

// DiscountRequest creates a new discount and saves it in the database
func (du *discountUseCase) DiscountRequest(ctx context.Context, dis *models.Discount) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "discountUseCase.Create")
	defer span.Finish()

	var err error
	var tx pgx.Tx

	defer func() {
		if err != nil && tx != nil {
			if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
				du.log.Errorf("tx.Rollback: %v", rollBackErr)
			}
		}

		if err == nil && tx != nil {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				du.log.Errorf("tx.Commit: %v", commitErr)
			}
		}

		if delErr := du.redisRepo.DeleteDiscount(ctx, dis); delErr != nil {
			du.log.Errorf("redisRepo.DeleteDiscount: %v", delErr)
		}
	}()

	ok, err := du.redisRepo.DiscountExist(ctx, dis)
	if err != nil {
		return errors.Wrap(err, "discountUsecase.DiscountRequest")
	}

	if ok {
		return fmt.Errorf("you have one request in progress")
	}

	if err = du.redisRepo.SetDiscount(ctx, dis); err != nil {
		du.log.Errorf("redisRepo.SetDiscount: %v", err)
	}

	// check for it gift code availability. only valid gift code that available in request time returned
	giftChargeObject, tx, err := du.checkGiftCode(ctx, dis, nil)
	if err != nil {
		return errors.Wrap(err, "discountUsecase.DiscountRequest")
	}

	// check for gift code usage allowed
	giftUsage, tx, err := du.reportPGRepo.CheckGiftUsageX(ctx, &models.CheckGiftUsage{
		GiftCode: dis.GiftCode,
		Mobile:   dis.Mobile,
	},
		tx)

	if err != nil {
		switch err.(type) {
		case nil:
		case *sql_errors.SqlNotFoundError:
		default:
			return errors.Wrap(err, "reportPGRepo.CheckGiftUsage")
		}
	}

	if giftUsage != nil {
		return fmt.Errorf("you received this gift before")
	}

	// send charge request to wallet
	_, err = du.walletGrpcClient.Charge(ctx, &walletService.ChargeReq{
		Mobile: dis.Mobile,
		Amount: giftChargeObject.Amount,
	})

	if err != nil {
		return errors.Wrap(err, "walletGrpcClient.Charge")
	}

	_, tx, err = du.reportPGRepo.CreateX(ctx, &models.CreateReportReq{
		GiftCode:   dis.GiftCode,
		Mobile:     dis.Mobile,
		Amount:     giftChargeObject.Amount,
		ReportTime: time.Now(),
	},
		tx)

	if err != nil {
		return errors.Wrap(err, "reportPGRepo.CreateX")
	}

	discountBytes, err := json.Marshal(giftChargeObject)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return du.publisher.Publish(createDiscountSubject, discountBytes)
}

func (du *discountUseCase) checkGiftCode(
	ctx context.Context,
	dis *models.Discount,
	tx pgx.Tx,
) (*models.GiftCharge, pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "discountUseCase.checkGiftCode")
	defer span.Finish()

	cached, err := du.giftChargeRedisRepo.GetGiftCharge(ctx, dis.GiftCode)
	if err != nil {
		du.log.Errorf("giftChargeRedisRepo.GetGiftCharge: %v", err)
	}

	if cached != nil {
		return cached, tx, nil
	}

	giftChargeObject, tx, err := du.giftChargePGRepo.GetByCodeX(ctx, dis.GiftCode, nil)
	if err != nil {
		return nil, tx, errors.Wrap(err, "discountUsecase.DiscountRequest")
	}

	if err = du.giftChargeRedisRepo.SetGiftCharge(ctx, giftChargeObject); err != nil {
		du.log.Errorf("giftChargeRedisRepo.SetGiftCharge: %v", err)
	}

	return giftChargeObject, tx, nil

}
