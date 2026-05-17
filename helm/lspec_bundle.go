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

// ResolveHelmLSpecPath returns HELM_LSPEC_PATH when set, otherwise the embedded spec path.
func ResolveHelmLSpecPath() (string, error) {
	return lspecembed.ResolvePath("HELM_LSPEC_PATH", HelmLspecBundle)
}
