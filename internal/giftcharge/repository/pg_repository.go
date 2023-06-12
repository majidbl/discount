package repository

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/majidbl/discount/internal/giftcharge"
	"github.com/majidbl/discount/internal/models"
)

type giftChargePGRepository struct {
	db *pgxpool.Pool
}

// NewGiftChargePGRepository GiftCharge postgresql discount_repository.go constructor
func NewGiftChargePGRepository(db *pgxpool.Pool) *giftChargePGRepository {
	return &giftChargePGRepository{db: db}
}

var _ giftcharge.PGRepository = &giftChargePGRepository{}

// Create  new giftCharge
func (gr *giftChargePGRepository) Create(
	ctx context.Context,
	createReq *models.GiftChargeCreateReq,
) (*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftchargePGRepository.Create")
	defer span.Finish()

	var g models.GiftChargeModel
	// Prepare the INSERT statement with the RETURNING clause
	err := gr.db.QueryRow(ctx,
		createGiftChargeQuery,
		createReq.Code,
		createReq.ValidityPeriodStart,
		createReq.ValidityPeriodEnd,
		createReq.Amount,
		createReq.MaxUsageCount,
		time.Now(),
	).Scan(
		&g.ID,
		&g.Code,
		&g.ValidityPeriodStart,
		&g.ValidityPeriodEnd,
		&g.Amount,
		&g.MaxUsageCount,
		&g.CreatedAt,
		&g.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return g.Entity(), nil
}

// GetByID get single giftCharge by id
func (gr *giftChargePGRepository) GetByID(ctx context.Context, id int64) (*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargePGRepository.GetByID")
	defer span.Finish()

	var g models.GiftChargeModel
	if err := gr.db.QueryRow(ctx, getByIDQuery, id).
		Scan(
			&g.ID,
			&g.Code,
			&g.ValidityPeriodStart,
			&g.ValidityPeriodEnd,
			&g.Amount,
			&g.MaxUsageCount,
			&g.CreatedAt,
			&g.UpdatedAt,
		); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return g.Entity(), nil
}

// GetByCode get single giftCharge by id
func (gr *giftChargePGRepository) GetByCode(ctx context.Context, code string) (*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargePGRepository.GetByID")
	defer span.Finish()

	var g models.GiftChargeModel
	currentTime := time.Now()

	if err := gr.db.QueryRow(ctx, getByCodeQuery, code, currentTime).
		Scan(
			&g.ID,
			&g.Code,
			&g.ValidityPeriodStart,
			&g.ValidityPeriodEnd,
			&g.Amount,
			&g.MaxUsageCount,
			&g.CreatedAt,
			&g.UpdatedAt,
		); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return g.Entity(), nil
}

// GetList get list giftCharge
func (gr *giftChargePGRepository) GetList(ctx context.Context) ([]*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargePGRepository.GetList")
	defer span.Finish()

	rows, err := gr.db.Query(ctx, getListQuery)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var res []*models.GiftCharge

	for rows.Next() {
		var g models.GiftChargeModel
		scanErr := rows.Scan(
			&g.ID,
			&g.Code,
			&g.ValidityPeriodStart,
			&g.ValidityPeriodEnd,
			&g.Amount,
			&g.MaxUsageCount,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if scanErr != nil {
			return nil, errors.Wrap(scanErr, "Scan")
		}

		res = append(res, g.Entity())
	}

	return res, nil
}

// GetValidList get list valid giftCharge
func (gr *giftChargePGRepository) GetValidList(ctx context.Context) ([]*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargePGRepository.GetList")
	defer span.Finish()

	currentTime := time.Now()

	rows, err := gr.db.Query(ctx, getValidListQuery, currentTime, currentTime)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var res []*models.GiftCharge

	for rows.Next() {
		var g models.GiftChargeModel
		scanErr := rows.Scan(
			&g.ID,
			&g.Code,
			&g.ValidityPeriodStart,
			&g.ValidityPeriodEnd,
			&g.Amount,
			&g.MaxUsageCount,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if scanErr != nil {
			return nil, errors.Wrap(scanErr, "Scan")
		}

		res = append(res, g.Entity())
	}

	return res, nil
}

// GetInValidList get list of inValid giftCharge
func (gr *giftChargePGRepository) GetInValidList(ctx context.Context) ([]*models.GiftCharge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargePGRepository.GetList")
	defer span.Finish()

	currentTime := time.Now()

	rows, err := gr.db.Query(ctx, getInValidListQuery, currentTime, currentTime)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var res []*models.GiftCharge

	for rows.Next() {
		var g models.GiftChargeModel
		scanErr := rows.Scan(
			&g.ID,
			&g.Code,
			&g.ValidityPeriodStart,
			&g.ValidityPeriodEnd,
			&g.Amount,
			&g.MaxUsageCount,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if scanErr != nil {
			return nil, errors.Wrap(scanErr, "Scan")
		}

		res = append(res, g.Entity())
	}

	return res, nil
}

// GetByCodeX get single giftCharge by id in a transaction
func (gr *giftChargePGRepository) GetByCodeX(ctx context.Context, code string, tx pgx.Tx) (*models.GiftCharge, pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "giftChargePGRepository.GetByCodeX")
	defer span.Finish()

	if tx == nil {
		var beginTxErr error

		tx, beginTxErr = gr.db.Begin(ctx)
		if beginTxErr != nil {
			return nil, nil, beginTxErr
		}
	}

	var g models.GiftChargeModel
	currentTime := time.Now()

	if err := gr.db.QueryRow(
		ctx,
		getByCodeQuery,
		code,
		currentTime,
	).
		Scan(
			&g.ID,
			&g.Code,
			&g.ValidityPeriodStart,
			&g.ValidityPeriodEnd,
			&g.Amount,
			&g.MaxUsageCount,
			&g.CreatedAt,
			&g.UpdatedAt,
		); err != nil {
		return nil, tx, errors.Wrap(err, "Scan")
	}

	return g.Entity(), tx, nil
}
