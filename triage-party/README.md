# Triage Party

This directory contains configuration and scripts that deploy an instance of
[triage party](https://github.com/google/triage-party) for Paketo's
repositories.

## Viewing/Using Triage Party
The current instance can be viewed here: https://triage-party-ngfy5sbarq-uc.a.run.app/s/weekly

## Deploying Triage Party

The triage party instance is hosted with Google Cloud Run. It uses a MySQL
backend to persist some data. The workflow
[`.github/workflows/deploy-triage-party.yml`](../.github/workflows/deploy-triage-party.yml)
builds a triage party Docker image using [`config.yml`](./config.yml), pushes
it to a Paketo GCR repository, and deploys the image to Google Cloud Run,
connected to its database backend.

To deploy a fresh Triage Party instance:
1. Create and push a new triage-party image.
1. Create Google Cloud SQL MySQL instance. Add an empty database for triage
   party data to that instance.
1. Run `./deploy.sh --init` with the appropriate environment variables
   configured to create the new service.
