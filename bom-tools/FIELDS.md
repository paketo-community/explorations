# Fields of interest in a BoM

## Directly Installed Dependencies
#### EX: Node-Engine or Go
Dependencies that directly provide runtimes and/or are tools used for compilation
* name
* sha256
* uri
* version
* cpe
* licenses
* source-uri
* source-sha


## Indirectly Installed Dependencies
#### EX: Go modules, node modules
Dependencies that are installed during the build process/vendored.
* name
* version

| Field    | CycloneDX|   SPDX        |
| ----------- | ----------- |  ----------- |
| name    | component name     |         package name |
| version | component version |         package version|
| uri | external references or properties |package information download location |
| sha256  | component hash             |package information checksum |
|    source uri | external references or properties  |package source information|
| source sha256  | external references of properties| package source information |
| cpe    | component CPE     |package external reference |
| licenses  | component license | package declared license|
