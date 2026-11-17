#!/bin/sh
# Gives this fork's own migrations new numbers after the upstream ones when
# an upstream update took their number. A fork migration is one whose file
# is not in the upstream tree ($1, for example upstream/main). The database
# code moves the record of an applied, renumbered migration by its name.
set -eu
upstream="$1"
changed=0
for dir in app/backend/internal/core/database/migrations/*/; do
	up="$(git ls-tree --name-only "$upstream" "$dir" | xargs -n1 basename 2>/dev/null || true)"
	for f in $(ls "$dir" | grep '\.sql$' | sort); do
		echo "$up" | grep -qx "$f" && continue
		num="${f%%_*}"
		if echo "$up" | grep -q "^${num}_"; then
			last="$(ls "$dir" | grep '\.sql$' | sed 's/_.*//' | sort -n | tail -1)"
			next="$(printf '%04d' $(( $(echo "$last" | sed 's/^0*//') + 1 )))"
			git mv "$dir$f" "$dir${next}_${f#*_}"
			echo "renumbered $dir$f -> ${next}_${f#*_}"
			changed=1
		fi
	done
done
echo "changed=$changed"
