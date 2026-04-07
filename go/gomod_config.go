package goDef

import (
	"autarch/pattern"
	"foundation/domain"
	"langspec/editor"
	"langspec/toolchain"
	. "lingua/go/artifacts"
	"lingua/internal/generation"
)

type EditorCtx = editor.EditorCtx[rune, uint32, uint32, string, uint32]

type EditorOverride = editor.EditorOverride[rune, uint32, uint32, string, uint32, toolchain.SublimeContext]

var runeFactory = pattern.RegulaASTFactoryCreate(domain.DiscreteDomainRuneCreate())

func GoModConfig() generation.RunnerConfig[GoModToken, GoModNode] {
	return generation.RunnerConfig[GoModToken, GoModNode]{
		SpecPath:        "gomod.lspec",
		FileExtensions:  []string{"go.mod"},
		ScopeExtension:  ".gomod",
		Manifest:        goModSemanticManifest,
		OverrideFactory: goModOverrideProducer,
	}
}

var goModSemanticManifest = toolchain.SemanticManifest[GoModToken, GoModNode]{
	InvalidScope: "invalid.illegal",
	BaseTokenScopes: map[GoModToken]string{
		// ------------------------------------------------
		// Comments & Trivia
		// ------------------------------------------------
		GoMod_TokLineComment: "comment.line.double-slash",
		GoMod_TokWhitespace:  "",
		GoMod_TokEOF:         "",

		// ------------------------------------------------
		// Keywords
		// ------------------------------------------------
		GoMod_TokKWModule:    "keyword.declaration.module",
		GoMod_TokKWGo:        "keyword.declaration.go",
		GoMod_TokKWRequire:   "keyword.control.require",
		GoMod_TokKWReplace:   "keyword.control.replace",
		GoMod_TokKWExclude:   "keyword.control.exclude",
		GoMod_TokKWRetract:   "keyword.control.retract",
		GoMod_TokKWToolchain: "keyword.declaration.toolchain",
		GoMod_TokKWGodebug:   "keyword.declaration.godebug",
		GoMod_TokKWTool:      "keyword.declaration.tool",

		// ------------------------------------------------
		// Literals & Identifiers
		// ------------------------------------------------
		GoMod_TokFloat:      "constant.numeric.float",
		GoMod_TokVersion:    "constant.numeric.version",
		GoMod_TokModulePath: "string.unquoted.module-path",

		// ------------------------------------------------
		// Punctuation & Operators
		// ------------------------------------------------
		GoMod_TokArrow:        "keyword.operator.arrow",
		GoMod_TokEquals:       "keyword.operator.assignment",
		GoMod_TokComma:        "punctuation.separator.comma",
		GoMod_TokParenOpen:    "punctuation.section.parens.begin",
		GoMod_TokParenClose:   "punctuation.section.parens.end",
		GoMod_TokBracketOpen:  "punctuation.section.brackets.begin",
		GoMod_TokBracketClose: "punctuation.section.brackets.end",
	},
	NodeBindings: map[GoModNode]toolchain.NodeBinding[GoModToken]{
		// ------------------------------------------------
		// Module Declaration
		// ------------------------------------------------
		GoMod_NodeModuleStatement: {Scopes: []string{"meta.declaration.module"}},
		GoMod_NodeModuleName:      {Scopes: []string{"entity.name.module"}},

		// ------------------------------------------------
		// Go Version Declaration
		// ------------------------------------------------
		GoMod_STD__NodeGoDecl:    {Scopes: []string{"meta.declaration.go.version"}},
		GoMod_STD__NodeGoVersion: {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Require Block
		// ------------------------------------------------
		GoMod_NodeRequireStatement:       {Scopes: []string{"meta.declaration.require"}},
		GoMod_NodeModuleReference:        {Scopes: []string{"meta.reference.module"}},
		GoMod_NodeModuleReferenceIdent:   {Scopes: []string{"entity.name.reference.module"}},
		GoMod_NodeModuleReferenceVersion: {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Replace Block
		// ------------------------------------------------
		GoMod_STD__NodeReplaceStatement:       {Scopes: []string{"meta.declaration.replace"}},
		GoMod_STD__NodeReplaceReference:       {Scopes: []string{"meta.reference.replace"}},
		GoMod_STD__NodeReplaceOriginal:        {Scopes: []string{"entity.name.reference.module.original"}},
		GoMod_STD__NodeReplaceOriginalVersion: {Scopes: []string{"constant.numeric.version"}},
		GoMod_STD__NodeReplaceTarget:          {Scopes: []string{"entity.name.reference.module.target"}},
		GoMod_STD__NodeReplaceTargetVersion:   {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Exclude Block
		// ------------------------------------------------
		GoMod_NodeExcludeStatement: {Scopes: []string{"meta.declaration.exclude"}},

		// ------------------------------------------------
		// Retract Block
		// ------------------------------------------------
		GoMod_NodeRetractStatement: {Scopes: []string{"meta.declaration.retract"}},
		GoMod_NodeRetractReference: {Scopes: []string{"meta.reference.retract"}},
		GoMod_NodeRetractSingle:    {Scopes: []string{"constant.numeric.version"}},
		GoMod_NodeRetractRange:     {Scopes: []string{"meta.range.retract"}},
		GoMod_NodeRetractLow:       {Scopes: []string{"constant.numeric.version.low"}},
		GoMod_NodeRetractHigh:      {Scopes: []string{"constant.numeric.version.high"}},

		// ------------------------------------------------
		// Toolchain Block
		// ------------------------------------------------
		GoMod_STD__NodeToolchainStatement: {Scopes: []string{"meta.declaration.toolchain"}},
		GoMod_STD__NodeToolchainName:      {Scopes: []string{"entity.name.reference.toolchain"}},

		// ------------------------------------------------
		// Godebug Block
		// ------------------------------------------------
		GoMod_STD__NodeGodebugStatement: {Scopes: []string{"meta.declaration.godebug"}},
		GoMod_STD__NodeGodebugReference: {Scopes: []string{"meta.reference.godebug"}},
		GoMod_STD__NodeGodebugKey:       {Scopes: []string{"support.type.property-name"}},
		GoMod_STD__NodeGodebugValue:     {Scopes: []string{"constant.language.value"}},

		// ------------------------------------------------
		// Tool Block
		// ------------------------------------------------
		GoMod_NodeToolStatement: {Scopes: []string{"meta.declaration.tool"}},
		GoMod_NodeToolPath:      {Scopes: []string{"entity.name.reference.tool"}},
	},
}

func goModOverrideProducer(
	_ *editor.LexingRuleSet[rune, uint32, uint32],
	_ func(ctx *EditorCtx) toolchain.SublimeContext,
) func(ec *EditorCtx) []*EditorOverride {
	registry := editor.NewOverrideRegistry[rune, uint32, uint32, string, uint32, toolchain.SublimeContext]()
	patterns := editor.NewTextPatternBuilder[uint32, uint32, string, uint32, toolchain.SublimeContext](runeFactory)

	registry.Register(uint32(GoMod_TokLineComment), patterns.LineComment(
		"//",
		toolchain.SublimeContext{Scope: "comment.line.double-slash"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment"},
	))

	return registry.Producer()
}
