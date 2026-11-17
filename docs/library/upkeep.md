---
title: Upkeep
description: Keep the Library tidy on its own: list or remove movies nobody watched and movies with no file, each rule automatic, suggested or off.
level: Setup
---

Upkeep finds movies you probably don't need any more. Each rule acts on its own, only lists what it finds, or is off. Whatever goes is moved to **Recently deleted** first.

## Before you start

- You can manage the Library. Changing the rules needs the settings permission.
- For **Movies nobody watched**, turn on **Read watch history** under **Settings > Media servers**, so HolyArrs knows what was played.

## The rules

| Rule | Finds |
|---|---|
| **Movies nobody watched** | Movies in the Library for at least the set days (365 by default) and never played by anyone on a media server |
| **Movies without a file** | Movies with no file that no Version keeps up to date |

Each rule has a mode:

| Mode | Does |
|---|---|
| **Automatic** | The daily run moves what it finds to Recently deleted |
| **Suggest** | Lists it on the Upkeep page; nothing goes unless you choose |
| **Off** | Finds nothing |

## Set the rules

1. Open **Library > Upkeep**.
2. Change a mode or the days. The save bar appears at the bottom.
3. Select **Save**, or **Discard** to go back.

**Pause upkeep: automatic rules only suggest** makes every automatic rule only suggest until you lift it.

## Remove what is listed

1. Tick the movies, or the box above the list for all.
2. Select **Move to Recently deleted**.

Protected copies and movies already gone are skipped, and the note says how many.

## Get it back

Open **Library > Recently deleted**. Files wait there for the holding days before they are gone for good.

## Related

- [Protected copies](protected-copies.md)
- [Media servers](media-servers.md)
