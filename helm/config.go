package helm

import (
	"autarch/pattern"
	"foundation/domain"
	"langspec/editor"
	"langspec/toolchain"
	. "lingua/helm/artifacts"
	"lingua/internal/generation"
)

// ------------------------------------------------------------------ ALIASES

type EditorCtx = editor.EditorCtx[rune, uint32, uint32, string, uint32]

type EditorOverride = editor.EditorOverride[rune, uint32, uint32, string, uint32, toolchain.SublimeContext]

var runeFactory = pattern.RegulaASTFactoryCreate(domain.DiscreteDomainRuneCreate())

// ------------------------------------------------------------------ CONFIG

func Config() generation.RunnerConfig[Token, Node] {
	return generation.RunnerConfig[Token, Node]{
		SpecPath:        "helm.lspec",
		FileExtensions:  []string{".helm"},
		ScopeExtension:  ".helm",
		Manifest:        helmSemanticManifest,
		OverrideFactory: helmOverrideProducer,
	}
}

var helmSemanticManifest = toolchain.SemanticManifest[Token, Node]{
	InvalidScope: "invalid.illegal",
	BaseTokenScopes: map[Token]string{
		TokNewLine:             "whitespace.newline",
		TokEquals:              "keyword.operator.assignment",
		TokStringLiteral:       "string.quoted.double",
		TokStringStart:         "punctuation.definition.string.begin",
		TokStringEnd:           "punctuation.definition.string.end",
		TokStringInterpolation: "constant.other.placeholder",

		// Variable fallback
		TokIdentifier: "variable.other.readwrite",
	},
	NodeBindings: map[Node]toolchain.NodeBinding[Token]{},
}

func helmOverrideProducer(
	ruleset *editor.LexingRuleSet[rune, uint32, uint32],
	ctxProducer func(ctx *EditorCtx) toolchain.SublimeContext,
) func(ec *EditorCtx) []*EditorOverride {
	registry := editor.NewOverrideRegistry[rune, uint32, uint32, string, uint32, toolchain.SublimeContext]()
	patterns := editor.NewTextPatternBuilder[uint32, uint32, string, uint32, toolchain.SublimeContext](runeFactory)

	registry.Register(uint32(TokLineComment), patterns.LineComment(
		"//",
		toolchain.SublimeContext{Scope: "comment.line.double-slash"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment"},
	))

	registry.Register(uint32(TokBlockComment), patterns.BlockComment(
		"/*",
		"*/",
		toolchain.SublimeContext{MetaScope: "comment.block"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment.begin"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment.end"},
	))

	return registry.Producer()
}
