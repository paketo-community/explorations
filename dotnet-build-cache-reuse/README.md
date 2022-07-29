# Build cache issue reproduction

## Testing locally

Run
```
APP_SOURCE=<path to console app> go run main.go
```

## Testing in a docker container
Use the official .NET SDK docker container for .NET 3.1 on top of ubuntu 18.04:
1. Compile the go program for linux: `GOOS=linux go build .`
3. Run the test in a container:
  ```
docker run -it -v $(pwd):/test \
               -v $APP_SOURCE:/source \
               --env APP_SOURCE=/source \
               mcr.microsoft.com/dotnet/sdk:3.1-bionic \
               /test/dotnet-build-cache
  ```

## Testing `dotnet publish` behaviour manually (on Linux)
On Linux, it appears that the issue only occurs when orchestrated through
Golang. When run manually or through a bash script, the failure does not crop
up.

Running the following appears to have the expected effects, in which the app
serves the updated content:
```
./manual.sh --dir <path to some temporary directory"
```
