package main

import (
	"log"
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/sbom"

	yarninstall "github.com/paketo-buildpacks/yarn-install"
)

type SBOMGenerator struct{}

func (s SBOMGenerator) Generate(path string) (sbom.SBOM, error) {
	return sbom.Generate(path)
}

func main() {
	packageJSONParser := yarninstall.NewPackageJSONParser()
	projectPathParser := yarninstall.NewProjectPathParser()
	home, err := os.UserHomeDir()
	if err != nil {
		// not tested
		log.Fatal(err)
	}

	packit.Run(
		yarninstall.Detect(
			projectPathParser,
			packageJSONParser,
		),
		yarninstall.Build(projectPathParser, home),
	)
}
