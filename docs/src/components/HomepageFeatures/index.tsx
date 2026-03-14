import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  emoji: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: '卡密自动发货',
    emoji: '🔑',
    description: (
      <>
        支持批量导入卡密，用户支付成功后自动分配发货，零人工干预，7×24 小时售卖。
      </>
    ),
  },
  {
    title: '10+ 支付方式',
    emoji: '💳',
    description: (
      <>
        支付宝、微信、Stripe、PayPal、USDT 等 10 种支付渠道，覆盖国内外用户，后台一键开关。
      </>
    ),
  },
  {
    title: '第三方登录',
    emoji: '🔐',
    description: (
      <>
        GitHub / Google OAuth 一键登录，管理后台动态配置，无需修改代码和重启服务。
      </>
    ),
  },
  {
    title: '全栈技术',
    emoji: '⚡',
    description: (
      <>
        Go + Gin 高性能后端、Next.js 买家前端、React + Ant Design 管理后台，全方位覆盖。
      </>
    ),
  },
  {
    title: '安全防护',
    emoji: '🛡️',
    description: (
      <>
        JWT 认证、防重放攻击、接口限流、bcrypt 密码加密，企业级安全防护开箱即用。
      </>
    ),
  },
  {
    title: '开箱即用',
    emoji: '📦',
    description: (
      <>
        一键启动脚本、Docker Compose 部署、完整的示例数据和文档，5 分钟跑起来。
      </>
    ),
  },
];

function Feature({title, emoji, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center" style={{fontSize: 48, marginBottom: 8}}>
        {emoji}
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
