# Running Tern on an NPM app

## True SPDX JSON format
Time:
real    2m5.101s
user    0m0.234s
sys     0m0.072s

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
* Includes an entry for npm, but has no version (line 184757)
