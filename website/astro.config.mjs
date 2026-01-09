// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://budimanr3101.github.io',
	base: '/nahkoda',
	integrations: [
		starlight({
			title: 'Nahkoda',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/budimanr3101/nahkoda' }],
			sidebar: [
				{
					label: 'Mulai Dari Sini',
					items: [
						{ label: 'Pengenalan', slug: 'index' },
						{ label: 'Instalasi', slug: 'instalasi' },
					],
				},
				{
					label: 'Panduan Perintah',
					items: [
						{ label: 'Navigasi (Context)', slug: 'perintah/navigasi' },
						{ label: 'Monitoring (Pods/Nodes)', slug: 'perintah/monitoring' },
						{ label: 'Debugging', slug: 'perintah/debugging' },
						{ label: 'Operasi & Metrik', slug: 'perintah/operasi' },
					],
				},
				{
					label: 'Dokumentasi Teknis',
					autogenerate: { directory: 'teknis' },
				},
			],
		}),
	],
});
