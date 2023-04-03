package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	dependencybuildpack "github.com/paketo-community/explorations/decoupled-dependencies/dependency-buildpack"
)

func main() {
	packit.Run(
		dependencybuildpack.Detect(),
		dependencybuildpack.Build(),
	)
}
