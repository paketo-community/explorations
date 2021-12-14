# Investigation into SBOM generation for .NET Core apps

Per [RFC 0038](https://github.com/paketo-buildpacks/rfcs/blob/932d89a970bea06c002496dbd6c7e3812ccab788/text/0038-cdx-syft-sbom.md)
we aim to support CycloneDX and Syft SBOMs in Paketo Buildpacks. We recently
began implementing this in the Node.js buildpacks, following [this set of
issues](https://github.com/paketo-buildpacks/nodejs/issues/496, following [this
set of issues](https://github.com/paketo-buildpacks/nodejs/issues/496))

For the .NET Core language family we will have two flavours of SBOMs:

## Dependency SBOMs
* dotnet-core-runtime
* dotnet-core-sdk
* dotnet-core-aspnet

The SBOM for these buildpacks can be created similarly to how they are for
other dependency-providing buildpacks like
[Yarn](https://github.com/paketo-buildpacks/yarn/pull/182).
I think that these should be pretty straightforward, using the new SBOM format
from packit v2.

## App Dependency SBOM
* dotnet-pubish

The SBOM for this buildpack will be much more complicated. This type of
buildpack is responsible for installing application dependencies.

In the Node.js language family, the NPM install and Yarn install SBOMs could be
generated via [Syft](https://github.com/anchore/syft), since the tool knew how
to parse node modules.

The .NET Publish buildpack is repsonsible for running the `dotnet publish`
command for an app, and as a result outputs a `.deps.json` file (that contains
project dependencies) as well as copies dependencies from the NuGet cache into
the output folder.

The issue is that the Syft tool does not know how to deal with NuGet
dependencies, so we may have to explore other avenues to generate the SBOM.

Options include:

- There is an [open issue](https://github.com/anchore/syft/issues/373) to support cataloguing NuGet packages.
  - We can wait to implement a .NET Core SBOM until this is supported.
  - We could collaborate and contribute to Syft to get this feature into the tool.
  - Benefits of one of these approaches would be that we are able to use the
    same Syft tool as the other buildpacks, and can leverage schema-specific
    SBOM logic without maintaining it ourselves.
  - Cons of this approach are that we don't know how long it might take to see
    this feature become a part of Syft, or the full complexity associated with
    implementing support for NuGet packages.

- Write and maintain our own code to catalogue all NuGet packages
  - We could create Syft catalog types, and the leverage Syft as a library to get the
    data in various formats. The `obj/project/assets.json` file created from
    `dotnet publish` could be a good starting place to look into.

- Use https://github.com/CycloneDX/cyclonedx-dotnet to generate a CycloneDX SBOM
  - Can we unmarshal the CycloneDX JSON into Syft struct to get ourselves an SBOM?
  - Note: I haven't been able to successfully build a substantial BOM with this tool yet.

Questions that might guide our decision (from discussions with @fg-j)
- Who are the main stakeholders that are interested in a .NET Core SBOM?
- Are there other offerings in the .NET Core space that provide an SBOM? Or
  would the buildpacks-provided options be the first in the area?
