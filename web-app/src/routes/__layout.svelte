<script lang="ts">
	import { browser } from '$app/env';
	import { API_SPEC_VERSION } from '../../../proto/versions/current/ts/version';
	import {
		PubGolfServiceClientImpl,
		GrpcWebImpl,
		ClientVersionResponse_VersionStatus
	} from '@rpc/pubgolf';
import { API_BASE } from 'src/_env';

	let message: string;
	let urgency: 'error' | 'warning' | null;

	if (browser) {
		const c = new PubGolfServiceClientImpl(new GrpcWebImpl(`${API_BASE}`, {}));
		c.ClientVersion({ clientVersion: API_SPEC_VERSION })
			.then((res) => {
				if (res.versionStatus === ClientVersionResponse_VersionStatus.VERSION_STATUS_INCOMPATIBLE) {
					message = 'Proto definition incompatible. Pull the latest code down and rebuild client.';
					urgency = 'error';
				}

				if (res.versionStatus === ClientVersionResponse_VersionStatus.VERSION_STATUS_OUTDATED) {
					message = 'Proto definition out of date. Pull the latest code down and rebuild client.';
					urgency = 'warning';
				}
			})
			.catch(() => {
				message = 'Unable to connect to API Server.';
				urgency = 'error';
			});
	}
</script>

<header>
	{#if message}
		<div class="message-box {urgency}">{message}</div>
	{/if}

	<nav>
		<ul>
      <li><a href="/">Home</a></li>
      <li><a href="/about">About</a></li>
    </ul>
	</nav>
</header>

<main>
	<slot />
</main>

<style lang="scss">
	.message-box {
		width: 100%;
		padding: 0.75rem 1.25rem;
		margin-bottom: 1rem;
		border-radius: 0.5rem;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		text-align: center;

		&.error {
			background-color: palevioletred;
			color: darkred;
		}

		&.warning {
			background-color: palegoldenrod;
			color: darkgoldenrod;
		}
	}
</style>
