# Running Syft on an NPM app

## CycloneDX XML
Time:
real    0m5.286s
user    0m1.747s
sys     0m1.363s

* no CPEs, SHAs, layer locations

## Enriched JSON output (not Cyclonedx)
Gets full information on all *indirectly installed node modules* that come with Node:
* name
* version
* url (source?)
* source version (for some, not all)
* location within layer
* licenses (for some, not all)
* CPEs

Notes:
* Only includes our indirectly installed node modules
* I don't see any entry for Node.JS or NPM
