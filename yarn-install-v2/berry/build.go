package berry

import (
	"os"
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/sbom"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		projectPath, err := pathParser.Get(context.WorkingDir)
		if err != nil {
			return packit.BuildResult{}, err
		}

		launch, build := entryResolver.MergeLayerTypes(PlanDependencyNodeModules, context.Plan.Entries)

		var layers []packit.Layer
		var currentModLayer string
		if build {
			layer, err := context.Layers.Get("build-modules")
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Process("Resolving installation process")

			layer.SharedEnv.Override("COREPACK_HOME", context.WorkingDir)
			// layer.SharedEnv.Override("COREPACK_ENABLE_NETWORK", "0")
			logger.EnvironmentVariables(layer)
			run, sha, err := installProcess.ShouldRun(projectPath, layer.Metadata)
			if err != nil {
				return packit.BuildResult{}, err
			}

			if run {
				logger.Subprocess("Selected default build process: 'yarn install'")
				logger.Break()
				logger.Process("Executing build environment install process")

				layer, err = layer.Reset()
				if err != nil {
					return packit.BuildResult{}, err
				}

				currentModLayer, err = installProcess.SetupModules(context.WorkingDir, currentModLayer, layer.Path)
				if err != nil {
					return packit.BuildResult{}, err
				}

				duration, err := clock.Measure(func() error {
					return installProcess.Execute(projectPath, layer.Path, false)
				})
				if err != nil {
					return packit.BuildResult{}, err
				}

				logger.Action("Completed in %s", duration.Round(time.Millisecond))
				logger.Break()

				layer.Metadata = map[string]interface{}{
					"built_at":  clock.Now().Format(time.RFC3339Nano),
					"cache_sha": sha,
				}

				path := filepath.Join(layer.Path, "node_modules", ".bin")
				layer.BuildEnv.Append("PATH", path, string(os.PathListSeparator))
				layer.BuildEnv.Override("NODE_ENV", "development")
				layer.BuildEnv.Override("NODE_ENV", "development")
				layer.BuildEnv.Override("NODE_ENV", "development")

				logger.EnvironmentVariables(layer)

				logger.GeneratingSBOM(layer.Path)
				var sbomContent sbom.SBOM
				duration, err = clock.Measure(func() error {
					sbomContent, err = sbomGenerator.Generate(context.WorkingDir)
					return err
				})
				if err != nil {
					return packit.BuildResult{}, err
				}
				logger.Action("Completed in %s", duration.Round(time.Millisecond))
				logger.Break()

				logger.FormattingSBOM(context.BuildpackInfo.SBOMFormats...)
				layer.SBOM, err = sbomContent.InFormats(context.BuildpackInfo.SBOMFormats...)
				if err != nil {
					return packit.BuildResult{}, err
				}
			} else {
				logger.Process("Reusing cached layer %s", layer.Path)

				err := os.RemoveAll(filepath.Join(projectPath, "node_modules"))
				if err != nil {
					return packit.BuildResult{}, err
				}

				err = symlinker.Link(filepath.Join(layer.Path, "node_modules"), filepath.Join(projectPath, "node_modules"))
				if err != nil {
					return packit.BuildResult{}, err
				}
			}

			layer.Build = true
			layer.Cache = true

			layers = append(layers, layer)
		}

		if launch {
			layer, err := context.Layers.Get("launch-modules")
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Process("Resolving installation process")

			layer.SharedEnv.Override("COREPACK_HOME", context.WorkingDir)
			// layer.SharedEnv.Override("COREPACK_ENABLE_NETWORK", "0")
			logger.EnvironmentVariables(layer)
			run, sha, err := installProcess.ShouldRun(projectPath, layer.Metadata)
			if err != nil {
				return packit.BuildResult{}, err
			}

			if run {
				logger.Subprocess("Selected default build process: 'yarn install'")
				logger.Break()
				logger.Process("Executing launch environment install process")

				layer, err = layer.Reset()
				if err != nil {
					return packit.BuildResult{}, err
				}

				_, err = installProcess.SetupModules(context.WorkingDir, currentModLayer, layer.Path)
				if err != nil {
					return packit.BuildResult{}, err
				}

				duration, err := clock.Measure(func() error {
					return installProcess.Execute(projectPath, layer.Path, true)
				})
				if err != nil {
					return packit.BuildResult{}, err
				}

				logger.Action("Completed in %s", duration.Round(time.Millisecond))
				logger.Break()

				layer.Metadata = map[string]interface{}{
					"built_at":  clock.Now().Format(time.RFC3339Nano),
					"cache_sha": sha,
				}

				path := filepath.Join(layer.Path, "node_modules", ".bin")
				layer.LaunchEnv.Append("PATH", path, string(os.PathListSeparator))

				logger.EnvironmentVariables(layer)

				logger.GeneratingSBOM(layer.Path)
				var sbomContent sbom.SBOM
				duration, err = clock.Measure(func() error {
					sbomContent, err = sbomGenerator.Generate(context.WorkingDir)
					return err
				})
				if err != nil {
					return packit.BuildResult{}, err
				}
				logger.Action("Completed in %s", duration.Round(time.Millisecond))
				logger.Break()

				logger.FormattingSBOM(context.BuildpackInfo.SBOMFormats...)
				layer.SBOM, err = sbomContent.InFormats(context.BuildpackInfo.SBOMFormats...)
				if err != nil {
					return packit.BuildResult{}, err
				}

				layer.ExecD = []string{filepath.Join(context.CNBPath, "bin", "setup-symlinks")}
			} else {
				logger.Process("Reusing cached layer %s", layer.Path)
				if !build {
					err := os.RemoveAll(filepath.Join(projectPath, "node_modules"))
					if err != nil {
						return packit.BuildResult{}, err
					}

					err = symlinker.Link(filepath.Join(layer.Path, "node_modules"), filepath.Join(projectPath, "node_modules"))
					if err != nil {
						return packit.BuildResult{}, err
					}
				}
			}

			layer.Launch = true

			layers = append(layers, layer)
		}

		err = symlinker.Unlink(filepath.Join(homeDir, ".npmrc"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		err = symlinker.Unlink(filepath.Join(homeDir, ".yarnrc"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		return packit.BuildResult{
			Layers: layers,
		}, nil
	}
}
