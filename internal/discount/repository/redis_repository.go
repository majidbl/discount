package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/discount/internal/discount"
	"github.com/majidbl/discount/internal/models"
)

const (
	prefix          = "discounts"
	codeUsagePrefix = "code"
	expiration      = time.Second * 3600
)

type discountRedisRepository struct {
	redis *redis.Client
}

// NewDiscountRedisRepository discounts redis discount_repository.go constructor
func NewDiscountRedisRepository(redis *redis.Client) *discountRedisRepository {
	return &discountRedisRepository{redis: redis}
}

var _ discount.RedisRepository = &discountRedisRepository{}

func (e *discountRedisRepository) SetDiscount(ctx context.Context, discount *models.Discount) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "discountRedisRepository.SetDiscount")
	defer span.Finish()

	discountBytes, err := json.Marshal(discount)
	if err != nil {
		return errors.Wrap(err, "discountRedisRepository.Marshal.SetDiscount")
	}

	return e.redis.SetNX(ctx, e.createKey(discount), string(discountBytes), expiration).Err()
}

func (e *discountRedisRepository) DiscountExist(ctx context.Context, dis *models.Discount) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "discountRedisRepository.ExistDiscount")
	defer span.Finish()

	result, err := e.redis.Get(ctx, e.createKey(dis)).Bytes()
	if err != nil {
		if redis.Nil == err {
			return false, nil
		}
		return false, errors.Wrap(err, "discountRedisRepository.ExistDiscount")
	}

	var res *models.Discount
	if err = json.Unmarshal(result, &res); err != nil {
		return false, errors.Wrap(err, "json.Unmarshal")
	}

	return res.GiftCode == dis.GiftCode && res.Mobile == dis.Mobile, nil
}

func (e *discountRedisRepository) DeleteDiscount(ctx context.Context, dis *models.Discount) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "discountRedisRepository.DeleteDiscount")
	defer span.Finish()
	return e.redis.Del(ctx, e.createKey(dis)).Err()
}

func (e *discountRedisRepository) createKey(d *models.Discount) string {
	return fmt.Sprintf("%s:%s:%s", prefix, d.Mobile, d.GiftCode)
}
