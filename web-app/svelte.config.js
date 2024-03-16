import { dirname, join } from 'path';

import adapter from '@sveltejs/adapter-static';
import { fileURLToPath } from 'url';
import { vitePreprocess } from '@sveltejs/kit/vite';

const __dirname = dirname(fileURLToPath(import.meta.url));

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: [
		vitePreprocess({
			style: {
				css: {
					postcss: join(__dirname, 'postcss.config.cjs')
				}
			}
		})
	],
	kit: {
		adapter: adapter({
			fallback: 'index.html'
		}),
		paths: {
			assets: 'https://assets.pubgolf.co'
		},
		prerender: {
			entries: [
				'/',
				'/event/nyc-2024',
				'/about/contact',
				'/about/privacy',
				'/admin',
				'/admin/nyc-2023',
				'/admin/nyc-2023/alerts',
				'/admin/nyc-2023/event-info',
				'/admin/nyc-2023/players',
				'/admin/nyc-2023/schedule',
				'/admin/nyc-2023/scores',
				'/admin/nyc-2024',
				'/admin/nyc-2024/alerts',
				'/admin/nyc-2024/event-info',
				'/admin/nyc-2024/players',
				'/admin/nyc-2024/schedule',
				'/admin/nyc-2024/scores',
				'/auth/login'
			]
		}
	}
};

export default config;
