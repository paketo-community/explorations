# Running Syft on a Node.JS app with no package manager

## CycloneDX XML
Time:
real    0m4.816s
user    0m1.380s
sys     0m1.307s

* no CPEs, SHAs, layer locations

## Enriched JSON output (not Cyclonedx)
**Note: The JSON output of Syft is NOT CycloneDX format, but rather
it's a superset of all metadata that can be retrieved.**

Gets full information on all *indirectly installed node modules* that come with Node:
* name
* version
* url (source?)
* source SHA (for some, not all)
* location within layer
* licenses (for some, not all)
* CPEs

Notes:
* Only includes our indirectly install node modules
* I don't see any entry for Node.JS itself, or Node version 14.17 that is installed

# CycloneDX XML format
* Does not appear to include CPEs, location within layer, hashes
