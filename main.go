package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--verbose] CONFIG_FILE\n", os.Args[0])
		os.Exit(1)
	}

	configFile := flag.Arg(0)

	cfg, err := LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
		os.Exit(1)
	}

	// Execute comparison
	differ := NewDiffer(cfg, *verbose)
	success := differ.RunComparison()

	if !success {
		os.Exit(1)
	}
}
