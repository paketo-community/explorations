# Kerberos Support

The MIT Kerberos distribution can be found at https://kerberos.org/dist/.

## Which libraries, exactly, are prerequisites for installing Kerberos on Linux?

It has already been validated that what is considered full Kerberos support can
be achieved by building applcation images using the Full stack. That stack
installs a package called
[`krb5-user`](https://packages.ubuntu.com/jammy/krb5-user) for this support.

## Where can the required libraries be obtained?

The source code for all versions of MIT Kerberos can be found at
https://kerberos.org/dist/krb5/.

## Do these libraries have to be installed by a root user in order to work properly? (i.e. Do they need to be installed in the stack, or can a buildpack install them?)

The library portions we want are only those that should be runnable by any
user, specifically not those executables generated in the output that appear in
the `/sbin` directory. For this reason, we should be able to drop any parts of
the compilation output that are not addressed specifically as supported layer
paths in the [Buildpack
specification](https://github.com/buildpacks/spec/blob/main/buildpack.md#environment)
(`/bin`, `/lib`, `/include`, `/pkgconfig`). This makes installation using a
buildpack a viable option.

## Building the library

1. Run the following:

```
docker build -t krb5:latest .
docker run -it \
  -v "${PWD}/output:/output" \
  -v "${PWD}/build.sh:/build.sh" \
  krb5:latest \
    /build.sh --output /output
```

2. Find the built library in `./output`.
