import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  tutorialSidebar: [
    'intro',
    {
      type: 'category',
      label: '快速开始',
      collapsed: false,
      items: [
        'getting-started/prerequisites',
        'getting-started/installation',
        'getting-started/configuration',
        'getting-started/first-run',
      ],
    },
    {
      type: 'category',
      label: '项目架构',
      items: [
        'architecture/overview',
        'architecture/tech-stack',
        'architecture/directory-structure',
        'architecture/database-design',
      ],
    },
    {
      type: 'category',
      label: '后端开发',
      items: [
        'backend/overview',
        'backend/config',
        'backend/models',
        'backend/router',
        'backend/middleware',
        'backend/payment',
        'backend/oauth',
      ],
    },
    {
      type: 'category',
      label: '前端开发',
      items: [
        'frontend/store',
        'frontend/admin',
      ],
    },
    {
      type: 'category',
      label: '部署运维',
      items: [
        'deployment/docker',
        'deployment/production',
      ],
    },
  ],

  apiSidebar: [
    'api/overview',
    {
      type: 'category',
      label: '接口列表',
      collapsed: false,
      items: [
        'api/auth',
        'api/products',
        'api/orders',
        'api/payment',
        'api/admin',
      ],
    },
  ],
};

export default sidebars;
