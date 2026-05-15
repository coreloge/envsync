package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/report"
)

func runCompare(args []string) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	sourceEnv := fs.String("source-env", "source", "label for the source environment")
	targetEnv := fs.String("target-env", "target", "label for the target environment")
	format := fs.String("format", "text", "output format: text or json")
	ignoreValues := fs.Bool("ignore-values", false, "only compare key presence, not values")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) < 2 {
		return fmt.Errorf("usage: envsync compare [flags] <source.env> <target.env>")
	}

	source, err := envfile.Parse(positional[0])
	if err != nil {
		return fmt.Errorf("parsing source: %w", err)
	}

	target, err := envfile.Parse(positional[1])
	if err != nil {
		return fmt.Errorf("parsing target: %w", err)
	}

	opts := envfile.CompareOptions{
		IgnoreValues: *ignoreValues,
	}

	result := envfile.Compare(*sourceEnv, *targetEnv, source, target, opts)

	if err := report.WriteCompare(os.Stdout, result, *format); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	if result.HasDrift() {
		os.Exit(1)
	}
	return nil
}
