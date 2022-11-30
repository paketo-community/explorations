# Dotnet Core Sample App using ASPNet

Copied from [paketo-buildpacks/samples](https://github.com/paketo-buildpacks/samples/blob/cd2e9027eddeb6e96c6e3610195940424438b992/dotnet-core/aspnet)

## Building with Watchexec

```bash
pack build dotnet-watchexec --descriptor project.toml```

## Running

```bash
docker run --name dotnet-watchexec-container \
           --entrypoint launcher --interactive --tty \
           --env PORT=8080 --publish 8080:8080 \
           dotnet-watchexec \
           watchexec --verbose --restart --watch /workspace '/workspace/aspnet --urls http://0.0.0.0:${PORT:-8080}'
```

## Viewing

`curl http://localhost:8080`

## Test Server Restart
Since `watchexec` is watching all of the `/workspace` directory, it will
restart the server when an empty file is copied into `/workspace`

```bash
touch hello.txt
docker cp ./hello.txt  dotnet-watchexec-container:/workspace/hello.txt
```

## Rebuild the app locally (before copying into container)
With the .NET SDK installed on your local workstation, run
```bash
dotnet publish --configuration Release --runtime ubuntu.18.04-x64 --self-contained false --output ./build
```

This will cross-compile the binary for linux in the same way it happens in the
build container. The output will be dumped in `./build`, which can then be
copied into the container.

```bash
docker cp ./build/ dotnet-watchexec-container:/build
```
`docker cp` does not currently support copying multiple files at once.
Therefore, it is necessary to run the following command to copy the content of build
into the workspace.

```bash
docker exec -it dotnet-watchexec-container sh -c "cp "/build/*" /workspace"
```

# NOTES

* Given that the start command currently set by the dotnet-execute buildpack
  requires shell variable expansion, running watchexec with the `--shell=none`
  option fails.

# Hot reloading of an FDD app:

1. Build this sample as an FDD locally:
```bash
dotnet publish --configuration Release -p:UseAppHost=false --runtime ubuntu.18.04-x64 --self-contained false --output ./build
```

2. Create a file `build/plan.toml` that contains:

```toml
[[requires]]
  name = "watchexec"

  [requires.metadata]
    launch = true

```

3. Build the FDD into a container
```bash
pack build dotnet-fdd-watchexec --descriptor fdd_project.toml --path ./build
```

4. Start the container with the watchexec process matching the buildpack-set start command:
```bash
docker run --name dotnet-fdd-container \
           --entrypoint launcher --interactive --tty \
           --env PORT=8080 --publish 8080:8080 \
           dotnet-fdd-watchexec \
           watchexec --verbose --restart --watch /workspace 'dotnet /workspace/aspnet.dll --urls http://0.0.0.0:${PORT:-8080}'
```

5. Make a change to the source code of the FDD

6. Recompile the app
```bash
rm -rf ./build && dotnet publish --configuration Release -p:UseAppHost=false --runtime ubuntu.18.04-x64 --self-contained false --output ./build
```

7. Copy the compiled app into the workspace:
```bash
docker cp ./build dotnet-fdd-container:/build
docker exec -it dotnet-fdd-container sh -c "cp /build/* /workspace/"
```

8. See that watchexec has restarted the server, changes to source code are reflected in running app.

## Watchexec + .NET Buildpack (SCD)

### Build the app locally (before `pack build`)
The .NET buildpack does not support building a self-contained deployment from
source. Thus, we must build the app on our local system and `pack build` using
the built artifact.

With the .NET SDK installed on your local workstation, run
```bash
dotnet publish --configuration Release --runtime ubuntu.18.04-x64 --self-contained true --output ./build
```

### Prepare SCD for building

Modify the project.toml to match the following:
```
[[build.buildpacks]]
uri = "gcr.io/paketo-buildpacks/watchexec"

[[build.buildpacks]]
uri = "gcr.io/paketo-buildpacks/dotnet-core"

[[build.buildpacks]]
uri = "gcr.io/paketo-community/build-plan"

```

Copy `plan.toml` into the `build` directory:
```bash
cp plan.toml ./build/plan.toml```

## Build with Watchexec

From the `build` directory (output from the previous command), run:
```bash
pack build dotnet-watchexec-scd --descriptor ../project.toml
```

### Running

```bash
docker run --name dotnet-watchexec-container-scd \
           --entrypoint launcher --interactive --tty \
           --env PORT=8080 --publish 8080:8080 \
           dotnet-watchexec-scd \
           watchexec --verbose --restart --watch /workspace '/workspace/aspnet --urls http://0.0.0.0:${PORT:-8080}'
```

### Viewing

`curl http://localhost:8080`


### Rebuild the app locally (before copying into container)
If you'd like to modify the app and rebuild it, first remove the output
directory (`build`), and re-run `dotnet publish`:

```bash
dotnet publish --configuration Release --runtime ubuntu.18.04-x64 --self-contained true --output ./build
```

Alternatively, you may specify a separate build directory (e.g. build2) in the
above command. This will, however, result in a second `build` directory nested
within your output directory.

```bash
docker cp ./build/ dotnet-watchexec-container-scd:/build
```

`docker cp` does not currently support copying multiple files at once.
Therefore, it is necessary to run the following command to copy the content of build
into the workspace.

```bash
docker exec -it dotnet-watchexec-container-scd sh -c "cp "/build/*" /workspace"
```
