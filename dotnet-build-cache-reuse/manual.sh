#!/usr/bin/env bash

set -eu
set -o pipefail
function main() {
  local working_dir

  while [ "${#}" != 0 ]; do
    case "${1}" in
      --dir)
        working_dir="${2}"
        shift 2
        ;;

      "")
        shift
        ;;

      *)
        echo "unknown argument \"${1}\""
        exit 1
    esac
  done

  pushd "${working_dir}" > /dev/null

  printf "Setting up nuget cache"
  rm -rf nuget-cache
  mkdir nuget-cache
  export NUGET_PACKAGES=$PWD/nuget-cache
  echo $NUGET_PACKAGES

  printf "Creating a .NET console app in %s\n" "${working_dir}"
  dotnet new console -o console_app --force

  printf "Running dotnet publish on console app\n"
  dotnet publish console_app/ --configuration Release --runtime ubuntu.18.04-x64 --self-contained false --output first-publish

  printf "Setting up intermediate build cache\n"
  rm -rf build-cache
  mkdir build-cache
  printf "Copying console_app/obj to build cache\n"
  cp -r console_app/obj build-cache

  printf "Recreating console_app\n"
  dotnet new console -o console_app --force

  printf "Modifying console_app/Program.cs console statement\n"
  sed -i 's/Hello/Goodbye/g' console_app/Program.cs

  printf "Copying in cached obj directory\n"
  cp -r build-cache/obj console_app/obj

  printf "Running dotnet publish on MODIFIED console app\n"
  dotnet publish console_app/ --configuration Release --runtime ubuntu.18.04-x64 --self-contained false --output second-publish

  printf "Testing the first app, expect to see Hello World!"
  ./first-publish/console_app

  printf "Testing the second app, expect to see Goodbye World!"
  ./second-publish/console_app
  popd > /dev/null
}

main "${@:-}"
