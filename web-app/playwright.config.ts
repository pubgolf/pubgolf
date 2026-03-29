import type { PlaywrightTestConfig } from '@playwright/test';

const config: PlaywrightTestConfig = {
	webServer: {
		command: 'npm run build && npm run preview',
		port: 4173,
		env: {
			SVELTE_ASSETS_PATH: ''
		}
	},
	testDir: 'e2e-tests'
};

export default config;
