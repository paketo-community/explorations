package yarninstall

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/sbom"
	"github.com/paketo-community/yarn-install-v2/berry"
	"github.com/paketo-community/yarn-install-v2/classic"
)

//go:generate faux --interface SymlinkManager --output fakes/symlink_manager.go
type SymlinkManager interface {
	Link(oldname, newname string) error
	Unlink(path string) error
}

//go:generate faux --interface InstallProcess --output fakes/install_process.go
type InstallProcess interface {
	ShouldRun(workingDir string, metadata map[string]interface{}) (run bool, sha string, err error)
	SetupModules(workingDir, currentModulesLayerPath, nextModulesLayerPath string) (string, error)
	Execute(workingDir, modulesLayerPath string, launch bool) error
}

//go:generate faux --interface EntryResolver --output fakes/entry_resolver.go
type EntryResolver interface {
	MergeLayerTypes(string, []packit.BuildpackPlanEntry) (launch, build bool)
}

//go:generate faux --interface SBOMGenerator --output fakes/sbom_generator.go
type SBOMGenerator interface {
	Generate(dir string) (sbom.SBOM, error)
}

//go:generate faux --interface ConfigurationManager --output fakes/configuration_manager.go
type ConfigurationManager interface {
	DeterminePath(typ, platformDir, entry string) (path string, err error)
}

func Build(pathParser PathParser, homeDir string) packit.BuildFunc {

	projectPath, err := projectPathParser.Get(context.WorkingDir)
	if err != nil {
		return packit.DetectResult{}, err
	}

	runBerry := true

	_, err = os.Stat(filepath.Join(projectPath, ".yarnrc.yml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			runBerry = false
		}
		return packit.DetectResult{}, err
	}

	if runBerry {
		return berry.Build()
	}

	return classic.Build()
}
