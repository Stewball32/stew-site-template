import { LayoutDashboardIcon, CircleUserIcon, SettingsIcon } from '@lucide/svelte';
import type { Component } from 'svelte';

export interface NavLink {
	label: string;
	href: string;
	icon: Component;
}

export const navLinks: NavLink[] = [
	{ label: 'Dashboard', href: '/', icon: LayoutDashboardIcon },
	{ label: 'Profile', href: '/profile/', icon: CircleUserIcon },
	{ label: 'Settings', href: '/settings/', icon: SettingsIcon }
];
