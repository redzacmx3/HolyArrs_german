package importer

import (
	"strings"
	"testing"

	"github.com/HolyArrs/holyarrs/app/backend/internal/core/releasename"
	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

// A show on a chosen episode order names its files by that order's numbers.
func TestEpisodePathByOrder(t *testing.T) {
	ns := mediafiles.NamingDefaults()
	ns.RenameOnImport = true
	eps := []Episode{{Ref: mediafiles.EpisodeRef{Season: 3, Episode: 4, AltSeason: 3, AltEpisode: 3}, Name: "Episode 4"}}
	name := "Trapped.S03E03.1080p.WEB.H264-GRP"
	p := releasename.Parse(name, releasename.Hints{Show: true})
	tgt := ShowTarget{FolderName: "Trapped (2015)", Show: mediafiles.TitleRef{Name: "Trapped", Year: 2015}, Order: mediafiles.OrderAired, Naming: &ns}
	if got := tgt.episodePath(p, name+".mkv", eps); !strings.Contains(got, "S03E04") {
		t.Fatalf("TMDB numbers: %s", got)
	}
	tgt.ByOrder = true
	if got := tgt.episodePath(p, name+".mkv", eps); !strings.Contains(got, "S03E03") {
		t.Fatalf("order numbers: %s", got)
	}
	if eps[0].Ref.Episode != 4 {
		t.Fatalf("the episode list is left as it was: %+v", eps)
	}
}

// A file for an episode the order leaves out is refused when its TMDB
// numbers would give it the name of another episode in the order.
func TestPlanRefusesOrderNameClash(t *testing.T) {
	ns := mediafiles.NamingDefaults()
	ns.RenameOnImport = true
	eps := []Episode{
		{Ref: mediafiles.EpisodeRef{ItemID: 7, Season: 3, Episode: 7, AirDay: "2019-12-07"}, Name: "Episode 7", SlotID: 17, Wanted: true},
		{Ref: mediafiles.EpisodeRef{ItemID: 8, Season: 3, Episode: 8, AltSeason: 3, AltEpisode: 7, AirDay: "2019-12-08"}, Name: "Episode 8", SlotID: 18, Wanted: true},
		{Ref: mediafiles.EpisodeRef{ItemID: 9, Season: 3, Episode: 9, AirDay: "2019-12-09"}, Name: "Episode 9", SlotID: 19, Wanted: true},
	}
	tgt := ShowTarget{LibraryPath: "/lib", FolderName: "Trapped (2015)", Show: mediafiles.TitleRef{Name: "Trapped", Year: 2015}, Order: mediafiles.OrderAired,
		Naming: &ns, Episodes: eps, ByOrder: true}
	steps := PlanEpisodes([]Candidate{{Path: "/dl/Trapped.2019.12.07.1080p.WEB.H264-GRP.mkv", Size: 1}, {Path: "/dl/Trapped.2019.12.09.1080p.WEB.H264-GRP.mkv", Size: 1}}, tgt)
	if steps[0].Method != "refuse" || !strings.Contains(strings.Join(steps[0].Why, "; "), "not in the chosen episode order") {
		t.Fatalf("clash: %+v", steps[0])
	}
	if steps[1].Method == "refuse" || !strings.Contains(steps[1].TargetPath, "S03E09") {
		t.Fatalf("no clash: %+v", steps[1])
	}
}

// With release names read by TMDB's numbers (ReleaseNames tmdb), a file named
// S02E46 fills the episode TMDB numbers S02E46 and is named by the chosen
// order (S03E47), not the episode the order calls S02E46.
func TestPlanTMDBNames(t *testing.T) {
	ns := mediafiles.NamingDefaults()
	ns.RenameOnImport = true
	eps := []Episode{
		{Ref: mediafiles.EpisodeRef{ItemID: 1, Season: 2, Episode: 46, AltSeason: 3, AltEpisode: 47}, Name: "Pups Save a Hoot", SlotID: 11, Wanted: true},
		{Ref: mediafiles.EpisodeRef{ItemID: 2, Season: 2, Episode: 45, AltSeason: 2, AltEpisode: 46}, Name: "Pups Save the Bunnies", SlotID: 12, Wanted: true},
	}
	for _, tmdb := range []bool{true, false} {
		names := ""
		if tmdb {
			names = mediafiles.ReleaseNumbersTMDB
		}
		tgt := ShowTarget{LibraryPath: "/lib", FolderName: "PAW Patrol (2013)", Show: mediafiles.TitleRef{Name: "PAW Patrol", Year: 2013}, Order: mediafiles.OrderAired,
			Naming: &ns, Episodes: eps, ByOrder: true, ReleaseNames: names}
		steps := PlanEpisodes([]Candidate{{Path: "/dl/PAW.Patrol.S02E46.German.AC3.DL.1080p.WebHD.x265-FuN.mkv", Size: 1}}, tgt)
		want := "S02E46"
		if tmdb {
			want = "S03E47"
		}
		if steps[0].Method == "refuse" || steps[0].Method == "skip" || !strings.Contains(steps[0].TargetPath, want) {
			t.Fatalf("tmdb %v: want %s, got %+v", tmdb, want, steps[0])
		}
	}
}

// With release names read by TVmaze's numbers (ReleaseNames tvmaze), a
// FuN double episode S02E40-E41 fills the two episodes TVmaze numbers so
// and is named by the chosen order; an episode without TVmaze numbers is
// not matched by its other numbers.
func TestPlanTVmazeNames(t *testing.T) {
	ns := mediafiles.NamingDefaults()
	ns.RenameOnImport = true
	eps := []Episode{
		{Ref: mediafiles.EpisodeRef{ItemID: 1, Season: 2, Episode: 46, AltSeason: 3, AltEpisode: 47, NameSeason: 2, NameEpisode: 41}, Name: "Pups Save an Eagle", SlotID: 11, Wanted: true},
		{Ref: mediafiles.EpisodeRef{ItemID: 2, Season: 2, Episode: 45, AltSeason: 3, AltEpisode: 46, NameSeason: 2, NameEpisode: 40}, Name: "Pups Save the Woof and Roll Show", SlotID: 12, Wanted: true},
		{Ref: mediafiles.EpisodeRef{ItemID: 3, Season: 2, Episode: 41, AltSeason: 2, AltEpisode: 41}, Name: "Not on TVmaze", SlotID: 13, Wanted: true},
	}
	tgt := ShowTarget{LibraryPath: "/lib", FolderName: "PAW Patrol (2013)", Show: mediafiles.TitleRef{Name: "PAW Patrol", Year: 2013}, Order: mediafiles.OrderAired,
		Naming: &ns, Episodes: eps, ByOrder: true, ReleaseNames: mediafiles.ReleaseNumbersTVmaze}
	steps := PlanEpisodes([]Candidate{{Path: "/dl/PAW.Patrol.S02E40-E41.German.AC3.DL.1080p.WebHD.x265-FuN.mkv", Size: 1}}, tgt)
	if steps[0].Method == "refuse" || steps[0].Method == "skip" || len(steps[0].TargetSlotIDs) != 2 || !strings.Contains(steps[0].TargetPath, "S03E46") {
		t.Fatalf("want S03E46-E47 for slots 12 and 11, got %+v", steps[0])
	}
	for _, id := range steps[0].TargetSlotIDs {
		if id == 13 {
			t.Fatalf("the episode without TVmaze numbers was matched: %+v", steps[0])
		}
	}
}
