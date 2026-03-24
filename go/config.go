package goDef

import (
	"autarch/pattern"
	"foundation/domain"
	"langspec/editor"
	"langspec/toolchain"
	"lexarch"
	. "lingua/go/artifacts"
	"lingua/internal/generation"
)

//go:generate go run ../internal/generation/bootstrap/main.go -spec=gomod.lspec
//go:generate go run generate/gen.go

type EditorCtx = editor.EditorCtx[rune, string, string, string, string]

type EditorOverride = editor.EditorOverride[rune, string, string, string, string, toolchain.SublimeContext]

var runeFactory = pattern.RegulaASTFactoryCreate(domain.DiscreteDomainRuneCreate())

func Config() generation.RunnerConfig[Token, Node] {
	return generation.RunnerConfig[Token, Node]{
		SpecPath:        "gomod.lspec",
		FileExtensions:  []string{"go.mod"},
		ScopeExtension:  ".gomod",
		Manifest:        goModSemanticManifest,
		OverrideFactory: goModOverrideProducer,
	}
}

var goModSemanticManifest = toolchain.SemanticManifest[Token, Node]{
	InvalidScope: "invalid.illegal",
	BaseTokenScopes: map[Token]string{
		// ------------------------------------------------
		// Comments & Trivia
		// ------------------------------------------------
		TokLineComment: "comment.line.double-slash",
		TokWhitespace:  "",

		// ------------------------------------------------
		// Keywords
		// ------------------------------------------------
		TokKWModule:    "keyword.declaration.module",
		TokKWGo:        "keyword.declaration.go",
		TokKWRequire:   "keyword.control.require",
		TokKWReplace:   "keyword.control.replace",
		TokKWExclude:   "keyword.control.exclude",
		TokKWRetract:   "keyword.control.retract",
		TokKWToolchain: "keyword.declaration.toolchain",
		TokKWGodebug:   "keyword.declaration.godebug",
		TokKWTool:      "keyword.declaration.tool",

		// ------------------------------------------------
		// Literals & Identifiers
		// ------------------------------------------------
		TokFloat:      "constant.numeric.float",
		TokVersion:    "constant.numeric.version",
		TokModulePath: "string.unquoted.module-path",

		// ------------------------------------------------
		// Punctuation & Operators
		// ------------------------------------------------
		TokArrow:        "keyword.operator.arrow",
		TokEquals:       "keyword.operator.assignment",
		TokComma:        "punctuation.separator.comma",
		TokParenOpen:    "punctuation.section.parens.begin",
		TokParenClose:   "punctuation.section.parens.end",
		TokBracketOpen:  "punctuation.section.brackets.begin",
		TokBracketClose: "punctuation.section.brackets.end",
	},
	NodeBindings: map[Node]toolchain.NodeBinding[Token]{
		// ------------------------------------------------
		// Module Declaration
		// ------------------------------------------------
		NodeModuleStatement: {Scopes: []string{"meta.declaration.module"}},
		NodeModuleName:      {Scopes: []string{"entity.name.module"}},

		// ------------------------------------------------
		// Go Version Declaration
		// ------------------------------------------------
		NodeGoDecl:    {Scopes: []string{"meta.declaration.go.version"}},
		NodeGoVersion: {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Require Block
		// ------------------------------------------------
		NodeRequireStatement:       {Scopes: []string{"meta.declaration.require"}},
		NodeModuleReference:        {Scopes: []string{"meta.reference.module"}},
		NodeModuleReferenceIdent:   {Scopes: []string{"entity.name.reference.module"}},
		NodeModuleReferenceVersion: {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Replace Block
		// ------------------------------------------------
		NodeReplaceStatement:       {Scopes: []string{"meta.declaration.replace"}},
		NodeReplaceReference:       {Scopes: []string{"meta.reference.replace"}},
		NodeReplaceOriginal:        {Scopes: []string{"entity.name.reference.module.original"}},
		NodeReplaceOriginalVersion: {Scopes: []string{"constant.numeric.version"}},
		NodeReplaceTarget:          {Scopes: []string{"entity.name.reference.module.target"}},
		NodeReplaceTargetVersion:   {Scopes: []string{"constant.numeric.version"}},

		// ------------------------------------------------
		// Exclude Block
		// ------------------------------------------------
		NodeExcludeStatement: {Scopes: []string{"meta.declaration.exclude"}},

		// ------------------------------------------------
		// Retract Block
		// ------------------------------------------------
		NodeRetractStatement: {Scopes: []string{"meta.declaration.retract"}},
		NodeRetractReference: {Scopes: []string{"meta.reference.retract"}},
		NodeRetractSingle:    {Scopes: []string{"constant.numeric.version"}},
		NodeRetractRange:     {Scopes: []string{"meta.range.retract"}},
		NodeRetractLow:       {Scopes: []string{"constant.numeric.version.low"}},
		NodeRetractHigh:      {Scopes: []string{"constant.numeric.version.high"}},

		// ------------------------------------------------
		// Toolchain Block
		// ------------------------------------------------
		NodeToolchainStatement: {Scopes: []string{"meta.declaration.toolchain"}},
		NodeToolchainName:      {Scopes: []string{"entity.name.reference.toolchain"}},

		// ------------------------------------------------
		// Godebug Block
		// ------------------------------------------------
		NodeGodebugStatement: {Scopes: []string{"meta.declaration.godebug"}},
		NodeGodebugReference: {Scopes: []string{"meta.reference.godebug"}},
		NodeGodebugKey:       {Scopes: []string{"support.type.property-name"}},
		NodeGodebugValue:     {Scopes: []string{"constant.language.value"}},

		// ------------------------------------------------
		// Tool Block
		// ------------------------------------------------
		NodeToolStatement: {Scopes: []string{"meta.declaration.tool"}},
		NodeToolPath:      {Scopes: []string{"entity.name.reference.tool"}},
	},
}

func goModOverrideProducer(
	ruleset *lexarch.LexingRuleset[rune, string, string],
	ctxProducer func(ctx *EditorCtx) toolchain.SublimeContext,
) func(ec *EditorCtx) []*EditorOverride {
	registry := editor.NewOverrideRegistry[rune, string, string, string, string, toolchain.SublimeContext]()
	patterns := editor.NewTextPatternBuilder[string, string, string, string, toolchain.SublimeContext](runeFactory)

	registry.Register(string(TokLineComment), patterns.LineComment(
		"//",
		toolchain.SublimeContext{Scope: "comment.line.double-slash"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment"},
	))

	return registry.Producer()
}
