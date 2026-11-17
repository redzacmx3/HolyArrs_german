package library

import "testing"

// A Radarr group such as TRaSH's "Merged QPs" ranks WEB and Bluray, 720p
// and 1080p by score alone, so the imported rules let points lead.
func TestOldProfileMergedGroupsRankByPoints(t *testing.T) {
	q := []string{"WEBDL-720p", "WEBRip-720p", "Bluray-720p", "WEBDL-1080p", "WEBRip-1080p", "Bluray-1080p"}
	if !MergesSteps(q) {
		t.Fatal("WEB and Bluray at 720p and 1080p should span several steps")
	}
	if MergesSteps([]string{"WEBDL-1080p", "WEBRip-1080p"}) {
		t.Fatal("WEBDL and WEBRip 1080p are one step")
	}
	body, _ := OldProfileRules(OldProfile{Qualities: q, CutoffQualities: q, UpgradeAllowed: true, MergedGroups: true, Language: "Any"})
	if body.Order != "points_then_quality" {
		t.Fatalf("merged groups: order %q, want points_then_quality", body.Order)
	}
	body, _ = OldProfileRules(OldProfile{Qualities: q, CutoffQualities: q[3:], UpgradeAllowed: true, Language: "Any"})
	if body.Order != "quality_then_points" {
		t.Fatalf("separate qualities: order %q, want quality_then_points", body.Order)
	}
}
