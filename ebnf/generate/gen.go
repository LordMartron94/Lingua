//go:build ignore

package main

import (
	"lingua/ebnf"
	"lingua/internal/generation"
)

func main() {
	generation.ExecuteFull(ebnf.Config())
}
