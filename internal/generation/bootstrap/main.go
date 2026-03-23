package main

import (
	"flag"
	"fmt"
	"langspec/bootstrap"
	"langspec/cliutil"
	"os"
)

func main() {
	specPath := flag.String("spec", "", "path to the .lspec file")
	flag.Parse()

	if *specPath == "" {
		fmt.Fprintln(os.Stderr, "bootstrap: -spec flag required")
		os.Exit(1)
	}

	g := cliutil.NewGenerator(*specPath)
	defer g.Close()

	// Only run the bindings toolchain
	opts := []bootstrap.Option{
		bootstrap.WithToolchainFilter("go_bindings"),
	}

	if err := g.RunToolchains(opts...); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap failed: %v\n", err)
		os.Exit(1)
	}
}
