import { redirect } from '@sveltejs/kit';

export const load = () => {
	redirect(302, '/admin/nyc-2026/scores/');
};
