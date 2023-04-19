import { AdminService } from '$lib/proto/api/v1/admin_connect';
import { PubGolfService } from '$lib/proto/api/v1/pubgolf_connect';
import { createConnectTransport } from '@bufbuild/connect-web';
import { createPromiseClient, type Interceptor, type PromiseClient } from '@bufbuild/connect';
import { getAPIToken, USER_NOT_AUTHORIZED_ERROR } from '$lib/auth/client';

export const PubGolfServiceClient = createPromiseClient(
	PubGolfService,
	createConnectTransport({
		baseUrl: '/rpc/'
	})
);

export type AdminServiceClient = PromiseClient<typeof AdminService>;

let adminToken = '';
const adminServiceClients: { [authToken: string]: AdminServiceClient } = {};

export async function getAdminServiceClient(): Promise<AdminServiceClient> {
	if (!adminToken) {
		const tokenResp = await getAPIToken();

		if (tokenResp.error) {
			if (tokenResp.error.type === USER_NOT_AUTHORIZED_ERROR) {
				throw USER_NOT_AUTHORIZED_ERROR;
			}
			throw tokenResp.error;
		}

		adminToken = tokenResp.token;
	}

	if (!adminServiceClients[adminToken]) {
		adminServiceClients[adminToken] = createPromiseClient(
			AdminService,
			createConnectTransport({
				baseUrl: '/rpc/',
				interceptors: [addAuthHeader(adminToken)]
			})
		);
	}

	return adminServiceClients[adminToken];
}

function addAuthHeader(authToken: string): Interceptor {
	return (next) => async (req) => {
		if (authToken) {
			req.header.set('X-PubGolf-Auth', authToken);
		}
		return await next(req);
	};
}
