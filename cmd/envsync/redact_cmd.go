package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/envfile"
	"github.com/yourorg/envsync/internal/report"
)

// runRedact executes the redact subcommand.
// Usage: envsync redact --file <path> [--format text|json] [--output <path>]
func runRedact(args []string) {
	fs := flag.NewFlagSet("redact", flag.ExitOnError)
	filePath := fs.String("file", "", "path to the .env file to redact")
	format := fs.String("format", "text", "output format: text or json")
	outPath := fs.String("output", "", "write redacted file to this path (optional)")
	_ = fs.Parse(args)

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "error: --file is required")
		fs.Usage()
		os.Exit(1)
	}

	vars, err := envfile.Parse(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing file: %v\n", err)
		os.Exit(1)
	}

	opts := envfile.DefaultRedactOptions()
	redacted := envfile.Redact(vars, opts)

	summary := report.BuildRedactSummary(vars, redacted)
	if err := report.WriteRedactSummary(os.Stdout, summary, *format); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if *outPath != "" {
		if err := envfile.Write(*outPath, redacted); err != nil {
			fmt.Fprintf(os.Stderr, "error writing redacted file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "\nRedacted file written to: %s\n", *outPath)
	}
}
