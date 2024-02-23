package main

import (
	"fmt"
	"log"
	"os"
	"path"

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

	defs, err := wut.Load(*flagDefinitionsFile, *flagMaxDistance)
	if err != nil {
		log.Fatalf("Failed to load acronyms: %v", err)
	}

	def, closeMatches := defs.Get(acronym)
	if def != "" {
		fmt.Println(def)
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
