package goDef

import "lingua/lspecembed"

var goModLspecBundle = lspecembed.BundleCreate(
	"go/gomod.lspec",
	lspecembed.File{RelPath: "go/gomod.lspec", Data: embeddedGoModLspec},
)

var goWorkLspecBundle = lspecembed.BundleCreate(
	"go/gowork.lspec",
	lspecembed.File{RelPath: "go/gowork.lspec", Data: embeddedGoWorkLspec},
)

// EmbeddedGoModLSpecPath returns a filesystem path to the embedded gomod.lspec.
func EmbeddedGoModLSpecPath() (string, error) {
	return lspecembed.BundleRootPath(goModLspecBundle)
}

// EmbeddedGoWorkLSpecPath returns a filesystem path to the embedded gowork.lspec.
func EmbeddedGoWorkLSpecPath() (string, error) {
	return lspecembed.BundleRootPath(goWorkLspecBundle)
}

// ResolveGoModLSpecPath returns GOMOD_LSPEC_PATH when set, otherwise the embedded spec path.
func ResolveGoModLSpecPath() (path string, release func(), err error) {
	return lspecembed.ResolvePath("GOMOD_LSPEC_PATH", goModLspecBundle)
}

// ResolveGoWorkLSpecPath returns GOWORK_LSPEC_PATH when set, otherwise acquires the embedded spec.
func ResolveGoWorkLSpecPath() (path string, release func(), err error) {
	return lspecembed.ResolvePath("GOWORK_LSPEC_PATH", goWorkLspecBundle)
}
