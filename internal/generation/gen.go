package generation

import (
	"fmt"
	"langspec/bootstrap"
	"langspec/cliutil"
	"langspec/dsl"
	"langspec/toolchain"
	"os"
	"reflect"
	"strconv"
)

/*
TokenConstraint and NodeConstraint are token/node key types from LangSpec go_bindings:
either numeric IDs with String() (e.g. ruleforge) or string enums (e.g. gomod/gowork).
*/
type TokenConstraint interface {
	comparable
}

type NodeConstraint interface {
	comparable
}

// manifestNodeKey turns a node key into the string form consumed by
// SemanticManifestRemapFromStrings. Node enums should keep symbolic names.
func manifestNodeKey[T comparable](k T) string {
	if s, ok := any(k).(fmt.Stringer); ok {
		return s.String()
	}
	v := reflect.ValueOf(k)
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	}
	return fmt.Sprint(k)
}

// manifestTokenKey turns a token key into string form consumed by
// SemanticManifestRemapFromStrings. Token keys intentionally prefer numeric IDs
// over String() to avoid global TokenKind resolver collisions across languages.
func manifestTokenKey[T comparable](k T) string {
	v := reflect.ValueOf(k)
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	}
	if s, ok := any(k).(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(k)
}

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
		out.BaseTokenScopes[manifestTokenKey(k)] = v
	}

	for k, v := range in.NodeBindings {
		binding := toolchain.NodeBinding[string]{
			Scopes:      v.Scopes,
			MetaScope:   v.MetaScope,
			TokenScopes: make(map[string][]string),
		}
		for tk, tv := range v.TokenScopes {
			binding.TokenScopes[manifestTokenKey(tk)] = tv
		}
		out.NodeBindings[manifestNodeKey(k)] = binding
	}
	return out
}
