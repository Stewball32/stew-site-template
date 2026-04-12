import { LayoutDashboardIcon, SettingsIcon } from '@lucide/svelte';
import type { Component } from 'svelte';

export interface NavLink {
	label: string;
	href: string;
	icon: Component;
}

export const mainLinks: NavLink[] = [
	{ label: 'Dashboard', href: '/examples/dashboard/', icon: LayoutDashboardIcon }
];

export const footerLinks: NavLink[] = [
	{ label: 'Settings', href: '/settings/', icon: SettingsIcon }
];
