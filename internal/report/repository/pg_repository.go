package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/majidbl/discount/internal/report"
	"github.com/majidbl/discount/pkg/sql_errors"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/discount/internal/models"
)

type reportPGRepository struct {
	db *pgxpool.Pool
}

// NewReportPGRepository Report postgresql discount_repository.go constructor
func NewReportPGRepository(db *pgxpool.Pool) *reportPGRepository {
	return &reportPGRepository{db: db}
}

var _ report.PGRepository = &reportPGRepository{}

// Create  new report
func (rr *reportPGRepository) Create(ctx context.Context, createReq *models.CreateReportReq) (*models.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.Create")
	defer span.Finish()

	var w models.Report
	if err := rr.db.QueryRow(
		ctx,
		createReportQuery,
		&createReq.Mobile,
		&createReq.GiftCode,
		&createReq.Amount,
		&createReq.ReportTime,
		time.Now(),
	).Scan(
		&w.Id,
		&w.GiftCode,
		&w.Mobile,
		&w.Amount,
		&w.ReportTime,
		&w.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &w, nil
}

// CreateX  new report
func (rr *reportPGRepository) CreateX(ctx context.Context, createReq *models.CreateReportReq, tx pgx.Tx) (*models.Report, pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.CreateX")
	defer span.Finish()

	if tx == nil {
		var beginTxErr error

		tx, beginTxErr = rr.db.Begin(ctx)
		if beginTxErr != nil {
			return nil, nil, beginTxErr
		}
	}

	var w models.Report
	if err := tx.QueryRow(
		ctx,
		createReportQuery,
		&createReq.Mobile,
		&createReq.GiftCode,
		&createReq.Amount,
		&createReq.ReportTime,
		time.Now(),
	).Scan(
		&w.Id,
		&w.GiftCode,
		&w.Mobile,
		&w.Amount,
		&w.ReportTime,
		&w.CreatedAt,
	); err != nil {
		return nil, tx, errors.Wrap(err, "Scan")
	}

	return &w, tx, nil
}

// GetByGiftCode get all report by gift code
func (rr *reportPGRepository) GetByGiftCode(ctx context.Context, giftCode string) ([]*models.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.GetByGiftCode")
	defer span.Finish()

	var reports []*models.Report
	rows, err := rr.db.Query(ctx, getByGiftCodeQuery, giftCode)
	if err != nil {
		return nil, errors.Wrap(err, "Query")
	}

	for rows.Next() {
		var r models.Report
		err = rows.Scan(
			&r.Id,
			&r.GiftCode,
			&r.Mobile,
			&r.Amount,
			&r.ReportTime,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "Scan")
		}

		reports = append(reports, &r)
	}

	return reports, nil
}

// GetByMobile get all report by gift code
func (rr *reportPGRepository) GetByMobile(ctx context.Context, mobile string) ([]*models.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.GetByMobile")
	defer span.Finish()

	var reports []*models.Report
	rows, err := rr.db.Query(ctx, getByMobileQuery, mobile)
	if err != nil {
		return nil, errors.Wrap(err, "Query")
	}

	for rows.Next() {
		var r models.Report
		err = rows.Scan(
			&r.Id,
			&r.GiftCode,
			&r.Mobile,
			&r.Amount,
			&r.ReportTime,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "Scan")
		}

		reports = append(reports, &r)
	}

	return reports, nil
}

func (rr *reportPGRepository) CheckGiftUsageX(ctx context.Context, usage *models.CheckGiftUsage, tx pgx.Tx) (*models.Report, pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.CheckGiftUsageX")
	defer span.Finish()

	if tx == nil {
		var beginTxErr error

		tx, beginTxErr = rr.db.Begin(ctx)
		if beginTxErr != nil {
			return nil, nil, beginTxErr
		}
	}

	var r models.ReportModel
	err := tx.QueryRow(ctx, getByGiftCodeAndMobileQuery, usage.Mobile, usage.GiftCode).
		Scan(
			&r.Id,
			&r.GiftCode,
			&r.Mobile,
			&r.Amount,
			&r.ReportTime,
			&r.CreatedAt,
		)
	if err != nil {
		return nil, tx, sql_errors.ParseSqlErrors(err)
	}

	return r.Entity(), tx, nil
}

func (rr *reportPGRepository) CheckGiftUsage(ctx context.Context, usage *models.CheckGiftUsage) (*models.Report, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.CheckGiftUsage")
	defer span.Finish()

	rows, err := rr.db.Query(ctx, getByGiftCodeAndMobileQuery, usage.GiftCode, usage.Mobile)
	if err != nil {
		return nil, errors.Wrap(err, "Query")
	}

	var r models.Report
	for rows.Next() {
		err = rows.Scan(
			&r.Id,
			&r.GiftCode,
			&r.Mobile,
			&r.Amount,
			&r.ReportTime,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "Scan")
		}

	}

	return &r, nil
}

// GetCountByGiftCode get count usage report
func (rr *reportPGRepository) GetCountByGiftCode(ctx context.Context, mobile string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.GetCountByGiftCode")
	defer span.Finish()

	var totalUsage sql.NullInt64
	if err := rr.db.QueryRow(ctx, getcountByGiftCodeQuery, mobile).
		Scan(
			&totalUsage,
		); err != nil {
		return 0, errors.Wrap(err, "Scan")
	}

	return totalUsage.Int64, nil
}

// GetCountByGiftCodeX get count usage report
func (rr *reportPGRepository) GetCountByGiftCodeX(ctx context.Context, mobile string, tx pgx.Tx) (int64, pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.GetCountByGiftCodeX")
	defer span.Finish()

	if tx == nil {
		var beginTxErr error

		tx, beginTxErr = rr.db.Begin(ctx)
		if beginTxErr != nil {
			return 0, tx, beginTxErr
		}
	}

	var totalUsage sql.NullInt64
	if err := rr.db.QueryRow(ctx, getcountByGiftCodeQuery, mobile).
		Scan(
			&totalUsage,
		); err != nil {
		return 0, tx, errors.Wrap(err, "Scan")
	}

	return totalUsage.Int64, tx, nil
}

// RollBack rolls back a transaction
func (rr *reportPGRepository) RollBack(ctx context.Context, tx pgx.Tx) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.RollBack")
	defer span.Finish()

	if tx == nil {
		return fmt.Errorf("transactions not begin")
	}

	return tx.Rollback(ctx)
}

// Commit commits a transaction
func (rr *reportPGRepository) Commit(ctx context.Context, tx pgx.Tx) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "reportPGRepository.Commit")
	defer span.Finish()

	if tx == nil {
		return fmt.Errorf("transactions not begin")
	}

	return tx.Commit(ctx)
}
