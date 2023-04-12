import { redirect } from '@sveltejs/kit';

export const load = () => {
	throw redirect(302, '/admin/nyc-2023/scores/');
};
