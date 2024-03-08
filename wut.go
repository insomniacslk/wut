package wut

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Defs map[string][]string

type Wut struct {
	// TODO: support multiple definitions
	defs Defs
	// used to do fuzzy matching
	acronyms []string
	// max Levenshtein distance for fuzzy matching
	maxDistance uint
}

// Get returns an acronym's definition. If no exact match is found, it returns an empty
// string, and a list of top results using fuzzy matching. This list can be empty if no
// match is close enough.
func (w *Wut) Get(acronym string) ([]string, []string) {
	def, ok := w.defs[acronym]
	if ok {
		return def, nil
	}
	if w.maxDistance == 0 {
		// only exact-matching is requested
		return nil, nil
	}
	// try fuzzy matching, this is slower because we iterate all the keys, plus
	// the fuzzy matching itself.
	matches := fuzzy.RankFindNormalized(acronym, w.acronyms)
	sort.Sort(matches)
	closeMatches := make([]string, 0)
	for _, m := range matches {
		if uint(m.Distance) <= w.maxDistance {
			closeMatches = append(closeMatches, m.Target)
		}
	}
	return nil, closeMatches
}

func Load(filename string, maxDistance uint) (*Wut, error) {
	data, err := os.ReadFile(filename)
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
		// normalize the case so the match can be case-insensitive.
		lk := strings.ToLower(k)
		for _, def := range v {
			defs[lk] = append(defs[lk], strings.TrimSpace(def))
		}
		// get all the keys in a slice to do fuzzy matching if necessary.
		keys = append(keys, lk)
	}
	return &Wut{defs: defs, acronyms: keys, maxDistance: maxDistance}, nil
}
