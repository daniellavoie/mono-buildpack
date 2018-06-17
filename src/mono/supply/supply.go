package supply

import (
	"io"

	"github.com/cloudfoundry/libbuildpack"
	"path/filepath"
)

type Stager interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/stager.go
	AddBinDependencyLink(string, string) error
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
	Execute(dir string, stdout io.Writer, stderr io.Writer, program string, args ...string) error
	Output(dir string, program string, args ...string) (string, error)
}

type Supplier struct {
	Manifest         Manifest
	Stager           Stager
	Command          Command
	Log              *libbuildpack.Logger
	MonoVersion      string
	XcRuntimeVersion string
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Supplying dependencies")

	err := s.InstallMono()

	if err != nil {
		return err
	}

	return s.InstallXCRuntime()
}

func (s *Supplier) InstallMono() error {
	monoInstallDir := filepath.Join(s.Stager.DepDir(), "mono"+s.MonoVersion)

	dep := libbuildpack.Dependency{Name: "mono", Version: s.MonoVersion}
	if err := s.Manifest.InstallDependency(dep, monoInstallDir); err != nil {
		s.Log.Error("Error while installing mono %s", err)

		return err
	}

	err := s.Stager.AddBinDependencyLink(filepath.Join(monoInstallDir, "mono", "bin", "mono"), "mono")
	if err != nil {
		return err
	}

	return err
}

func (s *Supplier) InstallXCRuntime() error {
	xcRuntimekDir := filepath.Join(s.Stager.BuildDir(), ".xc-buildpack", "xc-runtime")

	dep := libbuildpack.Dependency{Name: "xcruntime", Version: s.XcRuntimeVersion}
	if err := s.Manifest.InstallDependency(dep, xcRuntimekDir); err != nil {
		s.Log.Error("Error while installing xcruntime %s", err)

		return err
	}

	return nil
}
