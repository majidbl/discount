package v1

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/majidbl/discount/internal/giftcharge"
	"github.com/majidbl/discount/internal/models"
	httpErrors "github.com/majidbl/discount/pkg/http_errors"
	"github.com/majidbl/discount/pkg/logger"
)

type giftHandlers struct {
	group    *echo.Group
	giftUC   giftcharge.UseCase
	log      logger.Logger
	validate *validator.Validate
}

// NewGiftChargeHandlers giftHandlers constructor
func NewGiftChargeHandlers(
	group *echo.Group,
	giftUC giftcharge.UseCase,
	log logger.Logger,
	validate *validator.Validate,
) *giftHandlers {
	return &giftHandlers{group: group, giftUC: giftUC, log: log, validate: validate}
}

// Create New GiftCharge
// @Tags GiftCharge
// @Summary Create new giftcharge
// @Description Create new giftcharge and send it
// @Accept json
// @Produce json
// @Success 200 {object} models.GiftCharge
// @Router /giftCharge [post]
func (h *giftHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "giftHandlers.Create")
		defer span.Finish()
		createRequests.Inc()

		var giftChargeRequest models.GiftChargeCreateReq
		if err := c.Bind(&giftChargeRequest); err != nil {
			errorRequests.Inc()
			h.log.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &giftChargeRequest); err != nil {
			errorRequests.Inc()
			h.log.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		err := h.giftUC.Create(ctx, &giftChargeRequest)
		if err != nil {
			span.SetTag("err", err)
			errorRequests.Inc()
			h.log.Errorf("giftUC.Create: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.NoContent(http.StatusCreated)
	}
}

// GetByID Get GiftCharge by ID
// @Tags GiftCharge
// @Summary Get giftCharge by id
// @Description Get giftCharge by giftCharge uuid
// @Accept json
// @Produce json
// @Param gift_id path string true "gift_id"
// @Success 200 {object} models.GiftCharge
// @Router /giftCharge/{gift_id} [get]
func (h *giftHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "giftHandlers.GetByID")
		defer span.Finish()
		getByIdRequests.Inc()

		giftID, err := strconv.ParseInt(c.Param("gift_id"), 10, 64)
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("strconv.ParseInt: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		m, err := h.giftUC.GetByID(ctx, giftID)
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("giftUC.GetByID: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, m)
	}
}

// GetList Get GiftCharge List
// @Tags GiftCharge
// @Summary Get giftCharge list information
// @Description Get giftCharge list
// @Accept json
// @Produce json
// @Param   isActive  query  boolean  false  "giftCharge status"  Enums(true, false)
// @Success 200 {object} []models.GiftCharge
// @Router /giftCharge [get]
func (h *giftHandlers) GetList() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "giftHandlers.GetList")
		defer span.Finish()
		getByIdRequests.Inc()

		var gifts []*models.GiftCharge
		var err error

		isValid := c.QueryParam("isValid")
		switch isValid {
		case "true":
			gifts, err = h.giftUC.GetValidList(ctx)
		case "false":
			gifts, err = h.giftUC.GetInValidList(ctx)
		default:
			gifts, err = h.giftUC.GetList(ctx)
		}

		if err != nil {
			errorRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, gifts)
	}
}
