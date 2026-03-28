package goDef

import (
	"langspec/editor"
	"langspec/toolchain"
	"lexarch"
	. "lingua/go/artifacts"
	"lingua/internal/generation"
)

func GoWorkConfig() generation.RunnerConfig[GoWorkToken, GoWorkNode] {
	return generation.RunnerConfig[GoWorkToken, GoWorkNode]{
		SpecPath:        "gowork.lspec",
		FileExtensions:  []string{"go.work"},
		ScopeExtension:  ".gowork",
		Manifest:        goWorkSemanticManifest,
		OverrideFactory: goWorkOverrideProducer,
	}
}

var goWorkSemanticManifest = toolchain.SemanticManifest[GoWorkToken, GoWorkNode]{
	InvalidScope: "invalid.illegal",
	BaseTokenScopes: map[GoWorkToken]string{
		// ------------------------------------------------
		// Comments & Trivia
		// ------------------------------------------------
		GoWork_TokLineComment: "comment.line.double-slash",
		GoWork_TokWhitespace:  "",
		GoWork_TokEOF:         "",

		// ------------------------------------------------
		// Keywords
		// ------------------------------------------------
		GoWork_TokKWGo:        "keyword.declaration.go",
		GoWork_TokKWUse:       "keyword.control.use",
		GoWork_TokKWReplace:   "keyword.control.replace",
		GoWork_TokKWToolchain: "keyword.declaration.toolchain",
		GoWork_TokKWGodebug:   "keyword.declaration.godebug",

		// ------------------------------------------------
		// Literals & Identifiers
		// ------------------------------------------------
		GoWork_TokFloat:      "constant.numeric.float",
		GoWork_TokVersion:    "constant.numeric.version",
		GoWork_TokModulePath: "string.unquoted.module-path",

		// ------------------------------------------------
		// Punctuation & Operators
		// ------------------------------------------------
		GoWork_TokArrow:      "keyword.operator.arrow",
		GoWork_TokEquals:     "keyword.operator.assignment",
		GoWork_TokParenOpen:  "punctuation.section.parens.begin",
		GoWork_TokParenClose: "punctuation.section.parens.end",
	},
	NodeBindings: map[GoWorkNode]toolchain.NodeBinding[GoWorkToken]{
		// ------------------------------------------------
		// Go Version Declaration
		// ------------------------------------------------
		GoWork_NodeGoDecl:    {Scopes: []string{"meta.declaration.go.version"}},
		GoWork_NodeGoVersion: {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Use Block (Workspace specific)
		// ------------------------------------------------
		GoWork_NodeUseStatement: {Scopes: []string{"meta.declaration.use"}},
		GoWork_NodeUsePath:      {Scopes: []string{"entity.name.reference.module.use"}},

		// ------------------------------------------------
		// Replace Block
		// ------------------------------------------------
		GoWork_NodeReplaceStatement:       {Scopes: []string{"meta.declaration.replace"}},
		GoWork_NodeReplaceReference:       {Scopes: []string{"meta.reference.replace"}},
		GoWork_NodeReplaceOriginal:        {Scopes: []string{"entity.name.reference.module.original"}},
		GoWork_NodeReplaceOriginalVersion: {Scopes: []string{"constant.numeric.version"}},
		GoWork_NodeReplaceTarget:          {Scopes: []string{"entity.name.reference.module.target"}},
		GoWork_NodeReplaceTargetVersion:   {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Toolchain Block
		// ------------------------------------------------
		GoWork_NodeToolchainStatement: {Scopes: []string{"meta.declaration.toolchain"}},
		GoWork_NodeToolchainName:      {Scopes: []string{"entity.name.reference.toolchain"}},

		// ------------------------------------------------
		// Godebug Block
		// ------------------------------------------------
		GoWork_NodeGodebugStatement: {Scopes: []string{"meta.declaration.godebug"}},
		GoWork_NodeGodebugReference: {Scopes: []string{"meta.reference.godebug"}},
		GoWork_NodeGodebugKey:       {Scopes: []string{"support.type.property-name"}},
		GoWork_NodeGodebugValue:     {Scopes: []string{"constant.language.value"}},
	},
}

func goWorkOverrideProducer(
	_ *lexarch.LexingRuleset[rune, uint32, uint32],
	_ func(ctx *EditorCtx) toolchain.SublimeContext,
) func(ec *EditorCtx) []*EditorOverride {
	registry := editor.NewOverrideRegistry[rune, uint32, uint32, string, uint32, toolchain.SublimeContext]()
	patterns := editor.NewTextPatternBuilder[uint32, uint32, string, uint32, toolchain.SublimeContext](runeFactory)

	registry.Register(uint32(GoWork_TokLineComment), patterns.LineComment(
		"//",
		toolchain.SublimeContext{Scope: "comment.line.double-slash"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment"},
	))

	return registry.Producer()
}
