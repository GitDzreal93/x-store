import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'X-Store',
  tagline: '开源数字商品自动化售卖系统',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://x-store.dev',
  baseUrl: '/',

  organizationName: 'x-store',
  projectName: 'x-store',

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'zh-Hans',
    locales: ['zh-Hans'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/x-store/x-store/tree/main/docs/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/x-store-social-card.png',
    colorMode: {
      defaultMode: 'light',
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'X-Store',
      logo: {
        alt: 'X-Store Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: '📖 教程',
        },
        {
          type: 'docSidebar',
          sidebarId: 'apiSidebar',
          position: 'left',
          label: '🔌 API',
        },
        {
          href: 'https://github.com/x-store/x-store',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: '文档',
          items: [
            { label: '快速开始', to: '/docs/intro' },
            { label: '项目架构', to: '/docs/architecture/overview' },
            { label: 'API 参考', to: '/docs/api/overview' },
          ],
        },
        {
          title: '技术栈',
          items: [
            { label: 'Go + Gin', href: 'https://gin-gonic.com/' },
            { label: 'Next.js', href: 'https://nextjs.org/' },
            { label: 'React + Ant Design', href: 'https://ant.design/' },
            { label: 'PostgreSQL', href: 'https://www.postgresql.org/' },
          ],
        },
        {
          title: '更多',
          items: [
            { label: 'GitHub', href: 'https://github.com/x-store/x-store' },
            { label: 'Docusaurus', href: 'https://docusaurus.io/' },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} X-Store. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'go', 'yaml', 'sql', 'typescript', 'json'],
    },
    tableOfContents: {
      minHeadingLevel: 2,
      maxHeadingLevel: 4,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
