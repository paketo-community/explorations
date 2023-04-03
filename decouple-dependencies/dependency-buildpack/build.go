package dependencybuildpack

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/vacation"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		layer, err := context.Layers.Get("deps")
		if err != nil {
			return packit.BuildResult{}, err
		}

		layer, err = layer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		layer.Launch, layer.Build, layer.Cache = false, true, true

		f, err := os.Open(filepath.Join(context.CNBPath, "dependency.tgz"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		err = vacation.NewArchive(f).Decompress(layer.Path)
		if err != nil {
			return packit.BuildResult{}, err
		}

		layer.BuildEnv.Default("BP_DEPENDENCY_METADATA", layer.Path)

		return packit.BuildResult{
			Layers: []packit.Layer{layer},
		}, nil
	}
}
