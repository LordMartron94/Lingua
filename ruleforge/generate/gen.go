//go:build ignore

package main

import (
	"lingua/internal/generation"
	"lingua/ruleforge"
)

func main() {
	generation.ExecuteFull(ruleforge.Config())
}
