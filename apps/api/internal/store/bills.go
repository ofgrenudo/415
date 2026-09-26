// Package store provides database-backed CRUD operations for domain models.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ofgrenudo/415/apps/api/internal/models"
)

// ErrNotFound is returned when a requested bill does not exist.
var ErrNotFound = errors.New("bill not found")

// BillStore provides CRUD access to the bills table.
type BillStore struct {
	db *sql.DB
}

// NewBillStore returns a BillStore backed by conn.
func NewBillStore(conn *sql.DB) *BillStore {
	return &BillStore{db: conn}
}

// List returns all bills ordered by most recently created.
func (s *BillStore) List(ctx context.Context) ([]models.Bill, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, description, cost, date_paid, created_at, updated_at
		FROM bills
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list bills: %w", err)
	}
	defer func() { _ = rows.Close() }()

	bills := []models.Bill{}
	for rows.Next() {
		var b models.Bill
		var datePaid sql.NullTime
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Cost, &datePaid, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan bill: %w", err)
		}
		b.DatePaid = datePaid.Time
		bills = append(bills, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list bills: %w", err)
	}

	return bills, nil
}

// Get returns the bill with the given id.
func (s *BillStore) Get(ctx context.Context, id int64) (models.Bill, error) {
	var b models.Bill
	var datePaid sql.NullTime

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, description, cost, date_paid, created_at, updated_at
		FROM bills
		WHERE id = $1`, id,
	).Scan(&b.ID, &b.Name, &b.Description, &b.Cost, &datePaid, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Bill{}, ErrNotFound
	}
	if err != nil {
		return models.Bill{}, fmt.Errorf("get bill: %w", err)
	}
	b.DatePaid = datePaid.Time

	return b, nil
}

// Create inserts a new bill and returns it with its generated id and timestamps.
func (s *BillStore) Create(ctx context.Context, b models.Bill) (models.Bill, error) {
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO bills (name, description, cost, date_paid)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		b.Name, b.Description, b.Cost, nullableTime(b.DatePaid),
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return models.Bill{}, fmt.Errorf("create bill: %w", err)
	}

	return b, nil
}

// Update replaces the fields of the bill with the given id.
func (s *BillStore) Update(ctx context.Context, id int64, b models.Bill) (models.Bill, error) {
	var datePaid sql.NullTime

	err := s.db.QueryRowContext(ctx, `
		UPDATE bills
		SET name = $1, description = $2, cost = $3, date_paid = $4, updated_at = now()
		WHERE id = $5
		RETURNING id, name, description, cost, date_paid, created_at, updated_at`,
		b.Name, b.Description, b.Cost, nullableTime(b.DatePaid), id,
	).Scan(&b.ID, &b.Name, &b.Description, &b.Cost, &datePaid, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Bill{}, ErrNotFound
	}
	if err != nil {
		return models.Bill{}, fmt.Errorf("update bill: %w", err)
	}
	b.DatePaid = datePaid.Time

	return b, nil
}

// Delete removes the bill with the given id.
func (s *BillStore) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM bills WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete bill: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete bill: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func nullableTime(t interface{ IsZero() bool }) any {
	if t.IsZero() {
		return nil
	}
	return t
}
