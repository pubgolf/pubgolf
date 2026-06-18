import { expect, test } from '@playwright/test';
import { makePlayer, makeStage, makeStageScore, makeVenue } from '../helpers/fixtures';
import { ADMIN_BASE, captureRPC, EVENT_KEY, mockRPC, setupAuth } from '../helpers/mock-api';

const player1 = makePlayer();
const player2 = makePlayer({
	id: 'player-2',
	data: { name: 'Bob Test' },
	events: [{ eventKey: EVENT_KEY, scoringCategory: 2 }]
});
const stage1 = makeStage();
const stage2 = makeStage({
	id: 'stage-2',
	venue: makeVenue({ id: 'venue-2', name: 'The Brewery', address: '456 Oak Ave' }),
	rank: 2
});

const score1 = makeStageScore();
const score2 = makeStageScore({
	stageId: 'stage-2',
	playerId: 'player-2',
	score: { id: 'score-2', data: { value: 5 } },
	adjustments: [{ id: 'adj-1', data: { value: 1, label: 'Wrong drink' } }],
	isVerified: true
});

test.describe('scores page', () => {
	test.beforeEach(async ({ page }) => {
		await setupAuth(page);
		await mockRPC(page, 'ListPlayers', { players: [player1, player2] });
		await mockRPC(page, 'ListEventStages', { stages: [stage1, stage2] });
		await mockRPC(page, 'ListStageScores', { scores: [score1, score2] });
	});

	test('table renders with player and venue names', async ({ page }) => {
		await page.goto(`${ADMIN_BASE}/scores/`);

		await expect(page.getByRole('heading', { name: 'Scores' })).toBeVisible();

		const table = page.locator('table');
		await expect(table).toBeVisible();

		await expect(table.getByText('1: The Anchor')).toBeVisible();
		await expect(table.getByText('2: The Brewery')).toBeVisible();

		await expect(table.getByText('Alice Test (Pro)')).toBeVisible();
		await expect(table.getByText('Bob Test (Semi-Pro)')).toBeVisible();

		await expect(table.getByText('Sips: 3')).toBeVisible();
		await expect(table.getByText('Sips: 5')).toBeVisible();

		await expect(table.getByText('+1: Wrong drink')).toBeVisible();
	});

	test('verify button calls updateStageScore', async ({ page }) => {
		const requests = await captureRPC(page, 'UpdateStageScore');

		await page.goto(`${ADMIN_BASE}/scores/`);
		await expect(page.locator('table')).toBeVisible();

		const verifyBtn = page.getByRole('button', { name: 'Verify' });
		await expect(verifyBtn).toHaveCount(1);
		await verifyBtn.click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as { score: { stageId: string; isVerified: boolean } };
		expect(req.score.stageId).toBe('stage-1');
	});

	test('delete score with confirmation dialog', async ({ page }) => {
		const requests = await captureRPC(page, 'DeleteStageScore');

		await page.goto(`${ADMIN_BASE}/scores/`);
		await expect(page.locator('table')).toBeVisible();

		await page.getByRole('button', { name: 'Delete' }).first().click();

		await expect(page.getByText('Confirm Deletion')).toBeVisible();
		await expect(page.getByText('Are you sure you wish to delete')).toBeVisible();

		await page.getByRole('button', { name: 'Confirm' }).click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as { stageId: string; playerId: string };
		// Scores are displayed in reverse order, so first Delete button is on the last score
		expect(req.stageId).toBe('stage-2');
		expect(req.playerId).toBe('player-2');
	});

	test('enter score modal submits createStageScore', async ({ page }) => {
		const requests = await captureRPC(page, 'CreateStageScore', {
			score: makeStageScore({ score: { id: 'new-score', data: { value: 4 } } })
		});

		await page.goto(`${ADMIN_BASE}/scores/`);
		await expect(page.getByRole('heading', { name: 'Scores' })).toBeVisible();

		await page.getByRole('button', { name: 'Enter Score' }).click();

		await expect(page.getByText('Enter a Score')).toBeVisible();

		await page.locator('select').first().selectOption('player-1');
		await page.locator('select').nth(1).selectOption('stage-1');
		await page.locator('input[type="number"]').first().fill('4');

		await page.getByRole('button', { name: 'Create Score' }).click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as { data: { score: { value: number } } };
		expect(req.data.score.value).toBe(4);
	});
});
