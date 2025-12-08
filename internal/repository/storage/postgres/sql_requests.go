package postgres

const (
	insertOrderSQL = `
		INSERT INTO orders (
			id, user_id, status, created_at, 
			expires_at, issued_at, returned_at
		) VALUES (
			$1,$2,$3,$4,
			$5,$6,$7
		);`

	updateOrderSQL = `
		UPDATE orders SET
			user_id      = $2,
			status       = $3,
			expires_at   = $4,
			issued_at    = $5,
			returned_at  = $6
		WHERE id = $1;`

	selectOrderSQL = `
		SELECT
			id, user_id, status, created_at,
			expires_at, issued_at, returned_at
		FROM orders
		WHERE id = $1;`

	deleteOrderSQL = `DELETE FROM orders WHERE id = $1;`

	listAllSQL = `
		SELECT
			id, user_id, status, expires_at,
			issued_at, returned_at, created_at
		FROM orders
		ORDER BY created_at ASC, id ASC;`

	listByUserSQL = `
		SELECT
			id, user_id, status, expires_at,
			issued_at, returned_at, created_at
		FROM orders
		WHERE user_id = $1`
)
