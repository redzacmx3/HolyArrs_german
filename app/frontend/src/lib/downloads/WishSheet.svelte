<!-- Language and quality assistant: lays a small wish over the Pick rules for
every movie or show, one title, one season or one episode, with a preview of
what changes before it is saved. Empty choices keep what the level above says. -->
<script lang="ts">
	import Sheet from '#lib/ui/Sheet.svelte';
	import Button from '#lib/ui/Button.svelte';
	import { api, message } from '#lib/core/api.ts';
	import { toasts } from '#lib/core/toasts.svelte.ts';
	import { t } from '#lib/core/i18n/i18n.svelte.ts';
	import { languages, stepLabel } from './pick-rules.ts';
	import { wishBody, wishFrom, type LangChoice, type Wish, type WishRules } from './wish.ts';

	let {
		open = $bindable(false),
		kind,
		titleId,
		season,
		itemId,
		label,
		onsaved
	}: {
		open?: boolean;
		kind: 'movie' | 'show';
		/** Left out: the wish for every movie or show. */
		titleId?: number;
		season?: number;
		itemId?: number;
		/** What the wish is for, as the person reads it. */
		label: string;
		onsaved?: () => void;
	} = $props();

	type VersionOut = { version_id: number; label: string; before: WishRules; after: WishRules; files: number; good_enough_before: number; good_enough_after: number };
	type SetOut = { pick_rules_id: number; titles: number; before: WishRules; after: WishRules };
	type Preview = { versions?: VersionOut[]; sets?: SetOut[]; items?: number };

	let lang = $state<LangChoice>('keep');
	let code = $state('de');
	let quality = $state('');
	let loading = $state(false);
	let saving = $state(false);
	let preview = $state<Preview | null>(null);
	let error = $state('');
	let had = $state(false);

	const base = $derived(titleId ? `/titles/${titleId}/wishes` : `/settings/wishes/${kind}`);
	const wish = $derived<Wish>(wishBody(lang, code, quality));
	const body = $derived(titleId ? { season, item_id: itemId, wish } : wish);
	const qualities = [
		{ value: '', label: () => t('library.wish.qualityKeep') },
		{ value: '720p', label: () => '720p' },
		{ value: '1080p', label: () => '1080p' },
		{ value: '2160p', label: () => t('library.wish.quality4k') }
	];
	const langChoices: { value: LangChoice; label: () => string }[] = [
		{ value: 'keep', label: () => t('library.wish.langKeep') },
		{ value: 'with_original', label: () => t('library.wish.langWithOriginal') },
		{ value: 'only', label: () => t('library.wish.langOnly') },
		{ value: 'original', label: () => t('library.wish.langOriginal') }
	];

	// Reads the wish saved at this level each time the sheet opens.
	$effect(() => {
		if (!open) return;
		void load();
	});

	async function load() {
		loading = true;
		error = '';
		preview = null;
		try {
			let w: Wish | undefined;
			if (titleId) {
				const r = await api.get<{ wishes: { part_id?: number; item_id?: number; season?: number; wish: Wish }[] }>(base);
				w = r.wishes.find((x) => (itemId ? x.item_id === itemId : season !== undefined ? !x.item_id && x.season === season : !x.item_id && !x.part_id))?.wish;
			} else {
				w = (await api.get<{ wish: Wish }>(base)).wish;
			}
			const f = wishFrom(w);
			lang = f.lang;
			code = f.code || 'de';
			quality = f.quality;
			had = !!w && (w.languages != null || w.original_ok != null || !!w.quality);
		} catch (e) {
			error = message(e);
		} finally {
			loading = false;
		}
	}

	// The preview follows every change.
	$effect(() => {
		if (!open || loading) return;
		const sent = JSON.stringify(body);
		const timer = setTimeout(() => void look(sent), 250);
		return () => clearTimeout(timer);
	});

	async function look(sent: string) {
		try {
			preview = await api.post<Preview>(`${base}/preview`, JSON.parse(sent));
			error = '';
		} catch (e) {
			error = message(e);
		}
	}

	async function save(reset = false) {
		if (saving) return;
		saving = true;
		try {
			const send = titleId ? { season, item_id: itemId, wish: reset ? {} : wish } : reset ? {} : wish;
			await api.put(base, send);
			toasts.push(t(reset ? 'library.wish.resetDone' : 'library.wish.saved', { name: label }), 'ok');
			open = false;
			onsaved?.();
		} catch (e) {
			toasts.push(t('library.wish.saveFailed'), 'fail', { body: message(e) });
		} finally {
			saving = false;
		}
	}

	function steps(r: WishRules): string {
		return r.allowed.length ? r.allowed.map(stepLabel).join(', ') : t('library.wish.none');
	}
	function langs(r: WishRules): string {
		const named = r.languages.map((c) => languages.find(([k]) => k === c)?.[1] ?? c);
		if (named.length > 6) return t('library.wish.anyLanguage');
		const words = [...named, ...(r.original_ok ? [t('library.wish.original')] : [])];
		return words.length ? words.join(', ') : t('library.wish.none');
	}
	function same(a: WishRules, b: WishRules): boolean {
		return JSON.stringify(a) === JSON.stringify(b);
	}
</script>

<Sheet bind:open title={t('library.wish.title')} description={label} size="md">
	<form id="wish-form" class="wish" onsubmit={(e) => { e.preventDefault(); void save(); }} data-testid="wish-sheet">
		<p class="hint">{t(titleId ? 'library.wish.introTitle' : 'library.wish.introKind')}</p>
		<fieldset>
			<legend class="label">{t('library.wish.language')}</legend>
			{#each langChoices as c (c.value)}
				<label class="choice"><input type="radio" name="wish-lang" value={c.value} bind:group={lang} /> {c.label()}</label>
			{/each}
			{#if lang === 'with_original' || lang === 'only'}
				<select class="input" bind:value={code} aria-label={t('library.wish.whichLanguage')} data-testid="wish-code">
					{#each languages as [k, name] (k)}<option value={k}>{name}</option>{/each}
				</select>
			{/if}
		</fieldset>
		<fieldset>
			<legend class="label">{t('library.wish.quality')}</legend>
			<div class="seg">
				{#each qualities as q (q.value)}
					<label class="choice"><input type="radio" name="wish-quality" value={q.value} bind:group={quality} /> {q.label()}</label>
				{/each}
			</div>
		</fieldset>

		<section class="preview" aria-live="polite" data-testid="wish-preview">
			<h3 class="label">{t('library.wish.preview')}</h3>
			{#if error}
				<p class="fail" role="alert">{error}</p>
			{:else if !preview}
				<p class="hint">{t('library.wish.looking')}</p>
			{:else}
				{#each preview.versions ?? [] as v (v.version_id)}
					<div class="row">
						<b>{v.label} · {v.after.name}</b>
						<span>{t('library.wish.steps')}: {steps(v.after)}</span>
						<span>{t('library.wish.goodEnough')}: {stepLabel(v.after.good_enough_at)}</span>
						<span>{t('library.wish.languages')}: {langs(v.after)}</span>
						{#if v.files}
							<span class:warn={v.good_enough_after < v.good_enough_before}>{t('library.wish.files', { ok: v.good_enough_after, n: v.files, before: v.good_enough_before })}</span>
						{/if}
						{#if same(v.before, v.after)}<span class="hint">{t('library.wish.noChange')}</span>{/if}
					</div>
				{/each}
				{#each preview.sets ?? [] as s (s.pick_rules_id)}
					<div class="row">
						<b>{s.after.name} · {t('library.wish.titles', { n: s.titles })}</b>
						<span>{t('library.wish.steps')}: {steps(s.after)}</span>
						<span>{t('library.wish.goodEnough')}: {stepLabel(s.after.good_enough_at)}</span>
						<span>{t('library.wish.languages')}: {langs(s.after)}</span>
						{#if same(s.before, s.after)}<span class="hint">{t('library.wish.noChange')}</span>{/if}
					</div>
				{:else}
					{#if !preview.versions?.length}<p class="hint">{t('library.wish.nothingYet')}</p>{/if}
				{/each}
			{/if}
		</section>
	</form>
	{#snippet footer()}
		{#if had}<Button variant="plain" onclick={() => save(true)} disabled={saving}>{t('library.wish.reset')}</Button>{/if}
		<Button variant="plain" onclick={() => (open = false)}>{t('mediafiles.dialogs.cancel')}</Button>
		<Button variant="primary" type="submit" form="wish-form" busy={saving} disabled={loading}>{t('library.wish.save')}</Button>
	{/snippet}
</Sheet>

<style>
	.wish { display: flex; flex-direction: column; gap: var(--s-4); }
	fieldset { border: 0; margin: 0; padding: 0; display: flex; flex-direction: column; gap: var(--s-2); }
	.seg { display: flex; flex-wrap: wrap; gap: var(--s-3); }
	.choice { display: flex; align-items: center; gap: var(--s-2); }
	.label { font: var(--t-headline); margin: 0; }
	.hint { color: var(--text-2); font: var(--t-callout); margin: 0; }
	.preview { display: flex; flex-direction: column; gap: var(--s-3); }
	.row { display: flex; flex-direction: column; gap: 2px; font: var(--t-callout); }
	.warn { color: var(--warn); }
	.fail { color: var(--fail); margin: 0; }
</style>
