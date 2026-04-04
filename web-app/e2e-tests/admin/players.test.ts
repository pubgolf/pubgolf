import { expect, test } from '@playwright/test';
import { makePlayer } from '../helpers/fixtures';
import { ADMIN_BASE, captureRPC, EVENT_KEY, mockRPC, setupAuth } from '../helpers/mock-api';

const player1 = makePlayer();
const player2 = makePlayer({
	id: 'player-2',
	data: { name: 'Bob Test' },
	events: [{ eventKey: EVENT_KEY, scoringCategory: 2 }]
});

test.describe('players page', () => {
	test.beforeEach(async ({ page }) => {
		await setupAuth(page);
		await mockRPC(page, 'ListPlayers', { players: [player1, player2] });
	});

	test('table renders with player names and leagues', async ({ page }) => {
		await page.goto(`${ADMIN_BASE}/players/`);

		await expect(page.getByRole('heading', { name: 'Players' })).toBeVisible();

		const table = page.locator('table');
		await expect(table).toBeVisible();

		await expect(table.getByText('Alice Test')).toBeVisible();
		await expect(table.getByText('Bob Test')).toBeVisible();

		// Use exact match to avoid "Pro" matching "Semi-Pro"
		await expect(table.getByText('Pro', { exact: true })).toBeVisible();
		await expect(table.getByText('Semi-Pro', { exact: true })).toBeVisible();
	});

	test('register new player submits createPlayer', async ({ page }) => {
		const requests = await captureRPC(page, 'CreatePlayer', {
			player: makePlayer({ id: 'player-new' })
		});

		await page.goto(`${ADMIN_BASE}/players/`);
		await expect(page.getByRole('heading', { name: 'Players' })).toBeVisible();

		await page.getByRole('button', { name: 'New Player' }).click();

		await expect(page.getByText('Register New Player')).toBeVisible();

		await page.getByPlaceholder('Enter name...').fill('Charlie Test');
		await page.locator('input[type="tel"]').fill('5551234567');
		await page.locator('select').selectOption('1');

		await page.getByRole('button', { name: 'Register Player' }).click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as {
			playerData: { name: string };
			phoneNumber: string;
			registration: { eventKey: string; scoringCategory: number };
		};
		expect(req.playerData.name).toBe('Charlie Test');
		expect(req.phoneNumber).toContain('5551234567');
		expect(req.registration.eventKey).toBe(EVENT_KEY);
	});

	test('edit player submits updatePlayer', async ({ page }) => {
		const requests = await captureRPC(page, 'UpdatePlayer', {
			player: makePlayer()
		});

		await page.goto(`${ADMIN_BASE}/players/`);
		await expect(page.locator('table')).toBeVisible();

		await page.getByRole('button', { name: 'Edit' }).first().click();

		await expect(page.locator('span.text-2xl')).toContainText('Update Player');

		await page.locator('select').selectOption('3');

		await page.getByRole('button', { name: 'Update Player' }).click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as {
			playerId: string;
			registration: { scoringCategory: string | number };
		};
		expect(req.playerId).toBe('player-1');
		// Proto3 JSON serializes enums as string names
		expect(req.registration.scoringCategory).toBe('SCORING_CATEGORY_PUB_GOLF_CHALLENGES');
	});
});
