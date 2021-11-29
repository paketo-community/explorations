# Triage Party Deployment

Requirements:
- Local clone of the `google/triage-party` repo
- Github token with `public_repo` permissions (read-only)

Steps:
1. `echo <YOUR_GENERATED_TOKEN> > $HOME/.github-token`
1. Run:
`go run <path/to/local/triage-party/repo/cmd/server/main.go> --github-token-file=$HOME/.github-token --config ./triage-party-config.yml`




