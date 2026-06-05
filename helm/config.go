package helm

import (
	"fmt"

	"autarch/pattern"
	"foundation/domain"
	"langspec/dsl"
	"langspec/editor"
	"langspec/toolchain"
	. "lingua/helm/artifacts"
	"lingua/internal/generation"
)

// ------------------------------------------------------------------ ALIASES

type EditorCtx = editor.EditorCtx[rune, uint32, uint32, string, dsl.LangSpecParserNodeKind]

type EditorOverride = editor.EditorOverride[rune, uint32, uint32, string, dsl.LangSpecParserNodeKind, toolchain.SublimeContext]

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
		TokNewLine:              "whitespace.newline",
		TokEquals:               "keyword.operator.assignment",
		TokStringLiteral:        "string.quoted.double",
		TokEmptyString:          "string.quoted.double",
		TokStringStart:          "punctuation.definition.string.begin",
		TokStringEnd:            "punctuation.definition.string.end",
		TokStringEscape:         "constant.character.escape",
		TokMultilineStringStart: "punctuation.definition.string.begin",
		TokMultilineStringEnd:   "punctuation.definition.string.end",
		TokMultilineStringText:  "string.quoted.double",
		TokInterpolationStart:   "punctuation.section.interpolation.begin",
		TokInterpolationEnd:     "punctuation.section.interpolation.end",
		TokBraceOpen:            "punctuation.section.braces.begin",
		TokBraceClose:           "punctuation.section.braces.end",
		TokParenOpen:            "punctuation.section.parens.begin",
		TokParenClose:           "punctuation.section.parens.end",
		TokBracketOpen:          "punctuation.section.brackets.begin",
		TokBracketClose:         "punctuation.section.brackets.close",
		TokQuestion:             "punctuation.question",
		TokKWTarget:             "keyword.declaration.target",
		TokComma:                "punctuation.separator.comma",

		TokKWLet: "keyword.declaration.let",

		TokKWExport: "keyword.declaration.export",

		TokKWCollect: "support.function.builtin.collect",

		TokKWPath: "support.function.builtin.path",
		TokKWGlob: "support.function.builtin.glob",

		TokKWDefined:    "support.function.builtin.defined",
		TokKWNotDefined: "support.function.builtin.not-defined",
		TokKWEquals:     "support.function.builtin.equals",
		TokKWNotEquals:  "support.function.builtin.not-equals",

		TokKWWhen: "keyword.control.condition.when",

		TokKWRun:       "keyword.control.execution",
		TokKWDependsOn: "keyword.control.depends-on",

		TokKWHelp:        "support.type.property-name.help",
		TokKWInteractive: "support.type.property-name.interactive",
		TokKWHidden:      "support.type.property-name.hidden",
		TokKWArtifacts:   "keyword.declaration.artifacts",
		TokKWMatrix:      "keyword.declaration.matrix",
		TokKWAliases:     "support.type.property-name.aliases",
		TokKWIn:          "keyword.operator.iteration.in",

		TokKWInputs:   "support.type.property-name.inputs",
		TokKWOutputs:  "support.type.property-name.outputs",
		TokKWVolatile: "support.type.property-name.volatile",
		TokKWDynamic:  "support.type.property-name.dynamic",

		TokKWConfirm:  "support.type.property-name.confirm",
		TokKWOptional: "support.type.property-name.optional",
		TokKWParam:    "keyword.operator.parameter-splice",
		TokKWParams:   "support.type.property-name.params",

		TokKWTrue:  "constant.language.boolean.true",
		TokKWFalse: "constant.language.boolean.false",

		TokKWEnv:        "keyword.declaration.env",
		TokKWWorkingDir: "support.type.property-name.workdir",

		// Variable fallback
		TokIdentifier: "variable.other.readwrite",
	},
	NodeBindings: map[Node]toolchain.NodeBinding[Token]{
		NodeStringLiteral: {
			MetaScope: "meta.string",
		},
		NodeMultilineString: {
			MetaScope: "meta.string.multiline",
		},
		NodeStringInterpolation: {
			MetaScope: "meta.interpolation",
		},
		NodeStringEscape: {
			Scopes: []string{"constant.character.escape"},
		},
		NodeStringText: {
			Scopes: []string{"string.quoted.double"},
		},
		NodeInterpolatedVariable: {
			Scopes: []string{"variable.other.interpolated"},
		},
		NodeTargetParams: {
			MetaScope: "meta.target.parameters",
		},
		NodeTargetParameter: {
			MetaScope: "meta.target.parameter",
		},
		NodeTargetParameterIdentifier: {
			Scopes: []string{"variable.parameter"},
		},
		NodeTargetParameterOptional: {
			Scopes: []string{"keyword.operator.optional"},
		},
		NodeTargetIdentifier: {
			Scopes: []string{"entity.name.function.target"},
		},
		NodeTargetBody: {
			MetaScope: "meta.target.body",
		},
		NodeVariableReference: {
			Scopes: []string{"variable.other.reference"},
		},
		NodeGlobKwargIdentifier: {
			Scopes: []string{"variable.parameter"},
		},
		NodeConditionalParameter: {
			Scopes: []string{"variable.parameter"},
		},
		NodeStringDollar: {
			Scopes: []string{"string.quoted.double"},
		},
		NodeNumber: {
			Scopes: []string{"constant.language.numeric"},
		},
		NodeInvokeTarget: {
			Scopes: []string{"meta.target-reference", "variable.other.target.reference"},
		},
		NodeInputArray: {
			MetaScope: "meta.array.inputs",
		},
		NodeMatrixArray: {
			MetaScope: "meta.array.matrix",
		},
		NodeOutputArray: {
			MetaScope: "meta.array.outputs",
		},
		NodeDependencyParameterName: {
			Scopes: []string{"variable.other.property.dependency-parameter"},
		},
		NodeDependsOnParameterRef: {
			MetaScope: "meta.depends-on.parameter-ref",
		},
		NodeRunArray: {
			MetaScope: "meta.run.array",
		},
		NodeRunParameterRef: {
			MetaScope: "meta.run.parameter-ref",
		},
		NodeRunParameterName: {
			Scopes: []string{"variable.parameter"},
		},
		NodeRunAbsPathCall: {
			MetaScope: "meta.run.abs-path-call",
		},
		NodeRunAbsPathCallName: {
			Scopes: []string{"support.function.builtin.abs-path"},
		},
		NodeLetBinding: {
			MetaScope: "meta.let-binding",
		},
		NodeLetBindingName: {
			Scopes: []string{"entity.name.binding.let"},
		},
		NodePathVarOrCall: {
			MetaScope: "meta.path-expr",
		},
		NodePathVarOrCallName: {
			Scopes: []string{"variable.other.reference"},
		},
		NodePathCallArgs: {
			MetaScope: "meta.path-call-args.body",
		},
		// --- Environment & Properties ---
		NodeEnvKey: {
			Scopes: []string{"variable.other.property.env"},
		},

		// --- Structural Meta Boundaries ---
		NodeArtifactsBlock: {
			MetaScope: "meta.block.artifacts",
		},
		NodeEnvDeclaration: {
			MetaScope: "meta.block.env",
		},
		NodeExportDeclaration: {
			MetaScope: "meta.block.export",
		},
		NodeExportKey: {
			Scopes: []string{"variable.other.property.export"},
		},
		NodeStringListArray: {
			MetaScope: "meta.array.string-list",
		},
		NodeStringListParamRef: {
			MetaScope: "meta.string-list.parameter-ref",
		},
		NodeStringListCollectCall: {
			MetaScope: "meta.string-list.collect-call",
		},
		NodeStringListCollectArgs: {
			MetaScope: "meta.string-list.collect-args",
		},
		NodeStringListCollectDepsParam: {
			Scopes: []string{"variable.parameter"},
		},
		NodeParamValueArray: {
			MetaScope: "meta.array.param-value",
		},
		NodeParamValueParamRef: {
			MetaScope: "meta.param-value.parameter-ref",
		},
		NodeParamValueParamName: {
			Scopes: []string{"variable.parameter"},
		},
		NodeTargetDepends: {
			MetaScope: "meta.block.depends-on",
		},
	},
}

func helmOverrideProducer(
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

	registry.Register(uint32(TokIdentifier), helmPathExprIdentifierOverride(ruleset))

	return registry.Producer()
}

func helmPathExprIdentifierOverride(
	ruleset *editor.LexingRuleSet[rune, uint32, uint32],
) editor.OverrideHandler[rune, uint32, uint32, string, dsl.LangSpecParserNodeKind, toolchain.SublimeContext] {
	identRegex := helmTokenRegex(ruleset, TokIdentifier)
	parenRegex := helmTokenRegex(ruleset, TokParenOpen)

	pathBuiltinRegex := fmt.Sprintf(
		`(?:map_ext|join_prefix|rebase_dir)(?=\s*%s)`,
		parenRegex,
	)
	pathBuiltinCtx := toolchain.SublimeContext{Scope: "support.function.builtin.path-expr"}

	pathVarOrCallName := dsl.LangSpecParserNodeKind(NodePathVarOrCallName)

	return func(ctx *EditorCtx) []*EditorOverride {
		if ctx.NodeKind == nil || *ctx.NodeKind != pathVarOrCallName {
			return nil
		}

		return []*EditorOverride{
			{PatternRegex: &pathBuiltinRegex, MatchContext: &pathBuiltinCtx},
			{
				PatternRegex: &identRegex,
				MatchContext: &toolchain.SublimeContext{Scope: "variable.other.reference"},
			},
		}
	}
}

func helmTokenRegex(
	ruleset *editor.LexingRuleSet[rune, uint32, uint32],
	token Token,
) string {
	for _, rule := range editor.LexingRuleSetGetRules(ruleset) {
		if rule.Token == uint32(token) {
			regex, err := rule.Pattern.ToRegEx()
			if err != nil {
				panic(fmt.Sprintf("helm: token %v regex: %v", token, err))
			}
			return regex
		}
	}
	panic(fmt.Sprintf("helm: token %v not found in lexing ruleset", token))
}
