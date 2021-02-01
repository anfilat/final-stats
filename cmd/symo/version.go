package main

import (
	"flag"
	"fmt"
)

var (
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func printVersion() {
	fmt.Printf("System monitor built on %s git %s\n", buildDate, gitHash)
}

func isVersionCommand() bool {
	for _, name := range flag.Args() {
		if name == "version" {
			return true
		}
	}
	return false
}
