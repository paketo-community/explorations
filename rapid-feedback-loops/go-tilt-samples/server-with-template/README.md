# Go Sample App With Template Files

This sample app shows how a simple Go project that serves some templated files
can work with the watchexec buildpack and Tilt.

## Setting up
### Prerequisites
1. `tilt` + a connected k8s cluster
2. `pack` CLI

### Stand up the deployment
`./start.sh`

## Exploring Hot Reload

Make file changes to `main.go` or to the contents of `web/` and see Tilt
automatically recompile the source code and/or copy relevant files into the
container. See `main.go` and `web/templates/index.tmpl` for some potential
changes to demo.

## Viewing

`curl http://localhost:8080` or open the page in your browser.

## Caveats and Notes
- Note that in the Tiltfile's `pack` resource, we `sync()` all of the current
  working directory into `/workspace`, then separately `run()` a force-copy of
  the built binary into the correct layer that `watchexec` is watching. This is
  to avoid a scenario where `watchexec` triggers on intermediate file changes
  that occur when the `go` tool rebuilds the binary.
- Note also that we _must_ sync all of the working directory with `/workspace`
  rather than only `./build`, because otherwise `pack` will automatically
  re-trigger a full build on changes to files in the working directory (e.g.
  `main.go`). To mitigate this, we'd need the ability to configure the `pack`
  resource so that full rebuilds only trigger manually, or rebuilds ignore certain files/directories.
- Note also that the `local_resource` depends explicitly on all of the source
  code files in the directory. The `deps` parameter does not accept glob
  patterns like `*.go`. It would be less brittle if all source code were in a
  `./src` directory. However, the Paketo Go buildpack doesn't currently
  support building projects in a subdirectory.
    - Another way to mitigate this: if Tilt could `sync` files from _outside_
      the build context into the container. pack build with `--path` flag,
      restricting build context to `./src`. `local_resource` watches `./src`
      and drops rebuilt binaries in `./build`. Sync `./build` (now outside the
      build context) into the container.
