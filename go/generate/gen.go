//go:build ignore

package main

import (
	goDef "lingua/go"
	"lingua/internal/generation"
)

func main() {
	generation.ExecuteFull(goDef.Config())
}
