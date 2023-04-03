package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/postal"
	decoupledbuildpack "github.com/paketo-community/explorations/decoupled-dependencies/decoupled-buildpack"
)

func main() {

	dependencyManager := postal.NewService(cargo.NewTransport())
	packit.Run(
		decoupledbuildpack.Detect(),
		decoupledbuildpack.Build(dependencyManager),
	)
}
