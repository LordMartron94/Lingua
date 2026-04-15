package ebnf

import (
	"autarch/pattern"
	"foundation/domain"
	"langspec/dsl"
	"langspec/editor"
	"langspec/toolchain"
	"lingua/internal/generation"

	. "lingua/ebnf/artifacts"
)

type EditorCtx = editor.EditorCtx[rune, uint32, uint32, string, dsl.LangSpecParserNodeKind]

type EditorOverride = editor.EditorOverride[rune, uint32, uint32, string, dsl.LangSpecParserNodeKind, toolchain.SublimeContext]

var runeFactory = pattern.RegulaASTFactoryCreate(domain.DiscreteDomainRuneCreate())

func Config() generation.RunnerConfig[Token, Node] {
	return generation.RunnerConfig[Token, Node]{
		SpecPath:        "ebnf.lspec",
		FileExtensions:  []string{".ebnf"},
		ScopeExtension:  ".ebnf",
		Manifest:        ebnfSemanticManifest,
		OverrideFactory: ebnfOverrideProducer,
	}
}

var ebnfSemanticManifest = toolchain.SemanticManifest[Token, Node]{
	InvalidScope: "invalid.illegal",
	BaseTokenScopes: map[Token]string{
		TokAssignment:  "keyword.operator.assignment",
		TokIdentifier:  "variable.other",
		TokSemi:        "punctuation.terminator",
		TokParenOpen:   "punctuation.section.parens.begin",
		TokParenClose:  "punctuation.section.parens.close",
		TokQuestion:    "keyword.operator.optional",
		TokStar:        "keyword.operator.star",
		TokPlus:        "keyword.operator.plus",
		TokPipe:        "keyword.operator.alternation",
		TokRange:       "keyword.operator.range",
		TokTilde:       "keyword.operator.negation",
		TokExclamation: "keyword.operator.negation",

		TokStringSingleStart: "punctuation.definition.string.begin",
		TokStringSingleEnd:   "punctuation.definition.string.end",
		TokStringDoubleStart: "punctuation.definition.string.begin",
		TokStringDoubleEnd:   "punctuation.definition.string.end",
		TokStringEscape:      "constant.character.escape",
		TokStringSingleText:  "string.quoted.single",
		TokStringDoubleText:  "string.quoted.double",
	},
	NodeBindings: map[Node]toolchain.NodeBinding[Token]{
		NodeProductionGroup: {
			TokenScopes: map[Token][]string{
				TokParenOpen:  {"keyword.operator.group", "punctuation.section.parens.begin"},
				TokParenClose: {"keyword.operator.group", "punctuation.section.parens.close"},
			},
		},
		NodeProductionReference: {
			TokenScopes: map[Token][]string{
				TokIdentifier: {"variable.function"},
			},
		},
		NodeProductionName: {
			TokenScopes: map[Token][]string{
				TokIdentifier: {"meta.production", "entity.name.function"},
			},
		},
		NodeCharRange: {
			MetaScope: "meta.character-range",
		},
		NodeStringLiteral: {
			MetaScope: "meta.string",
		},
	},
}

func ebnfOverrideProducer(
	ruleset *editor.LexingRuleSet[rune, uint32, uint32],
	ctxProducer func(ctx *EditorCtx) toolchain.SublimeContext,
) func(ec *EditorCtx) []*EditorOverride {
	registry := editor.NewOverrideRegistry[rune, uint32, uint32, string, dsl.LangSpecParserNodeKind, toolchain.SublimeContext]()
	patterns := editor.NewTextPatternBuilder[uint32, uint32, string, dsl.LangSpecParserNodeKind, toolchain.SublimeContext](runeFactory)

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
