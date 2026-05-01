package generation

import (
	"fmt"
	"langspec/cliutil"
	"langspec/toolchain"
	"os"
)

/*
TokenConstraint and NodeConstraint bound token/node key types from LangSpec go_bindings:
either numeric IDs with String() (e.g. ruleforge) or string enums (e.g. gomod/gowork).
*/
type TokenConstraint = cliutil.TokenConstraint

type NodeConstraint = cliutil.NodeConstraint

// RunnerConfig is the Lingua-facing name for cliutil.InMemorySublimeRunnerConfig.
type RunnerConfig[T TokenConstraint, N NodeConstraint] = cliutil.InMemorySublimeRunnerConfig[T, N]

// ExecuteFull runs LangSpec toolchains for a Go-defined in-memory Sublime configuration.
func ExecuteFull[T TokenConstraint, N NodeConstraint](cfg RunnerConfig[T, N]) {
	if err := cliutil.RunInMemorySublimeToolchains(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "generation failed: %v\n", err)
		os.Exit(1)
	}
}

/*
ExportSemanticManifest converts a typed semantic manifest to string keys for
toolchains and tests (same mapping as AdaptInMemorySublimeManifest).
*/
func ExportSemanticManifest[T TokenConstraint, N NodeConstraint](in toolchain.SemanticManifest[T, N]) toolchain.SemanticManifest[string, string] {
	return cliutil.AdaptInMemorySublimeManifest(in)
}
