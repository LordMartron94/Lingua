//go:build ignore

package main

import (
	"lingua/helm"
	"lingua/internal/generation"
)

func main() {
	generation.ExecuteFull(helm.Config())
}
