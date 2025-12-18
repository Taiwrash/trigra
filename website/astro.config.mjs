// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://taiwrash.github.io',
	base: '/trigra',
	integrations: [
		starlight({
			title: 'TRIGRA',
			description: 'Lightweight GitOps for Kubernetes - Simple, Fast, Reliable',
			logo: {
				light: './src/assets/logo-light.svg',
				dark: './src/assets/logo-dark.svg',
				replacesTitle: false,
			},
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/Taiwrash/trigra' },
			],
			customCss: [
				'./src/styles/custom.css',
			],
			components: {
				Footer: './src/components/Footer.astro',
			},
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Quick Start', slug: 'getting-started/quickstart' },
						{ label: 'Installation', slug: 'getting-started/installation' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'GitHub Webhooks', slug: 'guides/github-webhooks' },
						{ label: 'Cloudflare Tunnel', slug: 'guides/cloudflare-tunnel' },
						{ label: 'Deploy Examples', slug: 'guides/deploy-examples' },
						{ label: 'Multi-Node Clusters', slug: 'guides/multi-node' },
					],
				},
				{
					label: 'Configuration',
					items: [
						{ label: 'Helm Values', slug: 'configuration/helm-values' },
						{ label: 'Environment Variables', slug: 'configuration/environment' },
						{ label: 'RBAC & Security', slug: 'configuration/security' },
					],
				},
				{
					label: 'Reference',
					autogenerate: { directory: 'reference' },
				},
			],
			head: [
				{
					tag: 'meta',
					attrs: {
						property: 'og:image',
						content: 'https://trigra.dev/og-image.png',
					},
				},
			],
		}),
	],
});
