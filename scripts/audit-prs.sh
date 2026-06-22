#!/usr/bin/env bash
# Audit open PRs and generate audit-report.md
# Usage: ./scripts/audit-prs.sh [--limit N]
# Output: audit-report.md + .pr-close-list.txt in repo root

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REPORT="$REPO_ROOT/audit-report.md"
TMP_JSON="$(mktemp /tmp/prs-XXXXXX.json)"
LIMIT=300

if [[ "${1:-}" == "--limit" && -n "${2:-}" ]]; then
  LIMIT="$2"
fi

echo "Fetching open PRs (limit=$LIMIT)..."
# Omit statusCheckRollup to avoid GraphQL timeout on large result sets.
gh pr list --limit "$LIMIT" --state open \
  --json number,title,author,mergeable,reviewDecision,\
comments,createdAt,isDraft,headRefName,changedFiles,additions,deletions \
  > "$TMP_JSON"

TOTAL=$(jq 'length' "$TMP_JSON")
echo "Fetched $TOTAL PRs."

# Detect bots via GitHub API: check each unique author's account type.
# GitHub Apps use "app/<name>" login format — always bots, no API call needed.
echo "Detecting bot accounts via GitHub API..."
BOT_ACCOUNTS=()
UNIQUE_AUTHORS=$(jq -r '[.[].author.login] | unique[]' "$TMP_JSON")
while IFS= read -r login; do
  if [[ "$login" == app/* ]]; then
    BOT_ACCOUNTS+=("$login")
    echo "  $login → Bot (GitHub App)"
  else
    acct_type=$(gh api "users/$login" --jq '.type' 2>/dev/null || echo "Unknown")
    echo "  $login → $acct_type"
    if [[ "$acct_type" == "Bot" ]]; then
      BOT_ACCOUNTS+=("$login")
    fi
  fi
done <<< "$UNIQUE_AUTHORS"

echo "Bot accounts detected: ${BOT_ACCOUNTS[*]:-none}"

# Build jq-compatible JSON array of bot logins
BOT_LIST=$(printf '"%s",' "${BOT_ACCOUNTS[@]:-}" | sed 's/,$//')
[[ -z "$BOT_LIST" ]] && BOT_LIST=""

jq --argjson bots "[$BOT_LIST]" '
  def is_bot: .author.login as $a | $bots | any(. == $a);

  def has_engagement:
    (.comments | length) > 0 or
    (.reviewDecision != null and .reviewDecision != "");

  def close_reason:
    if .mergeable == "CONFLICTING" then "merge conflict"
    elif (is_bot and (has_engagement | not)) then "bot/AI PR, no human engagement"
    elif (.isDraft and (has_engagement | not)) then "stale draft, no engagement"
    else "other"
    end;

  def should_close:
    (.mergeable == "CONFLICTING") or
    (is_bot and (has_engagement | not)) or
    (.isDraft and (has_engagement | not));

  . as $all |
  {
    auto_close: [
      $all[] | select(should_close) |
      {number, title, author: .author.login, reason: close_reason,
       mergeable, isDraft, changedFiles, additions, deletions,
       createdAt: (.createdAt | split("T")[0])}
    ],
    keep_open: [
      $all[] | select(should_close | not) |
      {number, title, author: .author.login,
       reason: (if (is_bot | not) then "human contributor" else "bot PR with engagement" end),
       mergeable, isDraft, changedFiles, additions, deletions,
       createdAt: (.createdAt | split("T")[0])}
    ]
  }
' "$TMP_JSON" > /tmp/categorized.json

AUTO_CLOSE_COUNT=$(jq '.auto_close | length' /tmp/categorized.json)
KEEP_COUNT=$(jq '.keep_open | length' /tmp/categorized.json)

echo "Auto-close candidates: $AUTO_CLOSE_COUNT"
echo "Keep open (manual review): $KEEP_COUNT"
echo "Generating $REPORT ..."

# Markdown header
cat > "$REPORT" << HEADER
# PR Audit Report

Generated: $(date -u '+%Y-%m-%d %H:%M UTC')
Total open PRs fetched: **$TOTAL**

## Summary

| Category | Count |
|----------|-------|
| Auto-close candidates | **$AUTO_CLOSE_COUNT** |
| Keep open (manual review) | **$KEEP_COUNT** |

**Auto-close criteria:**
- Has merge conflicts
- Bot/AI account (detected via GitHub API type=Bot or app/ prefix) with zero human engagement
- Stale draft PR with zero engagement

**Bot accounts detected:** ${BOT_ACCOUNTS[*]:-none}

**Keep criteria (any one is enough to preserve):**
- Author is not a known bot account
- PR has at least one comment or review

> **Next steps:**
> 1. Review the Auto-Close list below, remove any lines you want to preserve
> 2. Dry-run: \`./scripts/close-prs.sh --dry-run\`
> 3. Execute: \`./scripts/close-prs.sh\`

---

## Auto-Close List ($AUTO_CLOSE_COUNT PRs)

| # | Title | Author | Reason | Date | Files |
|---|-------|--------|--------|------|-------|
HEADER

jq -r '.auto_close[] | "| #\(.number) | \(.title | .[0:55]) | \(.author) | \(.reason) | \(.createdAt) | \(.changedFiles) |"' \
  /tmp/categorized.json >> "$REPORT"

cat >> "$REPORT" << 'KEEP_HEADER'

---

## Keep Open List (manual review)

| # | Title | Author | Reason | Date | Files |
|---|-------|--------|--------|------|-------|
KEEP_HEADER

jq -r '.keep_open[] | "| #\(.number) | \(.title | .[0:55]) | \(.author) | \(.reason) | \(.createdAt) | \(.changedFiles) |"' \
  /tmp/categorized.json >> "$REPORT"

# Write number list for close-prs.sh
jq -r '.auto_close[].number' /tmp/categorized.json > "$REPO_ROOT/.pr-close-list.txt"

rm -f "$TMP_JSON" /tmp/categorized.json
echo "Done."
echo "  Report:     $REPORT"
echo "  Close list: $REPO_ROOT/.pr-close-list.txt ($AUTO_CLOSE_COUNT entries)"
