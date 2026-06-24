#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
source_workflow="$repo_root/automation/n8n/workflows/pawsear-message-webhook.json"
image="docker.n8n.io/n8nio/n8n:2.26.8"
api_port="${PAWSEAR_SMOKE_API_PORT:-18081}"
n8n_port="${PAWSEAR_SMOKE_N8N_PORT:-15678}"
container_name="pawsear-n8n-smoke-$$"
volume_name="pawsear-n8n-smoke-$$"
temp_dir="$(mktemp -d)"
workflow="$temp_dir/workflow-source.json"
api_pid=""

cleanup() {
	if [[ -n "$api_pid" ]]; then
		kill "$api_pid" 2>/dev/null || true
		wait "$api_pid" 2>/dev/null || true
	fi
	docker rm -f "$container_name" >/dev/null 2>&1 || true
	docker volume rm "$volume_name" >/dev/null 2>&1 || true
	rm -rf "$temp_dir"
}
trap cleanup EXIT

wait_for_url() {
	local url="$1"
	for _ in $(seq 1 60); do
		if curl --silent --fail "$url" >/dev/null; then
			return 0
		fi
		sleep 1
	done
	return 1
}

assert_duplicate() {
	python3 - "$1" "$2" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as response_file:
    response = json.load(response_file)
expected = sys.argv[2] == "true"
if response.get("duplicate") is not expected:
    raise SystemExit(f"expected duplicate={expected}, got {response}")
if response.get("detectedRequest", {}).get("status") != "needs_review":
    raise SystemExit(f"expected needs_review request, got {response}")
PY
}

post_webhook() {
	local output_file="$1"
	local payload="$2"
	local status
	status="$(curl --silent --show-error \
		-X POST "http://127.0.0.1:$n8n_port/webhook/pawsear-message" \
		-H 'Content-Type: application/json' \
		-d "$payload" \
		-o "$output_file" \
		-w '%{http_code}')"
	if [[ "$status" -lt 200 || "$status" -ge 300 ]]; then
		echo "n8n webhook returned HTTP $status" >&2
		cat "$output_file" >&2
		docker logs "$container_name" --tail 80 >&2
		return 1
	fi
}

echo "Starting temporary Pawsear API on port $api_port"
python3 - "$source_workflow" "$workflow" "$api_port" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as source_file:
    workflow = json.load(source_file)
for node in workflow["nodes"]:
    if node["name"] == "Import into Pawsear":
        node["parameters"]["url"] = f"http://host.docker.internal:{sys.argv[3]}/api/message-imports"
with open(sys.argv[2], "w", encoding="utf-8") as output_file:
    json.dump(workflow, output_file)
PY
(
	cd "$repo_root/apps/api"
	PAWSEAR_HTTP_ADDR=":$api_port" \
		PAWSEAR_DB_PATH="$temp_dir/pawsear.db" \
		GOCACHE="$repo_root/apps/api/.gocache" \
		go run ./cmd/server
) >"$temp_dir/api.log" 2>&1 &
api_pid=$!
wait_for_url "http://127.0.0.1:$api_port/health"

docker volume create "$volume_name" >/dev/null
docker run --rm \
	-e N8N_ENCRYPTION_KEY=pawsear-smoke-validation-only \
	-e N8N_LOG_LEVEL=error \
	-v "$volume_name:/home/node/.n8n" \
	-v "$workflow:/workflow.json:ro" \
	"$image" import:workflow --input=/workflow.json >/dev/null
docker run --rm \
	-e N8N_ENCRYPTION_KEY=pawsear-smoke-validation-only \
	-e N8N_LOG_LEVEL=error \
	-v "$volume_name:/home/node/.n8n" \
	"$image" publish:workflow --id=d58e3965-1869-4266-a3ca-29c5e9012952 >/dev/null
docker run --rm \
	-e N8N_ENCRYPTION_KEY=pawsear-smoke-validation-only \
	-e N8N_LOG_LEVEL=error \
	-v "$volume_name:/home/node/.n8n" \
	"$image" update:workflow --id=d58e3965-1869-4266-a3ca-29c5e9012952 --active=true >/dev/null
docker run --rm \
	-e N8N_ENCRYPTION_KEY=pawsear-smoke-validation-only \
	-e N8N_LOG_LEVEL=error \
	-v "$volume_name:/home/node/.n8n" \
	-v "$temp_dir:/output" \
	"$image" export:workflow \
		--id=d58e3965-1869-4266-a3ca-29c5e9012952 \
		--output=/output/workflow.json >/dev/null
python3 - "$temp_dir/workflow.json" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as workflow_file:
    workflow = json.load(workflow_file)[0]
if workflow.get("active") is not True or not workflow.get("activeVersionId"):
    raise SystemExit(f"workflow was not activated: {workflow}")
PY

docker run -d --rm \
	--name "$container_name" \
	--add-host host.docker.internal:host-gateway \
	-p "127.0.0.1:$n8n_port:5678" \
	-e N8N_ENCRYPTION_KEY=pawsear-smoke-validation-only \
	-e N8N_DIAGNOSTICS_ENABLED=false \
	-e N8N_PERSONALIZATION_ENABLED=false \
	-e N8N_SECURE_COOKIE=false \
	-v "$volume_name:/home/node/.n8n" \
	"$image" start >/dev/null
wait_for_url "http://127.0.0.1:$n8n_port/healthz/readiness"

payload='{"channel":"telegram","externalConversationId":"smoke-chat","externalMessageId":"smoke-message-1","senderExternalId":"smoke-sender","body":"Paseo mañana a las 8","sentAt":"2026-06-21T14:00:00Z"}'
post_webhook "$temp_dir/first.json" "$payload"
assert_duplicate "$temp_dir/first.json" false

post_webhook "$temp_dir/second.json" "$payload"
assert_duplicate "$temp_dir/second.json" true

echo "n8n webhook smoke test passed"
