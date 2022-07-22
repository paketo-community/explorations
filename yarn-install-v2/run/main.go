package main

import (
	"log"
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/sbom"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	yarninstall "github.com/paketo-community/yarn-install-v2"
	"github.com/paketo-community/yarn-install-v2/berry"
	"github.com/paketo-community/yarn-install-v2/classic"
	"github.com/paketo-community/yarn-install-v2/common"
)

type SBOMGenerator struct{}

func (s SBOMGenerator) Generate(path string) (sbom.SBOM, error) {
	return sbom.Generate(path)
}

func main() {
	packageJSONParser := common.NewPackageJSONParser()
	emitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	projectPathParser := common.NewProjectPathParser()
	sbomGenerator := SBOMGenerator{}
	symlinker := common.NewSymlinker()
	home, err := os.UserHomeDir()
	if err != nil {
		// not tested
		log.Fatal(err)
	}
	classicInstallProcess := classic.NewYarnInstallProcess(fs.NewChecksumCalculator(), scribe.NewLogger(os.Stdout), pexec.NewExecutable("yarn"), pexec.NewExecutable("corepack"))
	berryInstallProcess := berry.NewYarnInstallProcess(fs.NewChecksumCalculator(), scribe.NewLogger(os.Stdout), pexec.NewExecutable("yarn"), pexec.NewExecutable("corepack"))

	packit.Run(
		yarninstall.Detect(
			projectPathParser,
			packageJSONParser,
		),
		yarninstall.Build(
			projectPathParser,
			home,
			symlinker,
			classicInstallProcess,
			berryInstallProcess,
			sbomGenerator,
			chronos.DefaultClock,
			emitter),
	)
}
