import pb from '$lib/pocketbase';
import type { AuthModel } from 'pocketbase';

function createAuthStore() {
	let user = $state<AuthModel | null>(pb.authStore.record);
	let token = $state(pb.authStore.token);
	let isValid = $state(pb.authStore.isValid);
	let isLoggedIn = $derived(isValid && user !== null);

	pb.authStore.onChange((newToken, record) => {
		user = record ?? null;
		token = newToken;
		isValid = !!newToken;
	});

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
		async login(email: string, password: string) {
			await pb.collection('users').authWithPassword(email, password);
		},
		async loginWithOAuth(provider: string) {
			await pb.collection('users').authWithOAuth2({ provider });
		},
		logout() {
			pb.authStore.clear();
		}
	};
}

export const auth = createAuthStore();
