package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/occam"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/pexec"
)

type Executable interface {
	Execute(pexec.Execution) error
}

type DotnetPublishProcess struct {
	executable Executable
}

func NewDotnetPublishProcess(executable Executable) DotnetPublishProcess {
	return DotnetPublishProcess{
		executable: executable,
	}
}

func (p DotnetPublishProcess) Execute(workingDir, root, nugetCachePath, intermediateBuildCachePath, projectPath, outputPath string, flags []string) error {
	err := loadBuildCache(workingDir, projectPath, intermediateBuildCachePath)
	if err != nil {
		return fmt.Errorf("failed to load build cache: %w", err)
	}

	args := []string{
		"publish",
		filepath.Join(workingDir, projectPath),
	}

	args = append(args, "--configuration", "Release")
	// args = append(args, "--runtime", "ubuntu.18.04-x64")
	args = append(args, "--self-contained", "false")
	args = append(args, "--output", outputPath)

	args = append(args, flags...)

	err = p.executable.Execute(pexec.Execution{
		Args:   args,
		Dir:    workingDir,
		Env:    append(os.Environ(), fmt.Sprintf("PATH=%s:%s", root, os.Getenv("PATH")), fmt.Sprintf("NUGET_PACKAGES=%s", nugetCachePath)),
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to execute 'dotnet publish': %w", err)
	}

	err = recreateBuildCache(workingDir, projectPath, intermediateBuildCachePath)
	if err != nil {
		return err
	}

	return nil
}

func loadBuildCache(workingDir, projectPath, cachePath string) error {
	obj, err := fs.Exists(filepath.Join(cachePath, "obj"))
	if err != nil {
		return err
	}

	if obj {
		// RemoveAll to clear the contents of the directory, which fs.Copy won't do
		err = os.RemoveAll(filepath.Join(workingDir, projectPath, "obj"))
		if err != nil {
			return err
		}
		err = fs.Copy(filepath.Join(cachePath, "obj"), filepath.Join(workingDir, projectPath, "obj"))
		if err != nil {
			return err
		}
	}
	return nil
}

func recreateBuildCache(workingDir, projectPath, cachePath string) error {
	obj, err := fs.Exists(filepath.Join(workingDir, projectPath, "obj"))
	if err != nil {
		return fmt.Errorf("failed to locate build cache: %w", err)
	}

	if obj {
		// RemoveAll to clear the contents of the directory, which fs.Copy won't do
		err = os.RemoveAll(filepath.Join(cachePath, "obj"))
		if err != nil {
			// not tested
			return fmt.Errorf("failed to reset build cache: %w", err)
		}
		err = os.MkdirAll(filepath.Join(cachePath, "obj"), os.ModePerm)
		if err != nil {
			// not tested
			return fmt.Errorf("failed to reset build cache: %w", err)
		}
		err = fs.Copy(filepath.Join(workingDir, projectPath, "obj"), filepath.Join(cachePath, "obj"))
		if err != nil {
			return fmt.Errorf("failed to store build cache: %w", err)
		}
	}
	return nil
}
func main() {
	if os.Getenv("APP_SOURCE") == "" {
		log.Fatal("Must set $APP_SOURCE")
	}

	source, err := occam.Source(os.Getenv("APP_SOURCE"))
	if err != nil {
		log.Fatal(err)
	}

	dotnet := NewDotnetPublishProcess(pexec.NewExecutable("dotnet"))

	dotnetRoot := ""
	nugetCache, err := os.MkdirTemp("", "nuget")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(nugetCache, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	buildCache, err := os.MkdirTemp("", "build")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(buildCache, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	output1, err := os.MkdirTemp("", "output1")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(output1, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = dotnet.Execute(source, dotnetRoot, nugetCache, buildCache, "", output1, []string{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output1)

	file, err := os.Open(filepath.Join(source, "Program.cs"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	contents = bytes.Replace(contents, []byte("Hello World!"), []byte("Hello Moon!"), 1)

	err = os.WriteFile(filepath.Join(source, "Program.cs"), contents, os.ModePerm)
	file.Close()

	if err != nil {
		log.Fatal(err)
	}

	output2, err := os.MkdirTemp("", "output2")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(output2, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = dotnet.Execute(source, dotnetRoot, nugetCache, buildCache, "", output2, []string{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output2)

	buf1 := bytes.NewBuffer(nil)
	err = pexec.NewExecutable(filepath.Join(output1, "console_app")).Execute(pexec.Execution{
		Stdout: buf1,
		Stderr: buf1,
	})
	if err != nil {
		log.Fatal(err)
	}

	buf2 := bytes.NewBuffer(nil)
	err = pexec.NewExecutable(filepath.Join(output2, "console_app")).Execute(pexec.Execution{
		Stdout: buf2,
		Stderr: buf2,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Build 1 output: ", buf1.String())
	fmt.Println("Build 2 output: ", buf2.String())
}
