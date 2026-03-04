#!/bin/sh
#
# Generate a CHANGELOG.md from conventional commits.
# Groups commits by tag (release) and by type (feat/fix/chore/etc).
#

set -e

echo "# Changelog"
echo ""

# Get all tags sorted by version (newest first), plus HEAD for unreleased.
tags=$(git tag --sort=-version:refname 2>/dev/null)

if [ -n "$tags" ]; then
    ranges="HEAD...$(echo "$tags" | head -1)"
    prev=""
    for tag in $tags; do
        if [ -n "$prev" ]; then
            ranges="$ranges
${prev}...${tag}"
        fi
        prev="$tag"
    done
    # Add the oldest tag (from beginning).
    ranges="$ranges
${prev}"
else
    ranges="HEAD"
fi

print_section() {
    type_label="$1"
    type_prefix="$2"
    commits="$3"

    matched=$(echo "$commits" | grep "^[a-f0-9]* ${type_prefix}" 2>/dev/null || true)
    if [ -z "$matched" ]; then
        return
    fi

    echo "### $type_label"
    echo ""
    echo "$matched" | while IFS= read -r line; do
        hash=$(echo "$line" | cut -d' ' -f1)
        msg=$(echo "$line" | sed "s/^[a-f0-9]* ${type_prefix}: //")
        echo "- $msg (\`$hash\`)"
    done
    echo ""
}

first=true
echo "$ranges" | while IFS= read -r range; do
    # Determine the label and git log range.
    case "$range" in
        HEAD...*|HEAD)
            label="Unreleased"
            if echo "$range" | grep -q '\.\.\.'; then
                ref="$range"
            else
                ref="HEAD"
            fi
            ;;
        *...*)
            label=$(echo "$range" | cut -d'.' -f1)
            ref="$range"
            ;;
        *)
            label="$range"
            ref="$range"
            ;;
    esac

    # Get commits for this range (skip bot commits).
    if echo "$ref" | grep -q '\.\.\.'; then
        commits=$(git log --oneline "$ref" --no-merges \
            --grep="^chore: rebuild" --grep="^chore: update changelog" --invert-grep 2>/dev/null || true)
    else
        commits=$(git log --oneline "$ref" --no-merges \
            --grep="^chore: rebuild" --grep="^chore: update changelog" --invert-grep 2>/dev/null || true)
    fi

    if [ -z "$commits" ]; then
        continue
    fi

    # Get date for label.
    if [ "$label" = "Unreleased" ]; then
        date_str=$(date +%Y-%m-%d)
    else
        date_str=$(git log -1 --format="%cs" "$label" 2>/dev/null || date +%Y-%m-%d)
    fi

    echo "## $label ($date_str)"
    echo ""

    print_section "Features" "feat" "$commits"
    print_section "Bug Fixes" "fix" "$commits"
    print_section "Chores" "chore" "$commits"
    print_section "Docs" "docs" "$commits"
    print_section "Refactor" "refactor" "$commits"

    # Catch anything else.
    other=$(echo "$commits" | grep -v "^[a-f0-9]* feat:" | grep -v "^[a-f0-9]* fix:" \
        | grep -v "^[a-f0-9]* chore:" | grep -v "^[a-f0-9]* docs:" \
        | grep -v "^[a-f0-9]* refactor:" 2>/dev/null || true)
    if [ -n "$other" ]; then
        echo "### Other"
        echo ""
        echo "$other" | while IFS= read -r line; do
            hash=$(echo "$line" | cut -d' ' -f1)
            msg=$(echo "$line" | cut -d' ' -f2-)
            echo "- $msg (\`$hash\`)"
        done
        echo ""
    fi
done
