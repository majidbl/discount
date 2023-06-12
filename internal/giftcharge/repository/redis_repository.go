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
)

const (
	prefix     = "gifts"
	expiration = time.Second * 3600
)

type giftChargeRedisRepository struct {
	redis *redis.Client
}

// NewGiftChargeRedisRepository giftCharges redis discount_repository.go constructor
func NewGiftChargeRedisRepository(redis *redis.Client) *giftChargeRedisRepository {
	return &giftChargeRedisRepository{redis: redis}
}

func (e *giftChargeRedisRepository) SetGiftCharge(ctx context.Context, payload *models.GiftCharge) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftchargeRedisRepository.SetGiftCharge")
	defer span.Finish()

	giftChargeBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "giftChargeRedisRepository.Marshal")
	}

	currentTime := time.Now()
	return e.redis.SetEX(
		ctx,
		e.createKey(payload.Code),
		string(giftChargeBytes),
		payload.ValidityPeriodEnd.Sub(currentTime)).Err()
}

func (e *giftChargeRedisRepository) GetGiftCharge(ctx context.Context, gId string) (*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargeRedisRepository.GetGiftCharge")
	defer span.Finish()

	result, err := e.redis.Get(ctx, e.createKey(gId)).Bytes()
	if err != nil {
		if redis.Nil == err {
			return nil, nil
		}

		return nil, errors.Wrap(err, "giftChargeRedisRepository.redis.Get")
	}

	var res models.GiftCharge
	if err = json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return &res, nil
}

func (e *giftChargeRedisRepository) DeleteGiftCharge(ctx context.Context, code string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftchargeRedisRepository.DeleteGiftCharge")
	defer span.Finish()
	return e.redis.Del(ctx, e.createKey(code)).Err()
}

func (e *giftChargeRedisRepository) createKey(id string) string {
	return fmt.Sprintf("%s: %s", prefix, id)
}
