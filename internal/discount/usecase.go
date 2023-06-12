package discount

import (
	"context"

	"github.com/majidbl/discount/internal/models"
)

// UseCase Discount usecase interface
type UseCase interface {
	DiscountRequest(ctx context.Context, req *models.Discount) error
}
