package finalize

import (
	"io"

	"github.com/cloudfoundry/libbuildpack"
	"path"
)

type Stager interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/stager.go
	BuildDir() string
	DepDir() string
	DepsIdx() string
	DepsDir() string
}

type Manifest interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/manifest.go
	AllDependencyVersions(string) []string
	DefaultVersion(string) (libbuildpack.Dependency, error)
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type Command interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/command.go
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(dir string, program string, args ...string) (string, error)
}

type Finalizer struct {
	Manifest Manifest
	Stager   Stager
	Command  Command
	Log      *libbuildpack.Logger
}

func (f *Finalizer) Run() error {
	f.Log.BeginStep("Configuring xcruntime")

	f.Log.Info("Testing xcruntime install.")
	versionOutput, err := f.Command.Output(".", "mono", path.Join(f.Stager.BuildDir(), ".xc-buildpack", "xc-runtime", "xcruntime.exe"))
	if err != nil {
		return err
	}
	f.Log.Info("mono xcruntime.exe output : %s", versionOutput)

	return nil
}
