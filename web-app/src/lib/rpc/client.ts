import { AdminService } from '$lib/proto/api/v1/admin_connect';
import { PubGolfService } from '$lib/proto/api/v1/pubgolf_connect';
import { createConnectTransport } from '@bufbuild/connect-web';
import { createPromiseClient } from '@bufbuild/connect';

export const PubGolfClient = createPromiseClient(
	PubGolfService,
	createConnectTransport({
		baseUrl: '/rpc/'
	})
);

export const AdminClient = createPromiseClient(
	AdminService,
	createConnectTransport({
		baseUrl: '/rpc/'
	})
);
