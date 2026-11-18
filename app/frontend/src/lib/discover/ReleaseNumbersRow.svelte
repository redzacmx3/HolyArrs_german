<script lang="ts">
	// Release numbers row on a show's Library tab: release names count by the
	// show's order, by TMDB's own numbering (only while a chosen episode order
	// is set) or by TVmaze's, which follows TVDB. Files keep the order's
	// numbers either way.
	import { t } from '#lib/core/i18n/i18n.svelte.ts';
	import Icon from '#lib/ui/Icon.svelte';
	import Menu, { type MenuItem } from '#lib/ui/Menu.svelte';

	let { value = '', grouped = false, busy = false, onpick }: { value?: string; grouped?: boolean; busy?: boolean; onpick: (v: string) => void } = $props();

	const words = (v: string) =>
		v === 'tmdb' ? t('discover.titleLibrary.releaseNumbersTmdb') : v === 'tvmaze' ? t('discover.titleLibrary.releaseNumbersTvmaze') : t('discover.titleLibrary.releaseNumbersOrder');
	const menu = $derived<MenuItem[]>((grouped ? ['', 'tmdb', 'tvmaze'] : ['', 'tvmaze']).map((v) => ({ label: words(v), checked: value === v, onselect: () => value !== v && onpick(v) })));
</script>

<div class="lib-row static" data-testid="release-numbers-row">
	<span><b>{t('discover.titleLibrary.releaseNumbers')}</b><small class="wrap">{t('discover.titleLibrary.releaseNumbersHint')}</small></span>
	<div class="lib-pick"><Menu label={t('discover.titleLibrary.releaseNumbers')} variant="button" items={menu} width={240}>{#snippet trigger()}<span>{busy ? t('discover.titleLibrary.saving') : words(value)}</span><Icon name="chevronDown" size={14} />{/snippet}</Menu></div>
</div>
