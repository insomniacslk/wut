package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/kirsle/configdir"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/pflag"
)

const progname = "wut"

var defaultAcronymsFile = path.Join(configdir.LocalConfig(progname), "acronyms.json")

var (
	flagDefinitionsFile = pflag.StringP("definitions-file", "f", defaultAcronymsFile, "JSON file containing the acronym definitions")
	flagMaxDistance     = pflag.UintP("max-distance", "d", 1, "Maximum Levenshtein distance for fuzzy matching when exact matching fails")
)

type Defs map[string]string

func main() {
	pflag.Parse()
	acronym := pflag.Arg(0)
	if acronym == "" {
		log.Fatalf("No acronym specified")
	}

	data, err := os.ReadFile(*flagDefinitionsFile)
	if err != nil {
		log.Fatalf("Failed to read definitions file: %v", err)
	}

	var (
		tmp  Defs
		defs = make(Defs)
	)
	if err := json.Unmarshal(data, &tmp); err != nil {
		log.Fatalf("Failed to unmarshal JSON file: %v", err)
	}
	keys := make([]string, 0, len(defs))
	for k, v := range tmp {
		lk := strings.ToLower(k)
		// normalize the case so the match can be case-insensitive.
		defs[lk] = v
		// get all the keys in a slice to do fuzzy matching if necessary.
		keys = append(keys, lk)
	}

	def, ok := defs[acronym]
	if ok {
		fmt.Println(def)
		return
	}
	// try fuzzy matching, this is slower because we iterate all the keys, plus
	// the fuzzy matching itself.
	matches := fuzzy.RankFindNormalized(acronym, keys)
	sort.Sort(matches)
	closeMatches := make([]string, 0)
	for _, m := range matches {
		if uint(m.Distance) <= *flagMaxDistance {
			closeMatches = append(closeMatches, m.Target)
		}
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
