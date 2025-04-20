import { getAPIToken, USER_NOT_AUTHORIZED_ERROR } from '$lib/auth/client';
import { AdminService } from '$lib/proto/api/v1/admin_pb';
import { PubGolfService } from '$lib/proto/api/v1/pubgolf_pb';
import { createClient, type Client, type Interceptor } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';

export const PubGolfServiceClient = createClient(
	PubGolfService,
	createConnectTransport({
		baseUrl: '/rpc/'
	})
);

export type AdminServiceClient = Client<typeof AdminService>;

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
		adminServiceClients[adminToken] = createClient(
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
			req.header.set('X-Pubgolf-Authtoken', authToken);
		}
		return await next(req);
	};
}
