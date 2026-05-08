package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/report"
)

func main() {
	var (
		sourceFile = flag.String("source", "", "path to source .env file (reference)")
		targetFile = flag.String("target", "", "path to target .env file")
		format     = flag.String("format", "text", "output format: text or json")
	)
	flag.Parse()

	if *sourceFile == "" || *targetFile == "" {
		fmt.Fprintln(os.Stderr, "error: --source and --target are required")
		flag.Usage()
		os.Exit(1)
	}

	source, err := envfile.Parse(*sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source file: %v\n", err)
		os.Exit(1)
	}

	target, err := envfile.Parse(*targetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target file: %v\n", err)
		os.Exit(1)
	}

	diffResult := envfile.Diff(source, target)

	fmt := report.Format(*format)
	if fmt != report.FormatText && fmt != report.FormatJSON {
		fmt.Fprintf(os.Stderr, "error: unknown format %q, use 'text' or 'json'\n", *format)
		os.Exit(1)
	}

	w := report.NewWriter(os.Stdout, fmt)
	if err := w.Write(*sourceFile, *targetFile, diffResult); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if diffResult.HasDrift() {
		os.Exit(2)
	}
}
