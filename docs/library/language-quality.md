---
title: Language and quality
description: Change the language and resolution HolyArrs looks for, for every movie or show, one title, one season or one episode, without editing Pick rules.
level: Setup
---

Change what HolyArrs looks for after a title is added: a different language, or a higher or lower resolution. The change lies over the Pick rules, so the sets themselves, and any TRaSH Guides sync, stay as they are.

## Before you start

- You can manage the Library (for a title, season or episode) or change settings (for every movie or show).
- Music is not covered.

## Where to open it

| For | Open |
|---|---|
| Every movie or every show | **Settings > Movies and TV defaults**, then **Language & quality for all movies** or **Language & quality for all shows** |
| One movie | The movie page, **⋯**, then **Language & quality** |
| One show | The show page, **⋯**, then **Language & quality** |
| One season | The season's menu on the show page, then **Language & quality** |
| One episode | The episode's menu, then **Language & quality** |

## Choose

1. Under **Language**, pick one:
   - **As in the rules**: keep what the level above says.
   - **This language or the original**: pick the language. A Release in that language or in the title's original language passes. German DL counts, since it holds both.
   - **Only this language**: only Releases in that language pass.
   - **Original language only**: only the title's original language passes.
2. Under **Quality**, pick **As in the rules**, **720p**, **1080p** or **4K**.
3. Read **What changes**. It shows, for each Version (or each Pick rules set in use), the qualities allowed, where it is good enough, the languages, and how many files are good enough before and after.
4. Select **Save**.

You should see a short note that the language and quality were saved.

## What a quality choice does

- Steps above the chosen resolution are no longer allowed.
- When the Pick rules allow nothing at that resolution, its web and Blu-ray steps are allowed. A 1080p set asked for 4K also gets 4K size limits.
- **Good enough at** moves to the best allowed step at that resolution (not remux).

## How the levels add up

From the top down: Pick rules, then every movie or show, then the title, the season and the episode. Each level only changes what it names; **As in the rules** passes the level above through. An episode wish beats its season's, and a season's beats the show's.

A season pack that holds episodes with different wishes must pass each of them.

## After saving

HolyArrs checks the files again at once. A file that is no longer good enough is searched for a better copy, if the title is kept up to date. A file that now meets the wish stops being searched.

To go back, open the same level and select **Reset**.

## Related

- [TRaSH Guides scores](trash-guides.md)
- [Audio versions](audio-versions.md)
- [Searching for movies](../movies/search.md) and [shows](../shows/search.md)
