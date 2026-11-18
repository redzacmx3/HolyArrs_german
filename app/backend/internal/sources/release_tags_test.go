package sources

import (
	"reflect"
	"testing"
)

func TestNameTagsReadLanguagePictureSound(t *testing.T) {
	for _, c := range []struct {
		name string
		want *NameTags
	}{
		{"A.Knight.of.the.Seven.Kingdoms.S01E06.2160p.HMAX.WEB-DL.DDP5.1.Atmos.DV.HDR.H.265-FLUX",
			&NameTags{Resolution: 2160, Source: "web_dl", Codec: "h265", HDR: []string{"dv", "hdr10"}, Audio: "eac3", Channels: "5.1", Atmos: true}},
		{"A.Knight.of.the.Seven.Kingdoms.2026.S01E06.The.Morrow.German.DL.1080p.HMAX.WEB-DL.DDP5.1.H.264-ZeroTwo",
			&NameTags{Languages: []string{"de"}, Dual: true, Resolution: 1080, Source: "web_dl", Codec: "h264", Audio: "eac3", Channels: "5.1"}},
		{"Some.Show.S01E01.German.Subbed.720p.WEB.x264-GRP", &NameTags{Subs: []string{"de"}, Resolution: 720, Source: "web_dl", Codec: "h264"}},
		{"Just a Title", nil},
	} {
		if got := nameTags(c.name); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s:\n got %+v\nwant %+v", c.name, got, c.want)
		}
	}
}
