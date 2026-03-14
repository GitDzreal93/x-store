import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Tag, Spin, Empty } from 'antd';
import {
  CreditCardOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import request from '../utils/request';
import type { PaymentChannel } from '../types';

const PaymentChannels: React.FC = () => {
  const [channels, setChannels] = useState<PaymentChannel[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchChannels = async () => {
    setLoading(true);
    try {
      const res = await request.get<any, { code: number; data: PaymentChannel[] }>('/payment-channels');
      setChannels(res.data || []);
    } catch { /* interceptor */ } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchChannels(); }, []);

  const providerLabel: Record<string, { label: string; color: string }> = {
    mock: { label: '模拟支付', color: 'purple' },
    stripe: { label: 'Stripe', color: 'blue' },
    alipay: { label: '支付宝', color: 'cyan' },
    wechat: { label: '微信支付', color: 'green' },
    paypal: { label: 'PayPal', color: 'geekblue' },
    epusdt: { label: 'USDT', color: 'lime' },
    yipay: { label: '易支付', color: 'orange' },
    payjs: { label: 'Payjs', color: 'volcano' },
    xunhupay: { label: '虎皮椒', color: 'gold' },
    vmqpay: { label: 'V免签', color: 'magenta' },
  };

  if (loading) {
    return <div style={{ textAlign: 'center', padding: 100 }}><Spin size="large" /></div>;
  }

  if (channels.length === 0) {
    return <Empty description="暂无支付渠道配置" />;
  }

  return (
    <div>
      <h3 style={{ marginBottom: 16 }}>支付渠道管理</h3>
      <Row gutter={[16, 16]}>
        {channels.map((ch) => {
          const provider = providerLabel[ch.provider_type] || { label: ch.provider_type, color: 'default' };
          return (
            <Col xs={24} sm={12} lg={8} key={ch.id}>
              <Card
                hoverable
                style={{ borderLeft: `4px solid ${ch.is_active ? '#52c41a' : '#d9d9d9'}` }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <CreditCardOutlined style={{ fontSize: 20, color: '#1890ff' }} />
                    <span style={{ fontWeight: 600, fontSize: 16 }}>{ch.name}</span>
                  </div>
                  {ch.is_active ? (
                    <Tag icon={<CheckCircleOutlined />} color="success">已启用</Tag>
                  ) : (
                    <Tag icon={<CloseCircleOutlined />} color="default">已禁用</Tag>
                  )}
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 6, color: '#666', fontSize: 13 }}>
                  <div>
                    <span style={{ color: '#999' }}>提供商：</span>
                    <Tag color={provider.color}>{provider.label}</Tag>
                  </div>
                  <div>
                    <span style={{ color: '#999' }}>渠道类型：</span>
                    <span>{ch.channel_type}</span>
                  </div>
                  <div>
                    <span style={{ color: '#999' }}>交互模式：</span>
                    <span>{ch.interaction_mode || '-'}</span>
                  </div>
                  <div>
                    <span style={{ color: '#999' }}>手续费：</span>
                    <span>{ch.fee_rate > 0 ? `${ch.fee_rate}%` : '无'}{ch.fixed_fee > 0 ? ` + ¥${ch.fixed_fee}` : ''}</span>
                  </div>
                  <div>
                    <span style={{ color: '#999' }}>排序：</span>
                    <span>{ch.sort}</span>
                  </div>
                </div>
              </Card>
            </Col>
          );
        })}
      </Row>
    </div>
  );
};

export default PaymentChannels;
