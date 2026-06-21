package household

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pawsear/pawsear-platform/apps/api/internal/platform/id"
)

var (
	ErrNotFound   = errors.New("household not found")
	ErrInvalid    = errors.New("invalid household")
	ErrNoChanges  = errors.New("no household changes provided")
	ErrStoreIssue = errors.New("household store error")
)

type Household struct {
	ID           string  `json:"id"`
	DisplayName  string  `json:"displayName"`
	AddressLine1 *string `json:"addressLine1,omitempty"`
	AddressLine2 *string `json:"addressLine2,omitempty"`
	Neighborhood *string `json:"neighborhood,omitempty"`
	City         *string `json:"city,omitempty"`
	Timezone     *string `json:"timezone,omitempty"`
	Notes        *string `json:"notes,omitempty"`
	Active       bool    `json:"active"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type CreateInput struct {
	DisplayName  string  `json:"displayName"`
	AddressLine1 *string `json:"addressLine1"`
	AddressLine2 *string `json:"addressLine2"`
	Neighborhood *string `json:"neighborhood"`
	City         *string `json:"city"`
	Timezone     *string `json:"timezone"`
	Notes        *string `json:"notes"`
}

type UpdateInput struct {
	DisplayName  *string `json:"displayName"`
	AddressLine1 *string `json:"addressLine1"`
	AddressLine2 *string `json:"addressLine2"`
	Neighborhood *string `json:"neighborhood"`
	City         *string `json:"city"`
	Timezone     *string `json:"timezone"`
	Notes        *string `json:"notes"`
	Active       *bool   `json:"active"`
}

type Service struct {
	store Store
	now   func() time.Time
}

func NewService(database *sql.DB) *Service {
	return &Service{
		store: NewSQLStore(database),
		now:   func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Household, error) {
	normalized, err := normalizeCreate(input)
	if err != nil {
		return Household{}, err
	}

	newID, err := id.New()
	if err != nil {
		return Household{}, fmt.Errorf("%w: create id: %v", ErrStoreIssue, err)
	}

	timestamp := s.now().Format(time.RFC3339Nano)
	params := createParams{
		ID:           newID,
		DisplayName:  normalized.DisplayName,
		AddressLine1: nullableString(normalized.AddressLine1),
		AddressLine2: nullableString(normalized.AddressLine2),
		Neighborhood: nullableString(normalized.Neighborhood),
		City:         nullableString(normalized.City),
		Timezone:     nullableString(normalized.Timezone),
		Notes:        nullableString(normalized.Notes),
		CreatedAt:    timestamp,
		UpdatedAt:    timestamp,
	}

	return s.store.Create(ctx, params)
}

func (s *Service) Get(ctx context.Context, householdID string) (Household, error) {
	householdID = strings.TrimSpace(householdID)
	if householdID == "" {
		return Household{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}

	return s.store.Get(ctx, householdID)
}

func (s *Service) List(ctx context.Context) ([]Household, error) {
	return s.store.List(ctx)
}

func (s *Service) Update(ctx context.Context, householdID string, input UpdateInput) (Household, error) {
	householdID = strings.TrimSpace(householdID)
	if householdID == "" {
		return Household{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}

	current, err := s.store.Get(ctx, householdID)
	if err != nil {
		return Household{}, err
	}

	updated, changed, err := applyUpdate(current, input)
	if err != nil {
		return Household{}, err
	}
	if !changed {
		return Household{}, ErrNoChanges
	}

	params := updateParams{
		ID:           householdID,
		DisplayName:  updated.DisplayName,
		AddressLine1: nullableString(updated.AddressLine1),
		AddressLine2: nullableString(updated.AddressLine2),
		Neighborhood: nullableString(updated.Neighborhood),
		City:         nullableString(updated.City),
		Timezone:     nullableString(updated.Timezone),
		Notes:        nullableString(updated.Notes),
		Active:       boolToInt(updated.Active),
		UpdatedAt:    s.now().Format(time.RFC3339Nano),
	}

	return s.store.Update(ctx, params)
}

func normalizeCreate(input CreateInput) (CreateInput, error) {
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	if input.DisplayName == "" {
		return CreateInput{}, fmt.Errorf("%w: displayName is required", ErrInvalid)
	}

	input.AddressLine1 = cleanOptional(input.AddressLine1)
	input.AddressLine2 = cleanOptional(input.AddressLine2)
	input.Neighborhood = cleanOptional(input.Neighborhood)
	input.City = cleanOptional(input.City)
	input.Timezone = cleanOptional(input.Timezone)
	input.Notes = cleanOptional(input.Notes)

	return input, nil
}

func applyUpdate(current Household, input UpdateInput) (Household, bool, error) {
	updated := current
	changed := false

	if input.DisplayName != nil {
		value := strings.TrimSpace(*input.DisplayName)
		if value == "" {
			return Household{}, false, fmt.Errorf("%w: displayName cannot be empty", ErrInvalid)
		}
		if value != current.DisplayName {
			updated.DisplayName = value
			changed = true
		}
	}

	changed = applyOptionalString(&updated.AddressLine1, input.AddressLine1) || changed
	changed = applyOptionalString(&updated.AddressLine2, input.AddressLine2) || changed
	changed = applyOptionalString(&updated.Neighborhood, input.Neighborhood) || changed
	changed = applyOptionalString(&updated.City, input.City) || changed
	changed = applyOptionalString(&updated.Timezone, input.Timezone) || changed
	changed = applyOptionalString(&updated.Notes, input.Notes) || changed

	if input.Active != nil && *input.Active != current.Active {
		updated.Active = *input.Active
		changed = true
	}

	return updated, changed, nil
}

func applyOptionalString(target **string, input *string) bool {
	if input == nil {
		return false
	}

	cleaned := cleanOptional(input)
	if sameOptionalString(*target, cleaned) {
		return false
	}

	*target = cleaned
	return true
}

func cleanOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func nullableString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func optionalString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func sameOptionalString(left *string, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func intToBool(value int) bool {
	return value != 0
}
