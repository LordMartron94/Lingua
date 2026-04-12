package main

import (
	"flag"
	"fmt"
	"langspec/bootstrap"
	"langspec/cliutil"
	"langspec/dsl"
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
	opts := []bootstrap.Option[dsl.LangSpecParserNodeKind]{
		bootstrap.WithToolchainFilter[dsl.LangSpecParserNodeKind]("go_bindings"),
		bootstrap.WithDiagnosticSink[dsl.LangSpecParserNodeKind](dsl.DefaultLangSpecDiagnosticSink()),
	}

	if err := g.RunToolchains(opts...); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap failed: %v\n", err)
		os.Exit(1)
	}
}
