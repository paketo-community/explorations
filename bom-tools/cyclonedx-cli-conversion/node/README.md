## Conversion

1. Generate CycloneDX XML (syft/node/from-image/syft-sbom.cdx.xml)
2. Use [CycloneDX CLI conversion
   tool](https://github.com/CycloneDX/cyclonedx-cli)
3. Get SPDX RDF format as an output

### Notes Via the CycloneDX CLI tool, there was no way to get SPDX
output in XML or JSON from what we could tell.

### Findings

Fields that were retained:
* purl
* name
* version
* licenses

Fields that are not retaiend:
* component type (probably not important?)

It seems most fields were retained, but perhaps because the XML didn't have a
huge amount of metadata available (compared to the Syft enriched JSON files).
