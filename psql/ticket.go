package psql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	tixer "github.com/mroobert/monorepo-tixer"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDbRecordNotFound = errors.New("db record not found")
	ErrDbEditConflict   = errors.New("db edit conflict")
)

const ticketsTable = "tickets"

// TicketRepository persists tickets in the database.
type TicketRepository struct {
	DB           *pgxpool.Pool
	QueryTimeout time.Duration
}

func NewTicketRepository(db *pgxpool.Pool, queryTimeout time.Duration) *TicketRepository {
	return &TicketRepository{
		DB:           db,
		QueryTimeout: queryTimeout,
	}
}

// Insert inserts a new ticket in the database.
func (tr *TicketRepository) Insert(ctx context.Context, ticket tixer.Ticket) (tixer.Ticket, error) {
	query := `INSERT INTO ` + ticketsTable +
		` (public_id, title, price) VALUES ($1, $2, $3)
        RETURNING id, public_id, title, price, version, created_at, updated_at`

	args := []any{ticket.PublicID, ticket.Title, ticket.Price}

	queryCtx, cancel := context.WithTimeout(ctx, tr.QueryTimeout)
	defer cancel()

	var createdTicket tixer.Ticket
	if err := tr.DB.QueryRow(queryCtx, query, args...).Scan(
		&createdTicket.ID,
		&createdTicket.PublicID,
		&createdTicket.Title,
		&createdTicket.Price,
		&createdTicket.Version,
		&createdTicket.CreatedAt,
		&createdTicket.UpdatedAt,
	); err != nil {
		return tixer.Ticket{}, fmt.Errorf("failed to insert ticket in database: %w", err)
	}

	return createdTicket, nil
}

// SelectOne reads a ticket from the database.
func (tr *TicketRepository) SelectOne(ctx context.Context, id tixer.PublicID) (tixer.Ticket, error) {
	query := `SELECT id, public_id, title, price, version, created_at, updated_at FROM ` + ticketsTable +
		` WHERE public_id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, tr.QueryTimeout)
	defer cancel()

	var ticket tixer.Ticket
	if err := tr.DB.QueryRow(queryCtx, query, id).Scan(
		&ticket.ID,
		&ticket.PublicID,
		&ticket.Title,
		&ticket.Price,
		&ticket.Version,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return tixer.Ticket{}, ErrDbRecordNotFound
		default:
			return tixer.Ticket{}, fmt.Errorf("failed to select ticket from database: %w", err)
		}
	}

	return ticket, nil
}

type TicketFilter struct {
	Title         string
	Limit         int
	Offset        int
	SortColumn    string
	SortDirection string
}

// SelectMultiple reads tickets based on filters from the database.
func (tr *TicketRepository) SelectMultiple(ctx context.Context, filter TicketFilter) ([]tixer.Ticket, Pagination, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, public_id, title, price, version, created_at, updated_at `+
		` FROM `+ticketsTable+
		` WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') `+
		` ORDER BY %s %s LIMIT $2 OFFSET $3`, filter.SortColumn, filter.SortDirection)

	args := []any{filter.Title, filter.Limit, filter.Offset}

	queryCtx, cancel := context.WithTimeout(ctx, tr.QueryTimeout)
	defer cancel()

	rows, err := tr.DB.Query(queryCtx, query, args...)
	if err != nil {
		return nil, Pagination{}, fmt.Errorf("failed to select tickets from database: %w", err)
	}

	defer rows.Close()

	totalRecords := 0
	tickets := []tixer.Ticket{}

	for rows.Next() {
		var ticket tixer.Ticket

		err := rows.Scan(
			&totalRecords,
			&ticket.ID,
			&ticket.PublicID,
			&ticket.Title,
			&ticket.Price,
			&ticket.Version,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			return nil, Pagination{}, fmt.Errorf("failed to scan row result: %w", err)
		}

		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, Pagination{}, fmt.Errorf("failed to iterate over rows result: %w", err)
	}

	pagination := calculatePagination(totalRecords, filter.Offset, filter.Limit)

	return tickets, pagination, nil
}

// Update updates a ticket in the database
func (tr *TicketRepository) Update(ctx context.Context, ticket *tixer.Ticket) error {
	query := `UPDATE ` + ticketsTable +
		` SET title = $1, price = $2, version = version + 1, updated_at = $3` +
		` WHERE public_id = $4 AND version = $5 RETURNING version`

	args := []any{ticket.Title, ticket.Price, time.Now(), ticket.PublicID, ticket.Version}

	queryCtx, cancel := context.WithTimeout(ctx, tr.QueryTimeout)
	defer cancel()

	if err := tr.DB.QueryRow(queryCtx, query, args...).Scan(&ticket.Version); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrDbEditConflict
		default:
			return fmt.Errorf("failed to update ticket in database: %w", err)
		}
	}

	return nil
}

// Delete deletes a ticket from the database.
func (tr *TicketRepository) Delete(ctx context.Context, id tixer.PublicID) error {
	query := `DELETE FROM ` + ticketsTable +
		` WHERE public_id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, tr.QueryTimeout)
	defer cancel()

	res, err := tr.DB.Exec(queryCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket from database: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrDbRecordNotFound
	}

	return nil
}
