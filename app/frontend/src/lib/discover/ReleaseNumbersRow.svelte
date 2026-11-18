<script lang="ts">
	// Release numbers row on a show's Library tab, shown while the show follows
	// a chosen episode order: release names count by that order or by TMDB's
	// own numbering. Files keep the order's numbers either way.
	import { t } from '#lib/core/i18n/i18n.svelte.ts';
	import Icon from '#lib/ui/Icon.svelte';
	import Menu, { type MenuItem } from '#lib/ui/Menu.svelte';

	let { value = '', busy = false, onpick }: { value?: string; busy?: boolean; onpick: (v: string) => void } = $props();

	const words = (v: string) => (v === 'tmdb' ? t('discover.titleLibrary.releaseNumbersTmdb') : t('discover.titleLibrary.releaseNumbersOrder'));
	const menu = $derived<MenuItem[]>(['', 'tmdb'].map((v) => ({ label: words(v), checked: value === v, onselect: () => value !== v && onpick(v) })));
</script>

<div class="lib-row static" data-testid="release-numbers-row">
	<span><b>{t('discover.titleLibrary.releaseNumbers')}</b><small class="wrap">{t('discover.titleLibrary.releaseNumbersHint')}</small></span>
	<div class="lib-pick"><Menu label={t('discover.titleLibrary.releaseNumbers')} variant="button" items={menu} width={240}>{#snippet trigger()}<span>{busy ? t('discover.titleLibrary.saving') : words(value)}</span><Icon name="chevronDown" size={14} />{/snippet}</Menu></div>
</div>
