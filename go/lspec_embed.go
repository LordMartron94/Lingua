package goDef

import _ "embed"

//go:embed gomod.lspec
var embeddedGoModLspec []byte

//go:embed gowork.lspec
var embeddedGoWorkLspec []byte
