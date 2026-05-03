#!/usr/bin/env bash
# Rename the Go module path and the container image name throughout the
# template. Run once after cloning.
#
#   ./scripts/rename-module.sh github.com/you/yourapp [image-name]
#
# If image-name is omitted it defaults to the basename of the module path.

set -euo pipefail

if [[ $# -lt 1 || $# -gt 2 ]]; then
    echo "Usage: $0 <new-module-path> [image-name]"
    echo "Example: $0 github.com/alice/widgets widgets"
    exit 1
fi

NEW_MODULE="$1"
NEW_IMAGE="${2:-$(basename "$NEW_MODULE")}"
OLD_MODULE="github.com/youruser/yourproject"
OLD_IMAGE="yourproject"

cd "$(dirname "$0")/.."

if ! grep -q "^module $OLD_MODULE$" go.mod; then
    echo "go.mod does not contain $OLD_MODULE — already renamed?" >&2
    exit 1
fi

echo "Module: $OLD_MODULE -> $NEW_MODULE"
echo "Image:  $OLD_IMAGE -> $NEW_IMAGE"
echo

go mod edit -module="$NEW_MODULE"

# Files that legitimately reference the placeholder. Keep the list explicit
# rather than scanning everything: avoids touching pb_data, lockfiles, etc.
mapfile -t TARGETS < <(grep -rl --null \
    --include='*.go' \
    --include='Taskfile.yml' \
    --include='Containerfile' \
    --include='compose.yml' \
    --include='*.md' \
    -e "$OLD_MODULE" -e "\b$OLD_IMAGE\b" . | tr '\0' '\n' | grep -v '^./scripts/rename-module.sh$' || true)

for f in "${TARGETS[@]}"; do
    [[ -z "$f" ]] && continue
    sed -i \
        -e "s|$OLD_MODULE|$NEW_MODULE|g" \
        -e "s|\\b$OLD_IMAGE\\b|$NEW_IMAGE|g" \
        "$f"
done

echo "Updated ${#TARGETS[@]} file(s)."
echo
echo "Next steps:"
echo "  1. go mod tidy"
echo "  2. git diff   # review changes"
echo "  3. git add -A && git commit -m 'rename module to $NEW_MODULE'"
