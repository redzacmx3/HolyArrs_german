package metadata

import "testing"

// PAW Patrol: TMDB keeps "Pups Save an Eagle" as S02E46 and TVmaze (like
// TVDB and the German FuN releases) as S02E41.
func TestSetMazeNumbers(t *testing.T) {
	eps := []Episode{
		{Season: 2, Number: 46, Name: "Pups Save an Eagle", AirDay: "2015-09-04"},
		{Season: 2, Number: 45, Name: "Pups Save the Woof & Roll Show", AirDay: "2015-09-04"},
		{Season: 2, Number: 30, Name: "Die Fellfreunde retten einen Schnupfen", AirDay: "2015-06-01"},
		{Season: 2, Number: 31, Name: "Ein anderer Titel", AirDay: "2015-06-01"},
		{Season: 2, Number: 50, Name: "Nowhere on TVmaze", AirDay: "2016-01-01"},
		{Season: 0, Number: 1, Name: "A Special", AirDay: "2015-01-01"},
	}
	maze := Show{Episodes: []Episode{
		{Season: 2, Number: 40, Name: "Pups Save the Woof and Roll Show", AirDay: "2015-09-04"},
		{Season: 2, Number: 41, Name: "Pups Save an Eagle", AirDay: "2015-09-04"},
		{Season: 2, Number: 31, Name: "Pups Save a Sniffle", AirDay: "2015-06-01"},
		{Season: 2, Number: 32, Name: "Pups Save a Hoot", AirDay: "2015-06-01"},
		{Season: 0, Number: 1, Name: "A Special", AirDay: "2015-01-01"},
	}}
	if n := SetMazeNumbers(eps, maze); n != 4 {
		t.Fatalf("numbered %d, want 4", n)
	}
	want := [][2]int{{2, 41}, {2, 40}, {2, 31}, {2, 32}, {0, 0}, {0, 0}}
	for i, w := range want {
		if got := [2]int{eps[i].NameSeason, eps[i].NameEpisode}; got != w {
			t.Errorf("%s: got %v, want %v", eps[i].Name, got, w)
		}
	}
}

// Two episodes of the same name are not told apart by name.
func TestSetMazeNumbersSameName(t *testing.T) {
	eps := []Episode{{Season: 1, Number: 1, Name: "Pilot"}, {Season: 2, Number: 1, Name: "Pilot"}}
	maze := Show{Episodes: []Episode{{Season: 1, Number: 1, Name: "Pilot"}}}
	if n := SetMazeNumbers(eps, maze); n != 0 {
		t.Fatalf("numbered %d, want 0", n)
	}
}
