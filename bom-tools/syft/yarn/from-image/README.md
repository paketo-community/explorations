# Running Syft on a yarn app

## CycloneDX XML
Time:
real    0m5.009s
user    0m1.740s
sys     0m1.329s

* no CPEs, SHAs, layer locations

## Enriched JSON output (not Cyclonedx)
Gets full information on all *indirectly installed node modules*:
* name
* version
* url (source?)
* source version (for some, not all)
* location within layer
* licenses (for some, not all)
* CPEs

Notes:
* Only includes our indirectly install node modules
* I don't see any entry for yarn itself
