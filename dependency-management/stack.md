With the addition of Jammy stacks, multi-architecture stacks, we will need to
support multiple versions of each dependency in the buildpacks.

//TODO check output images in oci archive output from jam create-stack
If stacks of different architectures are named differently, the buildpack
should declare its compatibility with all stacks in the `buildpack.toml`:
```
[[stacks]]
  id = "io.buildpacks.stacks.bionic"

[[stacks]]
  id = "io.buildpacks.stacks.jammy"
```

//TODO: verify naming
And when other architectures are supported:
```
[[stacks]]
  id = "io.buildpacks.stacks.bionic.amd64"

[[stacks]]
  id = "io.buildpacks.stacks.bionic.arm64"

[[stacks]]
  id = "io.buildpacks.stacks.jammy.amd64"

[[stacks]]
  id = "io.buildpacks.stacks.jammy.arm64"
```

If a buildpack supports jammy and bionic stacks, and ARM64/AMD64, the build
workflow should retrieve all versions of the dependency. For Golang v 1.17.10
for example (no compilation needed):

Workflow triggered by get-new-versions for AMD64 and retrieves go.17.10.linux-amd64
Workflow triggered separately by get-new-versions for ARM64 and retrieves go.17.10.linux-arm64

Each dependency will have the associated stack listed in it's metadata.

Dependencies in the `buildpack.toml` of implementation buildpacks contain a
`stacks` field, so the two dependnecies will show up as:
```
  [[metadata.dependencies]]
    cpe = "cpe:2.3:a:golang:go:1.17.10:*:*:*:*:*:*:*"
    id = "go"
    licenses = ["BSD-3-Clause"]
    name = "Go"
    purl = "some-purl"
    sha256 = "some-sha"
    source = "https://go.dev/dl/go1.17.10.linux-amd64.tar.gz"
    source_sha256 = "some-sha"
    stacks = ["io.buildpacks.stacks.bionic", "io.paketo.stacks.tiny"]
    uri = "https://go.dev/dl/go1.17.10.linux-amd64.tar.gz"
    version = "1.17.10"

  [[metadata.dependencies]]
    cpe = "cpe:2.3:a:golang:go:1.17.10:*:*:*:*:*:*:*"
    id = "go"
    licenses = ["BSD-3-Clause"]
    name = "Go"
    purl = "some-purl"
    sha256 = "some-sha"
    source = "https://go.dev/dl/go1.17.10.darwin-arm64.tar.gz"
    source_sha256 = "some-sha"
    stacks = ["io.buildpacks.stacks.bionic", "io.paketo.stacks.tiny"]
    uri = "https://go.dev/dl/go1.17.10.darwin-arm64.tar.gz"
    version = "1.17.10"
```

The buildpack should know which to use


## If the dependency requires compilation:
TODO: Jammy vs Bionic vs UBI? what happens in code now?
TODO: Explore ARM64 VM
