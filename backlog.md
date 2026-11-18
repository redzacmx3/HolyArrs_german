# Backlog (German fork)

Ideas and gaps from using the fork, with their status. New items go at the end of **Open**, and finished ones move to **Done** with their PR.

Status: **open** (noted, not started), **in progress**, **done**.

## Open

### 1. Filter list in "Release auswählen"
- **Status:** open (noted 2026-11-18)
- **Problem:** The dialog lists many releases (for example "Alle 92"). Only "Gute Treffer / Alle" narrows them, so finding the right one takes long.
- **Wish:** Drop-down filters above the list for language (German, German DL, English …), picture (720p/1080p/2160p, HDR/DV, x264/x265, WEB/BluRay) and sound (DDP5.1, Atmos, TrueHD …). The filters could combine and use the facts the parser already reads, which the dialog shows as chips.

### 2. Bulk actions for problem downloads
- **Status:** open (noted 2026-11-18)
- **Problem:** On Downloads › Probleme you can select every row ("43 ausgewählt"), but the actions at the top left (Starten, Pausieren, Entfernen) stay greyed out. Nothing can be done with the selection.
- **Wish:** Actions on the selection: remove (with or without files), remove and search again, and add the release to "abgewiesen" and search again. This would clear many blocked downloads at once and load them fresh.

### 3. Make "Prüfen / Import blockiert" understandable
- **Status:** open (noted 2026-11-18)
- **Problem:** Downloads in "Prüfen" with "Import blockiert" do not show clearly why they are blocked or what to do. There is no quick way to pick an existing file or a new file for the episode so the status clears.
- **Wish:** One clear reason per row, plus quick actions such as "use the file already in the library", "pick another file / release", "import anyway" and "remove and search again".
- **Note:** Many of the blocked rows are PAW Patrol (S06/S07, 0 %, from HolyArrs). Check whether item 4 below fixes them once "Release-Namen zählen: Nach TMDB" is set.

### 7. Quicker requests, modelled on Seerr
- **Status:** open (noted 2026-11-18), after the running features are stable
- **Problem:** Requesting a new title takes many steps. People who manage the Library always get "Zur Library hinzufügen" with up to ten fields (folder, Pick rules, seasons, order, search …), and a new show does not search by default.
- **Wish:** A short request like in Seerr: confirm a movie, tick seasons for a show, the rest comes from the request settings and sits under "Erweitert". Request button right in the search suggestions and on the cards, status badges on posters, and rules that pick folder and Pick rules by language or genre (German titles, kids' shows, anime).
- **Plan:** [docs/ui-seerr-plan.md](docs/ui-seerr-plan.md)

## Done

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
