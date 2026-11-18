# Plattform-Basis: Fork, eigener Weg oder Neubau?

Stand: 18.11.2026. Nur eine Bewertung, es wird nichts umgesetzt.

## Kurz

**Empfehlung: HolyArrs bleibt die Basis, und der Fork bleibt vorerst nah am Upstream.** Was allgemein nützlich ist, geht als Pull Request zu HolyArrs, damit der Fork klein bleibt. Ein eigener Weg (Weg b) ist der geplante Notausgang mit klaren Auslösern, aber noch nicht jetzt. Ein Neubau (Weg c) lohnt sich für eine Einzelperson nicht.

Der Grund in einem Satz: "nur das Nützliche übernehmen" klingt sparsam, ist aber die teuerste Art der Pflege, weil man dafür jede Upstream-Änderung lesen und bewerten muss. Mergen und das wenige Störende anpassen ist billiger, solange Upstream in die gleiche Richtung geht. Und das tut es zurzeit.

## Ausgangslage (geprüft im Repo und in der Upstream-Historie)

| | |
|---|---|
| **Upstream-Alter** | HolyArrs 1.0.0 erschien am 9.11.2026, 1.2.0 am 17.11.2026. Das Projekt ist erst neun Tage öffentlich. |
| **Tempo** | 46 Commits auf `main` in zehn Tagen, dazu viele Releases (1.0.0 bis 1.2.0). |
| **Wer pflegt** | 41 der 46 Commits stammen von einer Person (CoderLuii), einige wenige von Außenstehenden. Pull Requests von außen werden also angenommen. |
| **Größe** | Rund 450.000 Zeilen Go und 140.000 Zeilen Frontend (Svelte/TS), Tests eingeschlossen. |
| **Lizenz** | AGPL-3.0. Beiträge laufen über ein CLA, das dem Maintainer auch eine kommerzielle Lizenz erlaubt. Was einmal unter AGPL veröffentlicht ist, bleibt es. |
| **Metadaten** | Poster, Trailer und Bewertungen kommen standardmäßig von einem Server, den der Maintainer betreibt. Ein eigener TMDB-Schlüssel geht auch (seit heute auch ein v3-Schlüssel). |
| **Richtung** | "Six apps to watch a show. Now it's one", "Browse like a streaming app, ask for what you want". Das ist fast wörtlich dein Ziel: den arr-Stack bündeln und für Einsteiger bedienbar machen. |
| **Der Fork heute** | 475 geänderte Dateien, rund 12.300 Zeilen dazu, gut 18 PRs in zwei Tagen. Wöchentlicher Upstream-Sync. Der Sync auf 1.2.0 (PR #4) hatte Konflikte und musste von Hand gemergt werden, Migrationen werden umnummeriert. |

## Der Maßstab

Deine Ziele aus den bisherigen Threads:

1. Deutsche Releases werden richtig erkannt und bevorzugt (German DL, FuN, TRaSH-Profile).
2. Kinderserien laufen sauber, auch mit krummer Nummerierung (PAW Patrol, Doppelfolgen).
3. Einsteiger können es bedienen **und** verwalten.
4. Die Familie fragt am Fire TV und Google TV an.
5. Möglichst wenig Pflegeaufwand für dich. Du willst nicht von Hand pflegen.

## Die drei Wege

### (a) Fork bleibt nah am Upstream

Jeden Montag kommt Upstream per Sync-PR rein, der Fork legt seine Änderungen obendrauf.

- **Dafür:** Du bekommst jede Woche Fehlerbehebungen, Sicherheitsupdates und neue Funktionen umsonst. Allein heute kamen ein v3-TMDB-Schlüssel, "In Jellyfin öffnen", MDBList-Listen und ein Parser-Fix rein. Die Richtung passt zu Ziel 3. Rückweg zu (b) bleibt jederzeit offen.
- **Dagegen:** Bei einem jungen, schnellen Projekt gibt es häufig Konflikte, vor allem in großen Dateien (`openapi.json`, Parser, Title-Seiten nahe der 800-Zeilen-Grenze). Je größer der Fork wird, desto öfter knallt es.
- **Aufwand:** klein bis mittel, solange der Fork schlank bleibt. Ein Sync-PR prüfen dauert Minuten, ein Konflikt eine bis zwei Stunden mit mir.
- **Risiko:** gering. Du hängst an der Richtung eines Einzelnen, kannst aber jederzeit abbiegen.

### (b) Fork wird eigenständig, Upstream nur gezielt

Der Fork löst sich, Upstream-Änderungen werden einzeln gelesen und nur die nützlichen übernommen.

- **Dafür:** volle Freiheit bei Bedienung, Wortwahl und Aufbau. Keine Überraschungen durch Upstream-Entscheidungen, keine Konflikte im Wochentakt.
- **Dagegen:** Rund 600.000 Zeilen Code gehören dann dir. Jede Upstream-Änderung muss jemand lesen und entscheiden, das sind bei diesem Tempo 30 bis 50 Commits pro Woche. Rosinenpicken wird mit jeder Woche schwerer, weil die Codebasen auseinanderlaufen; nach ein paar Monaten passt kaum ein Upstream-Commit noch sauber. Sicherheitsfixes musst du dann selbst finden und portieren. Der Weg ist praktisch eine Einbahnstraße.
- **Aufwand:** mittel bis hoch, und er wächst mit der Zeit. Das widerspricht Ziel 5 am stärksten, auch wenn es anfangs nach weniger Arbeit aussieht.
- **Risiko:** mittel. Das Hauptrisiko ist nicht der Code von heute, sondern dass du in einem Jahr eine große, alleine gepflegte Codebasis hast.

### (c) Neubau, der den arr-Stack bündelt

Entweder alles neu (Parser, Indexer, Downloader, Import, Metadaten, Anfrage-Oberfläche) oder eine neue, einsteigerfreundliche Oberfläche über Sonarr, Radarr und Prowlarr.

- **Dafür:** sauberer Schnitt, genau deine Ideen von Anfang an, kein fremder Ballast.
- **Dagegen:** HolyArrs **ist** dieser Neubau, nur hat ihn schon jemand gemacht. Allein die Woche mit PAW Patrol (TMDB-Reihenfolgen, TVmaze, Doppelfolgen, WebHD, Microsized) zeigt, wie tief die Sonderfälle gehen; Sonarr und Radarr haben dafür über zehn Jahre gebraucht. Eine reine Oberfläche über die arr-Apps behält deren Verwaltung (sechs Apps, sechs Updates, Schlüssel zwischen den Apps), also genau das, was Einsteiger abschreckt.
- **Aufwand:** sehr hoch. Bis du auf dem Stand von heute bist, vergehen viele Monate, in denen das NAS mit dem alten Stack weiterläuft.
- **Risiko:** hoch. Halbfertige Neubauten sind der häufigste Ausgang solcher Projekte.

### Übersicht

| Ziel | (a) nah am Upstream | (b) eigenständig | (c) Neubau |
|---|---|---|---|
| Deutsche Releases | gut, schon gebaut | gut, schon gebaut | muss neu gebaut werden |
| Kinderserien / Nummerierung | gut, schon gebaut | gut, schon gebaut | muss neu gebaut werden |
| Einsteiger: Bedienung | Upstream arbeitet mit | du allein | du allein |
| Einsteiger: Verwaltung (eine App) | ja | ja | nur bei Vollneubau |
| Fire TV / Google TV | eigene kleine App, bei allen Wegen gleich | gleich | gleich |
| Pflegeaufwand | klein, wenn Fork schlank | wächst stetig | sehr hoch |
| Risiko | gering | mittel | hoch |

## Empfehlung im Einzelnen

Weg (a), mit drei Regeln, damit er billig bleibt:

1. **Upstream zuerst.** Was nicht nur für dich nützt, geht als Issue und dann als PR zu HolyArrs (Upstream verlangt erst ein Issue, dann einen kleinen PR, dazu einmal das CLA). Gute Kandidaten: German DL/ML/Subbed im Parser, TRaSH-Sprachtests, WebHD als Web-Quelle, Release-Namen nach TMDB/TVmaze zählen, die Profilprüfung, Sammelaktionen für Problem-Downloads, die Filterliste in "Release auswählen". Jeder angenommene PR ist Code, den du nie wieder mergen musst.
2. **Was im Fork bleibt, bleibt am Rand.** Eigene Logik in eigene Dateien und Module (wie schon `language_tags.go`), Upstream-Dateien nur mit kleinen Haken anfassen. Das hält Konflikte selten.
3. **Absicherung gegen den Einzelnen.** Einen eigenen TMDB-Schlüssel hinterlegen, damit Poster und Daten nicht vom Server des Maintainers abhängen. Die Docker-Images des Forks behalten, damit ein Rückweg immer geht.

**Wann auf (b) umsteigen**, dann bewusst und mit eingefrorenem Upstream-Stand:

- Upstream ist acht Wochen still.
- Neue Versionen erscheinen nicht mehr unter AGPL.
- Upstream lehnt die Richtung ab, die dir wichtig ist (Einsteiger, Familie, deutsche Releases), oder baut sie so um, dass deine Teile dauernd brechen.
- Die Syncs kosten im Schnitt mehr als etwa zwei Stunden im Monat.

Dann geht es weiter vom letzten AGPL-Stand, und von Upstream kommen nur noch Sicherheitsfixes.

**Unabhängig vom Weg:** Die Fire-TV-/Google-TV-App ist sowieso ein eigenes kleines Projekt (schlanke WebView-App, siehe [ui-seerr-plan.md](ui-seerr-plan.md)). Eine optionale lokale KI (llama.cpp, Qwen) für Empfehlungen und Fehlersuche lässt sich als Zusatz bauen, der nur Vorschläge macht, und gehört eher in den Fork als zu Upstream.

## Was ich nicht weiß

- Ob Upstream deutsche Sonderfälle annimmt. Das klärt das erste Issue, am besten mit dem kleinsten Kandidaten (WebHD oder German DL im Parser).
- Wie es um das Tempo in ein paar Monaten steht. Neun Tage Historie reichen nicht, um das sicher zu sagen; die Auslöser oben fangen das ab.

## Andere Projekte im Vergleich (18.11.2026)

Geprüft wurden die Repos selbst (Code, Historie, README), nicht die laufenden Apps.

| | **HolyArrs** (heute) | **Mydia** | **Cinephage** | **MediaManager** |
|---|---|---|---|---|
| Repo | HolyArrs/HolyArrs | getmydia/mydia | Cinephage/Cinephage | maxdorninger/MediaManager |
| Technik | Go + SvelteKit | Elixir (Phoenix) + Rust, Player in Flutter | TypeScript, SvelteKit, SQLite | Python + Svelte |
| Größe (Codezeilen, grob) | ~600.000 mit Tests | ~560.000 ohne Tests | ~400.000 ohne Tests | ~27.000 ohne Tests |
| Öffentlich seit | 11/2026 (Version 1.0) | 11/2025 | 12/2025 | 02/2025 |
| Letzter Commit | 18.11.2026 | 18.11.2026 | 03.11.2026 | 05.07.2026 |
| Tempo zuletzt | 46 Commits in 10 Tagen | 2.000 bis 2.700 Commits pro Monat | 150 bis 300 pro Monat, im November 5 | seit Juli nichts |
| Wer pflegt | fast nur eine Person | fast nur eine Person (10.300 von 10.500 Commits) | zwei Hauptleute, gut 20 insgesamt | eine Hauptperson, rund 35 insgesamt |
| Lizenz | AGPL-3.0 + CLA | AGPL-3.0 | GPL-3.0 | AGPL-3.0 |
| Ersetzt | Sonarr, Radarr, Prowlarr, Seerr, Bazarr, Lidarr | Sonarr, Radarr, Prowlarr, Overseerr | Sonarr, Radarr, Prowlarr, Bazarr, Overseerr, FlareSolverr, IPTV | Sonarr, Radarr (Grundfunktionen) |
| Deutsche Releases | im Fork gebaut (German DL/ML, TRaSH-Sprache, FuN, TVmaze) | Sprachwunsch wie `["de", "original"]`, Parser mit Radarr/Sonarr-Testkorpus | nur Sprachauswahl, kein German-DL-Parser gefunden | nichts gefunden |
| TRaSH-Profile | ja | ja (Custom Formats) | eigene Punkte, TRaSH-Benennung | nein |
| Anfragen / Familie | ja, Seerr-Plan im Fork | Gäste fragen an, Admin gibt frei | Listen, kaum Anfragen | einfache Anfragen |
| TV-Gerät | Web-App, eigene TV-App geplant | **eigener Player für Android TV**, iOS, Desktop, Web | Web-App | Web-App |

### Einschätzung

- **Mydia** ist die einzige ernsthafte Alternative. Es ist ein Jahr älter, bei deutschen Releases schon weiter als HolyArrs ohne Fork und hat mit dem eigenen Player für Android TV etwas, das zu Fire TV und Google TV passt. Dagegen spricht: Elixir und Rust sind eine Technik, in der ein Fork für dich schwerer zu pflegen ist als Go. Und es hängt genauso an einer einzigen Person wie HolyArrs, bei einem Tempo, das ein Fork kaum mitgehen kann.
- **Cinephage** setzt auf Streaming über `.strm`-Dateien aus Online-Quellen, das Umgehen von Cloudflare und das Absuchen von IPTV-Portalen nach MAC-Adressen. Das ist rechtlich und technisch wackelig und hat mit deinen Zielen (deutsche Releases, Familie, wenig Pflege) wenig zu tun. Im November ist es fast still.
- **MediaManager** ist klein und seit über vier Monaten ohne Commit. Es kann deutlich weniger als HolyArrs und kommt nicht in Frage.

### Folgerung

Ein Wechsel des Unterbaus lohnt sich jetzt nicht. Alles, was im Fork für deutsche Releases und Kinderserien gebaut ist (rund 20 PRs), müsste neu entstehen, und bei Mydia in einer anderen Sprache. Der Nutzen wäre vor allem der TV-Player.

Zwei Dinge lohnen sich trotzdem:

1. **Mydia als Plan B.** Falls einer der Auslöser oben eintritt (Upstream still, Lizenzwechsel, Richtung passt nicht mehr), ist Mydia der erste Kandidat und nicht ein eigener Neubau. Vorher einmal testen: eine Testinstanz mit deinen Indexern und einem PAW-Patrol-Fall.
2. **Mydia Player für den Fernseher ansehen.** Laut README spielt er auch ohne Mydia Server von Jellyfin ab. Falls er auf Fire TV gut läuft, spart das vielleicht einen Teil der eigenen TV-App. Anfragen bleiben dann in HolyArrs.
