# .NET Remote Debugging

Visual Studio Code can be configured to attach a remote debugging session into
a running container via `docker exec`. Presumably, this functionality could be
extended to work with Visual Studio and `kubectl exec`. These pieces of the
puzzle remain to be validated.

Contained in this exploration is a working prototype of the VSCode support for
debugging a .NET Core application built using buildpacks. This exploration
includes a simple buildpack that includes the `vsdbg` debugger in the
application container and includes extra flags and settings that enable remote
debugging of .NET Core applications in containers.

An RFC covering these buildpack features is [already
submitted](https://github.com/paketo-buildpacks/rfcs/pull/213) to the Paketo
project.

## Build the app

```
pack build debug-app \
  --path ./app \
  --buildpack ./buildpack \
  --buildpack paketo-buildpacks/dotnet-core \
  --env ASPNETCORE_ENVIRONMENT=Development \
  --env BPE_ASPNETCORE_ENVIRONMENT=Development \
  --env BP_DOTNET_PUBLISH_FLAGS="--configuration Debug"
```

## Setup VSCode

1. Add `.vscode/launch.json` to app source directory

```json
{
  "configurations": [
    {
      "name": ".NET Core Docker Attach",
      "type": "coreclr",
      "request": "attach",
      "processId": "${command:pickRemoteProcess}",
      "pipeTransport": {
        "pipeProgram": "docker",
        "pipeArgs": [ "exec", "-i", "<container id>" ],
        "debuggerPath": "/cnb/lifecycle/launcher vsdbg",
        "pipeCwd": "${workspaceRoot}",
        "quoteArgs": false
      },
      "sourceFileMap": {
        "/workspace": "${workspaceRoot}"
      }
    }
  ]
}
```

2. Install Microsoft C# Extension

## Start debugging

1. Run the app with `docker run -p 8080:8080 debug-app`
2. Open a browser window to `http://localhost:8080`
3. Update `<container id>` field in `launch.json` with actual container id
4. Add a breakpoint in `Controllers/WeatherForecastController.cs`
5. Click on the "Fetch data" link
6. Wait for the editor to load the breakpoint
7. Confirm that VSCode opens the `Controllers/WeatherForecastController.cs`
   file and stops at the line with a breakpoint
