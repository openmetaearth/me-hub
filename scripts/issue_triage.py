import os
import sys
import json
import time
import requests
from collections import defaultdict

REPO = os.environ["REPO"]
TOKEN = os.environ["GITHUB_TOKEN"]
DRY_RUN = os.environ.get("DRY_RUN", "false").lower() == "true"

HEADERS = {
    "Authorization": f"token {TOKEN}",
    "Accept": "application/vnd.github.v3+json",
    "X-GitHub-Api-Version": "2022-11-28",
}
BASE = "https://api.github.com"

def _request(method, path, **kwargs):
    resp = getattr(requests, method)(f"{BASE}{path}", headers=HEADERS, **kwargs)
    if resp.status_code == 429:
        reset = int(resp.headers.get("X-RateLimit-Reset", time.time() + 60))
        wait = max(reset - int(time.time()), 5)
        print(f"  Rate limited, waiting {wait}s...", flush=True)
        time.sleep(wait)
        resp = getattr(requests, method)(f"{BASE}{path}", headers=HEADERS, **kwargs)
    return resp

def gh_get(path, params=None):
    resp = _request("get", path, params=params)
    if resp.status_code not in (200, 201):
        print(f"  GET {path} => {resp.status_code}: {resp.text[:100]}", file=sys.stderr)
        return None
    return resp.json()

def gh_patch(path, body):
    if DRY_RUN:
        print(f"  [DRY RUN] PATCH {path}: {json.dumps(body)[:100]}", flush=True)
        return {}
    resp = _request("patch", path, json=body)
    if resp.status_code not in (200, 201):
        print(f"  PATCH {path} => {resp.status_code}: {resp.text[:100]}", file=sys.stderr)
        return None
    return resp.json()

def gh_post(path, body):
    if DRY_RUN:
        print(f"  [DRY RUN] POST {path}: {json.dumps(body)[:100]}", flush=True)
        return {}
    resp = _request("post", path, json=body)
    if resp.status_code not in (200, 201):
        print(f"  POST {path} => {resp.status_code}: {resp.text[:100]}", file=sys.stderr)
        return None
    return resp.json()

def set_labels(number, label_set):
    result = gh_patch(f"/repos/{REPO}/issues/{number}", {"labels": sorted(label_set)})
    if result is not None:
        names = sorted(l["name"] for l in result.get("labels", []))
        print(f"  #{number}: labels set to {names}", flush=True)

def add_comment(number, body):
    result = gh_post(f"/repos/{REPO}/issues/{number}/comments", {"body": body})
    if result is not None:
        print(f"  #{number}: comment posted", flush=True)

def get_all_issues():
    issues = []
    for state in ["open", "closed"]:
        page = 1
        while True:
            batch = gh_get(
                f"/repos/{REPO}/issues",
                params={"state": state, "per_page": 100, "page": page}
            )
            if not batch:
                break
            # Exclude pull requests
            batch = [i for i in batch if "pull_request" not in i]
            issues.extend(batch)
            print(f"  Fetched {state} page {page}: {len(batch)} issues", flush=True)
            if len(batch) < 100:
                break
            page += 1
            time.sleep(0.2)
    return issues

# ---- Classification rules ----

MEPASS_KW = [
    "me pass", "mepass", "daily check", "check-in", "checkin", "check in",
    "face verif", "facial verif", "human verif", "me 730",
    "invitation reward", "referral reward", "me app",
    "app crash", "app close", "login bug", "me id card",
    "mec token", "kyc verif", "identity verif",
    "verification loop", "daily reward", "check in reward",
    "telegram", "social account", "x account",
    "biometric", "camera issue", "notification bug", "daily streak",
    "unicorn", "p2p market", "me pass app", "android app", "ios app",
]

BLOCKCHAIN_KW = [
    "[bug bounty]", "fix(", "x/wstaking", "x/gravity", "x/sequencer",
    "x/megroup", "x/eibc", "x/kyc", "x/did", "x/rollapp", "x/wdistri",
    "x/wmint", "x/wnft", "x/wbank", "x/wgov", "x/dao", "x/bridgingfee",
    "[security]", "wstaking", "wmint", "wdistri", "wnft", "wbank", "wgov",
    "genesis import", "genesis export", "apphash", "appHash",
    "ante/", "wasm contract", "eip-712", "[high]", "[critical]",
    "[medium]", "[low]", "delayedack", "eibc", "bridge relayer",
    "bridge token", "attestation", "msgtransferregion", "msgremoveregion",
    "msgstake", "msgmint", "did query", "grpc query", "genesis state",
    "region treasury", "region share", "fixed deposit", "bridge fee",
    "double-sign", "slash", "rollapp", "ibc ", " ibc", "kvstore",
    "proposer", "relayer set", "bridge claim", "sequencer proposer",
]

INVALID_TITLES = {
    "s", "bug", "f", "kyc", "lemi", "me", "mec", "chat", "p2p",
    "[bug]:", "[bug]: ",
}

def get_labels(issue):
    return set(
        l["name"] if isinstance(l, dict) else str(l)
        for l in issue.get("labels", [])
    )

def fulltext(issue):
    return ((issue.get("title") or "") + " " + (issue.get("body") or "")).lower()

def is_blockchain(issue):
    t = fulltext(issue)
    return any(kw in t for kw in BLOCKCHAIN_KW)

def is_mepass(issue):
    t = fulltext(issue)
    return any(kw in t for kw in MEPASS_KW)

def is_spam(issue):
    # Minimum title length below which empty issues are considered non-actionable
    MIN_TITLE_LENGTH = 20
    title = (issue.get("title") or "").strip().lower()
    body = (issue.get("body") or "").strip()
    empty_body = not body
    return (
        title in INVALID_TITLES
        or len(title) <= 2
        or (empty_body and len(title) < MIN_TITLE_LENGTH
            and not title.startswith("[bug bounty]"))
    )

# ---- Main ----

print("Fetching all issues...", flush=True)
issues = get_all_issues()
print(f"Total: {len(issues)} issues\n", flush=True)

# Build title -> [issue_numbers] map for duplicate detection
title_map = defaultdict(list)
for i in issues:
    key = (i.get("title") or "").strip().lower()
    title_map[key].append(i["number"])

dup_groups = {t: sorted(nums) for t, nums in title_map.items() if len(nums) > 1}
print(f"Duplicate groups found: {len(dup_groups)}")
for t, nums in sorted(dup_groups.items()):
    print(f"  {nums}: '{t[:60]}'")
print()

updated = 0

for issue in issues:
    number = issue["number"]
    title = (issue.get("title") or "").strip()
    state = issue.get("state", "open")
    labels = get_labels(issue)
    title_key = title.lower()

    new_labels = set(labels)
    comments = []
    changes = []

    # --- Rule 1: Mark duplicates ---
    if title_key in dup_groups:
        group = dup_groups[title_key]
        if number != group[0] and "duplicate" not in labels:
            original = group[0]
            new_labels.add("duplicate")
            comments.append(
                f"This issue appears to be a duplicate of #{original}.\n\n"
                f"Labeling as **duplicate**. Please follow #{original} for updates."
            )
            changes.append("duplicate")

    # --- Rule 2: Mark spam/invalid ---
    if is_spam(issue) and "invalid" not in labels and "duplicate" not in new_labels:
        new_labels.add("invalid")
        comments.append(
            "This issue has been marked as **invalid** because it does not contain "
            "enough information to be actionable. "
            "Please reopen with a clear title and detailed description of the problem."
        )
        changes.append("invalid")

    # --- Rule 3: Add blockchain bug label ---
    if "invalid" not in new_labels and "duplicate" not in new_labels:
        chain = is_blockchain(issue)
        app = is_mepass(issue)

        if chain and "bug" not in labels:
            new_labels.add("bug")
            comments.append(
                "This issue has been triaged and labeled **`bug`** — "
                "it appears to report a defect in the ME Hub blockchain node. "
                "Thank you for your contribution!"
            )
            changes.append("bug")

        # --- Rule 4: Add mepass label ---
        if app and not chain and "mepass" not in labels:
            new_labels.add("mepass")
            if "bug" not in new_labels and "enhancement" not in new_labels:
                new_labels.add("bug")
            label_str = " and ".join(
                f"**`{l}`**"
                for l in sorted(new_labels - labels)
                if l not in {"duplicate", "invalid"}
            )
            if label_str:
                comments.append(
                    f"This issue has been triaged and labeled {label_str} — "
                    f"it appears to be related to the **ME Pass** mobile application.\n\n"
                    f"> **Note:** This repository (`me-hub`) hosts the ME Hub blockchain node. "
                    f"ME Pass app issues are tracked here with the `mepass` label for visibility."
                )
            changes.append("mepass")

    # --- Apply changes ---
    if changes:
        print(f"#{number} [{state}]: {title[:72]}", flush=True)
        print(f"  Changes: {changes}", flush=True)
        if new_labels != labels:
            set_labels(number, new_labels)
        for c in comments:
            add_comment(number, c)
            time.sleep(0.1)
        updated += 1

print(f"\n=== Triage complete: {updated} issues updated ===", flush=True)
if DRY_RUN:
    print("(Dry run mode — no actual changes were applied)")
