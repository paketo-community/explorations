# Go Sample App on Tiny Stack

This sample app shows how a simple Go project that uses Go modules
can work with the watchexec buildpack and Tilt. Instead of relying on code
changes to the Go buildpack to achieve hot reloading behaviour, we use the Build
Plan buildpack to ensure that `watchexec` is present in the runtime container
and use a Procfile to set the correct start command.

## Setting up
### Prerequisites
1. `tilt` + a connected k8s cluster
2. `pack` CLI

### Stand up the deployment
`./start.sh`

### See the build fail

The built app container fails to start with:
```bash
/layers/paketo-buildpacks_watchexec/watchexec/bin/watchexec: error while loading shared libraries: libgcc_s.so.1: cannot open shared object file: No such file or directory
```

The error indicates that `watchexec` relies on the shared library `libgcc1`,
which **is not installed in the tiny run image**. (See this [list of installed
packages in the tiny run
image](https://github.com/paketo-buildpacks/tiny-stack-release/blob/b1bf4317ffb78805384351d1a2a27f32af0f1328/run-receipt).)
Therefore, it seems that `watchexec` cannot be used on the Tiny stack without
expanding the set of packages installed on Tiny. It is unlikely that a new
package will be added to the Tiny stack for the sake of inner loop development,
given that the stack is designed for production use-cases where a minimal stack
is of paramount importance.

## Caveats and Notes
- Note that an unreleased version of the build-plan buildpack is downloaded
  here (in `.buildpacks`), because the latest release of the build plan
  buildpack does not support the Tiny stack.
