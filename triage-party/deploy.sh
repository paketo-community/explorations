#!/usr/bin/env bash

set -eu
set -o pipefail

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function main(){
  while [[ "${#}" != 0 ]]; do
    case "${1}" in
      --help|-h)
        shift 1
        usage
        exit 0
        ;;

      "")
        # skip if the argument is empty
        shift 1
        ;;

      *)
        util::print::error "unknown argument \"${1}\""
    esac
  done

  check_env_vars

  gcloud beta run deploy triage-party \
    --project "${PROJECT_ID}" \
    --image "${IMAGE}" \
    --set-env-vars="GITHUB_TOKEN=${GITHUB_TOKEN},PERSIST_BACKEND=cloudsql,PERSIST_PATH=${DB_USER}:${DB_PASS}@tcp(${PROJECT_ID}/us-central1/triage-party)/tp" \
        --allow-unauthenticated \
    --region us-central1 \
    --platform managed
}

function usage() {
  cat <<-USAGE
deploy.sh [OPTIONS]

Deploys the Paketo triage party instance with collocated config.
Script requires that the following environment variables are set:
  DB_PASS - The password to access the Google Cloud SQL instance
  DB_USER - The username to access the Google Cloud SQL instance
  GITHUB_TOKEN – The token that will be used to scrape Github data
  IMAGE - The URI of the triage party image to deploy

OPTIONS
  --help       -h  prints the command usage

USAGE
}

function check_env_vars() {
  if [[ -z "${GITHUB_TOKEN-}" ]]; then
    echo -e "Environment variable \$GITHUB_TOKEN is required" >&2
    exit 1
  fi

  if [[ -z "${DB_USER-}" ]]; then
    echo -e "Environment variable \$DB_USER is required" >&2
    exit 1
  fi

  if [[ -z "${DB_PASS-}" ]]; then
    echo -e "Environment variable \$DB_PASS is required" >&2
    exit 1
  fi

  if [[ -z "${IMAGE-}" ]]; then
    echo -e "Environment variable \$IMAGE is required" >&2
    exit 1
  fi
}

main "${@:-}"
