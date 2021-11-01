# Running Syft on a Maven app

## CycloneDX XML
Time:
real    0m6.261s
user    0m1.105s
sys     0m1.574s

* no CPEs, SHAs, layer locations

## Enriched JSON format (Not CycloneDX)

Gets full information on all indirectly installed javen archive dependencies and Debian packages:
* name
* version
* url (source?)
* source version (for some, not all)
* location within layer
* licenses (for some, not all)
* CPEs

Notes:
* Only includes our indirectly installed java archive dependencies
* There are no entries for the Java/Maven/etc dependencies our buildpacks directly install
