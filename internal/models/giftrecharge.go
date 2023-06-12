package models

import (
	"database/sql"
	"time"
)

// GiftCharge models
type GiftCharge struct {
	ID                  int64     `json:"id,omitempty"`
	Code                string    `json:"code,omitempty"`
	ValidityPeriodStart time.Time `json:"validityPeriodStart"`
	ValidityPeriodEnd   time.Time `json:"validityPeriodEnd"`
	Amount              int64     `json:"amount,omitempty"`
	MaxUsageCount       int       `json:"maxUsageCount,omitempty"`
	CreatedAt           time.Time `json:"createdAt,omitempty"`
	UpdatedAt           time.Time `json:"updatedAt,omitempty"`
}

// GiftChargeModel models
type GiftChargeModel struct {
	ID                  sql.NullInt64  `json:"ID"`
	Code                sql.NullString `json:"code"`
	ValidityPeriodStart sql.NullTime   `json:"validityPeriodStart"`
	ValidityPeriodEnd   sql.NullTime   `json:"validityPeriodEnd"`
	Amount              sql.NullInt64  `json:"amount"`
	MaxUsageCount       sql.NullInt64  `json:"maxUsageCount"`
	CreatedAt           sql.NullTime   `json:"createdAt"`
	UpdatedAt           sql.NullTime   `json:"updatedAt"`
}

// GiftChargeCreateReq Create request models
type GiftChargeCreateReq struct {
	Code                string    `json:"-"`
	ValidityPeriodStart time.Time `json:"validityPeriodStart"`
	ValidityPeriodEnd   time.Time `json:"validityPeriodEnd"`
	Amount              int64     `json:"amount,omitempty"`
	MaxUsageCount       int       `json:"maxUsageCount,omitempty"`
}

func (m GiftChargeModel) Entity() *GiftCharge {
	return &GiftCharge{
		ID:                  m.ID.Int64,
		Code:                m.Code.String,
		ValidityPeriodStart: m.ValidityPeriodStart.Time,
		ValidityPeriodEnd:   m.ValidityPeriodEnd.Time,
		Amount:              m.Amount.Int64,
		MaxUsageCount:       int(m.MaxUsageCount.Int64),
		CreatedAt:           m.CreatedAt.Time,
		UpdatedAt:           m.UpdatedAt.Time,
	}
}
