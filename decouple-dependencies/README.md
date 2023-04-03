# Decouple Dependencies
This is a Proof-of-concept to prove out some of the functionality that is laid
out in [this proposal](https://docs.google.com/document/d/1g5rRW-oE_v8Gdvz-CiCOK9z2rxg6L5XniKI25Zq2j6M/)
for decoupling buildpacks and the dependencies that they serve.

The goal of this POC is to show that the purpose API would work for the various
integration points that are given in the proposal.

## External Volume Mount
1. Build the decoupled buildpack:
```
./decoupled-buildpack/scripts/package.sh --version 1.2.3
```
1. Build a container and volume mount the dependecy metadata at the correct location:
```
pack build test -b decoupled-buildpack/build/buildpackage.cnb --volume $PWD/dependency:/platform/deps/metadata
```
1. Verify that the dependency was installed:
```
docker run -it --entrypoint=launcher test "go version"
```
Expected valued should look like:
```
go version go1.19.6 linux/amd64
```

## Via a Buildpack
1. Build the decoupled buildpack:
```
./decoupled-buildpack/scripts/package.sh --version 1.2.3
```
1. Build the dependency buildpack:
```
./dependency-buildpack/scripts/package.sh --version 1.2.3
```
1. Build a container by placing the dependency metadata loading buildpack first in the build order:
```
pack build test -b dependency-buildpack/build/buildpackage.cnb -b decoupled-buildpack/build/buildpackage.cnb
```
1. Verify that the dependency was installed:
```
docker run -it --entrypoint=launcher test "go version"
```
Expected valued should look like:
```
go version go1.19.6 linux/amd64
```

