package decoupledbuildpack

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/postal"
)

type DependencyManager interface {
	Deliver(dependency postal.Dependency, cnbPath, layerPath, platformPath string) error
}

type DependencyMetadata struct {
	Versions []struct {
		Cpes            []string `toml:"cpes"`
		Name            string   `toml:"name"`
		Purl            string   `toml:"purl"`
		Checksum        string   `toml:"checksum"`
		Arch            string   `toml:"arch"`
		Os              string   `toml:"os"`
		Distro          string   `toml:"distro"`
		URI             string   `toml:"uri"`
		Version         string   `toml:"version"`
		StripComponents int      `toml:"strip-components"`
	} `toml:"versions"`
}

func Build(dependencyManager DependencyManager) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		layer, err := context.Layers.Get("go")
		if err != nil {
			return packit.BuildResult{}, err
		}

		var depsPath string
		var ok bool
		depsPath, ok = os.LookupEnv("BP_DEPENDENCY_METADATA")
		if !ok {
			depsPath = "/platform/deps/metadata"
		}

		deps, err := searchPlatformForDeps(depsPath, "io.paketo.go")
		if err != nil {
			return packit.BuildResult{}, err
		}

		dependency := postal.Dependency{
			Checksum:        deps.Versions[0].Checksum,
			URI:             deps.Versions[0].URI,
			StripComponents: deps.Versions[0].StripComponents,
		}

		layer, err = layer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		layer.Launch, layer.Build, layer.Cache = true, true, true

		err = dependencyManager.Deliver(dependency, context.CNBPath, layer.Path, context.Platform.Path)
		if err != nil {
			return packit.BuildResult{}, err
		}

		return packit.BuildResult{
			Layers: []packit.Layer{layer},
		}, nil
	}
}

func searchPlatformForDeps(depsPath string, id string) (DependencyMetadata, error) {
	var deps DependencyMetadata

	f, err := os.Open(fmt.Sprintf("%s/%s.toml", depsPath, strings.ReplaceAll(id, ".", "/")))
	if err != nil {
		return deps, err
	}
	defer f.Close()

	_, err = toml.NewDecoder(f).Decode(&deps)
	if err != nil {
		return deps, err
	}

	return deps, nil
}
