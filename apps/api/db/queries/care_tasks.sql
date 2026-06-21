-- name: CreateCareTask :one
INSERT INTO care_tasks (
	id,
	booking_id,
	household_id,
	pet_id,
	task_type,
	title,
	instructions,
	due_at,
	status,
	assigned_staff_id,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetCareTask :one
SELECT *
FROM care_tasks
WHERE id = ?;

-- name: ListCareTasks :many
SELECT *
FROM care_tasks
ORDER BY due_at, created_at;

-- name: ListCareTasksByRange :many
SELECT *
FROM care_tasks
WHERE due_at >= ? AND due_at < ?
ORDER BY due_at, created_at;

-- name: ListCareTasksByBooking :many
SELECT *
FROM care_tasks
WHERE booking_id = ?
ORDER BY due_at, created_at;

-- name: ListCareTasksByHousehold :many
SELECT *
FROM care_tasks
WHERE household_id = ?
ORDER BY due_at, created_at;

-- name: UpdateCareTask :one
UPDATE care_tasks
SET
	booking_id = ?,
	household_id = ?,
	pet_id = ?,
	task_type = ?,
	title = ?,
	instructions = ?,
	due_at = ?,
	status = ?,
	assigned_staff_id = ?,
	completed_at = ?,
	completed_by_staff_id = ?,
	skipped_reason = ?,
	updated_at = ?
WHERE id = ?
RETURNING *;
