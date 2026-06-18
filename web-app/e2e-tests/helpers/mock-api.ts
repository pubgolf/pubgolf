import type { Page } from '@playwright/test';

const EVENT_KEY = 'nyc-2025';
const ADMIN_BASE = `/admin/${EVENT_KEY}`;
const RPC_BASE = '**/rpc/api.v1.AdminService/';

export { EVENT_KEY, ADMIN_BASE };

/**
 * Sets up auth for admin pages: injects cookie and mocks the token endpoint.
 */
export async function setupAuth(page: Page) {
	await page.context().addCookies([
		{
			name: 'web_admin_user_token',
			value: 'test-cookie-token',
			domain: 'localhost',
			path: '/'
		}
	]);

	await page.route('**/web-api/auth/generate-api-token', async (route) => {
		await route.fulfill({
			status: 200,
			contentType: 'application/json',
			body: JSON.stringify({ token: 'test-api-token' })
		});
	});
}

/**
 * Mocks a Connect-RPC AdminService method to return the given response.
 */
export async function mockRPC(page: Page, method: string, response: unknown) {
	await page.route(`${RPC_BASE}${method}`, async (route) => {
		await route.fulfill({
			status: 200,
			contentType: 'application/json',
			body: JSON.stringify(response)
		});
	});
}

/**
 * Mocks a Connect-RPC AdminService method to return an error.
 */
export async function mockRPCError(
	page: Page,
	method: string,
	code: string,
	message: string,
	status = 500
) {
	await page.route(`${RPC_BASE}${method}`, async (route) => {
		await route.fulfill({
			status,
			contentType: 'application/json',
			body: JSON.stringify({ code, message })
		});
	});
}

/**
 * Intercepts a Connect-RPC AdminService method, captures request bodies,
 * and optionally fulfills a response. Returns the captured requests array.
 */
export async function captureRPC(
	page: Page,
	method: string,
	response: unknown = {}
): Promise<unknown[]> {
	const requests: unknown[] = [];

	await page.route(`${RPC_BASE}${method}`, async (route) => {
		requests.push(route.request().postDataJSON());
		await route.fulfill({
			status: 200,
			contentType: 'application/json',
			body: JSON.stringify(response)
		});
	});

	return requests;
}
