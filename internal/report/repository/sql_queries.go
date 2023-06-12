package repository

const (
	createReportQuery = `INSERT INTO reports (mobile, gift_code, charge_amount, report_time, created_at)
                                               VALUES ($1, $2, $3, $4, $5) RETURNING id, gift_code, mobile, charge_amount, report_time, created_at`

	getByGiftCodeQuery = `SELECT id, gift_code, mobile, charge_amount, report_time, created_at
			                      FROM reports WHERE gift_code = $1`

	getByMobileQuery = `SELECT id, gift_code, mobile, charge_amount, report_time, created_at
			                      FROM reports WHERE mobile = $1 `

	getByGiftCodeAndMobileQuery = `SELECT id, gift_code, mobile, charge_amount, report_time, created_at
			                      FROM reports WHERE mobile = $1 AND gift_code = $2`

	getcountByGiftCodeQuery = `SELECT COUNT(*) AS usage_count FROM wallet_charge_report WHERE gift_code = $1`
)
