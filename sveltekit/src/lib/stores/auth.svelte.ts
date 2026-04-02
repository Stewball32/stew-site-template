import pb from '$lib/pocketbase';
import type { RecordModel } from 'pocketbase';

function createAuthStore() {
	let user = $state<RecordModel | null>(pb.authStore.record);
	let isLoggedIn = $derived(pb.authStore.isValid && user !== null);

	pb.authStore.onChange(() => {
		user = pb.authStore.record;
	});

	return {
		get user() {
			return user;
		},
		get isLoggedIn() {
			return isLoggedIn;
		},
		async login(email: string, password: string) {
			await pb.collection('users').authWithPassword(email, password);
		},
		logout() {
			pb.authStore.clear();
		}
	};
}

export const auth = createAuthStore();
