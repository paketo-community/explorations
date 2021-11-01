# Running Syft on a Go module app

## CycloneDX XML
Time:
real    0m2.659s
user    0m0.741s
sys     0m0.969s

* no CPEs, SHAs, layer locations

## Enriched JSON format (Not CycloneDX)
Gets full information on all OS level debian packages:
* name
* version
* url (source?)
* source version (for some, not all)
* location within layer
* licenses (for some, not all)
* CPEs

Notes:

* This BoM does not include any of our indirectly installed dependencies from
  go-modules since they are not present in the final app image.

* This BoM does not include any information about the Go version since it is
  also not present in the final app image.
