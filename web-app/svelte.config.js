import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/kit/vite';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: [vitePreprocess()],
	kit: {
		adapter: adapter({
			fallback: 'index.html'
		}),
        paths: {
            assets: 'https://assets.pubgolf.co',
        },
		prerender: {
			entries: [
				'/',
				'/event/nyc-2023',
                '/about/contact',
                '/about/privacy',
				'/admin',
				'/admin/nyc-2023',
				'/admin/nyc-2023/alerts',
				'/admin/nyc-2023/players',
				'/admin/nyc-2023/schedule',
				'/admin/nyc-2023/scores'
			]
		}
	}
};

export default config;
