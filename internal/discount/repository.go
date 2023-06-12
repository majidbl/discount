package discount

import (
	"context"

	"github.com/majidbl/discount/internal/models"
)

// RedisRepository redis discount discount_repository.go interface
type RedisRepository interface {
	SetDiscount(ctx context.Context, req *models.Discount) error
	DiscountExist(ctx context.Context, req *models.Discount) (bool, error)
	DeleteDiscount(ctx context.Context, req *models.Discount) error
}
