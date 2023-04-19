import type { DisplayError } from '$lib/components/ErrorBanner.svelte';

export const USER_NOT_AUTHORIZED_ERROR = 'Unauthorized Error';

type APIErrorResponse = {
	root_span_id: string;
	error_code: string;
	error_message: string;
};

export async function getAPIToken(): Promise<{ token: string; error: DisplayError }> {
	const { resp, error } = await makeAPIRequest('/web-api/auth/generate-api-token', {});

	if (error) {
		const errorResp = resp as APIErrorResponse | null;

		if (errorResp?.error_code === 'ERR_NOT_AUTHORIZED') {
			return {
				token: '',
				error: {
					message: errorResp?.error_message || 'User not logged in',
					type: USER_NOT_AUTHORIZED_ERROR
				}
			};
		}

		return { token: '', error };
	}

	const tokenResp = resp as { token: string };
	if (!tokenResp.token) {
		console.error('Response from /web-api/auth/generate-api-token did not include a token', {
			resp
		});
		return {
			token: '',
			error: {
				type: 'Server Error',
				message: 'Response did not include a token'
			}
		};
	}

	return { token: tokenResp.token, error: null };
}

export async function logIn(password: string): Promise<DisplayError> {
	const { resp, error } = await makeAPIRequest('/web-api/auth/login', { password });

	if (error) {
		return error;
	}

	const successResp = resp as { success: boolean };
	if (!successResp.success) {
		console.log('Response from /web-api/auth/login did not indicate success=true', { resp });
		return {
			type: 'Server Error',
			message: 'Login was unsuccessful'
		};
	}

	return null;
}

export async function logOut(): Promise<DisplayError> {
	const { resp, error } = await makeAPIRequest('/web-api/auth/logout', {});

	if (error) {
		return error;
	}

	const successResp = resp as { success: boolean };
	if (!successResp.success) {
		console.log('Response from /web-api/auth/logout did not indicate success=true', { resp });
		return {
			type: 'Server Error',
			message: 'Logout was unsuccessful'
		};
	}

	return null;
}

async function makeAPIRequest(
	url: string,
	bodyData: unknown
): Promise<{ resp: unknown; error: DisplayError }> {
	let resp;
	try {
		resp = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json'
			},
			body: JSON.stringify(bodyData)
		});
	} catch (error) {
		console.error('Error fetching ${url} endpoint:', { error });
		return {
			resp: null,
			error: {
				type: 'Application Error',
				message: 'Error connecting to API'
			}
		};
	}

	if (resp.status !== 200) {
		return parseErrorJSON(resp);
	}

	let respJSON;
	try {
		respJSON = await resp.json();
	} catch (error) {
		console.error('Error parsing API response as JSON:', { resp, error });
		return {
			resp: respJSON,
			error: {
				type: 'Application Error',
				message: 'Error parsing API response as JSON'
			}
		};
	}

	return { resp: respJSON, error: null };
}

async function parseErrorJSON(resp: Response): Promise<{ resp: unknown; error: DisplayError }> {
	let respJSON;
	try {
		respJSON = await resp.json();
	} catch (error) {
		console.error('Error parsing API response as JSON:', { resp, error });
		return {
			resp: null,
			error: {
				type: 'Application Error',
				message: 'Error parsing API response as JSON'
			}
		};
	}

	return {
		resp: respJSON,
		error: {
			type: 'Server Error',
			message: respJSON.error_message || 'Unknown error'
		}
	};
}
