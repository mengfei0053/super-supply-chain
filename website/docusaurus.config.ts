import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Super Supply Chain',
  tagline: '供应链管理后台文档',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  // 生产站点与管理后台同域。后台在 /super-supply-chain/，文档在主机根路径 /docs/。
  url: 'https://ssc.mengfei.tech',
  baseUrl: '/docs/',
  trailingSlash: true,

  organizationName: 'mengfei0053',
  projectName: 'super-supply-chain',

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
          routeBasePath: '/',
          editUrl:
            'https://github.com/mengfei0053/super-supply-chain/edit/master/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Super Supply Chain',
      logo: {
        alt: 'Super Supply Chain',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: '文档',
        },
        {
          type: 'custom-appLink',
          href: '/super-supply-chain/',
          label: '进入系统',
          position: 'right',
        },
        {
          href: 'https://github.com/mengfei0053/super-supply-chain',
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
            {label: '概览', to: '/intro'},
            {label: '本地运行', to: '/local-development'},
            {label: 'API 与 MCP', to: '/api-and-mcp'},
            {label: '部署', to: '/deployment'},
          ],
        },
        {
          title: '系统',
          items: [
            {
              html: '<a class="footer__link-item" href="/super-supply-chain/">管理后台</a>',
            },
            {
              label: 'GitHub',
              href: 'https://github.com/mengfei0053/super-supply-chain',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Super Supply Chain. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
