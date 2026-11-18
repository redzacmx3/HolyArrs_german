<script lang="ts">
	import Words from '#lib/ui/Words.svelte';
	import { t } from '#lib/core/i18n/i18n.svelte.ts';
	// Selection bars for ticked downloads: a toolbar above the list on wide screens, a floating sheet on phones.
	import Icon from '#lib/ui/Icon.svelte';
	import type { BulkAction, SteerMode } from './manage.ts';
	import { retryable } from './usenet.ts';
	import type { Decision, DownloadRow } from '../core/types.ts';
	import { api, message } from '../core/api.ts';
	import { toasts } from '../core/toasts.svelte.ts';
	import { ask } from '../ui/ask.svelte.ts';
	import { refreshDecisions } from './decisions.svelte.ts';

	type Targets = { rows: DownloadRow[] };
	let {
		pickedRows,
		bulkRows,
		steer,
		startT,
		pauseT,
		batchBusy = false,
		steerable,
		mayEdit,
		runOn,
		onbulk,
		onsteer,
		onretry,
		onremovemany,
		onall,
		onclear,
		decisionsOf,
		ondecided
	}: {
		pickedRows: DownloadRow[];
		/** Ticked torrents a torrent Downloader can steer by hash. */
		bulkRows: DownloadRow[];
		steer: { mode: SteerMode; label: string; short: string }[];
		startT: Targets;
		pauseT: Targets;
		batchBusy?: boolean;
		steerable: (r: DownloadRow) => boolean;
		mayEdit: (r: DownloadRow) => boolean;
		runOn: (t: Targets, one: 'resume' | 'pause', many: BulkAction) => void;
		onbulk?: (rows: DownloadRow[], action: BulkAction) => void;
		onsteer?: (rows: DownloadRow[], mode: SteerMode) => void;
		onretry?: (rows: DownloadRow[]) => void;
		onremovemany?: (rows: DownloadRow[], withFiles?: boolean) => void;
		onall: () => void;
		onclear: () => void;
		/** The decisions a ticked row waits on (Prüfen): the bulk fixes below settle them all at once. */
		decisionsOf?: (r: DownloadRow) => Decision[];
		ondecided?: () => void;
	} = $props();

	// Every decision behind the ticked rows, once each.
	const picks = $derived.by(() => {
		const m = new Map<number, Decision>();
		for (const r of pickedRows) for (const d of decisionsOf?.(r) ?? []) m.set(d.slot_id, d);
		return [...m.values()];
	});
	const withCopy = $derived(picks.filter((d) => d.has_copy));
	type Fix = 'search' | 'turn_down' | 'keep' | 'stop';
	let deciding = $state(false);
	// One call per decision, a few at a time; a row that settled meanwhile only counts as failed.
	async function decideAll(fix: Fix) {
		const list = fix === 'keep' ? withCopy : picks;
		if (!list.length || deciding) return;
		if (fix === 'stop' && !(await ask(t('downloads.bulkFix.askStop', { n: list.length }), { yes: t('ui.ask.stop') }))) return;
		if (fix === 'turn_down' && !(await ask(t('downloads.bulkFix.askTurnDown', { n: list.length })))) return;
		deciding = true;
		let ok = 0;
		let lastErr = '';
		const queue = [...list];
		const one = async () => {
			for (let d = queue.shift(); d; d = queue.shift()) {
				// A decision without a Release to turn down is searched again.
				const action = fix === 'keep' ? 'stop' : fix === 'turn_down' && !d.release_id ? 'search' : fix;
				try {
					await api.post(`/needs-decision/${d.slot_id}/decide`, { action });
					ok++;
				} catch (e) {
					lastErr = message(e);
				}
			}
		};
		await Promise.all([one(), one(), one(), one()]);
		deciding = false;
		const fail = list.length - ok;
		toasts.push(t('downloads.bulkFix.done', { n: ok }), fail ? 'fail' : 'ok', fail ? { body: t('downloads.bulkFix.failed', { n: fail, why: lastErr }) } : {});
		await refreshDecisions();
		ondecided?.();
		onclear();
	}
	const FIXES = $derived<{ fix: Fix; label: string; hint: string }[]>([
		{ fix: 'turn_down', label: t('downloads.bulkFix.turnDown'), hint: t('downloads.bulkFix.turnDownHint') },
		{ fix: 'search', label: t('downloads.dp.searchAgain'), hint: '' },
		...(withCopy.length ? [{ fix: 'keep' as const, label: t('downloads.bulkFix.keep', { n: withCopy.length }), hint: t('downloads.bulkFix.keepHint') }] : []),
		{ fix: 'stop', label: t('downloads.dp.stopLooking'), hint: '' }
	]);

	const PHONE_BULK = $derived<{ action: BulkAction; label: string; icon: 'play' | 'pause' | 'refresh' }[]>([{ action: 'start', label: t('downloads.act.start'), icon: 'play' }, { action: 'pause', label: t('downloads.act.pause'), icon: 'pause' }, { action: 'recheck', label: t('downloads.act.checkAgain'), icon: 'refresh' }]);
	const BULK = $derived<{ action: BulkAction; label: string }[]>([
		{ action: 'recheck', label: t('downloads.act.checkAgain') }, { action: 'reannounce', label: t('downloads.act.askTrackers') }, { action: 'top', label: t('downloads.act.top') },
		{ action: 'up', label: t('downloads.act.up') }, { action: 'down', label: t('downloads.act.down') }, { action: 'bottom', label: t('downloads.act.bottom') }, { action: 'force', label: t('downloads.act.forceStart') }
	]);
</script>

	<div class="selbar" role="toolbar" aria-label={t('downloads.sel.label')} data-testid="dl-selbar">
		<span class="sn" aria-live="polite"><Words key="downloads.sel.selectedCount" vars={{ n: pickedRows.length }} /></span>
		{#if picks.length}<span class="fx" data-testid="dl-bulk-fixes"><span class="fxl">{t('downloads.bulkFix.label', { n: picks.length })}</span>{#each FIXES as f (f.fix)}<button type="button" class="xbn" class:lead={f.fix === 'turn_down'} class:bad={f.fix === 'stop'} title={f.hint || undefined} disabled={deciding} onclick={() => decideAll(f.fix)} data-testid="dl-bulk-fix-{f.fix}">{f.label}</button>{/each}</span>{/if}
		{#if bulkRows.length}{#each BULK as b (b.action)}<button type="button" class="xbn" disabled={batchBusy} onclick={() => onbulk?.(bulkRows, b.action)} data-testid="dl-bulk-{b.action}">{b.label}</button>{/each}{/if}
		{#if pickedRows.some(steerable)}{#each steer as m (m.mode)}<button type="button" class="xbn" onclick={() => onsteer?.(pickedRows.filter(steerable), m.mode)} data-testid="dl-bulk-steer-{m.mode}">{m.short}</button>{/each}{/if}
		{#if pickedRows.some((r) => retryable(r) && mayEdit(r))}<button type="button" class="xbn" disabled={batchBusy} onclick={() => onretry?.(pickedRows.filter((r) => retryable(r) && mayEdit(r)))} data-testid="dl-bulk-retry">{t('downloads.act.again')}</button>{/if}
		<button type="button" class="xbn" onclick={onall} data-testid="dl-bulk-all">{t('downloads.sel.all')}</button>
		<button type="button" class="xbn" onclick={onclear} data-testid="dl-bulk-clear">{t('downloads.sel.clear')}</button>
	</div>
		<div class="psel" role="toolbar" aria-label={t('downloads.sel.label')} data-testid="dl-psel">
			<div class="pst"><span class="sn" aria-live="polite"><Words key="downloads.sel.selectedCount" vars={{ n: pickedRows.length }} /></span><button type="button" class="plk" onclick={onall} data-testid="dl-psel-all">{t('downloads.sel.all')}</button><button type="button" class="plk" onclick={onclear} data-testid="dl-psel-cancel">{t('downloads.sel.cancel')}</button></div>
			{#if picks.length}<div class="pfx">{#each FIXES as f (f.fix)}<button type="button" class="xbn" class:lead={f.fix === 'turn_down'} class:bad={f.fix === 'stop'} disabled={deciding} onclick={() => decideAll(f.fix)} data-testid="dl-psel-fix-{f.fix}">{f.label}</button>{/each}</div>{/if}
			<div class="psa">
				{#each PHONE_BULK.filter((b) => b.action !== 'recheck' || bulkRows.length) as b (b.action)}<button type="button" class="pa" disabled={batchBusy || (b.action === 'start' ? !startT.rows.length : b.action === 'pause' ? !pauseT.rows.length : false)} onclick={() => (b.action === 'start' ? runOn(startT, 'resume', 'start') : b.action === 'pause' ? runOn(pauseT, 'pause', 'pause') : onbulk?.(bulkRows, b.action))} data-testid="dl-psel-{b.action}"><Icon name={b.icon} size={20} /><span>{b.label}</span></button>{/each}
				<button type="button" class="pa bad" onclick={() => onremovemany?.(pickedRows)} data-testid="dl-psel-remove"><Icon name="trash" size={20} /><span>{t('downloads.act.remove')}</span></button>
			</div>
		</div>

<style>
	.selbar { position: sticky; top: 0; left: 0; z-index: 3; display: flex; flex-wrap: wrap; align-items: center; gap: 6px; padding: 8px var(--gutter, 12px); background: var(--surface-1); border-bottom: 1px solid var(--line); }
	.selbar .sn { font: var(--t-footnote); color: var(--text-muted); margin-right: 4px; }
	.selbar .sn :global(b) { color: var(--text); }
	.psel { display: none; }
	@keyframes pselIn { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: none; } }
	@media (prefers-reduced-motion: reduce) {
		.psel { animation: none; }
	}
	.xbn:focus-visible { outline: none; box-shadow: var(--focus-ring); }
	.xbn { min-height: 32px; padding: 0 14px; border: 1px solid var(--surface-3); border-radius: var(--r-btn); background: var(--surface-2); color: var(--text); font: var(--t-footnote); font-weight: 500; cursor: pointer; }
	.xbn:hover:not(:disabled) { background: var(--surface-3); }
	.xbn:disabled { opacity: 0.5; cursor: not-allowed; }
	.xbn.lead { border-color: transparent; background: var(--text); color: var(--bg); }
	.xbn.lead:hover:not(:disabled) { background: color-mix(in srgb, var(--text) 85%, transparent); }
	.xbn.bad { color: var(--fail-text); }
	.fx { display: inline-flex; flex-wrap: wrap; align-items: center; gap: 6px; padding-right: 8px; margin-right: 4px; border-right: 1px solid var(--line); }
	.fxl { font: var(--t-footnote); font-weight: 600; color: var(--text); }
	.pfx { display: flex; flex-wrap: wrap; gap: 4px; }
	@media (max-width: 640px) {
		.selbar { display: none; }
		.psel { position: fixed; left: 8px; right: 8px; bottom: calc(var(--tabbar, 56px) + env(safe-area-inset-bottom) + 8px); z-index: 45; display: grid; gap: 4px; padding: 6px 8px 8px; border: 1px solid var(--line-strong); border-radius: 18px; background: color-mix(in srgb, var(--surface-1) 97%, transparent); backdrop-filter: blur(20px) saturate(1.4); -webkit-backdrop-filter: blur(20px) saturate(1.4); box-shadow: var(--shadow-modal); animation: pselIn 220ms var(--ease-out); }
		.pst { display: flex; align-items: center; gap: 4px; }
		.pst .sn { flex: 1; padding-left: 8px; color: var(--text-muted); font: 400 13px/1 var(--font-ui); }
		.pst .sn :global(b) { color: var(--text); font-weight: 600; }
		.plk { min-height: 44px; padding: 0 12px; border: 0; border-radius: 10px; background: none; color: var(--accent-text, var(--text)); font: 600 14px/1 var(--font-ui); cursor: pointer; }
		.psa { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 4px; }
		.pa { display: grid; justify-items: center; align-content: center; gap: 4px; min-height: 52px; padding: 4px 2px; border: 0; border-radius: 12px; background: var(--surface-2); color: var(--text); font: 500 11px/1.1 var(--font-ui); cursor: pointer; -webkit-tap-highlight-color: transparent; }
		.pa span { max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
		.pa:active, .plk:active { transform: scale(0.97); }
		.pa.bad { color: var(--fail-text); }
		.pa:disabled { opacity: 0.45; }
		.pa:focus-visible, .plk:focus-visible { outline: none; box-shadow: var(--focus-ring); }
		.xbn { min-height: 32px; }
	}
</style>
