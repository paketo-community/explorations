FROM docker.io/paketobuildpacks/build-jammy-base:latest
ADD dependency/ /from-stack/deps/metadata
ADD BP_DEPENDENCY_METADATA.override /cnb/build-config/env/
