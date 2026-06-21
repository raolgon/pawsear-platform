package household

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type Store interface {
	Create(ctx context.Context, params createParams) (Household, error)
	Get(ctx context.Context, id string) (Household, error)
	List(ctx context.Context) ([]Household, error)
	Update(ctx context.Context, params updateParams) (Household, error)
}

type SQLStore struct {
	queries *dbqueries.Queries
}

func NewSQLStore(database *sql.DB) *SQLStore {
	return &SQLStore{queries: dbqueries.New(database)}
}

type createParams struct {
	ID           string
	DisplayName  string
	AddressLine1 sql.NullString
	AddressLine2 sql.NullString
	Neighborhood sql.NullString
	City         sql.NullString
	Timezone     sql.NullString
	Notes        sql.NullString
	CreatedAt    string
	UpdatedAt    string
}

type updateParams struct {
	ID           string
	DisplayName  string
	AddressLine1 sql.NullString
	AddressLine2 sql.NullString
	Neighborhood sql.NullString
	City         sql.NullString
	Timezone     sql.NullString
	Notes        sql.NullString
	Active       int
	UpdatedAt    string
}

func (s *SQLStore) Create(ctx context.Context, params createParams) (Household, error) {
	created, err := s.queries.CreateHousehold(ctx, dbqueries.CreateHouseholdParams{
		ID:           params.ID,
		DisplayName:  params.DisplayName,
		AddressLine1: params.AddressLine1,
		AddressLine2: params.AddressLine2,
		Neighborhood: params.Neighborhood,
		City:         params.City,
		Timezone:     params.Timezone,
		Notes:        params.Notes,
		CreatedAt:    params.CreatedAt,
		UpdatedAt:    params.UpdatedAt,
	})
	if err != nil {
		return Household{}, fmt.Errorf("%w: create household: %v", ErrStoreIssue, err)
	}

	return fromQueryHousehold(created), nil
}

func (s *SQLStore) Get(ctx context.Context, id string) (Household, error) {
	household, err := s.queries.GetHousehold(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Household{}, ErrNotFound
	}
	if err != nil {
		return Household{}, fmt.Errorf("%w: get household: %v", ErrStoreIssue, err)
	}

	return fromQueryHousehold(household), nil
}

func (s *SQLStore) List(ctx context.Context) ([]Household, error) {
	rows, err := s.queries.ListHouseholds(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: list households: %v", ErrStoreIssue, err)
	}

	households := make([]Household, 0, len(rows))
	for _, row := range rows {
		households = append(households, fromQueryHousehold(row))
	}

	return households, nil
}

func (s *SQLStore) Update(ctx context.Context, params updateParams) (Household, error) {
	updated, err := s.queries.UpdateHousehold(ctx, dbqueries.UpdateHouseholdParams{
		DisplayName:  params.DisplayName,
		AddressLine1: params.AddressLine1,
		AddressLine2: params.AddressLine2,
		Neighborhood: params.Neighborhood,
		City:         params.City,
		Timezone:     params.Timezone,
		Notes:        params.Notes,
		Active:       int64(params.Active),
		UpdatedAt:    params.UpdatedAt,
		ID:           params.ID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return Household{}, ErrNotFound
	}
	if err != nil {
		return Household{}, fmt.Errorf("%w: update household: %v", ErrStoreIssue, err)
	}

	return fromQueryHousehold(updated), nil
}

func fromQueryHousehold(row dbqueries.Household) Household {
	return Household{
		ID:           row.ID,
		DisplayName:  row.DisplayName,
		AddressLine1: optionalString(row.AddressLine1),
		AddressLine2: optionalString(row.AddressLine2),
		Neighborhood: optionalString(row.Neighborhood),
		City:         optionalString(row.City),
		Timezone:     optionalString(row.Timezone),
		Notes:        optionalString(row.Notes),
		Active:       intToBool(int(row.Active)),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
