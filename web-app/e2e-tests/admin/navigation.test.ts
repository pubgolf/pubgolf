import { expect, test } from '@playwright/test';
import { makePlayer, makeStage, makeStageScore } from '../helpers/fixtures';
import { ADMIN_BASE, mockRPC, mockRPCError, setupAuth } from '../helpers/mock-api';

test.describe('drawer navigation', () => {
	test.beforeEach(async ({ page }) => {
		await setupAuth(page);
		await mockRPC(page, 'ListPlayers', { players: [] });
		await mockRPC(page, 'ListEventStages', { stages: [] });
		await mockRPC(page, 'ListStageScores', { scores: [] });
		await mockRPC(page, 'ListVenues', { venues: [] });
		await mockRPC(page, 'ListAdjustmentTemplates', { templates: [] });
	});

	test('mobile drawer opens, navigates, and closes', async ({ page }) => {
		await page.setViewportSize({ width: 375, height: 667 });

		await page.goto(`${ADMIN_BASE}/scores/`);
		await expect(page.getByRole('heading', { name: 'Scores' })).toBeVisible();

		await page
			.locator('button')
			.filter({ has: page.locator('svg') })
			.first()
			.click();

		const drawer = page.locator('.drawer-backdrop');
		await expect(drawer).toBeVisible();

		await drawer.getByText('Players').click();

		await expect(page).toHaveURL(new RegExp(`${ADMIN_BASE}/players/`));
		await expect(page.getByRole('heading', { name: 'Players' })).toBeVisible();
	});
});

test.describe('error handling', () => {
	test.beforeEach(async ({ page }) => {
		await setupAuth(page);
	});

	test('toast appears on RPC error during action', async ({ page }) => {
		await mockRPC(page, 'ListPlayers', { players: [makePlayer()] });
		await mockRPC(page, 'ListEventStages', { stages: [makeStage()] });
		await mockRPC(page, 'ListStageScores', { scores: [makeStageScore()] });
		await mockRPCError(page, 'DeleteStageScore', 'internal', 'Something went wrong');

		await page.goto(`${ADMIN_BASE}/scores/`);
		await expect(page.locator('table')).toBeVisible();

		await page.getByRole('button', { name: 'Delete' }).click();
		await page.getByRole('button', { name: 'Confirm' }).click();

		await expect(page.locator('.toast-message, [class*="toast"]').first()).toBeVisible();
	});

	test('error banner shows on data fetch failure', async ({ page }) => {
		await mockRPCError(page, 'ListPlayers', 'internal', 'Database error');
		await mockRPC(page, 'ListEventStages', { stages: [] });
		await mockRPC(page, 'ListStageScores', { scores: [] });

		await page.goto(`${ADMIN_BASE}/scores/`);

		const banner = page.locator('.alert.variant-filled-error');
		await expect(banner).toBeVisible();
		await expect(banner.getByRole('button', { name: 'Retry' })).toBeVisible();
	});
});
