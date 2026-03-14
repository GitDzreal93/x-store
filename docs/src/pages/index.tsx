import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <p style={{fontSize: '1.1rem', opacity: 0.9, maxWidth: 600, margin: '0 auto 1.5rem'}}>
          Go + Gin + Next.js + React + Ant Design + PostgreSQL
        </p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro">
            快速开始 🚀
          </Link>
          <Link
            className="button button--outline button--lg"
            style={{marginLeft: 12, color: '#fff', borderColor: '#fff'}}
            to="/docs/api/overview">
            API 文档 📖
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title="开源数字商品自动化售卖系统"
      description="X-Store 是一个基于 Go + Next.js + React 的开源数字商品自动化售卖系统，支持卡密自动发货、多渠道支付、OAuth 登录">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
