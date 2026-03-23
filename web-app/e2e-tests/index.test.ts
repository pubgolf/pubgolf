import { expect, test } from ‘@playwright/test’;

test(‘index page has expected h1’, async ({ page }) => {
	await page.goto(‘/’);
	await expect(page.getByRole(‘heading’, { name: ‘NYC Bottle Open ‘25’ })).toBeVisible();
});

test(‘login page renders’, async ({ page }) => {
	await page.goto(‘/auth/login/’);
	await expect(page.locator(‘form’)).toBeVisible();
	await expect(page.getByText(‘Admin Log In’)).toBeVisible();
});

test(‘privacy page renders’, async ({ page }) => {
	await page.goto(‘/about/privacy/’);
	await expect(page.getByRole(‘heading’, { name: ‘Privacy Policy’ })).toBeVisible();
});

test(‘no console errors on home page load’, async ({ page }) => {
	const consoleErrors: string[] = [];
	page.on(‘console’, (msg) => {
		if (msg.type() === ‘error’) {
			consoleErrors.push(msg.text());
		}
	});

	await page.goto(‘/’);
	expect(consoleErrors).toHaveLength(0);
});
