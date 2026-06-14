// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import sitemap from '@astrojs/sitemap';
import icon from 'astro-icon';
import starlightLlmsTxt from 'starlight-llms-txt';
import starlightContextualMenu from 'starlight-contextual-menu';
import starlightLinksValidator from 'starlight-links-validator';
import rehypeExternalLinks from 'rehype-external-links';
import { unified } from '@astrojs/markdown-remark';

// https://astro.build/config
export default defineConfig({
  site: 'https://go-automate.timmo.dev',
  markdown: {
    processor: unified({
      rehypePlugins: [
        [rehypeExternalLinks, { target: '_blank', rel: ['noopener', 'noreferrer'] }],
      ],
    }),
  },
  integrations: [
    icon(),
    sitemap(),
    starlight({
      title: 'Go Automate',
      logo: {
        src: './src/assets/logo.svg',
        alt: 'Go Automate logo',
      },
      favicon: '/favicon.svg',
      customCss: ['./src/styles/starlight.css', './src/styles/landing.css'],
      editLink: {
        baseUrl: 'https://github.com/timmo001/go-automate/edit/main/docs/',
      },
      lastUpdated: true,
      head: [
        {
          tag: 'meta',
          attrs: { property: 'og:image', content: 'https://go-automate.timmo.dev/og.png' },
        },
        {
          tag: 'meta',
          attrs: { property: 'og:image:width', content: '1200' },
        },
        {
          tag: 'meta',
          attrs: { property: 'og:image:height', content: '630' },
        },
        {
          tag: 'meta',
          attrs: { property: 'og:image:alt', content: 'Go Automate' },
        },
        {
          tag: 'meta',
          attrs: { name: 'twitter:image', content: 'https://go-automate.timmo.dev/og.png' },
        },
      ],
      plugins: [
        starlightLinksValidator(),
        starlightLlmsTxt({
          projectName: 'Go Automate',
          description: 'A CLI utility to run common tasks and trigger Home Assistant automations.',
          promote: ['index*', 'overview*'],
        }),
        starlightContextualMenu({
          actions: ['copy', 'view'],
        }),
      ],
      components: {
        PageFrame: './src/components/PageFrame.astro',
      },
      social: [
        { icon: 'github', label: 'GitHub', href: 'https://github.com/timmo001/go-automate' },
      ],
      sidebar: [
        { label: 'Overview', slug: 'overview' },
        {
          label: 'Getting Started',
          items: [
            { label: 'Install', slug: 'install' },
            { label: 'Configuration', slug: 'configuration' },
            { label: 'Running the Bridge', slug: 'running' },
          ],
        },
        {
          label: 'Using',
          items: [
            { label: 'CLI', slug: 'using/cli' },
            { label: 'TUI', slug: 'using/tui' },
            { label: 'Home Assistant', slug: 'using/home-assistant' },
            { label: 'Watching Entities', slug: 'using/watching' },
            { label: 'Notifications', slug: 'using/notifications' },
          ],
        },
        {
          label: 'Reference',
          items: [
            { label: 'Overview', slug: 'reference' },
            { label: 'Commands', slug: 'reference/commands' },
            { label: 'Bridge Protocol', slug: 'reference/bridge' },
          ],
        },
        { label: 'LLMs', slug: 'llms' },
      ],
    }),
  ],
});
