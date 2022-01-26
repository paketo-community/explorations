#!/usr/bin/env bash

set -eu
set -o pipefail

readonly DB_INSTANCE="${DB_INSTANCE:-triage-party}"
readonly DB_NAME="${DB_NAME:-tpdb}"
readonly DB_REGION="${DB_REGION:-us-central1}"
readonly CLOUD_RUN_REGION="${CLOUD_RUN_REGION:-us-central1}"

function main(){
  local initial="false"
  while [[ "${#}" != 0 ]]; do
    case "${1}" in
      --help|-h)
        shift 1
        usage
        exit 0
        ;;

      --init|-i)
        shift 1
        initial="true"
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
  if [[ "${initial}" == "true" ]]; then
    gcloud run deploy triage-party \
      --project "${GCP_PROJECT_ID}" \
      --region "${CLOUD_RUN_REGION}" \
      --image "${IMAGE}" \
      --set-env-vars="GITHUB_TOKEN=${GITHUB_TOKEN},PERSIST_BACKEND=cloudsql,PERSIST_PATH=${DB_USER}:${DB_PASS}@tcp(${GCP_PROJECT_ID}/${DB_REGION}/${DB_INSTANCE})/${DB_NAME}" \
      --allow-unauthenticated \
      --platform managed
    exit 0
  fi

  gcloud run services update triage-party \
    --project "${GCP_PROJECT_ID}" \
    --region us-central1 \
    --image "${IMAGE}" \
    --set-env-vars="GITHUB_TOKEN=${GITHUB_TOKEN},PERSIST_BACKEND=cloudsql,PERSIST_PATH=${DB_USER}:${DB_PASS}@tcp(${GCP_PROJECT_ID}/${DB_REGION}/${DB_INSTANCE})/${DB_NAME}"
}

function usage() {
  cat <<-USAGE
deploy.sh [OPTIONS]

Updates the Paketo triage party instance with a prebuilt image.

Script requires that the following environment variables are set:
  DB_PASS - The password to access the Google Cloud SQL instance
  DB_USER - The username to access the Google Cloud SQL instance
  GITHUB_TOKEN – The token that will be used to scrape Github data
  IMAGE - The URI of the triage party image to deploy

A CloudSQL MySQL instance for persisting Triage Party data must already exist
with the name '${DB_INSTANCE}' and contain a database called '${DB_NAME}'.

OPTIONS
  --help       -h  prints the command usage
  --init       -i  performs an initial gcloud run deploy

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

  if [[ -z "${GCP_PROJECT_ID-}" ]]; then
    echo -e "Environment variable \$GCP_PROJECT_ID is required" >&2
    exit 1
  fi
}

main "${@:-}"
