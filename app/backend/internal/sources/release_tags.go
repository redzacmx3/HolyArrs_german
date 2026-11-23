package sources

import "github.com/HolyArrs/holyarrs/app/backend/internal/core/releasename"

// NameTags are what a Release name says about its languages, picture and
// sound, read the way the judge reads them. "Release auswählen" filters
// its list by them.
type NameTags struct {
	// Languages are the audio languages the name names (ISO 639-1); empty
	// when it names none, which mostly means the original language.
	Languages []string `json:"languages,omitempty"`
	// Dual says the name holds the original audio besides German
	// (German DL / ML).
	Dual       bool     `json:"dual,omitempty"`
	Subs       []string `json:"subs,omitempty"`       // subtitle languages ("German Subbed", "OmU")
	Resolution int      `json:"resolution,omitempty"` // 480, 720, 1080, 2160
	Source     string   `json:"source,omitempty"`     // web_dl, web_rip, bluray, remux, tv, dvd …
	Codec      string   `json:"codec,omitempty"`      // h264, h265, av1 …
	HDR        []string `json:"hdr,omitempty"`        // hdr10, hdr10plus, dv, hlg
	Audio      string   `json:"audio,omitempty"`      // aac, ac3, eac3, dts, dts_hd_ma, truehd, flac …
	Channels   string   `json:"channels,omitempty"`   // "5.1", "7.1"
	Atmos      bool     `json:"atmos,omitempty"`
}

// nameTags reads a Release name for the filters; nil when it says nothing.
func nameTags(name string) *NameTags {
	p := releasename.Parse(name, releasename.Hints{})
	t := NameTags{Languages: p.Languages, Dual: p.OriginalAudio, Subs: p.SubLanguages, Resolution: p.Resolution.Value,
		Source: p.Source.Value, Codec: p.Codec.Value, HDR: p.HDR.Value, Audio: p.Audio.Format, Channels: p.Audio.Channels, Atmos: p.Audio.Atmos}
	if p.Source.Conf == releasename.Unknown {
		t.Source = ""
	}
	if len(t.HDR) == 0 {
		t.HDR = nil
	}
	if len(t.Languages) == 0 {
		t.Languages = nil
	}
	if len(t.Languages) == 0 && !t.Dual && len(t.Subs) == 0 && t.Resolution == 0 && t.Source == "" && t.Codec == "" && len(t.HDR) == 0 && t.Audio == "" && t.Channels == "" && !t.Atmos {
		return nil
	}
	return &t
}
