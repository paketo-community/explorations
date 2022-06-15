package berry

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface Summer --output fakes/summer.go
type Summer interface {
	Sum(paths ...string) (string, error)
}

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) error
}

type YarnInstallProcess struct {
	executable []Executable
	summer     Summer
	logger     scribe.Logger
}

func NewYarnInstallProcess(summer Summer, logger scribe.Logger, executable ...Executable) YarnInstallProcess {
	return YarnInstallProcess{
		executable: executable,
		summer:     summer,
		logger:     logger,
	}
}

func (ip YarnInstallProcess) ShouldRun(workingDir string, metadata map[string]interface{}) (run bool, sha string, err error) {

	ip.logger.Subprocess("Process inputs:")

	env := os.Environ()
	env = append(env, fmt.Sprintf("COREPACK_HOME=%s", workingDir))

	cpBuffer := bytes.NewBuffer(nil)

	_, err = os.Stat(filepath.Join(workingDir, "corepack.tgz"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
		panic(err)
	}

	err = ip.executable[1].Execute(pexec.Execution{
		Args:   []string{"hydrate", "./corepack.tgz"},
		Stdout: cpBuffer,
		Stderr: cpBuffer,
		Dir:    workingDir,
		Env:    env,
	})
	if err != nil {
		return true, "", fmt.Errorf("failed to execute corepack hydrate output:\n%s\nerror: %s", cpBuffer.String(), err)
	}

	err = ip.executable[1].Execute(pexec.Execution{
		Args:   []string{"enable"},
		Stdout: cpBuffer,
		Stderr: cpBuffer,
		Dir:    workingDir,
	})
	if err != nil {
		return true, "", fmt.Errorf("failed to execute corepack enable output:\n%s\nerror: %s", cpBuffer.String(), err)
	}

	_, err = os.Stat(filepath.Join(workingDir, "yarn.lock"))
	if os.IsNotExist(err) {
		ip.logger.Action("yarn.lock -> Not found")
		ip.logger.Break()
		return true, "", nil
	} else if err != nil {
		return true, "", fmt.Errorf("unable to read yarn.lock file: %w", err)
	}

	ip.logger.Action("yarn.lock -> Found")
	ip.logger.Break()

	buffer := bytes.NewBuffer(nil)

	err = ip.executable[0].Execute(pexec.Execution{
		Args:   []string{"info", "--all"},
		Stdout: buffer,
		Stderr: buffer,
		Dir:    workingDir,
	})
	if err != nil {
		return true, "", fmt.Errorf("failed to execute yarn info output:\n%s\nerror: %s", buffer.String(), err)
	}

	nodeEnv := os.Getenv("NODE_ENV")
	buffer.WriteString(nodeEnv)

	file, err := os.CreateTemp("", "config-file")
	if err != nil {
		return true, "", fmt.Errorf("failed to create temp file for %s: %w", file.Name(), err)
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return true, "", fmt.Errorf("failed to write temp file for %s: %w", file.Name(), err)
	}

	sum, err := ip.summer.Sum(filepath.Join(workingDir, "yarn.lock"), file.Name())
	if err != nil {
		return true, "", fmt.Errorf("unable to sum config files: %w", err)
	}

	prevSHA, ok := metadata["cache_sha"].(string)
	if (ok && sum != prevSHA) || !ok {
		_, err = os.Stat(filepath.Join(workingDir, ".pnp.cjs"))
		if !os.IsNotExist(err) {
			ip.logger.Action(".pnp.cjs found, will not run `yarn install`")
			ip.logger.Break()
			return false, "", nil
		}
		return true, sum, nil
	}

	return false, "", nil
}

func (ip YarnInstallProcess) SetupModules(workingDir, currentModulesLayerPath, nextModulesLayerPath string) (string, error) {
	if currentModulesLayerPath != "" {
		err := fs.Copy(filepath.Join(currentModulesLayerPath, "node_modules"), filepath.Join(nextModulesLayerPath, "node_modules"))
		if err != nil {
			return "", fmt.Errorf("failed to copy node_modules directory: %w", err)
		}
	} else {
		err := os.MkdirAll(filepath.Join(workingDir, "node_modules"), os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create node_modules directory: %w", err)
		}

		err = fs.Move(filepath.Join(workingDir, "node_modules"), filepath.Join(nextModulesLayerPath, "node_modules"))
		if err != nil {
			return "", fmt.Errorf("failed to move node_modules directory to layer: %w", err)
		}

		err = os.Symlink(filepath.Join(nextModulesLayerPath, "node_modules"), filepath.Join(workingDir, "node_modules"))
		if err != nil {
			return "", fmt.Errorf("failed to symlink node_modules into working directory: %w", err)
		}
	}

	return nextModulesLayerPath, nil
}

// The build process here relies on yarn install ... --frozen-lockfile note that
// even if we provide a node_modules directory we must run a 'yarn install' as
// this is the ONLY way to rebuild native extensions.
func (ip YarnInstallProcess) Execute(workingDir, modulesLayerPath string, launch bool) error {
	environment := os.Environ()
	environment = append(environment, fmt.Sprintf("PATH=%s%c%s", os.Getenv("PATH"), os.PathListSeparator, filepath.Join("node_modules", ".bin")))

	buffer := bytes.NewBuffer(nil)

	installArgs := []string{"install"}

	// if !launch {
	// 	installArgs = append(installArgs, "--production", "false")
	// }

	// installArgs = append(installArgs, "--modules-folder", filepath.Join(modulesLayerPath, "node_modules"))
	ip.logger.Subprocess("Running yarn %s", strings.Join(installArgs, " "))

	buffer = bytes.NewBuffer(nil)
	err := ip.executable[0].Execute(pexec.Execution{
		Args:   installArgs,
		Env:    environment,
		Stdout: buffer,
		Stderr: buffer,
		Dir:    workingDir,
	})

	ip.logger.Action("%s", buffer)
	if err != nil {
		return fmt.Errorf("failed to execute yarn install: %w", err)
	}

	return nil
}
