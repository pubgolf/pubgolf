import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';
import tsconfigPaths from 'vite-tsconfig-paths';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: preprocess(),

	kit: {
		adapter: adapter(),
		target: '#svelte',
		// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
		vite: () => ({
			define: {
				'process.env.API_BASE': JSON.stringify(process.env.PUBGOLF_API_BASE)
			},
			// Needed to make Vite play nicely with import aliases (e.g. `import ... from '@rpc/pubgolf';`).
			plugins: [tsconfigPaths()],
			ssr: {
				// For some reason the generated client depending on `protobufjs/minimal` causes issues with either the build or dev command, depending on how we configure this. Luckily we can make it dynamic.
				noExternal: process.env.NODE_ENV === 'production' ? ['protobufjs'] : []
			}
		})
	}
};

export default config;
