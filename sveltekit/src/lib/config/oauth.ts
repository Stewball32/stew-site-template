export type OAuthProviderMeta = {
	label: string;
};

export const OAUTH_PROVIDERS: Record<string, OAuthProviderMeta> = {
	discord: { label: 'Discord' },
	google: { label: 'Google' },
	github: { label: 'GitHub' },
	apple: { label: 'Apple' },
	microsoft: { label: 'Microsoft' },
	facebook: { label: 'Facebook' },
	twitch: { label: 'Twitch' },
	spotify: { label: 'Spotify' },
	patreon: { label: 'Patreon' },
	instagram: { label: 'Instagram' }
};
