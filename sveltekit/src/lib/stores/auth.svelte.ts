import pb from '$lib/pocketbase';
import type { UsersResponse } from '$lib/types/pocketbase-types';

function createAuthStore() {
	let user = $state<UsersResponse | null>(pb.authStore.record as UsersResponse | null);
	let token = $state(pb.authStore.token);
	const isLoggedIn = $derived(!!token && user !== null);

	pb.authStore.onChange((newToken, record) => {
		user = (record as UsersResponse | null) ?? null;
		token = newToken;
	});

	// pb.authStore.isValid only checks JWT expiry locally; if the server rotated
	// signing keys (e.g. dev DB wipe on task dev restart), the cached token is
	// stale. authRefresh() round-trips to verify; on 401 clear() cascades through
	// onChange to reset the runes and snap the UI to logged-out.
	if (pb.authStore.isValid) {
		pb.collection('users')
			.authRefresh()
			.catch(() => pb.authStore.clear());
	}

	return {
		get user() {
			return user;
		},
		get token() {
			return token;
		},
		get isLoggedIn() {
			return isLoggedIn;
		},
		async register(email: string, password: string, passwordConfirm: string) {
			await pb.collection('users').create({ email, password, passwordConfirm });
			await pb.collection('users').authWithPassword(email, password);
		},
		async login(email: string, password: string) {
			await pb.collection('users').authWithPassword(email, password);
		},
		async loginWithOAuth(provider: string) {
			await pb.collection('users').authWithOAuth2({ provider });
		},
		async listExternalAuths(userId: string) {
			return await pb.collection('users').listExternalAuths(userId);
		},
		async linkOAuth(provider: string) {
			await pb.collection('users').authWithOAuth2({ provider });
		},
		async unlinkOAuth(userId: string, provider: string) {
			await pb.collection('users').unlinkExternalAuth(userId, provider);
		},
		logout() {
			pb.authStore.clear();
		},
		async requestPasswordReset(email: string) {
			await pb.collection('users').requestPasswordReset(email);
		},
		async confirmPasswordReset(token: string, password: string, passwordConfirm: string) {
			await pb.collection('users').confirmPasswordReset(token, password, passwordConfirm);
		},
		async requestVerification(email: string) {
			await pb.collection('users').requestVerification(email);
		},
		async confirmVerification(token: string) {
			await pb.collection('users').confirmVerification(token);
		}
	};
}

export const auth = createAuthStore();
