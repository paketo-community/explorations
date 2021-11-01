# Running Tern on a yarn app

## True SPDX JSON format
Time:
real    2m3.314s
user    0m0.233s
sys     0m0.074s

## Non-SPDX JSON format
Gets every file within the node_modules directory and reports:
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

* Includes an entry for node, but has no version (line 257417)
* Does not include an entry for yarn, but this is expected since yarn is not on
  the final app image (only required during build-time)
