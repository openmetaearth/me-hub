#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/common.sh"
case ${1:-all} in
 all)
  query params
  query proposals --limit 100
  query proposals --status voting --kind treasury --limit 100
  query proposals --kind replace-operator --limit 100
  query equity "$REGION_ID"
  query region-freeze "$REGION_ID"
  query equity-freezes "$REGION_ID"
  query pending-actions --limit 100
  query pending-actions --status pending --limit 100
  ;;
 proposal)
  [[ $# == 2 ]] || fail 'usage: queries.sh proposal ID'
  query proposal "$2"
  query votes "$2" --limit 100
  query electors "$2"
  ;;
 action)
  [[ $# == 2 ]] || fail 'usage: queries.sh action ID'
  query pending-action "$2"
  ;;
 *) query "$@" ;;
esac
