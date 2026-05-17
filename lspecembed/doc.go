// Package lspecembed materializes embedded .lspec trees to temporary filesystem paths
// so langspec can compile them (import resolution requires real paths).
//
// Adoption recipe for a Lingua language package:
//
//  1. Add //go:embed your.lspec in the language package.
//  2. Build a Bundle with BundleCreate (add StdFile() when the spec imports lspec_std).
//  3. Expose EmbeddedLSpecPath() that returns BundleRootPath(bundle).
//  4. In tooling, call ResolvePath("YOUR_LANG_LSPEC_PATH", bundle) for dev overrides.
//
// Materialized directories live for the process lifetime.
package lspecembed
