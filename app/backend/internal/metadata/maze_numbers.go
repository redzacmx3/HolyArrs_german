package metadata

import (
	"sort"
	"strings"
	"unicode"
)

// SetMazeNumbers gives each regular episode of eps the season and number
// TVmaze has for the same episode as NameSeason and NameEpisode (TVmaze
// numbers like TVDB, and German releases of some kids' shows count that
// way), and returns how many it numbered. maze must be the same show (see
// SameShow). An episode is the same when, in this order:
//   - its name is the name of exactly one TVmaze episode, and no other
//     episode of eps has that name;
//   - it aired on a day TVmaze lists episodes for, and one of them has the
//     most words of its name in common (at least half);
//   - the episodes still open on that day are as many as TVmaze's still
//     open ones: they pair up in their own order.
//
// Each TVmaze episode numbers one episode at most; an episode found no
// other way keeps no TVmaze numbers.
func SetMazeNumbers(eps []Episode, maze Show) int {
	var regular []Episode
	for _, m := range maze.Episodes {
		if m.Season > 0 && m.Number > 0 {
			regular = append(regular, m)
		}
	}
	for i := range eps {
		eps[i].NameSeason, eps[i].NameEpisode = 0, 0
	}
	used := make([]bool, len(regular))
	n := 0
	take := func(i, j int) {
		eps[i].NameSeason, eps[i].NameEpisode = regular[j].Season, regular[j].Number
		used[j] = true
		n++
	}
	open := func(i int) bool { return eps[i].Season > 0 && eps[i].NameSeason == 0 }

	// 1. One name, one episode on either side.
	ours, theirs := map[string][]int{}, map[string][]int{}
	for i, e := range eps {
		if k := nameKey(e.Name); k != "" && e.Season > 0 {
			ours[k] = append(ours[k], i)
		}
	}
	for j, m := range regular {
		if k := nameKey(m.Name); k != "" {
			theirs[k] = append(theirs[k], j)
		}
	}
	for k, is := range ours {
		if js := theirs[k]; len(is) == 1 && len(js) == 1 {
			take(is[0], js[0])
		}
	}

	// 2. The same day, the most words in common.
	byDay := map[string][]int{}
	for j, m := range regular {
		if m.AirDay != "" {
			byDay[m.AirDay] = append(byDay[m.AirDay], j)
		}
	}
	for i := range eps {
		if !open(i) || eps[i].AirDay == "" {
			continue
		}
		best, score, tie := -1, 0.0, false
		for _, j := range byDay[eps[i].AirDay] {
			if used[j] {
				continue
			}
			s := wordsInCommon(eps[i].Name, regular[j].Name)
			switch {
			case s > score:
				best, score, tie = j, s, false
			case s == score:
				tie = true
			}
		}
		if best >= 0 && score >= 0.5 && !tie {
			take(i, best)
		}
	}

	// 3. What is left of a day pairs up in order.
	ourDays := map[string][]int{}
	for i := range eps {
		if open(i) && eps[i].AirDay != "" {
			ourDays[eps[i].AirDay] = append(ourDays[eps[i].AirDay], i)
		}
	}
	for day, is := range ourDays {
		var js []int
		for _, j := range byDay[day] {
			if !used[j] {
				js = append(js, j)
			}
		}
		if len(js) != len(is) {
			continue
		}
		sort.Slice(is, func(a, b int) bool { return before(eps[is[a]], eps[is[b]]) })
		sort.Slice(js, func(a, b int) bool { return before(regular[js[a]], regular[js[b]]) })
		for k := range is {
			take(is[k], js[k])
		}
	}
	return n
}

func before(a, b Episode) bool {
	if a.Season != b.Season {
		return a.Season < b.Season
	}
	return a.Number < b.Number
}

// nameKey is an episode name reduced to its lower-case letters and digits.
func nameKey(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// wordsInCommon is the share of the words of the longer name that both
// names hold, from 0 to 1.
func wordsInCommon(a, b string) float64 {
	split := func(s string) map[string]bool {
		out := map[string]bool{}
		for _, w := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
			out[w] = true
		}
		return out
	}
	wa, wb := split(a), split(b)
	if len(wa) == 0 || len(wb) == 0 {
		return 0
	}
	same := 0
	for w := range wa {
		if wb[w] {
			same++
		}
	}
	return float64(same) / float64(max(len(wa), len(wb)))
}
