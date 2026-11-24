# Backlog (German fork)

Ideas and gaps from using the fork, with their status. New items go at the end of **Open**, and finished ones move to **Done** with their PR.

Status: **open** (noted, not started), **in progress**, **done**.

**Upstream** says whether an item should also go to HolyArrs as an issue and pull request, so the fork stays small (see [docs/plattform-basis.md](docs/plattform-basis.md)): **offer it**, **ask first** or **keep in the fork**. New code goes into files of its own; upstream files get only small changes.

## Open

### 7. Quicker requests, modelled on Seerr
- **Status:** on hold (Thomas, 2026-11-21: zunächst zurückstellen). Plan switched to a plugin on 2026-11-24 (Thomas: "punkt 7 auf plug-in umstellen"); see "Umstellung auf Plugin" in the plan.
- **Upstream:** plugin, not fork code. The rules (C1), searching right after adding (A2), a quick request page for the owner and the parents' approval by message become plugins. The parts plugins cannot reach (buttons in search suggestions, badges, season dialog, TV view) go to HolyArrs as issues. First issue: let the owner show a plugin page to other roles, because plugin pages are owner-only today.
- **Problem:** Requesting a new title takes many steps. People who manage the Library always get "Zur Library hinzufügen" with up to ten fields (folder, Pick rules, seasons, order, search …), and a new show does not search by default.
- **Wish:** A short request like in Seerr: confirm a movie, tick seasons for a show, the rest comes from the request settings and sits under "Erweitert". Request button right in the search suggestions and on the cards, status badges on posters, and rules that pick folder and Pick rules by language or genre (German titles, kids' shows, anime).
- **Plan:** [docs/ui-seerr-plan.md](docs/ui-seerr-plan.md)

### 8. Show each episode's release-name number
- **Status:** done, [PR #23](https://github.com/redzacmx3/HolyArrs_german/pull/23) (2026-11-21). For people who manage the Library, each episode row on a show's Library tab says "Im Release-Namen: S02E41" when the show reads release names by TMDB or TVmaze, or "Bei TVmaze nicht gefunden" when TVmaze has no match. Nothing shows for "Nach dieser Reihenfolge" and for specials.
- **Upstream:** offer it. Build in the fork, then propose it to HolyArrs (issue first).
- **Wish:** see per episode which number releases use (TVmaze, TMDB or the chosen order), so a wrong match can be spotted.

### 9. Rejection reasons in German, with the setting's UI name
- **Status:** done, [PR #24](https://github.com/redzacmx3/HolyArrs_german/pull/24) (2026-11-21). The reasons a Release is turned down or held are known server lines now (German translated, other languages English), and the season pack reasons name the setting by its UI name. Text passed through as it is (trait names, title match details) stays as the server sends it.
- **Upstream:** offer it. Wording in German is a locale fix; naming the setting by its UI name helps every language.
- **Problem:** reasons such as "not every episode would gain … (Settings > Search: every episode must gain)" are English and name a setting the German UI calls "Einstellungen › Suche › Staffelpakete › Jede Folge darin neu oder besser wäre".

### 10. Points breakdown per release in "Release auswählen"
- **Status:** done by upstream HolyArrs 1.3.0 (#539: the reasons list each trait with its points). Our [PR #25](https://github.com/redzacmx3/HolyArrs_german/pull/25) was closed unmerged (Thomas, 2026-11-23).
- **Upstream:** already there.
- **Problem:** a release shows only the sum ("−23300 points, needs at least 11000"); which trait cost −35000 is not visible at a glance.

### 11. Episode search finds double-episode releases
- **Status:** done, [PR #15](https://github.com/redzacmx3/HolyArrs_german/pull/15) (2026-11-18), for shows that read release names by TVmaze: an episode search without an answer also asks the release-name season. Open: whether all shows should do this.
- **Upstream:** offer it, together with the TVmaze numbering (item 6).

### 12. "German Microsized" by size per minute for short episodes
- **Status:** done, [PR #26](https://github.com/redzacmx3/HolyArrs_german/pull/26) (2026-11-21). New own trait kind "Pro Minute" (size per minute of runtime in MB). **To do on the NAS:** set "German Microsized" to 0 in the series Pick rules and give the new trait (for example "bis 10 MB/Min." = −35000) instead.
- **Upstream:** offer it (Thomas, 2026-11-21): build in the fork, then propose it to HolyArrs as a fix.
- **Problem:** 12-minute kids' episodes (377 MB, 1080p x265, FuN) get −35000.
- **Found:** TRaSH's "German Microsized" tests no size at all. It is a list of release groups (FuN, w00t, GTF, PaTroL, JellyfinPlex, BB, 266ers), so every FuN release gets −35000 whatever its size. A trait that reads MB per minute of runtime judges the file itself.

### 13. WebHD fix: compare with upstream 1.2.1 before removing ours
- **Status:** done (2026-11-23, upstream 1.3.0 sync). Upstream reads WebHD/Web-HD as a web download itself (`releasename/quality.go`, test `issue497_test.go` with the FuN name). Our `TestWebHDIsWebDownload` passes on upstream's code for all three names, so our parser change is gone and our test stays.
- **Upstream:** already there. [HolyArrs issue #497](https://github.com/HolyArrs/HolyArrs/issues/497) (WebHD/Web-HD not read as a source): the maintainer confirmed a fix in 1.2.1, together with a review of how release names are read for source and quality.
- **Ours:** [PR #16](https://github.com/redzacmx3/HolyArrs_german/pull/16) reads WebHD/Web-HD as a web download (`releasename/parse.go`, test `TestWebHDIsWebDownload` in `languages_test.go`).
- **To do when 1.2.1 comes:** before removing our fix, run both against the same release names, for example `PAW.Patrol.S02E31.German.1080p.WebHD.x264-FuN`, `PAW.Patrol.S02E31.German.AC3.DL.1080p.WebHD.x265-FuN` and `PAW.Patrol.S02E31.German.DL.1080p.Web-HD.x265-FuN`. Check source, quality and ladder step, see which one reads more names right, then decide: drop ours, keep ours, or offer the missing cases upstream. Keep our test either way.

### 14. AI search: find a title from a description (plugin)
- **Status:** open (Thomas' idea, 2026-11-24)
- **Upstream:** plugin, published through HolyArrs' signed plugin list.
- **Wish:** Sometimes you look for a film without knowing its name. A description ("Film, wo ein Astronaut allein auf dem Mars Kartoffeln anbaut") should find it directly. Typed, and above all spoken on Android TV.
- **Idea:**
  1. A language model turns the description into up to five guesses with title and year. It runs locally, for example llama.cpp with Qwen on the NAS, or optionally as a cloud model.
  2. The plugin checks each guess through HolyArrs' metadata search. Only real hits are shown, so made-up titles never appear.
  3. For each hit the plugin shows whether it is in the Library, and offers to request it as the person.
- **Where people use it:**
  - as a tool in AI apps over MCP (`assistant.tools`), for everyone with their own rights;
  - as a plugin page, for the owner only until HolyArrs lets other roles see plugin pages (see item 7);
  - on the TV: through the planned TV app with Android voice input, or through the microphone key of the TV keyboard in a normal search field.
- **To check:**
  - how fast the NAS runs a small model;
  - the plugin reaches its own server by name only, so the model server needs a host name.

## Done

### 1. Filter list in "Release auswählen"
- **Status:** done, [PR #18](https://github.com/redzacmx3/HolyArrs_german/pull/18) (2026-11-18). Since the upstream 1.3.0 sync (2026-11-23) the sheet uses upstream's own filters (resolution, Source, language, name text, sort); our HDR, codec and sound drop-downs are gone. The DE / DL / Sub chips stay (`readTags` in `mediafiles/releases.ts`).
- **Done:** drop-downs above the list for Sprache, Auflösung, Quelle, HDR, Codec and Ton, each with how many Releases fit. They combine with each other and with Gute Treffer / Alle, and the browser keeps the picks for the next title. The server now sends what each name says (language, German DL, subtitles, picture, sound) as read by the parser, and the rows show DE / DL chips. The sheet loads up to 300 Releases instead of 60.
- **Problem:** The dialog lists many releases (for example "Alle 92"). Only "Gute Treffer / Alle" narrows them, so finding the right one takes long.
- **Wish:** Drop-down filters above the list for language (German, German DL, English …), picture (720p/1080p/2160p, HDR/DV, x264/x265, WEB/BluRay) and sound (DDP5.1, Atmos, TrueHD …). The filters could combine and use the facts the parser already reads, which the dialog shows as chips.

### 2. Bulk actions for problem downloads
- **Status:** done, [PR #17](https://github.com/redzacmx3/HolyArrs_german/pull/17) (2026-11-18)
- **Why greyed out:** the rows in "Prüfen" wait on a decision. Most have no download left in the Downloader (PAW Patrol at 0 %, or SABnzbd already cleaned up), so Starten, Pausieren and Entfernen have nothing to act on.
- **Done:** with rows ticked, the selection bar shows "N zu prüfen:" with Verwerfen und neu suchen, Erneut suchen, Datei in der Bibliothek behalten (when a copy is there) and Nicht mehr suchen. Each runs for every ticked row at once.
- **Problem:** On Downloads › Probleme you can select every row ("43 ausgewählt"), but the actions at the top left (Starten, Pausieren, Entfernen) stay greyed out. Nothing can be done with the selection.
- **Wish:** Actions on the selection: remove (with or without files), remove and search again, and add the release to "abgewiesen" and search again. This would clear many blocked downloads at once and load them fresh.

### 3. Make "Prüfen / Import blockiert" understandable
- **Status:** done, [PR #17](https://github.com/redzacmx3/HolyArrs_german/pull/17) (2026-11-18)
- **Done:** the list says the reason in a few words instead of "Import blockiert" (for example "Keine passenden Releases mehr", "Keine Datei passt", "Gehört in eine andere Variante"). The row details show the quick fixes next to the main button: Release von Hand wählen opens "Release auswählen" for just this episode or movie right there, Verwerfen und neu suchen, and Datei in der Bibliothek behalten.
- **Problem:** Downloads in "Prüfen" with "Import blockiert" do not show clearly why they are blocked or what to do. There is no quick way to pick an existing file or a new file for the episode so the status clears.
- **Wish:** One clear reason per row, plus quick actions such as "use the file already in the library", "pick another file / release", "import anyway" and "remove and search again".
- **Note:** Many of the blocked rows are PAW Patrol (S06/S07, 0 %, from HolyArrs). Check whether item 4 below fixes them once "Release-Namen zählen: Nach TMDB" is set.

### 6. Read release names by TVDB numbering
- **Status:** done, [PR #8](https://github.com/redzacmx3/HolyArrs_german/pull/8) and [PR #9](https://github.com/redzacmx3/HolyArrs_german/pull/9) (2026-11-18)
- **Problem:** German groups (FuN) name PAW Patrol by TVDB. The library follows TMDB "Production". No TMDB order matches, and there is no TVDB key.
- **Solution:** "Release-Namen zählen: Nach TVmaze (TVDB)". TVmaze numbers like TVDB, needs no key; each episode is found there by name or air day. PR #8: a double episode picked by hand imports into the chosen episode.

### 4. German releases with other numbering and double episodes (PAW Patrol)
- **Status:** done, [PR #6](https://github.com/redzacmx3/HolyArrs_german/pull/6) (2026-11-18)
- **Problem:** The library follows "Production", but German releases (FuN) count by TMDB. `S02E46` named another episode, so every release was turned down.
- **Fix:** A per-show setting "Release-Namen zählen: Nach TMDB" on the Library tab. Search, matching and import read release names by TMDB's numbers, and files keep the chosen order's numbers.
- **Follow-up:** for PAW Patrol it is not enough. FuN releases count by **TVDB** ("Pups Save an Eagle" = TVDB S02E41, Production S03E47), not by TMDB. See item 6.

### 5. Profile check: Fix fills in the advice
- **Status:** done, [PR #5](https://github.com/redzacmx3/HolyArrs_german/pull/5) (2026-11-17)
