package repository

const (
	createGiftChargeQuery = `INSERT INTO gift_charges (code, validity_period_start, validity_period_end, amount, 
                                          max_usage_Count,created_at)  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, code, validity_period_start, validity_period_end, amount, max_usage_count, created_at, updated_at`

	getByIDQuery = `SELECT id, code, validity_period_start, validity_period_end, amount, max_usage_count, created_at, updated_at
                                 FROM gift_charges  WHERE id = $1;`

	getByCodeQuery = `SELECT id, code, validity_period_start, validity_period_end, amount, max_usage_count, created_at, updated_at 
			                  FROM gift_charges  WHERE code = $1 AND $2 BETWEEN validity_period_start AND validity_period_end;`

	getListQuery = `SELECT id, code, validity_period_start, validity_period_end, amount, max_usage_count, created_at, updated_at 
			                  FROM gift_charges;`
	getValidListQuery = `SELECT id, code, validity_period_start, validity_period_end, amount, max_usage_count, created_at, updated_at 
			                  FROM gift_charges WHERE validity_period_start < $2  AND validity_period_end > $1;`

	getInValidListQuery = `SELECT id, code, validity_period_start, validity_period_end, amount, max_usage_count, created_at, updated_at 
			                  FROM gift_charges WHERE $1 < validity_period_start OR $2 > validity_period_end;`
)
