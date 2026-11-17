package mediafiles

// Advice the profile check gives about Pick rules that will not pick as
// meant. The words are fixed so the web app can put them in the reader's
// language.
const (
	AdviceOrder      = "It ranks quality before points, so a Release in the wanted language loses to a better quality without it. Set the order to points first."
	AdviceGoodEnough = "Good enough is 0 points, so a copy without the language the traits score is kept and never replaced. Set good enough points to the language trait's points."
)

// languagePoints is the least a trait must score to count as a language
// trait, as the TRaSH language formats (German, French, …) do: thousands
// of points, far above any quality or codec trait.
const languagePoints = 1000

// AdviseRules lists what keeps rules from picking as their traits ask: a
// language trait worth thousands of points still loses on quality when the
// rules rank quality first, and good enough at 0 points stops upgrades
// before the language is reached. Rules without such a trait get none.
func AdviseRules(b PickRulesBody) []string {
	top := 0
	for _, t := range b.Traits {
		top = max(top, t.Points)
	}
	if top < languagePoints {
		return nil
	}
	var out []string
	allowed := 0
	for _, s := range b.Steps {
		if s.Allowed {
			allowed++
		}
	}
	if b.Order != "points_then_quality" && allowed > 1 {
		out = append(out, AdviceOrder)
	}
	if b.GoodEnoughPoints != nil && *b.GoodEnoughPoints <= 0 {
		out = append(out, AdviceGoodEnough)
	}
	return out
}
