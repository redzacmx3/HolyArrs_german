# HolyArrs German (Fork)

Dieser Fork von [HolyArrs/HolyArrs](https://github.com/HolyArrs/HolyArrs) ergänzt einen Filter für deutsche Usenet-Releases (German DL/ML, Subbed/OmU, TRaSH-Sprachformate, eigene Sprach-Traits).

## Test-Image

`.github/workflows/german-image.yml` baut bei jedem Push auf `main` oder `german-*` ein Docker-Image:

```
ghcr.io/redzacmx3/holyarrs_german:<branch>   # z. B. :main oder :german-usenet-filter
```

## Upstream-Stand übernehmen

`.github/workflows/upstream-sync.yml` läuft jeden Montag um 4:17 UTC (oder von Hand unter *Actions*):

1. holt `HolyArrs/HolyArrs` `main` und merged es in den Branch `upstream-sync` (von diesem `main` aus),
2. gibt eigenen Migrationen eine neue Nummer, wenn Upstream ihre Nummer belegt hat (`.github/scripts/renumber-fork-migrations.sh`),
3. baut und testet die Teile des Forks,
4. öffnet einen Pull Request „Upstream-Update …“, den du prüfst und mergst.

Bei einem Merge-Konflikt bricht der Lauf ab; die Zusammenfassung des Laufs nennt die Dateien.

Damit der Lauf auch geänderte Upstream-Workflows pushen kann und die Checks auf dem PR laufen, braucht er ein Repository-Secret `SYNC_TOKEN`: ein Fine-grained Token nur für dieses Repository mit *Contents*, *Pull requests* und *Workflows* (Lesen und Schreiben). Ohne das Secret klappt es, solange Upstream nichts unter `.github/workflows` ändert.

## Migrationen

HolyArrs führt jede noch nicht gelaufene Migration aus, auch eine mit kleinerer Nummer als die neueste schon gelaufene. Eine umnummerierte Fork-Migration wird über ihren Namen wiedererkannt und nicht doppelt ausgeführt (`app/backend/internal/core/database/migrate.go`, Test `TestForkMigrationRenumberedBehindUpstream`). So kann ein NAS mit dem Fork-Image jederzeit auf einen neueren Stand gehen.

## Von Hand

```
git remote add upstream https://github.com/HolyArrs/HolyArrs.git
git fetch upstream
git checkout -b upstream-sync main && git merge upstream/main
sh .github/scripts/renumber-fork-migrations.sh upstream/main
go test ./app/backend/tools/structurecheck -update
```

Als Patch: `git format-patch upstream/main..main` und auf einem frischen Stand `git am *.patch`.
