package ruleforge

import (
	"autarch/pattern"
	"fmt"
	"foundation/domain"
	"langspec/editor"
	"langspec/toolchain"
	"lexarch"
	"lingua/internal/generation"
	. "lingua/ruleforge/artifacts"
)

type EditorCtx = editor.EditorCtx[rune, string, string, string, string]

type EditorOverride = editor.EditorOverride[rune, string, string, string, string, toolchain.SublimeContext]

var runeFactory = pattern.RegulaASTFactoryCreate(domain.DiscreteDomainRuneCreate())

func Config() generation.RunnerConfig[Token, Node] {
	return generation.RunnerConfig[Token, Node]{
		SpecPath:        "ruleforge.lspec",
		FileExtensions:  []string{".rf"},
		ScopeExtension:  ".rf",
		Manifest:        ruleForgeSemanticManifest,
		OverrideFactory: ruleforgeOverrideProducer,
	}
}

var ruleForgeSemanticManifest = toolchain.SemanticManifest[Token, Node]{
	InvalidScope: "invalid.illegal",
	BaseTokenScopes: map[Token]string{

		// ------------------------------------------------
		// Comments / trivia
		// ------------------------------------------------
		TokBlockComment: "comment.block",
		TokLineComment:  "comment.line.double-slash",
		TokWhitespace:   "",
		"TokEOF":        "",

		// ------------------------------------------------
		// Core language declarations (The Blueprints)
		// ------------------------------------------------
		TokKWVersion: "keyword.other.version",
		TokKWModule:  "keyword.declaration.module",

		TokKWExport:  "storage.modifier.export",
		TokKWPrivate: "storage.modifier.private",

		TokKWExtends:    "storage.modifier.extends",
		TokKWImplements: "storage.modifier.implements",
		TokKWOverride:   "keyword.control.override",
		TokKWInherit:    "keyword.control.inherit",
		TokKWImport:     "keyword.control.import",
		TokKWAs:         "keyword.control.import.as",

		TokKWTemplate: "storage.type.function.template keyword.declaration.function.template",

		TokKWConfiguration: "storage.type.class.configuration keyword.declaration.configuration",
		TokKWAdapter:       "storage.type.class.adapter keyword.declaration.adapter",
		TokKWFunc:          "storage.type.function keyword.declaration.function",
		TokKWInclude:       "keyword.control.include",
		TokKWWith:          "keyword.control.with",

		// ------------------------------------------------
		// Ruleset & Rule Structural Declarations
		// ------------------------------------------------
		TokKWRuleset: "storage.type.class.ruleset keyword.declaration.class.ruleset",
		TokKWRule:    "storage.type.function.rule keyword.declaration.function.rule",

		TokKWDescription:    "support.type.property-name.description",
		TokKWCoreConditions: "keyword.control.block.core-conditions",

		// ------------------------------------------------
		// Theme DSL declarations
		// ------------------------------------------------
		TokKWThemeSchema: "storage.type.struct.theme-schema keyword.declaration.struct.theme-schema",
		TokKWTheme:       "storage.type.struct.theme keyword.declaration.struct.theme",
		TokKWColorSet:    "storage.type.struct.color-set keyword.declaration.color-set",

		// ------------------------------------------------
		// Block Control Flow (Template / Rule Structure)
		// ------------------------------------------------
		TokKWScope:    "keyword.control.block.scope",
		TokKWBaseLine: "keyword.control.block.baseline",
		TokKWVariants: "keyword.control.block.variants",
		TokKWWhen:     "keyword.control.conditional.when",

		// ------------------------------------------------
		// Properties (The "CSS" of Ruleforge)
		// ------------------------------------------------
		TokKWAction: "support.type.property-name.action",
		TokKWStyle:  "support.type.property-name.style",

		TokKWBackground: "support.type.property-name.background",
		TokKWBorder:     "support.type.property-name.border",
		TokKWTextColor:  "support.type.property-name.text-color",
		TokKWFontSize:   "support.type.property-name.font-size",

		TokKWSize:  "support.type.property-name.size",
		TokKWColor: "support.type.property-name.color",
		TokKWShape: "support.type.property-name.shape",

		// Effects
		TokKWEffectColor:     "support.type.property-name.effect-color",
		TokKWEffectTemporary: "support.type.property-name.effect-temporary",

		// Event hooks (Control flow branching based on game state)
		TokKWOnShow: "keyword.control.event.on-show",
		TokKWOnHide: "keyword.control.event.on-hide",

		// Sound & Minimap Blocks
		TokKWSound:            "support.type.property-name.sound",
		TokKWSoundID:          "support.type.property-name.sound-id",
		TokKWSoundPath:        "support.type.property-name.sound-path",
		TokKWSoundVolume:      "support.type.property-name.sound-volume",
		TokKWPositional:       "support.type.property-name.sound-positional",
		TokKWDropSoundEnabled: "support.type.property-name.sound-enabled",

		TokKWMinimap: "support.type.property-name.minimap",

		// ------------------------------------------------
		// Structural punctuation
		// ------------------------------------------------
		TokSemi:         "punctuation.terminator.statement",
		TokDot:          "punctuation.accessor.dot",
		TokBraceOpen:    "punctuation.section.braces.begin",
		TokBraceClose:   "punctuation.section.braces.end",
		TokEquals:       "keyword.operator.assignment",
		TokComma:        "punctuation.separator.comma",
		TokParenOpen:    "punctuation.section.parens.begin",
		TokParenClose:   "punctuation.section.parens.end",
		TokAt:           "punctuation.definition.variable.compiler",
		TokDollar:       "punctuation.definition.variable.parameter",
		TokBracketOpen:  "punctuation.section.brackets.begin",
		TokBracketClose: "punctuation.section.brackets.end",
		TokArrow:        "keyword.operator.arrow.conversion",
		TokColon:        "punctuation.separator.key-value",

		// ------------------------------------------------
		// Literals & Domain Constants
		// ------------------------------------------------
		TokInteger:          "constant.numeric.integer",
		TokFloat:            "constant.numeric.float",
		TokVersionIndicator: "constant.numeric.version",
		TokStringLiteral:    "string.quoted.double",
		TokPercentage:       "constant.numeric.percentage",

		TokKWTrue:  "constant.language.boolean.true",
		TokKWFalse: "constant.language.boolean.false",

		TokKWShow: "support.constant.action.show",
		TokKWHide: "support.constant.action.hide",

		// RGB channel literals
		TokKWRGB:    "support.function.color.rgb",
		TokKWRGBA:   "support.function.color.rgba",
		TokKWHSV:    "support.function.color.hsv",
		TokKWHSVA:   "support.function.color.hsva",
		TokKWHSL:    "support.function.color.hsl",
		TokKWHSLA:   "support.function.color.hsla",
		TokHexColor: "constant.other.color.rgb.hex",

		TokRGBRed:     "constant.numeric.color.channel.red",
		TokRGBGreen:   "constant.numeric.color.channel.green",
		TokRGBBlue:    "constant.numeric.color.channel.blue",
		TokRGBAlpha:   "constant.numeric.color.channel.alpha",
		TokHue:        "constant.numeric.color.channel.hue",
		TokSaturation: "constant.numeric.color.channel.saturation",
		TokValue:      "constant.numeric.color.channel.value",
		TokLightness:  "constant.numeric.color.channel.lightness",

		// ------------------------------------------------
		// Comparison / logical operators
		// ------------------------------------------------
		TokOpEq:  "keyword.operator.comparison.eq",
		TokOpNeq: "keyword.operator.comparison.neq",
		TokOpGt:  "keyword.operator.comparison.gt",
		TokOpGte: "keyword.operator.comparison.gte",
		TokOpLt:  "keyword.operator.comparison.lt",
		TokOpLte: "keyword.operator.comparison.lte",

		// ------------------------------------------------
		// Fallback identifier
		// ------------------------------------------------
		TokIdentifier: "variable.other.readwrite",
	},
	NodeBindings: map[Node]toolchain.NodeBinding[Token]{

		// ------------------------------------------------
		// Module system
		// ------------------------------------------------
		"NodeNamespaceSegment":         {Scopes: []string{"entity.name.namespace"}},
		"NodeModuleSegment":            {Scopes: []string{"entity.name.module"}},
		"NodeEntityReference":          {Scopes: []string{"support.class.reference.configuration"}},
		"NodeRulesetReferenceTerminal": {Scopes: []string{"support.class.ruleset"}},
		"NodeModuleReferenceSegment":   {Scopes: []string{"support.other.module"}},
		NodeModuleAlias:                {Scopes: []string{"entity.name.module.alias"}},

		// ------------------------------------------------
		// Theme / schema types (The Entities)
		// ------------------------------------------------
		NodeThemeSchemaName: {Scopes: []string{"entity.name.class.schema"}},
		NodeThemeName:       {Scopes: []string{"entity.name.class.theme"}},
		NodeColorSetName:    {Scopes: []string{"entity.name.class.color-set"}},

		NodeThemeReference:       {Scopes: []string{"entity.other.inherited-class.theme"}},
		NodeThemeSchemaReference: {Scopes: []string{"entity.other.inherited-class.schema"}},

		// ------------------------------------------------
		// Theme style blocks & Property Keys
		// ------------------------------------------------
		NodeThemeStyleName: {Scopes: []string{"entity.name.class.theme-style"}},
		NodeStyleKey:       {Scopes: []string{"meta.property.name", "support.type.property-name"}},
		NodeColorKey:       {Scopes: []string{"meta.property.name", "support.type.property-name"}},

		// Color References (e.g., Palette.DEFAULT)
		NodeColorSetReference: {Scopes: []string{"support.class.color-set"}},
		NodeColorKeyReference: {Scopes: []string{"variable.other.constant.property"}},

		// ------------------------------------------------
		// Numeric property values
		// ------------------------------------------------
		NodeTextSize:    {Scopes: []string{"constant.numeric.integer"}},
		NodeMinimapSize: {Scopes: []string{"constant.numeric.integer"}},
		NodeSoundID:     {Scopes: []string{"constant.numeric.integer"}},
		NodeSoundVolume: {Scopes: []string{"constant.numeric.integer"}},

		// ------------------------------------------------
		// String values
		// ------------------------------------------------
		NodeMinimapColor: {Scopes: []string{"string.quoted"}},
		NodeMinimapShape: {Scopes: []string{"string.quoted"}},
		NodeSoundPath:    {Scopes: []string{"string.quoted"}},
		NodeEffectColor:  {Scopes: []string{"string.quoted"}},

		// ------------------------------------------------
		// Boolean values
		// ------------------------------------------------
		NodeSoundEnabled:    {Scopes: []string{"constant.language.boolean"}},
		NodeSoundPositional: {Scopes: []string{"constant.language.boolean"}},
		NodeEffectTemporary: {Scopes: []string{"constant.language.boolean"}},

		// ------------------------------------------------
		// Color literals
		// ------------------------------------------------
		NodeHexColor: {Scopes: []string{"meta.color.hex", "constant.other.color.rgb.hex"}},

		NodeRGBRed:     {Scopes: []string{"constant.numeric.color.channel.red"}},
		NodeRGBGreen:   {Scopes: []string{"constant.numeric.color.channel.green"}},
		NodeRGBBlue:    {Scopes: []string{"constant.numeric.color.channel.blue"}},
		NodeRGBAlpha:   {Scopes: []string{"constant.numeric.color.channel.alpha"}},
		NodeHue:        {Scopes: []string{"constant.numeric.color.channel.hue"}},
		NodeSaturation: {Scopes: []string{"constant.numeric.color.channel.saturation"}},
		NodeValue:      {Scopes: []string{"constant.numeric.color.channel.value"}},
		NodeLightness:  {Scopes: []string{"constant.numeric.color.channel.lightness"}},

		// ------------------------------------------------
		// Actions / styles (Property Wrappers)
		// ------------------------------------------------
		NodeActionDef: {Scopes: []string{"meta.property.action"}},
		NodeStyleDef:  {Scopes: []string{"meta.property.style"}},

		// ------------------------------------------------
		// Templates & Functions
		// ------------------------------------------------
		NodeTemplateName:       {Scopes: []string{"entity.name.function.template"}},
		NodeFunctionArgument:   {Scopes: []string{"variable.parameter.function"}},
		NodeFunctionName:       {Scopes: []string{"entity.name.function"}},
		NodeRulesetOverrideKey: {Scopes: []string{"support.type.property-name.override"}},

		// ------------------------------------------------
		// Variables
		// ------------------------------------------------
		NodeVariableName:      {Scopes: []string{"variable.other.readwrite"}},
		NodeVariableRef:       {Scopes: []string{"variable.other.readwrite"}},
		NodeAssignTarget:      {Scopes: []string{"variable.other.assignment"}},
		NodeSchemaUsageTarget: {Scopes: []string{"support.class.schema"}},

		// Arguments
		NodeArgumentName: {Scopes: []string{"variable.parameter"}},
		NodeArgumentType: {Scopes: []string{"support.type"}},

		// ------------------------------------------------
		// Conditions / expressions
		// ------------------------------------------------
		NodeConditionExpr:    {Scopes: []string{"meta.expression.condition"}},
		NodeConditionTarget:  {Scopes: []string{"variable.other.property.game"}},
		NodeConditionValue:   {Scopes: []string{"constant.language.condition"}},
		NodeVariantCondition: {Scopes: []string{"meta.condition.variant"}},

		// Engine Directives (@strictness)
		NodeEngineDirective: {Scopes: []string{"variable.language.compiler.directive"}},

		// ------------------------------------------------
		// Structural Block Metas
		// ------------------------------------------------
		NodeStringArray: {MetaScope: "meta.sequence.list.string", TokenScopes: map[Token][]string{
			TokBracketOpen: {"punctuation.section.brackets.begin"},
		}},

		// ------------------------------------------------
		// Inheritance
		// ------------------------------------------------
		NodeInheritTarget: {Scopes: []string{"support.type.property-name"}},

		// ------------------------------------------------
		// Operations
		// ------------------------------------------------
		NodeOp: {Scopes: []string{"keyword.operator"}},

		// ------------------------------------------------
		// ColorSet entries
		// ------------------------------------------------
		NodeColorSetEntry: {Scopes: []string{"meta.mapping.entry.colorset"}},

		// ------------------------------------------------
		// Ruleset Architecture
		// ------------------------------------------------
		NodeRulesetName: {Scopes: []string{"entity.name.class.ruleset"}},
		NodeRulesetDescription: {
			Scopes: []string{"meta.property.description", "string.quoted.double"},
		},

		// ------------------------------------------------
		// Rules & Template Invocations
		// ------------------------------------------------
		NodeRule:     {Scopes: []string{"meta.function.rule"}},
		NodeRuleName: {Scopes: []string{"entity.name.function.rule"}},
		NodeRuleBlock: {MetaScope: "meta.block.rule", TokenScopes: map[Token][]string{
			TokBraceOpen: {"punctuation.section.brace.begin"},
		}},

		NodeTemplateReferenceTarget: {Scopes: []string{"meta.function-call entity.name.function.template"}},
		NodeStyleKeyRef:             {Scopes: []string{"variable.other.constant.property"}},

		// ------------------------------------------------
		// Configuration & Adapters
		// ------------------------------------------------
		NodeConfigurationKey: {Scopes: []string{"support.type.property-name.configuration"}},

		// References to external blueprints (e.g., FancySchema, core.BaseTiering)
		NodeOriginalSchemaReference: {Scopes: []string{"support.class.schema.original"}},
		NodeTargetSchemaReference:   {Scopes: []string{"support.class.schema.target"}},

		// Mapping specific properties (e.g., fancyMythic = godTier)
		NodeOriginalSchemaKey: {Scopes: []string{"variable.other.constant.property.original"}},
		NodeTargetSchemaKey:   {Scopes: []string{"variable.other.constant.property.target"}},

		// ------------------------------------------------
		// Literal / value abstraction
		// ------------------------------------------------
		NodeAssignmentValue: {MetaScope: "meta.value", TokenScopes: map[Token][]string{
			TokStringLiteral: {"string.quoted.double"},
		}},
	},
}

func ruleforgeOverrideProducer(
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

	registry.Register(string(TokBlockComment), patterns.BlockComment(
		"/*",
		"*/",
		toolchain.SublimeContext{MetaScope: "comment.block"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment.begin"},
		toolchain.SublimeContext{Scope: "punctuation.definition.comment.end"},
	))

	registry.Register(string(TokIdentifier), identifierOverride(
		ruleset,
		ctxProducer,
	))

	return registry.Producer()
}

func identifierOverride(
	ruleset *lexarch.LexingRuleset[rune, string, string],
	ctxProducer func(ctx *EditorCtx) toolchain.SublimeContext,
) editor.OverrideHandler[rune, string, string, string, string, toolchain.SublimeContext] {

	identRegex := getTokenRegex(ruleset, string(TokIdentifier))
	dotRegex := getTokenRegex(ruleset, string(TokDot))

	// The mechanical lookaheads
	prefixRegex := fmt.Sprintf(`%s(?=\s*%s)`, identRegex, dotRegex)
	terminalRegex := fmt.Sprintf(`%s(?!\s*%s)`, identRegex, dotRegex)

	// 1. Contexts for MODULE_PATH (Declarations & Imports: e.g., 'import core.test')
	nsNode := string(Node("NodeNamespaceSegment"))
	modNode := string(Node("NodeModuleSegment"))
	pathPrefixCtx := ctxProducer(&EditorCtx{NodeKind: &nsNode})
	pathTerminalCtx := ctxProducer(&EditorCtx{NodeKind: &modNode})

	// 2. Contexts for ENTITY_REF_WRAP (Pointers: e.g., 'core.BaseTiering')
	modRefNode := "NodeModuleReferenceSegment" // Synthetic node for pointers
	entityRefNode := string(Node("NodeEntityReference"))
	refPrefixCtx := ctxProducer(&EditorCtx{NodeKind: &modRefNode})
	refTerminalCtx := ctxProducer(&EditorCtx{NodeKind: &entityRefNode})

	// 3. Contexts for RULESET_REFERENCE (Pointers: e.g., 'include core.Currency')
	rulesetRefTerminalNode := "NodeRulesetReferenceTerminal"
	rulesetPrefixCtx := ctxProducer(&EditorCtx{NodeKind: &modRefNode}) // Uses the same pointer node
	rulesetTerminalCtx := ctxProducer(&EditorCtx{NodeKind: &rulesetRefTerminalNode})

	// Target AST groups we are overriding
	pathSeg := string(NodePathSegment)
	refSeg := string(NodeReferenceSegment)
	rulesetSeg := string(NodeRulesetReferenceSegment)

	return func(ctx *EditorCtx) []*EditorOverride {
		if ctx.NodeKind == nil {
			return nil
		}

		switch *ctx.NodeKind {
		case pathSeg:
			return []*EditorOverride{
				{PatternRegex: &prefixRegex, MatchContext: &pathPrefixCtx},
				{PatternRegex: &terminalRegex, MatchContext: &pathTerminalCtx},
			}
		case refSeg:
			return []*EditorOverride{
				{PatternRegex: &prefixRegex, MatchContext: &refPrefixCtx},
				{PatternRegex: &terminalRegex, MatchContext: &refTerminalCtx},
			}
		case rulesetSeg:
			return []*EditorOverride{
				{PatternRegex: &prefixRegex, MatchContext: &rulesetPrefixCtx},
				{PatternRegex: &terminalRegex, MatchContext: &rulesetTerminalCtx},
			}
		default:
			return nil
		}
	}
}

func getTokenRegex(
	ruleset *lexarch.LexingRuleset[rune, string, string],
	token string,
) string {
	for _, rule := range ruleset.GetRules() {
		if rule.Token == token {
			regex, err := rule.Pattern.ToRegEx()
			if err != nil {
				panic(fmt.Sprintf("Failed to convert token %v to regex: %v", token, err))
			}
			return regex
		}
	}
	panic(fmt.Sprintf("Token %v not found in lexing ruleset", token))
}
