# Pawsear n8n foundation

n8n transports external events into Pawsear. It must not write to the Pawsear SQLite database or create confirmed bookings directly.

## Start the complete local stack from WSL2

Use a Docker Engine available inside WSL2:

```sh
cd pawsear-platform
DOCKER_BUILDKIT=0 docker compose --env-file automation/n8n/.env -f compose.local.yaml up -d
```

Open Pawsear at `http://localhost:5173` and n8n at `http://localhost:5678`.
Complete the local n8n setup and import the workflow you need from `automation/n8n/workflows/`.

On first start, n8n generates its encryption key inside the persistent `n8n_data` volume. Keep that volume: removing it deletes the local n8n account, credentials, workflows, and encryption key. The service binds only to `127.0.0.1`.

Daily commands from the repository root:

```sh
docker compose --env-file automation/n8n/.env -f compose.local.yaml up -d
docker compose --env-file automation/n8n/.env -f compose.local.yaml ps
docker compose --env-file automation/n8n/.env -f compose.local.yaml logs -f n8n
docker compose --env-file automation/n8n/.env -f compose.local.yaml stop
```

Avoid `docker compose down -v`: the `-v` option deletes the persistent n8n volume.

`scripts/webhook_proxy.py` is a local safety boundary for tunnels. It forwards only n8n's production and test webhook paths and returns 404 for the editor and every other route.

The workflow calls Pawsear through the shared Docker network at
`http://api:8080/api/message-imports`. Do not use `localhost` in the HTTP Request
node: from n8n it means the n8n container, not Pawsear. The HTTP Request node reads
only the injected `PAWSEAR_AUTOMATION_TOKEN` to build its bearer header; the secret is
not stored in workflow JSON.

## Protect message imports

Generate or rotate the ignored local token from the repository root:

```sh
automation/n8n/scripts/configure_automation_token.sh
docker compose --env-file automation/n8n/.env -f compose.local.yaml up -d
```

Compose injects the same value into API, web, and n8n. Rotating it requires recreating
those containers with the second command. Requests without the bearer token receive
`401`; n8n supplies it automatically.

Never store the token in the workflow JSON or commit a real `.env` file.

## Test payload

After activating the workflow, send a unique external message ID:

```sh
curl -X POST http://localhost:5678/webhook/pawsear-message \
  -H 'Content-Type: application/json' \
  -d '{
    "channel": "telegram",
    "externalConversationId": "local-test-chat",
    "externalMessageId": "local-test-message-1",
    "senderExternalId": "local-test-sender",
    "body": "¿Puedes pasear a Luna mañana a las 8?",
    "sentAt": "2026-06-20T14:00:00Z"
  }'
```

Sending the same payload again is safe: Pawsear returns the existing message with `duplicate: true`.

For a disposable end-to-end check that starts a temporary Pawsear API and n8n instance, run from the repository root:

```sh
bash automation/n8n/scripts/smoke.sh
```

The script uses temporary SQLite and n8n data, verifies the first import and its duplicate retry, then removes its container, volume, and temporary files. It does not use or modify real credentials.

## Telegram ingestion

Import `workflows/pawsear-telegram-ingestion.json` and create a Telegram API credential in n8n using a bot token obtained from BotFather. Assign it to `Telegram Trigger`.

Telegram must reach n8n through a public HTTPS webhook. Keep the base Compose file local while developing. When a trusted HTTPS reverse proxy or tunnel is ready, set `N8N_WEBHOOK_URL` in the ignored `.env` file and start with the external-webhook override:

```sh
docker compose -f compose.yaml -f compose.telegram.yaml up -d
```

Do not activate the Telegram workflow until the public URL, TLS, Pawsear automation token, and n8n owner account are configured. The workflow ignores messages without text or a caption, retries API delivery three times, and relies on Pawsear idempotency for safe retries.

## Telegram replies

Import `workflows/pawsear-telegram-outbound.json`, assign the existing Telegram API
credential to `Send Telegram reply`, and activate it. Every minute it loads up to 20
operator-approved pending replies, sends them, and marks successful delivery in Pawsear.
The workflow contains no bot token or Pawsear token.

Reply content is selected in Pawsear, not n8n. A failed Telegram node leaves the item
pending so a later scheduled execution can retry it.

WhatsApp should follow only after choosing an official provider and documenting its webhook verification requirements.
