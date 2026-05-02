import { expect, test } from '@playwright/test';
import { makeAdjustmentTemplate, makeStage } from '../helpers/fixtures';
import { ADMIN_BASE, captureRPC, EVENT_KEY, mockRPC, setupAuth } from '../helpers/mock-api';

const stage1 = makeStage();
const template1 = makeAdjustmentTemplate();
const template2 = makeAdjustmentTemplate({
	id: 'template-2',
	data: {
		adjustment: { value: -2, label: 'Bonus: Speed round' },
		rank: 2,
		eventKey: EVENT_KEY,
		stageId: 'stage-1',
		isVisible: false
	}
});

test.describe('adjustments page', () => {
	test.beforeEach(async ({ page }) => {
		await setupAuth(page);
		await mockRPC(page, 'ListEventStages', { stages: [stage1] });
		await mockRPC(page, 'ListAdjustmentTemplates', { templates: [template1, template2] });
	});

	test('table renders with adjustment labels and values', async ({ page }) => {
		await page.goto(`${ADMIN_BASE}/adjustments/`);

		await expect(page.getByRole('heading', { name: 'Adjustments' })).toBeVisible();

		const table = page.locator('table');
		await expect(table).toBeVisible();

		await expect(table.getByText('Penalty: Wrong drink')).toBeVisible();
		await expect(table.getByText('Bonus: Speed round')).toBeVisible();

		await expect(page.getByRole('button', { name: 'Hide' })).toBeVisible();
		await expect(page.getByRole('button', { name: 'Show' })).toBeVisible();
	});

	test('create adjustment submits createAdjustmentTemplate', async ({ page }) => {
		const requests = await captureRPC(page, 'CreateAdjustmentTemplate', {
			template: makeAdjustmentTemplate({ id: 'template-new' })
		});

		await page.goto(`${ADMIN_BASE}/adjustments/`);
		await expect(page.getByRole('heading', { name: 'Adjustments' })).toBeVisible();

		await page.getByRole('button', { name: 'Create Adjustment' }).click();

		await expect(page.getByTestId('modal-component').locator('span.text-2xl')).toContainText(
			'Create Adjustment'
		);

		await page.getByPlaceholder('Enter label...').fill('Penalty: Spill');
		await page.getByPlaceholder('Enter value...').fill('2');
		await page.getByPlaceholder('Enter rank...').fill('3');

		await page
			.getByTestId('modal-component')
			.getByRole('button', { name: 'Create Adjustment' })
			.click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as {
			data: {
				adjustment: { value: number; label: string };
				rank: number;
				eventKey: string;
			};
		};
		expect(req.data.adjustment.label).toBe('Penalty: Spill');
		expect(req.data.adjustment.value).toBe(2);
		expect(req.data.rank).toBe(3);
	});

	test('toggle visibility shows confirm dialog and calls updateAdjustmentTemplate', async ({
		page
	}) => {
		const requests = await captureRPC(page, 'UpdateAdjustmentTemplate');

		await page.goto(`${ADMIN_BASE}/adjustments/`);
		await expect(page.locator('table')).toBeVisible();

		await page.getByRole('button', { name: 'Hide' }).click();

		await expect(page.getByText('Confirm Hide')).toBeVisible();
		await expect(page.getByText('Are you sure you wish to')).toBeVisible();

		await page.getByRole('button', { name: 'Confirm' }).click();

		expect(requests).toHaveLength(1);
		const req = requests[0] as {
			template: { id: string; data: { isVisible: boolean } };
		};
		expect(req.template.id).toBe('template-1');
		// Proto3 JSON omits false booleans, so isVisible will be undefined (not false)
		expect(req.template.data.isVisible).toBeFalsy();
	});
});
