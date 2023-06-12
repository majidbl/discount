package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/majidbl/discount/internal/report"
	httpErrors "github.com/majidbl/discount/pkg/http_errors"
	"github.com/majidbl/discount/pkg/logger"
)

type reportHandlers struct {
	group    *echo.Group
	reportUC report.UseCase
	log      logger.Logger
	validate *validator.Validate
}

// NewReportHandlers reportHandlers constructor
func NewReportHandlers(
	group *echo.Group,
	reportUC report.UseCase,
	log logger.Logger,
	validate *validator.Validate,
) *reportHandlers {
	return &reportHandlers{group: group, reportUC: reportUC, log: log, validate: validate}
}

// GetGiftCodeReport Get Gift Report by Gift Code
// @Tags Report
// @Summary Get reports by giftCode
// @Description return usage reports for a giftCode
// @Accept json
// @Produce json
// @Param gift_code path string true "gift_code"
// @Success 200 {object} []models.Report
// @Router /report/{gift_code} [get]
func (h *reportHandlers) GetGiftCodeReport() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "reportHandlers.GetGiftCodeReport")
		defer span.Finish()

		getByIdRequests.Inc()

		giftCode := c.Param("gift_code")

		m, err := h.reportUC.GetByGiftCode(ctx, giftCode)
		if err != nil {
			span.SetTag("err", err)
			errorRequests.Inc()
			h.log.Errorf("reportUC.GetByGiftCode: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, m)
	}
}

// GetMobileReport Get Gift Report by mobile
// @Tags Report
// @Summary Get reports by mobile
// @Description return usage reports for a mobile
// @Accept json
// @Produce json
// @Param gift_code path string true "gift_code"
// @Success 200 {object} []models.Report
// @Router /report/{mobile} [get]
func (h *reportHandlers) GetMobileReport() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "reportHandlers.GetMobileReport")
		defer span.Finish()

		getByIdRequests.Inc()

		mobile := c.Param("mobile")

		m, err := h.reportUC.GetByMobile(ctx, mobile)
		if err != nil {
			span.SetTag("err", err)
			errorRequests.Inc()
			h.log.Errorf("reportUC.GetByMobile: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, m)
	}
}
