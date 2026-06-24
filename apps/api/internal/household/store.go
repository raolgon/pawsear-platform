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
	Delete(ctx context.Context, id string) error
}

type SQLStore struct {
	database *sql.DB
	queries  *dbqueries.Queries
}

func NewSQLStore(database *sql.DB) *SQLStore {
	return &SQLStore{database: database, queries: dbqueries.New(database)}
}

func (s *SQLStore) Delete(ctx context.Context, id string) error {
	tx, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: begin household deletion: %v", ErrStoreIssue, err)
	}
	defer tx.Rollback()
	queries := s.queries.WithTx(tx)
	contactIDs, err := queries.ListHouseholdContactIDsForDeletion(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: list household contacts for deletion: %v", ErrStoreIssue, err)
	}
	nullableID := sql.NullString{String: id, Valid: true}
	steps := []func() error{
		func() error { return queries.DeleteHouseholdBookingSources(ctx, id) },
		func() error { return queries.DeleteHouseholdPaymentAllocations(ctx, id) },
		func() error { return queries.DeleteHouseholdOutboundMessages(ctx, nullableID) },
		func() error { return queries.DeleteHouseholdDetectedRequests(ctx, nullableID) },
		func() error { return queries.DeleteHouseholdMessages(ctx, nullableID) },
		func() error { return queries.DeleteHouseholdConversations(ctx, nullableID) },
		func() error { return queries.DeleteHouseholdCareTasks(ctx, id) },
		func() error { return queries.DeleteHouseholdCareRoutines(ctx, id) },
		func() error { return queries.DeleteHouseholdCharges(ctx, id) },
		func() error { return queries.DeleteHouseholdBookingPets(ctx, id) },
		func() error { return queries.DeleteHouseholdBookings(ctx, id) },
		func() error { return queries.DeleteHouseholdPetMedications(ctx, id) },
		func() error { return queries.DeleteHouseholdPetDiets(ctx, id) },
		func() error { return queries.DeleteHouseholdPets(ctx, id) },
		func() error { return queries.DeleteHouseholdContactLinks(ctx, id) },
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return fmt.Errorf("%w: delete household dependencies: %v", ErrStoreIssue, err)
		}
	}
	deleted, err := queries.DeleteHouseholdRecord(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: delete household: %v", ErrStoreIssue, err)
	}
	if deleted == 0 {
		return ErrNotFound
	}
	for _, contactID := range contactIDs {
		if err := queries.DeleteContactIdentitiesIfOtherwiseOrphaned(ctx, contactID); err != nil {
			return fmt.Errorf("%w: delete orphaned contact identities: %v", ErrStoreIssue, err)
		}
		if err := queries.DeleteContactIfOrphaned(ctx, contactID); err != nil {
			return fmt.Errorf("%w: delete orphaned contact: %v", ErrStoreIssue, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: commit household deletion: %v", ErrStoreIssue, err)
	}
	return nil
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
