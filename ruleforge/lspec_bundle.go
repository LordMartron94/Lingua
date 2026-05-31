package ruleforge

import "lingua/lspecembed"

var ruleforgeLspecBundle = lspecembed.BundleCreate(
	"ruleforge/ruleforge.lspec",
	lspecembed.File{RelPath: "ruleforge/ruleforge.lspec", Data: embeddedRuleforgeLspec},
)

// EmbeddedLSpecPath returns a filesystem path to the embedded ruleforge.lspec.
func EmbeddedLSpecPath() (string, error) {
	return lspecembed.BundleRootPath(ruleforgeLspecBundle)
}

// ResolveRuleforgeLSpecPath returns RULEFORGE_LSPEC_PATH when set, otherwise the embedded spec path.
func ResolveRuleforgeLSpecPath() (path string, release func(), err error) {
	return lspecembed.ResolvePath("RULEFORGE_LSPEC_PATH", ruleforgeLspecBundle)
}
