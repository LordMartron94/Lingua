package generation

import (
	"fmt"
	"langspec/bootstrap"
	"langspec/cliutil"
	"langspec/dsl"
	"langspec/toolchain"
	"os"
)

type TokenConstraint interface{ ~string }
type NodeConstraint interface{ ~string }

type RunnerConfig[T TokenConstraint, N NodeConstraint] struct {
	SpecPath        string
	Manifest        toolchain.SemanticManifest[T, N]
	OverrideFactory bootstrap.SublimeOverrideFactory
	FileExtensions  []string
	ScopeExtension  string
}

func ExecuteFull[T TokenConstraint, N NodeConstraint](cfg RunnerConfig[T, N]) {
	g := cliutil.NewGenerator(cfg.SpecPath)
	defer g.Close()

	// The engine handles the conversion once, forever.
	stringManifest := adaptManifest(cfg.Manifest)

	opts := []bootstrap.Option{
		bootstrap.WithDiagnosticSink(dsl.DefaultLangSpecDiagnosticSink()),
		bootstrap.WithSublimeToolchain(
			stringManifest,
			cfg.OverrideFactory,
			cfg.FileExtensions,
			cfg.ScopeExtension,
		),
	}

	if err := g.RunToolchains(opts...); err != nil {
		fmt.Fprintf(os.Stderr, "generation failed: %v\n", err)
		os.Exit(1)
	}
}

// adaptManifest performs the exact logic you previously had in ruleforge, but generically.
func adaptManifest[T TokenConstraint, N NodeConstraint](in toolchain.SemanticManifest[T, N]) toolchain.SemanticManifest[string, string] {
	out := toolchain.SemanticManifest[string, string]{
		InvalidScope:    in.InvalidScope,
		BaseTokenScopes: make(map[string]string),
		NodeBindings:    make(map[string]toolchain.NodeBinding[string]),
	}

	for k, v := range in.BaseTokenScopes {
		out.BaseTokenScopes[string(k)] = v
	}

	for k, v := range in.NodeBindings {
		binding := toolchain.NodeBinding[string]{
			Scopes:      v.Scopes,
			MetaScope:   v.MetaScope,
			TokenScopes: make(map[string][]string),
		}
		for tk, tv := range v.TokenScopes {
			binding.TokenScopes[string(tk)] = tv
		}
		out.NodeBindings[string(k)] = binding
	}
	return out
}
