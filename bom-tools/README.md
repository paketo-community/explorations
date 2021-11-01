# Investigation into available BOM tools and formats

For README section, we used the top level tool against a sample ap to generate a Bill of Materials.
Result Bills of Materials are located in the matching directories for each tool

### [FIELDS.md](FIELDS.md)
Contains the BoM fields of interest and how we believe we can represent them in
CycloneDX and SPDX.


### Tools Investigations

#### Pack
* [Node.JS example](pack/node) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/no-package-manager)
* [NPM example](pack/npm) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/npm)
* [Yarn example](pack/yarn) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/yarn)
* [Go Mod example](pack/go-mod) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/go/mod)
* [Maven example](pack/maven) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/maven)


#### [Syft](https://github.com/anchore/syft)
Command run for CycloneDX XML (from source): `syft packages <path-to-source> -o cyclonedx`

Command run for CycloneDX XML (from image): `syft packages <image-name> -o cyclonedx`

Includes examples in both Syft enriched JSON format (not CycloneDX) and CycloneDX (XML)

* [Node.JS example from image](syft/node/from-image) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/no-package-manager)
* [NPM example from image](syft/npm/from-image) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/npm)
* [NPM example from image](syft/yarn/from-image) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/yarn)
* [Go Mod example from image](syft/go-mod/from-image) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/go/mod)
* [Maven example from image](syft/maven/from-image) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/maven)
--------------------------------
* [Node.JS example from source](syft/node/from-source) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/no-package-manager)
* [NPM example from source](syft/npm/from-source) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/npm)
* [NPM example from source](syft/yarn/from-source) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/yarn)
* [Go Mod example from source](syft/go-mod/from-source) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/go/mod)
* [Maven example from source](syft/maven/from-source) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/maven)


#### [Conversion](https://github.com/CycloneDX/cyclonedx-cli)
* [Node.JS example from image](cyclonedx-cli-conversion/node)

#### Tern
Command run for true SPDX JSON format: `./docker_run.sh ternd "report -f spdxjson -i built-app-image:latest"`

Includes examples in both Tern enriched JSON format (not SPDX) and true SPDX (JSON)

* [Node.JS example](tern/node) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/no-package-manager)
* [NPM example](tern/npm) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/npm)
* [NPM example](tern/yarn) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/nodejs/yarn)
* [Go Mod example](tern/go-mod) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/go/mod)
* [Maven example](tern/maven) / [sample app used](https://github.com/paketo-buildpacks/samples/tree/main/maven)


#### Conclusions

**Time Averages for Scanning**
* Syft on a pre-built image: 4.8062s
* Syft on application source code: 0.7624s
* Tern on a pre-built image: 136.4412s

**Tern**
* We generated both a Tern-specific SPDX JSON file, as well as a "true" SPDX JSON file.
* It appears that while neither have *all* of the metadata we are looking for,
  the "true" SPDX format has the least.
  * It does not have CPEs, SHAs, and layer paths

**Syft**
* The JSON output of Syft is NOT CycloneDX format, but rather it's a superset
  of all metadata that can be retrieved.
* The real CycloneDX format from Syft (XML) is missing some information.
  * It does not have CPEs, SHAs, and layer paths

* [This issue](https://github.com/anchore/syft/issues/325) on Syft looks like a
  request for what we might want

* The enriched formats have all of the information we need, but they don't seem
  to align with CycloneDX or SPDX as well as we once thought

Syft enriched JSON seems to fit our use case
* Gets full information on all OS level and indirectly installed packages
  * CPEs, licenses, urls, shas, layer location, name, version
  * These are all fields we feel strongly about being first-class citizens
  * The SPDX output from Tern did not surface this metadata fully/or as clearly

* For all of our language modules (packages installed by go-mod and npm) can
  easily be retrieved with fully-fledged metadata
* It might be nice to integrate with this tooling rather than build out our own custom logic

* No information about the the actual dependencies we directly install (node,
  go, etc), but this is information that we can easily provide

* Does a better job of conveying the information we think is important
* The tooling ecosystem feels better fleshed out for the use cases we have (language module metadata collection)

**Format Concerns**

The fact that these scanning tools all seem to have enriched BOM outputs that
are outside of either SPDX or CycloneDX and then there translations to these
formats are sparse is with interesting and cause for pause.

* Why do these tools not just try and do everything in the existing BOM formats
  that they support?
* Why is the translation between their format and the "official" format so
  sparse?
* Is it sparse because they don't care about those formats or because they
  cannot get more specfic?

#### To Do
* Conversion tooling
* Offline environments?
* Check CycloneDX Scanning Tools
