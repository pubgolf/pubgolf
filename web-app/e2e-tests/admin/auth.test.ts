import { expect, test } from '@playwright/test';
import { ADMIN_BASE, mockRPC, setupAuth } from '../helpers/mock-api';

test.describe('admin auth guard', () => {
	test('unauthenticated visit redirects to login', async ({ page }) => {
		await page.route('**/web-api/auth/generate-api-token', async (route) => {
			await route.fulfill({
				status: 401,
				contentType: 'application/json',
				body: JSON.stringify({
					root_span_id: 'test',
					error_code: 'ERR_NOT_AUTHORIZED',
					error_message: 'User not logged in'
				})
			});
		});

		await page.goto(`${ADMIN_BASE}/scores/`);
		await page.waitForURL('**/auth/login**', { timeout: 10000 });
		expect(page.url()).toContain('/auth/login');
		expect(page.url()).toContain('redirect=');
	});

	test('server error shows error banner with retry', async ({ page }) => {
		let callCount = 0;

		await page.context().addCookies([
			{
				name: 'web_admin_user_token',
				value: 'test-cookie-token',
				domain: 'localhost',
				path: '/'
			}
		]);

		await page.route('**/web-api/auth/generate-api-token', async (route) => {
			callCount++;
			if (callCount === 1) {
				await route.fulfill({
					status: 500,
					contentType: 'application/json',
					body: JSON.stringify({
						root_span_id: 'test',
						error_code: 'ERR_INTERNAL',
						error_message: 'Internal server error'
					})
				});
			} else {
				await route.fulfill({
					status: 200,
					contentType: 'application/json',
					body: JSON.stringify({ token: 'test-api-token' })
				});
			}
		});

		await page.goto(`${ADMIN_BASE}/scores/`);

		const banner = page.locator('.alert.variant-filled-error');
		await expect(banner).toBeVisible();

		// Mock data RPCs that fire after successful auth retry
		await mockRPC(page, 'ListPlayers', { players: [] });
		await mockRPC(page, 'ListEventStages', { stages: [] });
		await mockRPC(page, 'ListStageScores', { scores: [] });

		await banner.getByRole('button', { name: 'Retry' }).click();

		await expect(page.getByRole('heading', { name: 'Scores' })).toBeVisible();
	});
});
