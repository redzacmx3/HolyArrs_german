# Backlog (German fork)

Ideas and gaps from using the fork, with their status. New items go at the end of **Open**, and finished ones move to **Done** with their PR.

Status: **open** (noted, not started), **in progress**, **done**.

**Upstream** says whether an item should also go to HolyArrs as an issue and pull request, so the fork stays small (see [docs/plattform-basis.md](docs/plattform-basis.md)): **offer it**, **ask first** or **keep in the fork**. New code goes into files of its own; upstream files get only small changes.

## Open

### 7. Quicker requests, modelled on Seerr
- **Status:** open (noted 2026-11-18), after the running features are stable
- **Upstream:** keep in the fork. Build it as own pages and modules and touch upstream files only with small hooks, because upstream changes these pages a lot.
- **Problem:** Requesting a new title takes many steps. People who manage the Library always get "Zur Library hinzufügen" with up to ten fields (folder, Pick rules, seasons, order, search …), and a new show does not search by default.
- **Wish:** A short request like in Seerr: confirm a movie, tick seasons for a show, the rest comes from the request settings and sits under "Erweitert". Request button right in the search suggestions and on the cards, status badges on posters, and rules that pick folder and Pick rules by language or genre (German titles, kids' shows, anime).
- **Plan:** [docs/ui-seerr-plan.md](docs/ui-seerr-plan.md)

### 8. Show each episode's release-name number
- **Status:** open (noted 2026-11-18). The episodes API lists `tvmaze_season`/`tvmaze_episode` since [PR #15](https://github.com/redzacmx3/HolyArrs_german/pull/15); the episode list in the UI does not show them yet.
- **Upstream:** offer it. Build in the fork, then propose it to HolyArrs (issue first).
- **Wish:** see per episode which number releases use (TVmaze, TMDB or the chosen order), so a wrong match can be spotted.

### 9. Rejection reasons in German, with the setting's UI name
- **Status:** open (noted 2026-11-18)
- **Upstream:** offer it. Wording in German is a locale fix; naming the setting by its UI name helps every language.
- **Problem:** reasons such as "not every episode would gain … (Settings > Search: every episode must gain)" are English and name a setting the German UI calls "Einstellungen › Suche › Staffelpakete › Jede Folge darin neu oder besser wäre".

### 10. Points breakdown per release in "Release auswählen"
- **Status:** open (noted 2026-11-18)
- **Upstream:** offer it. Build in the fork, then propose it to HolyArrs (issue first).
- **Problem:** a release shows only the sum ("−23300 points, needs at least 11000"); which trait cost −35000 is not visible at a glance.

### 11. Episode search finds double-episode releases
- **Status:** done, [PR #15](https://github.com/redzacmx3/HolyArrs_german/pull/15) (2026-11-18), for shows that read release names by TVmaze: an episode search without an answer also asks the release-name season. Open: whether all shows should do this.
- **Upstream:** offer it, together with the TVmaze numbering (item 6).

### 12. "German Microsized" by size per minute for short episodes
- **Status:** open (noted 2026-11-18)
- **Upstream:** ask first. If HolyArrs does not want it, it stays a small rule of its own in the fork.
- **Problem:** TRaSH's fixed size limit hits 12-minute kids' episodes (377 MB, 1080p x265) with −35000. Judge by MB per minute of runtime instead.

## Done

### 1. Filter list in "Release auswählen"
- **Status:** done, [PR #18](https://github.com/redzacmx3/HolyArrs_german/pull/18) (2026-11-18)
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
