package models

import (
	"database/sql"
	"time"
)

type Report struct {
	Id         int64     `json:"id"`
	GiftCode   string    `json:"giftCode"`
	Mobile     string    `json:"mobile"`
	Amount     int64     `json:"chargeAmount"`
	ReportTime time.Time `json:"reportTime"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ReportModel struct {
	Id         sql.NullInt64
	GiftCode   sql.NullString
	Mobile     sql.NullString
	Amount     sql.NullInt64
	ReportTime sql.NullTime
	CreatedAt  sql.NullTime
}

type CreateReportReq struct {
	GiftCode   string    `json:"giftCode"`
	Mobile     string    `json:"mobile"`
	Amount     int64     `json:"chargeAmount"`
	ReportTime time.Time `json:"reportTime"`
}

type CountUsage struct {
	GiftCode   string `json:"giftCode"`
	CountUsage int64  `json:"countUsage"`
}

type CheckGiftUsage struct {
	GiftCode string `json:"giftCode"`
	Mobile   string `json:"mobile"`
}

func (m ReportModel) Entity() *Report {
	return &Report{
		Id:         m.Id.Int64,
		GiftCode:   m.GiftCode.String,
		Mobile:     m.Mobile.String,
		Amount:     m.Amount.Int64,
		ReportTime: m.ReportTime.Time,
		CreatedAt:  m.CreatedAt.Time,
	}
}
