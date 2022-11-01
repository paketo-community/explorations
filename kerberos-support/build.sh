#!/usr/bin/env bash

set -e
set -u
set -o pipefail

function main() {
  local tmpDir outputDir
  tmpDir="$(mktemp -d)"

  while [[ "${#}" != 0 ]]; do
    case "${1}" in
      --output)
        outputDir="${2}"
        shift 2;
        ;;

      "")
        # skip if the argument is empty
        shift 1
        ;;

      *)
        util::print::error "unknown argument \"${1}\""
    esac
  done

  if [[ -z "${outputDir}" ]]; then
    echo "--output is a required argument"
    exit 1
  fi

  curl -sSLo "${tmpDir}/source.tgz" "https://kerberos.org/dist/krb5/1.20/krb5-1.20.tar.gz"
  pushd "${tmpDir}" > /dev/null || true
    tar --strip-components 1 -xzvf "${tmpDir}/source.tgz"
    rm "${tmpDir}/source.tgz"

    pushd "${tmpDir}/src" > /dev/null
      ./configure --prefix="${outputDir}"
      make
      make install
    popd > /dev/null || true
  popd > /dev/null || true

}

main "${@:-}"
