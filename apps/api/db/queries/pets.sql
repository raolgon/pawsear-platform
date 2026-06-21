-- name: CreatePet :one
INSERT INTO pets (
	id,
	household_id,
	name,
	species,
	breed,
	size,
	sex,
	birthdate,
	color_markings,
	behavior_notes,
	medical_notes,
	feeding_notes,
	vet_notes,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPet :one
SELECT *
FROM pets
WHERE id = ?;

-- name: ListPets :many
SELECT *
FROM pets
WHERE active = 1
ORDER BY name;

-- name: ListPetsByHousehold :many
SELECT *
FROM pets
WHERE active = 1 AND household_id = ?
ORDER BY name;

-- name: UpdatePet :one
UPDATE pets
SET
	household_id = ?,
	name = ?,
	species = ?,
	breed = ?,
	size = ?,
	sex = ?,
	birthdate = ?,
	color_markings = ?,
	behavior_notes = ?,
	medical_notes = ?,
	feeding_notes = ?,
	vet_notes = ?,
	active = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;
