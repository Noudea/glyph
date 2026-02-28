#!/usr/bin/env bash
set -euo pipefail

# Builds spellbooks/registry.json from individual spellbook.json manifests.
# Usage: ./scripts/build-registry.sh

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SPELLBOOKS_DIR="$REPO_ROOT/spellbooks"
REGISTRY="$SPELLBOOKS_DIR/registry.json"

# Collect all spellbook.json files into a single registry object.
registry='{}'
for manifest in "$SPELLBOOKS_DIR"/*/spellbook.json; do
  [ -f "$manifest" ] || continue
  id="$(basename "$(dirname "$manifest")")"
  registry="$(jq --arg id "$id" --slurpfile sb "$manifest" '.[$id] = $sb[0]' <<< "$registry")"
done

jq -S '{"spellbooks": .}' <<< "$registry" > "$REGISTRY"
echo "registry: wrote $REGISTRY ($(jq '.spellbooks | length' "$REGISTRY") spellbooks)"
