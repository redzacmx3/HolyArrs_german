<script lang="ts">
	import PageLoader from '#lib/ui/PageLoader.svelte';
	import Poster from '#lib/ui/Poster.svelte';
	import Button from '#lib/ui/Button.svelte';
	import SaveBar from '#lib/ui/SaveBar.svelte';
	import EmptyState from '#lib/ui/EmptyState.svelte';
	import LoadError from '#lib/ui/LoadError.svelte';
	import { titleArt, LibraryToolPage } from '#lib/mediafiles/index.ts';
	import { onMount } from 'svelte';
	import { to } from '#lib/core/paths.ts';
	import { api, message } from '#lib/core/api.ts';
	import { bytes, date } from '#lib/core/format.ts';
	import { toasts } from '#lib/core/toasts.svelte.ts';
	import { plainReason } from '#lib/core/plain.ts';
	import { t } from '#lib/core/i18n/i18n.svelte.ts';

	// Upkeep: movies the upkeep rules name (never watched, no file). Each rule
	// acts on its own (auto), only lists (suggest) or is off. What goes is
	// moved to Recently deleted, where it waits for the holding days.

	type Mode = 'auto' | 'suggest' | 'off';
	type Settings = { never_watched: Mode; never_watched_days: number; no_file: Mode; paused: boolean };
	type Item = { rule: 'never_watched' | 'no_file'; mode: Mode; title_id: number; name: string; year?: number; added_at: string; size: number; files: number };
	type Done = { removed: number[]; skipped: number[]; problems?: string[] };

	let settings = $state<Settings>({ never_watched: 'suggest', never_watched_days: 365, no_file: 'suggest', paused: false });
	let history = $state(true);
	let items = $state<Item[]>([]);
	let loaded = $state(false);
	let error = $state('');
	let picked = $state<Set<number>>(new Set());
	let busy = $state(false);
	let saving = $state(false);
	// The rules as last loaded or saved; the save bar shows while they differ.
	let saved = $state('');
	const dirty = $derived(loaded && !!saved && JSON.stringify(settings) !== saved);
	const all = $derived(items.length > 0 && picked.size === items.length);
	const modes: Mode[] = ['auto', 'suggest', 'off'];

	async function load() {
		try {
			const r = await api.get<{ settings: Settings; watch_history: boolean; items: Item[] }>('/upkeep');
			settings = r.settings;
			saved = JSON.stringify(r.settings);
			history = r.watch_history;
			items = r.items;
			picked = new Set([...picked].filter((id) => items.some((i) => i.title_id === id)));
			error = '';
		} catch (e) {
			error = message(e);
		} finally {
			loaded = true;
		}
	}
	let art = $state<Map<number, string>>(new Map());
	onMount(() => {
		void load();
		void titleArt().then((m) => (art = m));
	});

	function toggle(id: number) {
		const next = new Set(picked);
		if (!next.delete(id)) next.add(id);
		picked = next;
	}

	async function save() {
		if (saving) return;
		saving = true;
		try {
			settings = await api.put<Settings>('/settings/upkeep', { ...settings, never_watched_days: Number(settings.never_watched_days) });
			saved = JSON.stringify(settings);
			toasts.push(t('library.upkeep.saved'), 'ok');
			await load();
		} catch (e) {
			toasts.push(t('library.upkeep.saveFailed'), 'fail', { body: message(e) });
		} finally {
			saving = false;
		}
	}

	async function remove() {
		if (!picked.size || busy) return;
		busy = true;
		try {
			const r = await api.post<Done>('/upkeep/apply', { title_ids: [...picked] });
			if (r.removed.length) toasts.push(t('library.upkeep.removed', { n: r.removed.length }), 'ok');
			if (r.skipped.length) toasts.push(t('library.upkeep.skipped', { n: r.skipped.length }), 'fail', { body: r.problems?.join('\n') });
			picked = new Set();
			await load();
		} catch (e) {
			toasts.push(t('library.upkeep.removeFailed'), 'fail', { body: message(e) });
		} finally {
			busy = false;
		}
	}

	function why(i: Item): string {
		return i.rule === 'no_file' ? t('library.upkeep.whyNoFile') : t('library.upkeep.whyNeverWatched', { date: date(i.added_at) });
	}
</script>

<svelte:head><title>{t('library.upkeep.title')} · HolyArrs</title></svelte:head>

<LibraryToolPage title={t('library.upkeep.title')} count={loaded && !error && items.length ? t('library.upkeep.count', { n: items.length }) : ''}>
	{#snippet bar()}
		<span class="tp-gap"></span>
		<span class="tp-cap">
			<Button variant="primary" icon="trash" onclick={remove} disabled={!picked.size} {busy} data-testid="upkeep-remove">
				{t('library.upkeep.remove')}{picked.size ? ` ${picked.size}` : ''}
			</Button>
		</span>
	{/snippet}

{#if error}
	<LoadError hasData={false} message={plainReason(error).text} onretry={load} />
{:else}
<PageLoader loading={!loaded} label={t('library.upkeep.loading')}>
	<form id="upkeep-form" class="tp-sheet rules" onsubmit={(e) => { e.preventDefault(); void save(); }} data-testid="upkeep-rules">
		<div class="rule">
			<label for="never-watched">{t('library.upkeep.neverWatched')}</label>
			<select id="never-watched" class="input" bind:value={settings.never_watched}>
				{#each modes as m (m)}<option value={m}>{t('library.upkeep.mode.' + m)}</option>{/each}
			</select>
			<label class="days">{t('library.upkeep.days')} <input class="input narrow" type="number" min="30" max="3650" bind:value={settings.never_watched_days} onchange={() => (settings.never_watched_days = Number(settings.never_watched_days))} /></label>
			<span class="tp-sub">{t('library.upkeep.neverWatchedHint')}</span>
			{#if !history}<span class="tp-sub warn" data-testid="upkeep-no-history">{t('library.upkeep.noHistory')}</span>{/if}
		</div>
		<div class="rule">
			<label for="no-file">{t('library.upkeep.noFile')}</label>
			<select id="no-file" class="input" bind:value={settings.no_file}>
				{#each modes as m (m)}<option value={m}>{t('library.upkeep.mode.' + m)}</option>{/each}
			</select>
			<span class="tp-sub">{t('library.upkeep.noFileHint')}</span>
		</div>
		<label class="paused"><input type="checkbox" bind:checked={settings.paused} /> {t('library.upkeep.paused')}</label>
		<span class="tp-sub">{t('library.upkeep.holdingNote')}</span>
	</form>
	<SaveBar {dirty} form="upkeep-form" {saving} ondiscard={() => (settings = JSON.parse(saved))} />

	{#if items.length === 0}
		<EmptyState icon="tools" title={t('library.upkeep.nothing')} />
	{:else}
	<div class="tp-sheet">
		<label class="all">
			<input type="checkbox" checked={all} onchange={() => (picked = all ? new Set() : new Set(items.map((i) => i.title_id)))} />
			{t('library.deleted.selectAll')}
		</label>
		<ul class="tp-rows">
			{#each items as i (i.title_id)}
				<li class="row-item" class:picked={picked.has(i.title_id)} data-testid="upkeep-item">
					<label class="pick"><input type="checkbox" checked={picked.has(i.title_id)} onchange={() => toggle(i.title_id)} aria-label={t('mediafiles.list.select', { name: i.name })} /></label>
					<span class="tp-art" style="--poster-ratio: 2 / 3"><Poster src={art.get(i.title_id)} name={i.name} /></span>
					<div class="grow">
						<a class="name" href={to(`/library/${i.title_id}`)}>{i.name}{i.year ? ` (${i.year})` : ''}</a>
						<span class="tp-sub">{why(i)}{i.size ? ` · ${bytes(i.size)}` : ''}</span>
					</div>
					{#if i.mode === 'auto'}<span class="auto">{t('library.upkeep.goesAuto')}</span>{/if}
				</li>
			{/each}
		</ul>
	</div>
	{/if}
</PageLoader>
{/if}
</LibraryToolPage>

<style>
	.rules { display: grid; gap: 14px; padding: 16px 20px; margin-bottom: 16px; }
	.rule { display: grid; gap: 6px; }
	.rule > label:first-child { font-weight: 600; }
	.days { display: inline-flex; align-items: center; gap: 8px; }
	.narrow { width: 7em; }
	.paused { display: flex; align-items: center; gap: 10px; min-height: 44px; }
	.warn { color: var(--warn-text, var(--fail-text)); }
	.all { display: flex; align-items: center; gap: 10px; min-height: 44px; padding: 0 20px; color: var(--text-muted); font: 600 11px/1 var(--font-ui); letter-spacing: .06em; text-transform: uppercase; background: color-mix(in srgb, var(--text) 3%, transparent); border-bottom: 1px solid var(--line); cursor: pointer; }
	input[type='checkbox'] { width: 20px; height: 20px; margin: 0; accent-color: var(--accent); }
	.row-item { display: flex; align-items: center; gap: 12px; }
	.row-item.picked { background: color-mix(in srgb, var(--accent) 8%, transparent); }
	.pick { display: inline-flex; align-items: center; justify-content: center; flex: none; min-width: 44px; min-height: 44px; margin: -12px; cursor: pointer; }
	.name { color: var(--text); font-weight: 600; text-decoration: none; overflow-wrap: anywhere; }
	a.name:hover { text-decoration: underline; text-underline-offset: 3px; }
	.grow { flex: 1; min-width: 0; font: var(--t-callout); font-variant-numeric: tabular-nums; }
	.auto { flex: none; font: var(--t-footnote); color: var(--text-muted); }
	@media (max-width: 639px) {
		.rules { padding: 16px; }
		.all { padding: 0 16px; }
		.row-item { flex-wrap: wrap; align-items: flex-start; }
		.auto { flex-basis: 100%; padding-left: 44px; }
	}
	@media (pointer: coarse) { a.name { display: inline-block; padding-block: 12px; margin-block: -12px; } }
</style>
