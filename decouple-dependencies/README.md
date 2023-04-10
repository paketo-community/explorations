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
2. Build a container and volume mount the dependecy metadata at the correct location:
```
pack build test -b decoupled-buildpack/build/buildpackage.cnb --volume $PWD/dependency:/platform/deps/metadata
```
3. Verify that the dependency was installed:
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
2. Build the dependency buildpack:
```
./dependency-buildpack/scripts/package.sh --version 1.2.3
```
3. Build a container by placing the dependency metadata loading buildpack first in the build order:
```
pack build test -b dependency-buildpack/build/buildpackage.cnb -b decoupled-buildpack/build/buildpackage.cnb
```
4. Verify that the dependency was installed:
```
docker run -it --entrypoint=launcher test "go version"
```
Expected valued should look like:
```
go version go1.19.6 linux/amd64
```

## Via a Builder
NOTE: You need to be using a `pack` that supports platform API `0.11` which
currently is currently `0.30.0-pre1` and greater

1. Build extended build image:
```
cd custom-stack && docker build -f extend.Dockerfile -t extended-build . && cd ..
```
2. Build the decoupled buildpack:
```
./decoupled-buildpack/scripts/package.sh --version 1.2.3
```
3. Create the builder with embedded dependency metadata:
```
pack builder create test-builder --config builder/builder.toml
```
4. Build a container using the builder:
```
 pack build test --builder test-builder
```
5. Verify that the dependency was installed:
```
docker run -it --entrypoint=launcher test "go version"
```
Expected valued should look like:
```
go version go1.19.6 linux/amd64
```
