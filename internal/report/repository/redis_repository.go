package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/discount/internal/models"
	"github.com/majidbl/discount/internal/report"
)

const (
	prefix          = "reports"
	codeUsagePrefix = "code"
	expiration      = time.Second * 3600
)

type reportRedisRepository struct {
	redis *redis.Client
}

// NewReportRedisRepository reports redis discount_repository.go constructor
func NewReportRedisRepository(redis *redis.Client) *reportRedisRepository {
	return &reportRedisRepository{redis: redis}
}

var _ report.RedisRepository = &reportRedisRepository{}

func (e *reportRedisRepository) SetReport(ctx context.Context, report []*models.Report, giftCode string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportRedisRepository.SetReport")
	defer span.Finish()

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return errors.Wrap(err, "reportRedisRepository.Marshal.SetReport")
	}

	return e.redis.Set(ctx, e.createReportKey(giftCode), string(reportBytes), expiration).Err()
}

func (e *reportRedisRepository) GetReport(ctx context.Context, code string) ([]*models.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportRedisRepository.GetReport")
	defer span.Finish()

	result, err := e.redis.Get(ctx, e.createReportKey(code)).Bytes()
	if err != nil {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, errors.Wrap(err, "reportRedisRepository.redis.GetReport")
	}

	var res []*models.Report
	if err = json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return res, nil
}

func (e *reportRedisRepository) SetCountUsage(ctx context.Context, usage *models.CountUsage) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportRedisRepository.SetCountUsage")
	defer span.Finish()

	usageBytes, err := json.Marshal(usage)
	if err != nil {
		return errors.Wrap(err, "reportRedisRepository.Marshal.SetCountUsage")
	}

	return e.redis.Set(ctx, e.createReportKey(usage.GiftCode), string(usageBytes), expiration).Err()
}

func (e *reportRedisRepository) GetCountUsage(ctx context.Context, giftCode string) (*models.CountUsage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportRedisRepository.GetReport")
	defer span.Finish()

	result, err := e.redis.Get(ctx, e.createCountKey(giftCode)).Bytes()
	if err != nil {
		if redis.Nil == err {
			return nil, nil
		}
		return nil, errors.Wrap(err, "reportRedisRepository.redis.GetCountUsage")
	}

	var res models.CountUsage
	if err = json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return &res, nil
}

func (e *reportRedisRepository) createCountKey(giftCode string) string {
	return fmt.Sprintf("%s: %s", codeUsagePrefix, giftCode)
}

func (e *reportRedisRepository) createReportKey(code string) string {
	return fmt.Sprintf("%s: %s", prefix, code)
}
