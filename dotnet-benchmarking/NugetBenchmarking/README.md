# NugetBenchmarking

This simple .NET console app's sole purpose is to benchmark the `dotnet
restore` phase of the build. Its project file is padded with lots of NuGet
packages (sourced from the NuGet [package gallery](https://www.nuget.org/)) and
added with the `dotnet add <package name>` command. Since `dotnet` restores
packages that appear in the project file whether or not they're used, this app
is quick to compile (basically no business logic), but has lots of
dependencies. Building this app with the buildpack helps demonstrate
the efficacy of storing the NuGet cache for rebuild speeds.

## Building
To observe how long it takes to restore packages on an initial build, run:
```bash
pack build nuget-benchmarking -b paketo-buildpacks/dotnet-core
```
Note the log output from the `dotnet` CLI indicating how long the restore took. For instance:
```
[builder]         Determining projects to restore...
[builder]         Restored /workspace/NugetBenchmarking.csproj (in 14.24 sec).
```

Rebuild the app with the same `pack build` command, and see that the package restore phase
is significantly faster.

