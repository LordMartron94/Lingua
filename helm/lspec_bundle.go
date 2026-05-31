package helm

import "lingua/lspecembed"

// HelmLspecBundle is the embedded helm.lspec tree (including lspec_std for IMPORT).
var HelmLspecBundle = lspecembed.BundleCreate(
	"helm/helm.lspec",
	lspecembed.File{RelPath: "helm/helm.lspec", Data: embeddedHelmLspec},
	lspecembed.StdFile(),
)

// EmbeddedLSpecPath returns a filesystem path to the embedded helm.lspec.
func EmbeddedLSpecPath() (string, error) {
	return lspecembed.BundleRootPath(HelmLspecBundle)
}

// ResolveHelmLSpecPath returns HELM_LSPEC_PATH when set, otherwise acquires the embedded spec.
// Call release when the path is no longer needed (release is a no-op when HELM_LSPEC_PATH is set).
func ResolveHelmLSpecPath() (path string, release func(), err error) {
	return lspecembed.ResolvePath("HELM_LSPEC_PATH", HelmLspecBundle)
}
