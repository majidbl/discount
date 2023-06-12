package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/discount/internal/models"
	"github.com/majidbl/discount/internal/report"
	"github.com/majidbl/discount/internal/report/delivery/nats"
	"github.com/majidbl/discount/pkg/logger"
)

const (
	createReportSubject = "report:create"
)

type reportUseCase struct {
	log          logger.Logger
	reportPGRepo report.PGRepository
	publisher    nats.Publisher
	redisRepo    report.RedisRepository
}

// NewReportUseCase report usecase constructor
func NewReportUseCase(
	log logger.Logger,
	reportPGRepo report.PGRepository,
	publisher nats.Publisher,
	redisRepo report.RedisRepository,
) *reportUseCase {
	return &reportUseCase{
		log:          log,
		reportPGRepo: reportPGRepo,
		publisher:    publisher,
		redisRepo:    redisRepo,
	}
}

var _ report.UseCase = &reportUseCase{}

// Create creates a new report and saves it in the database
func (wu *reportUseCase) Create(ctx context.Context, report *models.Report) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportUseCase.Create")
	defer span.Finish()

	created, err := wu.reportPGRepo.Create(ctx, &models.CreateReportReq{
		GiftCode:   report.GiftCode,
		Mobile:     report.Mobile,
		Amount:     report.Amount,
		ReportTime: time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "reportPGRepo.Create")
	}

	reportBytes, err := json.Marshal(created)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return wu.publisher.Publish(createReportSubject, reportBytes)
}

// GetByGiftCode fnd report by id
func (wu *reportUseCase) GetByGiftCode(ctx context.Context, giftCode string) ([]*models.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportUseCase.GetByGiftCode")
	defer span.Finish()

	cached, err := wu.redisRepo.GetReport(ctx, giftCode)
	if err != nil && err != redis.Nil {
		wu.log.Errorf("redisRepo.GetByGiftCode: %v", err)
	}
	if cached != nil {
		return cached, nil
	}

	reports, err := wu.reportPGRepo.GetByGiftCode(ctx, giftCode)
	if err != nil {
		return nil, errors.Wrap(err, "reportPGRepo.GetByGiftCode")
	}

	if err = wu.redisRepo.SetReport(ctx, reports, giftCode); err != nil {
		wu.log.Errorf("redisRepo.SetReport: %v", err)
	}

	return reports, nil
}

// GetByMobile fnd report by mobile
func (wu *reportUseCase) GetByMobile(ctx context.Context, giftCode string) ([]*models.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportUseCase.GetByMobile")
	defer span.Finish()

	cached, err := wu.redisRepo.GetReport(ctx, giftCode)
	if err != nil && err != redis.Nil {
		wu.log.Errorf("redisRepo.GetByMobile: %v", err)
	}
	if cached != nil {
		return cached, nil
	}

	reports, err := wu.reportPGRepo.GetByMobile(ctx, giftCode)
	if err != nil {
		return nil, errors.Wrap(err, "reportPGRepo.GetByMobile")
	}

	if err = wu.redisRepo.SetReport(ctx, reports, giftCode); err != nil {
		wu.log.Errorf("redisRepo.SetReport: %v", err)
	}

	return reports, nil
}

// GetGiftCodeCountUsage return giftCode usage
func (wu *reportUseCase) GetGiftCodeCountUsage(ctx context.Context, giftCode string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportUseCase.GetGiftCodeCountUsage")
	defer span.Finish()

	cached, err := wu.redisRepo.GetCountUsage(ctx, giftCode)
	if err != nil && err != redis.Nil {
		wu.log.Errorf("redisRepo.GetByGiftCode: %v", err)
	}
	if cached != nil {
		return cached.CountUsage, nil
	}

	countUsage, err := wu.reportPGRepo.GetCountByGiftCode(ctx, giftCode)
	if err != nil {
		return 0, errors.Wrap(err, "reportPGRepo.GetByGiftCode")
	}

	if err = wu.redisRepo.SetCountUsage(ctx, &models.CountUsage{
		GiftCode:   giftCode,
		CountUsage: countUsage,
	}); err != nil {
		wu.log.Errorf("redisRepo.SetCountUsage: %v", err)
	}

	return countUsage, nil
}
