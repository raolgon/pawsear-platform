#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADDR="${PAWSEAR_SMOKE_ADDR:-127.0.0.1:18080}"
BASE_URL="http://${ADDR}"
DB_PATH="${PAWSEAR_SMOKE_DB:-$(mktemp -t pawsear-smoke-XXXXXX.db)}"
GOCACHE="${GOCACHE:-${ROOT_DIR}/.gocache}"

server_pid=""

cleanup() {
	if [[ -n "${server_pid}" ]]; then
		kill "${server_pid}" >/dev/null 2>&1 || true
		wait "${server_pid}" >/dev/null 2>&1 || true
	fi
	if [[ "${PAWSEAR_KEEP_SMOKE_DB:-0}" != "1" ]]; then
		rm -f "${DB_PATH}" "${DB_PATH}-shm" "${DB_PATH}-wal"
	fi
}
trap cleanup EXIT

request() {
	local method="$1"
	local path="$2"
	local body="${3:-}"

	if [[ -n "${body}" ]]; then
		curl -fsS -X "${method}" "${BASE_URL}${path}" \
			-H "Content-Type: application/json" \
			-d "${body}"
	else
		curl -fsS -X "${method}" "${BASE_URL}${path}"
	fi
}

expect_status() {
	local expected="$1"
	local method="$2"
	local path="$3"
	local body="${4:-}"
	local status

	if [[ -n "${body}" ]]; then
		status="$(curl -sS -o /dev/null -w "%{http_code}" -X "${method}" "${BASE_URL}${path}" \
			-H "Content-Type: application/json" \
			-d "${body}")"
	else
		status="$(curl -sS -o /dev/null -w "%{http_code}" -X "${method}" "${BASE_URL}${path}")"
	fi

	if [[ "${status}" != "${expected}" ]]; then
		echo "Expected ${method} ${path} to return ${expected}, got ${status}" >&2
		exit 1
	fi
}

json_id() {
	sed -n 's/.*"id":"\([^"]*\)".*/\1/p'
}

require_id() {
	local label="$1"
	local value="$2"
	if [[ -z "${value}" ]]; then
		echo "Missing id for ${label}" >&2
		exit 1
	fi
}

echo "Starting Pawsear API smoke server on ${ADDR}"
(
	cd "${ROOT_DIR}"
	PAWSEAR_HTTP_ADDR="${ADDR}" PAWSEAR_DB_PATH="${DB_PATH}" GOCACHE="${GOCACHE}" go run ./cmd/server
) &
server_pid="$!"

for _ in {1..40}; do
	if curl -fsS "${BASE_URL}/health" >/dev/null 2>&1; then
		break
	fi
	sleep 0.25
done

request GET /health >/dev/null
request GET /api/meta >/dev/null

household="$(request POST /api/households '{"displayName":"Casa Smoke","neighborhood":"Roma Norte","city":"CDMX"}')"
household_id="$(printf '%s' "${household}" | json_id)"
require_id "household" "${household_id}"

contact="$(request POST /api/contacts '{"displayName":"Ana Smoke","phone":"+525500000001"}')"
contact_id="$(printf '%s' "${contact}" | json_id)"
require_id "contact" "${contact_id}"

request POST "/api/households/${household_id}/contacts" "{\"contactId\":\"${contact_id}\",\"role\":\"owner\",\"isPrimary\":true}" >/dev/null
request GET "/api/households/${household_id}/contacts" | grep -q '"role":"owner"'

pet="$(request POST /api/pets "{\"householdId\":\"${household_id}\",\"name\":\"Mora Smoke\",\"species\":\"dog\"}")"
pet_id="$(printf '%s' "${pet}" | json_id)"
require_id "pet" "${pet_id}"

staff="$(request POST /api/staff '{"displayName":"Rafa Smoke","role":"walker"}')"
staff_id="$(printf '%s' "${staff}" | json_id)"
require_id "staff" "${staff_id}"

booking="$(request POST /api/bookings "{
	\"householdId\":\"${household_id}\",
	\"serviceType\":\"walk\",
	\"status\":\"confirmed\",
	\"startAt\":\"2026-06-05T15:00:00Z\",
	\"assignedStaffId\":\"${staff_id}\",
	\"petIds\":[\"${pet_id}\"]
}")"
booking_id="$(printf '%s' "${booking}" | json_id)"
require_id "booking" "${booking_id}"

request GET "/api/bookings/${booking_id}" | grep -q '"pets":'

task="$(request POST /api/care-tasks "{
	\"bookingId\":\"${booking_id}\",
	\"householdId\":\"${household_id}\",
	\"petId\":\"${pet_id}\",
	\"taskType\":\"walk\",
	\"title\":\"Paseo smoke\",
	\"dueAt\":\"2026-06-05T15:00:00Z\"
}")"
task_id="$(printf '%s' "${task}" | json_id)"
require_id "care task" "${task_id}"

charge="$(request POST /api/charges "{
	\"householdId\":\"${household_id}\",
	\"bookingId\":\"${booking_id}\",
	\"description\":\"Paseo smoke\",
	\"amountMinor\":20000
}")"
charge_id="$(printf '%s' "${charge}" | json_id)"
require_id "charge" "${charge_id}"

payment="$(request POST /api/payments "{
	\"payerContactId\":\"${contact_id}\",
	\"receivedAt\":\"2026-06-05T18:00:00Z\",
	\"amountMinor\":20000,
	\"method\":\"cash\",
	\"allocations\":[{\"chargeId\":\"${charge_id}\",\"amountMinor\":20000}]
}")"
payment_id="$(printf '%s' "${payment}" | json_id)"
require_id "payment" "${payment_id}"

request GET "/api/payments/${payment_id}" | grep -q '"allocations":'
request GET "/api/charges/${charge_id}" | grep -q '"status":"paid"'
request GET '/api/dashboard/today?date=2026-06-05' | grep -q '"bookings":'

expect_status 400 POST /api/households '{"displayName":""}'
expect_status 400 POST /api/bookings "{\"householdId\":\"${household_id}\",\"startAt\":\"2026-06-05T16:00:00Z\",\"petIds\":[\"missing_pet\"]}"
expect_status 400 POST /api/charges "{\"householdId\":\"${household_id}\",\"description\":\"bad charge\",\"amountMinor\":0}"
expect_status 400 POST /api/payments "{\"amountMinor\":100,\"method\":\"cash\",\"allocations\":[{\"chargeId\":\"${charge_id}\",\"amountMinor\":200}]}"

echo "Smoke test passed"
