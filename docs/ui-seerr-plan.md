# Plan: Schneller anfragen nach Seerr-Vorbild

Status: **offen**, erst nach der Stabilisierung der laufenden Features (Backlog-Punkte 1 bis 3). Noch keine Implementierung.
Stand: 2026-11-18. Backlog-Eintrag: Punkt 7 in [backlog.md](../backlog.md).

## Kurzfassung

HolyArrs hat für Mitglieder schon einen Seerr-ähnlichen Weg: **Anfragen** fragt nur nach Staffeln und eventuell Qualität. Ordner, Pick rules und „Suchen, wenn genehmigt“ kommen aus den Request-Einstellungen. Wer die Library verwaltet (Owner, Admins, also Thomas), bekommt diesen Weg aber nie. Er sieht immer **Zur Library hinzufügen** mit einem Formular aus bis zu zehn Feldern.

Der größte Gewinn ist daher kein Umbau, sondern: **auch Admins bekommen den kurzen Weg, das lange Formular wird zu „Erweitert“.** Dazu kommen ein Anfrage-Knopf direkt in der Suche und auf den Karten sowie Zuordnungsregeln nach Seerr-Vorbild, die Ordner und Pick rules von selbst wählen (zum Beispiel deutsche Titel oder Kinderserien).

## Wie die beiden heute anfragen

| Schritt | Seerr | HolyArrs (Admin) | HolyArrs (Mitglied) |
|---|---|---|---|
| Finden | Suchleiste oben, Ergebnisse als Posterraster | Suchfeld oben mit Vorschlägen, Enter öffnet Search | gleich |
| Aus der Suche anfragen | Knopf **Anfragen** direkt auf der Karte (Hover, auf dem Handy Tap) | Vorschlag öffnet die Titelseite; auf der Search-Seite öffnet das runde **+** ein Formular | wie Admin, Formular ist kurz |
| Film | Ein Dialog: Bestätigen. Server, Profil und Ordner stecken unter „Erweitert“ | **Hinzufügen**: Library-Ordner, Pick rules, Veröffentlichung abwarten, Aktuell halten, Suchen beim Hinzufügen, zweite Version (4K), Folgen | **Anfragen**: eventuell Qualität (Standard/4K), fertig |
| Serie | Liste aller Staffeln mit Schalter und Status (verfügbar, angefragt), „Alle Staffeln“, „Erweitert“ eingeklappt | Ordner, Pick rules, Staffeln, neue Staffeln, Episodenreihenfolge, Quelle für Details, Aktuell halten, Suchen, zweite Version, Folgen | Ganze Serie oder Staffeln ankreuzen, eventuell Qualität |
| Nach dem Absenden | Admin-Anfrage ist sofort genehmigt, Sonarr/Radarr sucht sofort | Serie: **Suchen beim Hinzufügen ist aus** (Standard), es passiert erst beim nächsten Feed | Genehmigung (automatisch oder durch Admin), dann Suche |
| Status sehen | Status-Badge auf jeder Karte, Download-Fortschritt in der Anfrage | Status-Pill auf der Titelseite, Requests-Seite | Requests-Seite |

Quellen im Code: `app/frontend/src/lib/discover/TitlePage.svelte` (Zeilen um 428 bis 446: Admins sehen nur „Hinzufügen“), `lib/requests/AddShowDialog.svelte`, `lib/requests/AddMovieDialog.svelte`, `lib/requests/QuickAsk.svelte` („never adds or asks in one click“), `lib/discover/live-search.ts` (`suggestHref`: Vorschläge sind nur Links). Request-Einstellungen: `docs/requests/settings.md`.

## Vorschläge

Sortiert nach Nutzen pro Aufwand. Jeder Schritt lässt sich einzeln ausliefern und testen.

### A. Schnelle Gewinne (klein, wenig Konfliktgefahr mit Upstream)

**A1. „Anfragen“ auch für Admins (Kern des Plans)**
- Auf Titelseite und Karten bekommt ein Admin denselben kurzen Dialog wie ein Mitglied: Film bestätigen, bei Serien Staffeln wählen.
- Ordner, Pick rules, 4K und „Suchen“ kommen aus den Request-Einstellungen („Wohin genehmigte Titel gehen“). Die Admin-Anfrage ist automatisch genehmigt, wie in Seerr.
- Unten im Dialog: **Erweitert…** öffnet das heutige Hinzufügen-Formular, vorausgefüllt. Nichts geht verloren.
- Neue Einstellung „Admins fragen schnell an“ (Standard an), damit der alte Weg zurückkommt, wenn jemand ihn lieber hat.
- Technik: Das Backend kann Anfragen schon anlegen und genehmigen. Zu prüfen ist, ob eine Anfrage eines Owners mit `manage_library` den Ask-Endpunkt nutzen darf oder ob der Dialog direkt „Hinzufügen“ mit den Request-Standards aufruft. Letzteres ist kleiner und braucht keine Backend-Änderung.

**A2. „Suchen beim Hinzufügen“ für Serien standardmäßig an**
- Heute ist es bei Serien aus (`search_on_add: false` in den eingebauten Standards). Wer eine Serie hinzufügt, wartet sonst auf den nächsten Feed und denkt, es passiert nichts.
- Seerr sucht sofort. Vorschlag: an, mit Hinweis „Sucht jetzt alle ausgestrahlten Folgen“. Die gespeicherte eigene Wahl bleibt erhalten.

**A3. Anfrage-Knopf direkt in den Suchvorschlägen**
- Jede Zeile im Dropdown der Suchleiste bekommt rechts einen kleinen Knopf: **+** (anfragen), **✓** (in der Library) oder den Status (angefragt, lädt).
- Ein Klick öffnet den kurzen Dialog aus A1, ohne die Titelseite zu laden. Für einen Film sind das drei Schritte: Name tippen, +, Bestätigen.

**A4. Status auf jeder Karte**
- Wie in Seerr ein Badge oben auf dem Poster: Verfügbar, Teilweise verfügbar (Serie), Angefragt, Lädt (mit Prozent), Fehlt. `CardBadges.svelte` hat schon eine Grundlage; zu prüfen ist, ob Download-Fortschritt dort schon ankommt.

### B. Mittlere Schritte

**B1. Staffel-Dialog wie in Seerr**
- Eine Liste aller Staffeln mit Schalter pro Staffel, Spalten „Folgen“ und „Status“ (in Library, angefragt, fehlt) und oben „Alle Staffeln“ plus „Neue Staffeln automatisch“.
- Episodenreihenfolge, Detail-Quelle und „Release-Namen zählen“ (TMDB/TVmaze) wandern unter „Erweitert“. Sie sind für die meisten Serien richtig auf Standard und verwirren beim schnellen Anfragen.

**B2. Ein Ort für „Anfragen“**
- Ein deutlicher Knopf **Neue Anfrage** auf Startseite und Requests-Seite, der direkt die Suche mit Fokus öffnet (heute: „Ask for something“ nur auf der Requests-Seite).
- Tastenkürzel `/` fokussiert die Suche (falls noch nicht vorhanden, prüfen).

**B3. Entdecken mit deutschen Reihen**
- Seerr erlaubt eine Region und Sprache für Entdecken sowie eigene Reihen (nach Genre, Studio, Sender, Stichwort).
- Vorschlag: Reihen „Beliebt in Deutschland“, „Deutsche Produktionen“ (TMDB `with_original_language=de`), „Neu bei Netflix/Prime/Disney+ in DE“. Die Region-Logik für Streaming-Toplisten gibt es schon (`ServiceTop(ctx, region)` in `cmd/holyarrs/services.go`).
- Option „Angefragte und vorhandene Titel ausblenden“ (Seerr 3.5 hat das neu), damit Entdecken nur zeigt, was man noch nicht hat.

**B4. Bedienung am Fernseher (für die Familie)**

Ziel: Ein Familienmitglied sitzt mit der Fernbedienung vor dem Fernseher und fragt einen Film oder eine Serie an, ohne Maus, Tastatur oder Handy. Seerr kann das übrigens auch nicht gut, es ist auf Maus und Touch gebaut. Hier können wir besser werden als das Vorbild.

*Wie die Familie HolyArrs am Fernseher öffnet*
- Über den Browser des Fernsehers oder der Box (Fire TV Silk, Android-TV-Browser, Samsung/LG-Browser). Eine eigene TV-App ist nicht Teil dieses Plans.
- Vorschlag: eine eigene Adresse `/tv` oder ein Schalter „Fernseher-Ansicht“ unter Einstellungen › Darstellung, der pro Gerät gespeichert wird. Automatisch erkennen ist unzuverlässig (TV-Browser melden sich oft wie ein Desktop), deshalb lieber ein fester Einstieg, den man einmal als Lesezeichen anlegt.
- Die TV-Ansicht zeigt nur, was man dort braucht: Entdecken, Suchen, Anfragen, Meine Anfragen. Downloads, Library-Pflege und Einstellungen bleiben auf PC und Handy.

*Anmelden am Fernseher*
- Passwort mit der Fernbedienung tippen ist mühsam. Vorschlag: Anmelden per QR-Code. Der Fernseher zeigt einen Code, man scannt ihn mit dem Handy, auf dem man schon angemeldet ist, und bestätigt dort. HolyArrs kann QR-Codes schon zeichnen (`lib/accounts/QrCode.svelte`, heute für die Zwei-Faktor-Einrichtung).
- Profilauswahl wie bei Netflix: Wer an diesem Fernseher schon angemeldet war, erscheint als Bild zum Antippen. „Konto wechseln“ gibt es schon (`switchAccount` im Layout), heute ist jeder Wechsel aber eine volle Anmeldung. Für den Fernseher sollte ein Wechsel zwischen gemerkten Profilen ohne Passwort gehen, bei Erwachsenenkonten optional mit 4-stelliger PIN, damit Kinder nicht als Eltern anfragen.

*Bedienung mit der Fernbedienung*
- Alles geht mit den Pfeiltasten, OK und Zurück. Die Zurück-Taste schließt zuerst den offenen Dialog, dann geht sie eine Seite zurück, nie aus der App heraus.
- Räumliche Navigation: Pfeil hoch und runter springt zwischen Reihen, links und rechts innerhalb einer Reihe. Die Grundlage gibt es schon: `lib/ui/rove.ts` macht in einer Kartenreihe genau eine Karte per Tab erreichbar und bewegt den Fokus mit den Pfeiltasten. Für den Fernseher fehlt der Sprung zwischen den Reihen und an die richtige Stelle (die Karte, die senkrecht darunter liegt).
- Der Fokus ist immer sichtbar und groß: Karte vergrößert sich leicht, dicker heller Rahmen, kein dünner Browser-Rahmen. Beim Betreten einer Seite liegt der Fokus sofort auf dem ersten sinnvollen Ziel (Suchfeld oder erste Karte), nie im Nichts.
- Nichts darf nur über Hover gehen. Der Anfrage-Knopf auf der Karte (A3/A4) erscheint, sobald die Karte den Fokus hat, nicht nur beim Überfahren mit der Maus. HolyArrs trennt das heute schon teilweise mit `@media (hover: hover)` (z. B. in `WideCard.svelte`); das muss für jede Aktion im Anfrage-Weg gelten.

*Anfrage-Ablauf am Fernseher*
1. Startseite: oben „Suchen“, darunter Reihen (Beliebt, Für Kinder, Neu bei Netflix/Prime/Disney+ …, siehe B3).
2. Suchen: Ein Bildschirm-Alphabet in einem Raster, Ergebnisse erscheinen schon nach zwei, drei Buchstaben neben dem Raster. Wo die Fernbedienung ein Mikrofon hat, kann der TV-Browser sprechen statt tippen; dafür reicht ein normales Eingabefeld, das wir nicht blockieren. Die zuletzt gesuchten und angefragten Titel stehen oben, damit man selten tippen muss.
3. Titel öffnen: große Titelseite mit Poster, kurzer Beschreibung, Altersfreigabe und genau einem Hauptknopf „Anfragen“ (oder „Schon da“, „Angefragt“, „Lädt 45 %“). Der Fokus liegt schon auf diesem Knopf.
4. Film: OK drücken, fertig. Ein kurzer Hinweis „Angefragt, du bekommst Bescheid“ erscheint und verschwindet von selbst.
5. Serie: Der Staffel-Dialog aus B1 in TV-Form. Oben „Alle Staffeln“ (vorausgewählt), darunter eine Zeile pro Staffel zum An- und Abwählen mit OK, unten „Anfragen“. Alles passt ohne Scrollen auf einen 1080p-Bildschirm; bei sehr vielen Staffeln scrollt nur die Staffelliste.
6. „Erweitert“ gibt es am Fernseher nicht. Ordner, Pick rules und Qualität kommen immer aus den Request-Einstellungen und den Zuordnungsregeln (C1).
- Ziel: Ein Film in höchstens fünf Tastendrücken nach dem Finden, eine ganze Serie in höchstens sechs.

*Kinder am Fernseher*
- Kinderkonten gibt es schon (mit Altersgrenze und eigener Kennzeichnung, `lib/accounts/kid.ts`, `KidMark.svelte`). Am Fernseher sehen sie nur Titel innerhalb ihrer Altersgrenze, nur den kurzen Weg und keine Qualitätswahl.
- Ihre Anfragen warten auf Freigabe durch die Eltern, wenn „Automatisch genehmigen“ für das Konto aus ist (gibt es schon). Vorschlag: Die Eltern bekommen dann eine Nachricht aufs Handy mit „Genehmigen“ und „Ablehnen“ direkt in der Nachricht, damit niemand an den PC muss.
- „Meine Anfragen“ zeigt Kindern einfache Bilder statt Zustandswörtern: Uhr (wartet), Haken (kommt), Fernseher (kann man schauen).

*Lesbarkeit aus drei Metern*
- Mindestens 24 px Text auf 1080p, große Poster, wenige Elemente pro Bildschirm, hoher Kontrast im dunklen Design.
- Sichere Ränder von etwa 5 % rundum, weil manche Fernseher den Bildrand abschneiden (Overscan).
- Keine kleinen Symbole ohne Text, keine Tooltips (die gibt es am Fernseher nicht).

*Prüfen und testen*
- Jeder Schritt aus A1, A3, A4 und B1 bekommt einen Test, der nur mit Tastatur bedient wird (Pfeile, Enter, Escape/Zurück): Titel finden, anfragen, Staffel wählen, Dialog schließen. Das läuft mit Playwright in der Testsuite.
- Einmal von Hand auf einem echten Gerät der Familie (welches, ist offen, siehe Fragen unten) prüfen, bevor B4 als fertig gilt.

*Reihenfolge innerhalb von B4*
1. Fokus sichtbar machen und Hover-only-Aktionen im Anfrage-Weg beseitigen (klein, hilft auch am PC).
2. Pfeil-Navigation zwischen Reihen und Zurück-Taste.
3. TV-Ansicht `/tv` mit Bildschirm-Alphabet und TV-Staffeldialog.
4. Anmelden per QR-Code und Profilauswahl mit PIN.
5. Freigabe für Eltern per Nachricht aufs Handy.

### C. Größere Schritte

**C1. Zuordnungsregeln (Seerr: „Override rules“)**
- Seerr wählt Server, Qualitätsprofil, Ordner und Tags per Regel, je nach Genre, Originalsprache, Stichwort oder Jahr. Beispiel: Originalsprache Japanisch und Genre Animation → Anime-Ordner.
- Für unseren Fork besonders nützlich:
  - Originalsprache Deutsch → Pick rules „German“ (ohne DL-Pflicht).
  - Genre Kinder/Familie oder Altersfreigabe bis 6 → Kinder-Ordner, Pick rules „German DL“, „Release-Namen zählen: Nach TVmaze (TVDB)“ (die PAW-Patrol-Lösung aus Punkt 6 würde dann für neue Kinderserien von selbst greifen).
  - Anime → Anime-Ordner und Anime-Namensschema.
- Die Regeln gelten für Anfragen, „Hinzufügen“, Plex-Watchlist und „Alle anfragen“. Der Dialog zeigt dann „Regel ‚Kinderserien‘ wählt: Kinder-Ordner, German DL“ statt Formularfeldern.
- Ort: Einstellungen › Anfragen › Zuordnungsregeln. Braucht Backend (Tabelle, Auswertung beim Anlegen) und ist der größte Punkt hier.

**C2. Anfrage-Assistent für den Sprachwunsch (optional)**
- Seerr hat keinen Sprachfilter pro Anfrage; die Sprache steckt im Profil. Bei uns könnte die Quality-Zeile im Dialog statt „Standard/4K“ auch „Deutsch / Deutsch DL / Original“ anbieten, als zusätzliche Zeilen aus den Request-Einstellungen („What people can ask for“ kann schon mehrere Zeilen). Das passt zur Idee „Sprach-Assistent“ aus dem Projekt.

## Was Seerr sonst noch hat, und ob es uns hilft

| Seerr-Funktion | Bei HolyArrs | Empfehlung |
|---|---|---|
| Override rules (Genre, Sprache, Stichwort, Jahr → Profil/Ordner) | fehlt | **übernehmen**, siehe C1 |
| Region und Sprache für Entdecken, eigene Entdecken-Reihen | teilweise (Streaming-Toplisten nach Region) | übernehmen, siehe B3 |
| „Angefragtes ausblenden“ beim Entdecken (3.5) | fehlt | übernehmen, klein |
| Admin-Anfrage sofort genehmigt und gesucht | Hinzufügen, Serien ohne Suche | übernehmen, siehe A1/A2 |
| Status-Badges und Download-Fortschritt auf Karten | teilweise | übernehmen, siehe A4 |
| Blocklist | „Ausblenden“ / „Never add“ | vorhanden |
| Wochenlimits pro Nutzer und Art | vorhanden | vorhanden |
| Plex-Watchlist automatisch anfragen | vorhanden | vorhanden |
| Probleme melden (Issues) | Reports | vorhanden |
| Ganze Sammlung anfragen | „Alle anfragen“ | vorhanden |
| Benachrichtigungen: Discord, E-Mail, Telegram, Pushover, Gotify, ntfy, Web-Push | Discord, E-Mail (laut Doku) | später prüfen, eigener Punkt |
| TVDB als Quelle für Staffeln und Folgen | gelöst über TVmaze (Punkt 6, PR #8/#9) | nicht nötig; Seerr bezieht Staffeldaten nicht aus TVDB für die Release-Zuordnung, die Arbeit macht dort Sonarr |
| Sprache im Release-Filter | Seerr hat keinen, das macht Sonarr/Radarr | unser Fork ist hier schon weiter |

## Was wir nicht übernehmen

- Seerr ist nur ein Anfrage-Frontend vor Sonarr/Radarr. HolyArrs' Library-, Download- und Einstellungsseiten bleiben, wie sie sind; der Plan betrifft nur den Weg „Titel finden und anfragen“.
- Kein eigener Seerr-Nachbau: Optik und Komponenten bleiben HolyArrs, übernommen wird nur der Ablauf.

## Risiko: Konflikte mit dem Upstream-Sync

Der Fork holt jeden Montag Upstream. Die betroffenen Dateien (`TitlePage.svelte`, `QuickAsk.svelte`, `Add*Dialog.svelte`) ändert Upstream häufig. Darum:
- neue Teile als eigene Komponenten (z. B. `lib/requests/FastAsk.svelte`) und nur kleine Einhänge-Stellen in bestehenden Dateien,
- Schalter in den Einstellungen, damit der alte Weg jederzeit zurückkommt,
- A1 bis A3 eignen sich auch als Vorschlag an Upstream (HolyArrs), dann fällt die Pflege im Fork weg.

## Reihenfolge

1. Stabilisierung abschließen (Backlog-Punkte 1 bis 3 und laufende Fehler).
2. A2 (eine Zeile, sofort spürbar), dann A1, dann A3.
3. B1, dann A4.
4. B2, B3, dann B4 Schritt 3 bis 5 (TV-Ansicht, Anmelden, Eltern-Freigabe). B4 Schritt 1 und 2 laufen schon bei A1 bis B1 mit.
5. C1 als eigenes Vorhaben mit eigenem Plan: Kinderserien, dann Anime, dann deutsche Originale.

## Entscheidungen von Thomas (2026-11-18)

- **A1:** wie in Seerr. Eine Admin-Anfrage erscheint auf der Requests-Seite, mit Verlauf und Status, und ist sofort genehmigt.
- **C1:** Reihenfolge der Regeln: zuerst Kinderserien, dann Anime, zum Schluss deutsche Originale.
- **Nutzung:** Im Endausbau fragt die ganze Familie an, auch am Fernseher. Folgen für den Plan:
  - B1 (Staffel-Dialog) rückt nach vorn, direkt hinter A3.
  - Neu **B4. Bedienung am Fernseher**, ausgeführt oben unter B. Die Grundregeln (Fokus, keine Aktion nur über Hover) werden schon bei A1, A3, A4 und B1 mitgeprüft.
  - Kinderkonten (gibt es schon, mit Altersgrenze) sollen den kurzen Weg ohne „Erweitert“ sehen.

## Offene Fragen an Thomas (B4)

- Welche Geräte nutzt die Familie am Fernseher (Fire TV, Android TV/Google TV, Samsung, LG, Apple TV)? Apple TV hat keinen Browser; dort ginge es nur über das Handy oder eine eigene App.
- Soll ein Kind am Fernseher überhaupt anfragen dürfen, oder nur wünschen, und die Eltern entscheiden?
