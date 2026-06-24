#!/bin/sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
env_file="$script_dir/../.env"
token=$(openssl rand -hex 32)
temporary_file=$(mktemp)
trap 'rm -f "$temporary_file"' EXIT

if [ -f "$env_file" ]; then
	awk -v token="$token" '
		BEGIN { replaced = 0 }
		/^PAWSEAR_AUTOMATION_TOKEN=/ {
			print "PAWSEAR_AUTOMATION_TOKEN=" token
			replaced = 1
			next
		}
		{ print }
		END {
			if (!replaced) print "PAWSEAR_AUTOMATION_TOKEN=" token
		}
	' "$env_file" > "$temporary_file"
else
	printf 'PAWSEAR_AUTOMATION_TOKEN=%s\n' "$token" > "$temporary_file"
fi

chmod 600 "$temporary_file"
mv "$temporary_file" "$env_file"
trap - EXIT
printf '%s\n' 'Pawsear automation token configured without displaying it.'
