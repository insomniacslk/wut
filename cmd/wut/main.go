package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/insomniacslk/wut"
	"github.com/kirsle/configdir"
	"github.com/spf13/pflag"
)

const progname = "wut"

var defaultAcronymsFile = path.Join(configdir.LocalConfig(progname), "acronyms.json")

var (
	flagDefinitionsFile = pflag.StringP("definitions-file", "f", defaultAcronymsFile, "JSON file containing the acronym definitions")
	flagMaxDistance     = pflag.UintP("max-distance", "d", 1, "Maximum Levenshtein distance for fuzzy matching when exact matching fails. 0 means exact match")
)

func main() {
	pflag.Parse()
	acronym := pflag.Arg(0)
	if acronym == "" {
		log.Fatalf("No acronym specified")
	}

	allDefs, err := wut.Load(*flagDefinitionsFile, *flagMaxDistance)
	if err != nil {
		log.Fatalf("Failed to load acronyms: %v", err)
	}

	defs, closeMatches := allDefs.Get(acronym)
	if defs != nil {
		separator := strings.Repeat("-", 80)
		output := strings.Join(defs, separator)
		fmt.Println(output)
		return
	}
	if len(closeMatches) > 0 {
		fmt.Printf("Acronym '%s' not found, did you mean one of the following?\n", acronym)
		for _, m := range closeMatches {
			fmt.Printf("  %s\n", m)
		}
	} else {
		fmt.Printf("No match found for '%s'.\n", acronym)
		os.Exit(1)
	}
}
