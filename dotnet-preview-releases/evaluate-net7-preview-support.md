## Evaluate support for prerelease versions of .net
### Validate:
Can a preview version of the .net sdk and runtime be supported by Paketo v3 buildpacks?

### Method:
Fork runtime and sdk buildpacks
* Add a dependency to a preview version of the buildpack
* Modify/add tests to use .net 7 preview, quick fixing issues until they pass

Fork execute buildpack
* Add an integration test to run a pre-compiled app using .net 7 preview
* Add an integration test to build and run an app using .net 7 preview
* Run unit and integration tests, quick fixing issues until they pass

Retreived latest SDK and Runtime from https://dotnet.microsoft.com/en-us/download/dotnet/7.0

### [dotnet-core-aspnet-runtime](https://github.com/KieranJeffreySmart/dotnet-core-aspnet-runtime/)
Adding a dependency involved updating the `buildpack.toml` file to include a uri and sha512 checksum for the microsoft download of .net 7 rc.2 runtime.
> **_Note_**: there is a section \[\[metadata.dependency-constraints\]\] which is present for .net core 3.1 and .net 6, however this does not need to be present for tests to pass

References to .net 6 were replaced with .net 7 rc.2, in `./integration/default_test.go` and `./integration/layer_reuse.go`. These were the only tests I found that directly referenced .net 6.

These tests failed initially due to the use of a preview version number, which led to changes being made to `./runtime_version_resolver.go` and its respective tests in `./runtime_version_resolver_test.go`

This highlighted a complexity around rolling forward on versions of the .net framework, as there is an assumption this can be done with an asterix, e.g. 7.0.*, however this will not work with prerelease versions as they have a lower precedence than 7.0.* as described [here](https://semver.org/#spec-item-9)

There is no specified way to check for versions of prerelease with wildcards. However it is possible to use the following condition `>= 7.0.0-0` as described [here](https://github.com/Masterminds/semver#readme)

Additional logic will be needed to better handle wildcards and rollforward, with specific attention to the inclusion of prereleases

> **_Note_**: there is a section `[[metadata.dependency-constraints]]` which is present for .net core 3.1 and .net 6, however this does not need to be present for tests to pass

### [dotnet-core-sdk](https://github.com/KieranJeffreySmart/dotnet-core-sdk)
Adding a dependency involved updating the `buildpack.toml` file to include a uri and sha512 checksum for the microsoft download of .net 7 rc.2 runtime.
Existing tests were duplicated and references to .net 6 were replaced with .net 7 rc.2, in `./integration/default_net7_test.go`, `./integration/layer_net7_reuse.go` and `./integrationoffline_net7_test.go`, excluding tests for using `BP_DOTNET_FRAMEWORK_VERSION`.
To test `BP_DOTNET_FRAMEWORK_VERSION` I added a test for .Net 7 which failed due to the specified version number, `7.0.100-rc.2.22477.23`, being reformatted to `7.0.*`. This was unexpected, as 7.0.100-rc should be considered a higher version than 7.0.0 but it seems the prerelease version was ignored. This led me to modifying `./build.go` to override the formatting

Setting the value of `[metadata.default-versions].dotnet-sdk` in `buildpack.toml` to `7.0.100-rc.2.22477.23`, allows the tests to run using the .net 7 framework, however this causes .net 6 tests to fail

> **_Note_**: Setting the value of `BP_DOTNET_FRAMEWORK_VERSION` or hard coding `./build.go` to use `7.0.0-0` failed to roll forward, but instead looked for an exact match.

> **_Note_**: there is a section `[[metadata.dependency-constraints]]` which is present for .net core 3.1 and .net 6, however this does not need to be present for tests to pass

### [dotnet-execute](https://github.com/KieranJeffreySmart/dotnet-execute)
To test the changes made to the dotnet-core-aspnet-runtime and dotnet-core-sdk buildpacks, I added 2 tests, the first “when building a .NET 7 FDE app that uses the ASP.NET Framework” in `./integration/fde_aspnet_test.go` and “when 'net7.0' is specified as the TargetFramework” in `./integration/source_app_test.go`

To run the tests I modified the file `./integration.json` to reference the local folders where the dotnet-core-aspnet-runtime and dotnet-core-sdk buildpacks could be found and added the line `WithVersion("1.2.3")` to `./integration/init_test.go`, where the buildpackStore is being used to retrieve the dotnet-core-aspnet-runtime and dotnet-core-sdk buildpacks.

The tests in `./integration/fde_aspnet_test.go` passed without further modification, however the tests in `./integration/source_app_test.go` failed due to the version number formatting performed while building an app from source code. This led me to adding a condition to the `ParseVersion` function in `./project_file_parser.go` to return `"7.0.0-0"` if `"net7"` was found as the target framework. To prevent reformatting when deciding the sdk version, I also added a clause to the function `getSDKVersion` found in `./detect.go` to hard return `"7.0.0-0"` if `7.0.0-0` was the detected framework returned by the ParseVersion function.

### Conclusion
The majority of effort is simple, however there are some complexities around handling preview versions, both when retrieving an exact version and when rolling forward to the latest version. There will need to be some discussion as to what the expected behavior should be and how we should provide a release that supports preview versions.
There also seems to be version management in many places, possibly doing the same thing but duplicated.
I have yet to find where the dotnet-sdk buildpack formats the detected version to be 7.0.* (the result of line 46 in file `build.go`) and fear it might be in a different repository all together
The method of retrieving the target framework from the csproj file could be seen as a little flakey and could be improved using the method described [here](https://github.com/cloudfoundry/dotnet-core-buildpack/issues/520)