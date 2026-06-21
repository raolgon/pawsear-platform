package httpapi

import (
	"fmt"
	"net/http"

	dbqueries "github.com/pawsear/pawsear-platform/apps/api/internal/db/queries"
)

type petInput struct {
	HouseholdID   string  `json:"householdId"`
	Name          string  `json:"name"`
	Species       string  `json:"species"`
	Breed         *string `json:"breed"`
	Size          *string `json:"size"`
	Sex           *string `json:"sex"`
	Birthdate     *string `json:"birthdate"`
	ColorMarkings *string `json:"colorMarkings"`
	BehaviorNotes *string `json:"behaviorNotes"`
	MedicalNotes  *string `json:"medicalNotes"`
	FeedingNotes  *string `json:"feedingNotes"`
	VetNotes      *string `json:"vetNotes"`
	Active        *bool   `json:"active"`
}

type staffInput struct {
	DisplayName string  `json:"displayName"`
	Phone       *string `json:"phone"`
	Role        string  `json:"role"`
	Active      *bool   `json:"active"`
}

func (h *mvpHandler) listPets(w http.ResponseWriter, r *http.Request) {
	householdID := r.URL.Query().Get("householdId")
	var rows []dbqueries.Pet
	var err error
	if householdID == "" {
		rows, err = h.queries.ListPets(r.Context())
	} else {
		rows, err = h.queries.ListPetsByHousehold(r.Context(), householdID)
	}
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"pets": petListResponse(rows)})
}

func (h *mvpHandler) createPet(w http.ResponseWriter, r *http.Request) {
	var input petInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.petCreateParams(input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	created, err := h.queries.CreatePet(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, petResponse(created))
}

func (h *mvpHandler) getPet(w http.ResponseWriter, r *http.Request) {
	pet, err := h.queries.GetPet(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, petResponse(pet))
}

func (h *mvpHandler) updatePet(w http.ResponseWriter, r *http.Request) {
	current, err := h.queries.GetPet(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	var input petInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.petUpdateParams(current, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	updated, err := h.queries.UpdatePet(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, petResponse(updated))
}

func (h *mvpHandler) listStaff(w http.ResponseWriter, r *http.Request) {
	rows, err := h.queries.ListStaffMembers(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"staff": staffListResponse(rows)})
}

func (h *mvpHandler) createStaff(w http.ResponseWriter, r *http.Request) {
	var input staffInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.staffCreateParams(input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	created, err := h.queries.CreateStaffMember(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, staffResponse(created))
}

func (h *mvpHandler) getStaff(w http.ResponseWriter, r *http.Request) {
	staff, err := h.queries.GetStaffMember(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, staffResponse(staff))
}

func (h *mvpHandler) updateStaff(w http.ResponseWriter, r *http.Request) {
	current, err := h.queries.GetStaffMember(r.Context(), r.PathValue("id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	var input staffInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	params, err := h.staffUpdateParams(current, input)
	if err != nil {
		writeInvalid(w, err)
		return
	}
	updated, err := h.queries.UpdateStaffMember(r.Context(), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, staffResponse(updated))
}

func (h *mvpHandler) petCreateParams(input petInput) (dbqueries.CreatePetParams, error) {
	householdID, err := requiredText(input.HouseholdID, "householdId")
	if err != nil {
		return dbqueries.CreatePetParams{}, err
	}
	name, species, err := validatePetNameSpecies(input.Name, input.Species)
	if err != nil {
		return dbqueries.CreatePetParams{}, err
	}
	recordID, err := newRecordID()
	if err != nil {
		return dbqueries.CreatePetParams{}, err
	}
	now := timestamp(h.now)
	return dbqueries.CreatePetParams{
		ID: recordID, HouseholdID: householdID, Name: name, Species: species,
		Breed: optionalText(input.Breed), Size: optionalText(input.Size), Sex: optionalText(input.Sex),
		Birthdate: optionalText(input.Birthdate), ColorMarkings: optionalText(input.ColorMarkings),
		BehaviorNotes: optionalText(input.BehaviorNotes), MedicalNotes: optionalText(input.MedicalNotes),
		FeedingNotes: optionalText(input.FeedingNotes), VetNotes: optionalText(input.VetNotes),
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (h *mvpHandler) petUpdateParams(current dbqueries.Pet, input petInput) (dbqueries.UpdatePetParams, error) {
	householdID := defaultText(input.HouseholdID, current.HouseholdID)
	name := defaultText(input.Name, current.Name)
	species := defaultText(input.Species, current.Species)
	name, species, err := validatePetNameSpecies(name, species)
	if err != nil {
		return dbqueries.UpdatePetParams{}, err
	}
	active := boolFlag(current.Active)
	if input.Active != nil {
		active = *input.Active
	}
	return dbqueries.UpdatePetParams{
		HouseholdID: householdID, Name: name, Species: species,
		Breed: patchText(input.Breed, current.Breed), Size: patchText(input.Size, current.Size),
		Sex: patchText(input.Sex, current.Sex), Birthdate: patchText(input.Birthdate, current.Birthdate),
		ColorMarkings: patchText(input.ColorMarkings, current.ColorMarkings),
		BehaviorNotes: patchText(input.BehaviorNotes, current.BehaviorNotes),
		MedicalNotes:  patchText(input.MedicalNotes, current.MedicalNotes),
		FeedingNotes:  patchText(input.FeedingNotes, current.FeedingNotes),
		VetNotes:      patchText(input.VetNotes, current.VetNotes),
		Active:        intFlag(active), UpdatedAt: timestamp(h.now), ID: current.ID,
	}, nil
}

func validatePetNameSpecies(name string, species string) (string, string, error) {
	cleanName, err := requiredText(name, "name")
	if err != nil {
		return "", "", err
	}
	cleanSpecies := defaultText(species, "dog")
	if !allowed(cleanSpecies, "dog", "cat", "other") {
		return "", "", fmt.Errorf("species is not supported")
	}
	return cleanName, cleanSpecies, nil
}

func (h *mvpHandler) staffCreateParams(input staffInput) (dbqueries.CreateStaffMemberParams, error) {
	displayName, err := requiredText(input.DisplayName, "displayName")
	if err != nil {
		return dbqueries.CreateStaffMemberParams{}, err
	}
	role := defaultText(input.Role, "walker")
	if !allowed(role, "owner_operator", "walker", "sitter", "admin") {
		return dbqueries.CreateStaffMemberParams{}, fmt.Errorf("role is not supported")
	}
	recordID, err := newRecordID()
	if err != nil {
		return dbqueries.CreateStaffMemberParams{}, err
	}
	now := timestamp(h.now)
	return dbqueries.CreateStaffMemberParams{
		ID: recordID, DisplayName: displayName, Phone: optionalText(input.Phone),
		Role: role, CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (h *mvpHandler) staffUpdateParams(current dbqueries.StaffMember, input staffInput) (dbqueries.UpdateStaffMemberParams, error) {
	displayName := defaultText(input.DisplayName, current.DisplayName)
	role := defaultText(input.Role, current.Role)
	if !allowed(role, "owner_operator", "walker", "sitter", "admin") {
		return dbqueries.UpdateStaffMemberParams{}, fmt.Errorf("role is not supported")
	}
	active := boolFlag(current.Active)
	if input.Active != nil {
		active = *input.Active
	}
	return dbqueries.UpdateStaffMemberParams{
		DisplayName: displayName, Phone: patchText(input.Phone, current.Phone), Role: role,
		Active: intFlag(active), UpdatedAt: timestamp(h.now), ID: current.ID,
	}, nil
}

func petListResponse(rows []dbqueries.Pet) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, petResponse(row))
	}
	return items
}

func petResponse(row dbqueries.Pet) map[string]any {
	return map[string]any{
		"id": row.ID, "householdId": row.HouseholdID, "name": row.Name, "species": row.Species,
		"breed": textValue(row.Breed), "size": textValue(row.Size), "sex": textValue(row.Sex),
		"birthdate": textValue(row.Birthdate), "colorMarkings": textValue(row.ColorMarkings),
		"behaviorNotes": textValue(row.BehaviorNotes), "medicalNotes": textValue(row.MedicalNotes),
		"feedingNotes": textValue(row.FeedingNotes), "vetNotes": textValue(row.VetNotes),
		"active": boolFlag(row.Active), "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}

func staffListResponse(rows []dbqueries.StaffMember) []map[string]any {
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, staffResponse(row))
	}
	return items
}

func staffResponse(row dbqueries.StaffMember) map[string]any {
	return map[string]any{
		"id": row.ID, "displayName": row.DisplayName, "phone": textValue(row.Phone),
		"role": row.Role, "active": boolFlag(row.Active), "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt,
	}
}
