-- name: ListHouseholdContactIDsForDeletion :many
SELECT DISTINCT contact_id
FROM household_contacts
WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdBookingSources :exec
DELETE FROM booking_sources
WHERE booking_id IN (SELECT bookings.id FROM bookings WHERE bookings.household_id = sqlc.arg(target_household_id))
   OR detected_request_id IN (
	SELECT dr.id
	FROM detected_requests AS dr
	WHERE dr.household_id = sqlc.arg(target_household_id)
	   OR dr.converted_booking_id IN (SELECT bookings.id FROM bookings WHERE bookings.household_id = sqlc.arg(target_household_id))
	   OR dr.message_id IN (
		SELECT m.id FROM messages AS m
		JOIN conversations AS c ON c.id = m.conversation_id
		WHERE c.household_id = sqlc.arg(target_household_id)
	   )
   )
   OR message_id IN (
	SELECT m.id FROM messages AS m
	JOIN conversations AS c ON c.id = m.conversation_id
	WHERE c.household_id = sqlc.arg(target_household_id)
   );

-- name: DeleteHouseholdPaymentAllocations :exec
DELETE FROM payment_allocations
WHERE charge_id IN (SELECT id FROM charges WHERE household_id = sqlc.arg(target_household_id));

-- name: DeleteHouseholdOutboundMessages :exec
DELETE FROM outbound_messages
WHERE detected_request_id IN (
	SELECT dr.id
	FROM detected_requests AS dr
	WHERE dr.household_id = sqlc.arg(target_household_id)
	   OR dr.converted_booking_id IN (SELECT bookings.id FROM bookings WHERE bookings.household_id = sqlc.arg(target_household_id))
	   OR dr.message_id IN (
		SELECT m.id FROM messages AS m
		JOIN conversations AS c ON c.id = m.conversation_id
		WHERE c.household_id = sqlc.arg(target_household_id)
	   )
);

-- name: DeleteHouseholdDetectedRequests :exec
DELETE FROM detected_requests
WHERE detected_requests.household_id = sqlc.arg(target_household_id)
   OR converted_booking_id IN (SELECT bookings.id FROM bookings WHERE bookings.household_id = sqlc.arg(target_household_id))
   OR message_id IN (
	SELECT m.id FROM messages AS m
	JOIN conversations AS c ON c.id = m.conversation_id
	WHERE c.household_id = sqlc.arg(target_household_id)
   );

-- name: DeleteHouseholdMessages :exec
DELETE FROM messages
WHERE conversation_id IN (SELECT id FROM conversations WHERE household_id = sqlc.arg(target_household_id));

-- name: DeleteHouseholdConversations :exec
DELETE FROM conversations WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdCareTasks :exec
DELETE FROM care_tasks WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdCareRoutines :exec
DELETE FROM care_routines WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdCharges :exec
DELETE FROM charges WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdBookingPets :exec
DELETE FROM booking_pets
WHERE booking_id IN (SELECT bookings.id FROM bookings WHERE bookings.household_id = sqlc.arg(target_household_id))
   OR pet_id IN (SELECT pets.id FROM pets WHERE pets.household_id = sqlc.arg(target_household_id));

-- name: DeleteHouseholdBookings :exec
DELETE FROM bookings WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdPetMedications :exec
DELETE FROM pet_medications
WHERE pet_id IN (SELECT id FROM pets WHERE household_id = sqlc.arg(target_household_id));

-- name: DeleteHouseholdPetDiets :exec
DELETE FROM pet_diets
WHERE pet_id IN (SELECT id FROM pets WHERE household_id = sqlc.arg(target_household_id));

-- name: DeleteHouseholdPets :exec
DELETE FROM pets WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdContactLinks :exec
DELETE FROM household_contacts WHERE household_id = sqlc.arg(target_household_id);

-- name: DeleteHouseholdRecord :execrows
DELETE FROM households WHERE id = sqlc.arg(target_household_id);

-- name: DeleteContactIfOrphaned :exec
DELETE FROM contacts
WHERE id = sqlc.arg(contact_id)
  AND NOT EXISTS (SELECT 1 FROM household_contacts WHERE contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM bookings WHERE requested_by_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM payments WHERE payer_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM conversations WHERE primary_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM messages WHERE sender_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM detected_requests WHERE contact_id = ?1);

-- name: DeleteContactIdentitiesIfOtherwiseOrphaned :exec
DELETE FROM contact_channel_identities
WHERE contact_id = ?1
  AND NOT EXISTS (SELECT 1 FROM household_contacts WHERE contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM bookings WHERE requested_by_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM payments WHERE payer_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM conversations WHERE primary_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM messages WHERE sender_contact_id = ?1)
  AND NOT EXISTS (SELECT 1 FROM detected_requests WHERE contact_id = ?1);
