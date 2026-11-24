package mediafiles

import (
	"fmt"
	"strings"
)

// Fork: a size per minute trait tests a Release's size against the runtime
// it covers (an episode, a season pack's episodes, a movie), written
// "over-upto" in MB a minute like a size test in GB. A 12-minute kids'
// episode of 377 MB (31 MB a minute) then reads like a 45-minute one of
// 1.4 GB, where a fixed size or a group list calls it too small.

// MatchSizePerMinute is the own trait match for a size per minute test.
const MatchSizePerMinute = "size_per_minute"

const bytesPerMB = 1 << 20

// perMinuteRange reads a size per minute trait's one term; bad is true when
// it cannot be read (CheckOwnTrait refuses those on save).
func perMinuteRange(terms []string) (*sizeRange, bool) {
	if len(terms) != 1 {
		return &sizeRange{}, true
	}
	over, upTo, ok := ParseSizeRange(terms[0])
	if !ok {
		return &sizeRange{}, true
	}
	return &sizeRange{over: int64(over * bytesPerMB), upTo: int64(upTo * bytesPerMB)}, false
}

// checkPerMinuteTerms cleans a size per minute trait's terms before save.
func checkPerMinuteTerms(raw []string) ([]string, error) {
	var terms []string
	for _, r := range raw {
		if s := strings.TrimSpace(r); s != "" {
			terms = append(terms, s)
		}
	}
	if len(terms) != 1 {
		return nil, fieldErr("terms", "Give one range in MB a minute, such as -10 or 10-40.")
	}
	over, upTo, ok := ParseSizeRange(terms[0])
	if !ok {
		return nil, fieldErr("terms", "Give an over or an up to size in MB a minute, and up to above over.")
	}
	return []string{SizeRangeValue(over, upTo)}, nil
}

// perMinuteMatches tests size / runtime against the range. Without a size
// or a runtime the trait does not fit, so it never takes points on a guess.
func perMinuteMatches(name string, r sizeRange, size int64, runtime int) (bool, string) {
	if size <= 0 {
		return false, name + ": the size is not known"
	}
	if runtime <= 0 {
		return false, name + ": the runtime is not known"
	}
	per := size / int64(runtime)
	words := fmt.Sprintf("%s: %.1f MB a minute (%s for %d min)", name, float64(per)/bytesPerMB, sizeWord(size), runtime)
	if r.fits(per) {
		return true, words
	}
	return false, words + " is outside it"
}

// hasPerMinute says whether the rules score a size per minute trait, so the
// file in place is scored again with the runtime.
func hasPerMinute(rules PickRulesBody) bool {
	all := ownTraits.Load()
	if all == nil {
		return false
	}
	for _, tp := range rules.Traits {
		if id, ok := OwnTraitID(tp.Trait); ok && (*all)[id].perMin != nil {
			return true
		}
	}
	return false
}
