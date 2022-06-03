Feedback:
- Let's make the Github Actions less confusing
- Let's get rid of known-versions.json and rely on the `buildpack.toml` as a
  source of truth
- Let's get rid of metadata.json and rely on the `buildpack.toml` as a source
  of truth
- Let's just leave the old dependencies where they are

- Means that the dep-server doesn't need to work for new versions, just needs to support old versions
- Means that the `jam update-dependency` command can go away

Example buildpack.toml for Go if we support using the dependency from upstream
directly, and also support other stacks:

```
...
[metadata]
  [[metadata.dependencies]]
    cpe = "cpe:2.3:a:golang:go:1.17.10:*:*:*:*:*:*:*"
    id = "go"
    licenses = ["BSD-3-Clause"]
    name = "Go"
    purl = "pkg:generic/go@go1.17.10?checksum=amd64-sha"
    sha256 = "amd64-source-sha"
    source = "https://go.dev/dl/go1.17.10.linux-amd64.tar.gz"
    source_sha256 = "amd64-source-sha"
    stacks = ["io.buildpacks.stacks.bionic", "io.paketo.stacks.tiny", "io.buildpacks.stacks.jammy.tiny"]
    uri = "https://go.dev/dl/go1.17.10.linux-amd64.tar.gz"
    version = "1.17.10"

  [[metadata.dependencies]]
    cpe = "cpe:2.3:a:golang:go:1.17.10:*:*:*:*:*:*:*"
    id = "go"
    licenses = ["BSD-3-Clause"]
    name = "Go"
    purl = "pkg:generic/go@go1.17.10?checksum=arm64-sha"
    sha256 = "arm64-source-sha"
    source = "https://go.dev/dl/go1.17.10.linux-arm64.tar.gz"
    source_sha256 = "arm64-source-sha"
    stacks = [io.buildpacks.stacks.bionic.arm64, io.buildpacks.stacks.jammy.arm64]
    uri = "https://go.dev/dl/go1.17.10.linux-arm64.tar.gz"
    version = "1.17.10"

  [[metadata.dependency-constraints]]
    constraint = "1.17.*"
    id = "go"
    patches = 2

[[stacks]]
  id = "*"
```

From this, we see that there are two entries for the same version, since the
source URI is different dependening on the stack architecture.

--------- 

Workflows: a single workflow with multiple jobs

Job 1: Get New Versions
- Use buildpack `version retrieval` code to discover new versions
- Get a list of all versions, sorted (Go example: https://golang.org/doc/devel/release.html)
- Some riff on the code to diff against known versions:
  - get versions in the constraints, determine if there's a version in the
    constraint (1.17.*) higher than what we have (greater than 1.17.10)
  - if a in-constraint higher version is found add it to a matrix
  - check for any versions  greater than the constraints we have (greater than 1.17.*)
    - if a higher version line exists: file an issue that a higher version line exists

- code for version retrieval lives in buildpack
- code for comparing output to buildpack toml and spitting out diff should be common

Output: list of new versions: [1.17.11, 1.17.12]

Job 2: Jam Update, takes in array of versions: [1.17.11, 1.17.12]
- `jam update-dependencies` updated to take in an array of versions, buildpack.toml path,
  and path to metadata-gathering code.
  - CPE, PURL, licenses, stacks, version, source URI, source SHA256, etc.
  - The URI and SHA256 will be left off if the dependency needs to be compiled
  - Updates the `buildpack.toml` if the version is newer than what's in the buildpack.toml
  - TODO: Will need to be modified to work for multiple stacks
- open a draft PR with `buildpack.toml` updates

- code to gather metadata live in buildpack
- jam code is common

Workflow: PR validation
- Check the PR contains the SHA256 and URI fields
- Un-draft the PR if the content is good to go

- If the PR does not contain SHA256 and URI:
- Get all versions, and source URIs in a 2-dimensional matrix
  - TODO account for compilation across multiple stacks
  - For each element of matrix run dependency compilation code from the
    buildpack, take in the version, source URI
  - Run a smoke test against the compiled dependency
  - Upload it to GCP Bucket
  - Get SHA256 of dependency
  - Update the PR with a commit to add the SHA256/bucket access URI

- Re-run validation workflow

- validation code is a common action
- dependency compilation code lives in buildpack
- smoke test lives in buildpack
- upload code is an action
- update the PR is an action
