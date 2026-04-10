import { expect, test } from '@playwright/test';
import { makeStage, makeVenue } from '../helpers/fixtures';
import { ADMIN_BASE, captureRPC, mockRPC, setupAuth } from '../helpers/mock-api';

const venue1 = makeVenue();
const venue2 = makeVenue({ id: 'venue-2', name: 'The Brewery', address: '456 Oak Ave' });
const stage1 = makeStage();
const stage2 = makeStage({
	id: 'stage-2',
	venue: venue2,
	rule: { id: 'rule-2', venueDescription: 'Par 5 - finish the pint' },
	rank: 2,
	durationMin: 60
});

test.describe('schedule page', () => {
	test.beforeEach(async ({ page }) => {
		await setupAuth(page);
		await mockRPC(page, 'ListEventStages', { stages: [stage1, stage2] });
		await mockRPC(page, 'ListVenues', { venues: [venue1, venue2] });
	});

	test('table renders with venue names, durations, and rules', async ({ page }) => {
		await page.goto(`${ADMIN_BASE}/schedule/`);

		await expect(page.getByRole('heading', { name: 'Schedule' })).toBeVisible();

		const table = page.locator('table');
		await expect(table).toBeVisible();

		await expect(table.getByText('The Anchor')).toBeVisible();
		await expect(table.getByText('The Brewery')).toBeVisible();

		await expect(table.getByText('45')).toBeVisible();
		await expect(table.getByText('60')).toBeVisible();

		await expect(table.getByText('Par 3 - drink in 3 sips')).toBeVisible();
		await expect(table.getByText('Par 5 - finish the pint')).toBeVisible();
	});

	test('edit stage submits updateStage', async ({ page }) => {
		const requests = await captureRPC(page, 'UpdateStage');

		await page.goto(`${ADMIN_BASE}/schedule/`);
		await expect(page.locator('table')).toBeVisible();

		await page.getByRole('button', { name: 'Edit' }).first().click();

		await expect(page.locator('span.text-2xl')).toContainText('Update Stage');

		await page.locator('input[type="number"]').fill('30');

		await page.getByRole('button', { name: 'Update Stage' }).click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as { stageId: string; durationMin: number };
		expect(req.stageId).toBe('stage-1');
		expect(req.durationMin).toBe(30);
	});
});
