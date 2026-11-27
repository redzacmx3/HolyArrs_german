// Fork fields on the upstream types, kept out of types.ts so it stays as
// upstream has it (file size cap, fewer conflicts on each sync).
export {};

declare module '#lib/core/types.ts' {
	interface ShowInfo {
		// How release names number episodes: '' by the order, 'tmdb' (chosen order only), 'tvmaze' (like TVDB).
		release_numbers?: string;
	}
	interface ReleaseView {
		// What the name says about its languages (sources/release_tags.go): audio languages, DL (the original too), subtitle languages.
		tags?: { languages?: string[]; dual?: boolean; subs?: string[] };
	}
	interface DownloadRow {
		// A row waiting on a decision: the few words saying what is wrong (decision-fix.ts shortWhy).
		decide_why?: string;
	}
}
