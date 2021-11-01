# Running Tern on a Go module app

## True SPDX JSON format
Time:
real    1m37.153s
user    0m0.227s
sys     0m0.071s

## Non-SPDX JSON format
Gets full information on all OS level debian packages:
* name
* version
* source sha (for some, not all)
* location within layer
* licenses

Notes:
* No CPEs
* Download URLs are present but empty for every entry
* This is file by file, rather than dependency by dependency. It is more fine
  grained than the Syft output.

* Includes image labels

* Only includes our indirectly install node modules

* This BoM does not include any of our indirectly installed dependencies from
  go-modules since they are not present in the final app image.
* This BoM does not include any information about the Go version since it is
  also not present in the final app image.
