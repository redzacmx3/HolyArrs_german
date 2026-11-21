import { describe, expect, it } from 'vitest';
import { pointsLines, signed } from './points-breakdown.ts';

describe('pointsLines', () => {
	it('lists the traits with their points, the costliest minus first', () => {
		const lines = pointsLines([
			{ code: 'points.trait:own:2', effect: 'add', text: 'German: fits (+10000 points)', points: 10000 },
			{ code: 'points.trait:x265', effect: 'add', text: 'video is HEVC (+1700 points)', points: 1700 },
			{ code: 'points.trait:own:7', effect: 'subtract', text: 'German Microsized: fits (-35000 points)', points: -35000 },
			{ code: 'points.trait:own:9', effect: 'subtract', text: 'German LQ: fits (-500 points)', points: -500 },
			{ code: 'points.below_min', effect: 'fail', text: '-23300 points, needs at least 11000' },
			{ code: 'hard.size', effect: 'pass', text: '0.37 GB within 0.10 to 2.00 GB' }
		]);
		expect(lines).toEqual([
			{ name: 'German Microsized', points: -35000 },
			{ name: 'German LQ', points: -500 },
			{ name: 'German', points: 10000 },
			{ name: 'video is HEVC', points: 1700 }
		]);
	});
	it('is empty without trait points', () => {
		expect(pointsLines([{ code: 'step.none', effect: 'fail', text: 'no ladder step: the source is unknown' }])).toEqual([]);
	});
});

describe('signed', () => {
	it('puts the sign in front', () => {
		expect(signed(-35000, 'de')).toBe('−35.000');
		expect(signed(11000, 'en')).toBe('+11,000');
	});
});
