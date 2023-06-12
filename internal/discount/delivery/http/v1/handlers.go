package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/majidbl/discount/internal/discount"
	"github.com/majidbl/discount/internal/models"
	httpErrors "github.com/majidbl/discount/pkg/http_errors"
	"github.com/majidbl/discount/pkg/logger"
)

type discountHandlers struct {
	group      *echo.Group
	discountUC discount.UseCase
	log        logger.Logger
	validate   *validator.Validate
}

// NewDiscountHandlers discountHandlers constructor
func NewDiscountHandlers(
	group *echo.Group,
	discountUC discount.UseCase,
	log logger.Logger,
	validate *validator.Validate,
) *discountHandlers {
	return &discountHandlers{group: group, discountUC: discountUC, log: log, validate: validate}
}

// DiscountRequest Get Gift Discount by Gift Code
// @Tags Discount
// @Summary Get discounts by giftCode
// @Description return usage discounts for a giftCode
// @Accept json
// @Produce json
// @Param gift_code path string true "gift_code"
// @Success 204
// @Router /discount [post]
func (h *discountHandlers) DiscountRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(
			c.Request().Context(),
			"discountHandlers.GetGiftCodeDiscount")
		defer span.Finish()

		createRequests.Inc()

		var discountRequest models.Discount
		if err := c.Bind(&discountRequest); err != nil {
			errorRequests.Inc()
			h.log.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &discountRequest); err != nil {
			errorRequests.Inc()
			h.log.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		err := h.discountUC.DiscountRequest(ctx, &discountRequest)
		if err != nil {
			span.SetTag("err", err)
			errorRequests.Inc()
			h.log.Errorf("discountUC.DiscountRequest: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, "ok")
	}
}
