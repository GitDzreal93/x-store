import React, { useEffect, useState } from 'react';
import {
  Table, Button, Space, Tag, Input, Card, Modal, Descriptions, Select,
} from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import { getOrders } from '../api/order';
import type { Order } from '../types';
import { ORDER_STATUS_MAP, CARD_KEY_STATUS_MAP } from '../types';
import dayjs from 'dayjs';

const Orders: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [size] = useState(20);
  const [loading, setLoading] = useState(false);
  const [keyword, setKeyword] = useState('');
  const [filterStatus, setFilterStatus] = useState<number | undefined>();
  const [detailOpen, setDetailOpen] = useState(false);
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null);

  const fetchOrders = async (p = page) => {
    setLoading(true);
    try {
      const res = await getOrders({
        page: p, size,
        keyword: keyword || undefined,
        status: filterStatus,
      });
      setOrders(res.data.list || []);
      setTotal(res.data.total);
    } catch {
      setOrders([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchOrders(page); }, [page, filterStatus]);

  const handleViewDetail = (record: Order) => {
    setCurrentOrder(record);
    setDetailOpen(true);
  };

  const columns = [
    { title: '订单号', dataIndex: 'order_no', width: 180 },
    { title: '邮箱', dataIndex: 'email', ellipsis: true },
    {
      title: '金额', dataIndex: 'amount', width: 100,
      render: (v: number) => <span style={{ fontWeight: 600 }}>¥{v.toFixed(2)}</span>,
    },
    {
      title: '状态', dataIndex: 'status', width: 90,
      render: (v: number) => {
        const s = ORDER_STATUS_MAP[v] || { label: '未知', color: 'default' };
        return <Tag color={s.color}>{s.label}</Tag>;
      },
    },
    { title: '支付方式', dataIndex: 'pay_method', width: 100, render: (v: string) => v || '-' },
    {
      title: '创建时间', dataIndex: 'created_at', width: 170,
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作', width: 80,
      render: (_: any, record: Order) => (
        <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => handleViewDetail(record)}>
          详情
        </Button>
      ),
    },
  ];

  return (
    <div>
      <Card size="small" style={{ marginBottom: 16 }}>
        <Space wrap>
          <Input.Search
            placeholder="搜索订单号或邮箱"
            allowClear
            style={{ width: 280 }}
            onSearch={(v) => { setKeyword(v); setPage(1); fetchOrders(1); }}
          />
          <Select
            placeholder="订单状态"
            allowClear
            style={{ width: 140 }}
            onChange={(v) => { setFilterStatus(v); setPage(1); }}
            options={Object.entries(ORDER_STATUS_MAP).map(([k, v]) => ({ value: Number(k), label: v.label }))}
          />
        </Space>
      </Card>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={orders}
        loading={loading}
        pagination={{
          current: page, pageSize: size, total,
          onChange: (p) => setPage(p),
          showTotal: (t) => `共 ${t} 条`,
          showSizeChanger: false,
        }}
        size="middle"
      />

      <Modal
        title="订单详情"
        open={detailOpen}
        onCancel={() => setDetailOpen(false)}
        footer={null}
        width={640}
      >
        {currentOrder && (
          <div>
            <Descriptions column={2} bordered size="small" style={{ marginBottom: 16 }}>
              <Descriptions.Item label="订单号" span={2}>{currentOrder.order_no}</Descriptions.Item>
              <Descriptions.Item label="邮箱">{currentOrder.email}</Descriptions.Item>
              <Descriptions.Item label="金额">¥{currentOrder.amount.toFixed(2)}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={ORDER_STATUS_MAP[currentOrder.status]?.color}>
                  {ORDER_STATUS_MAP[currentOrder.status]?.label}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="支付方式">{currentOrder.pay_method || '-'}</Descriptions.Item>
              <Descriptions.Item label="支付时间">
                {currentOrder.paid_at ? dayjs(currentOrder.paid_at).format('YYYY-MM-DD HH:mm:ss') : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="过期时间">
                {dayjs(currentOrder.expire_at).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间" span={2}>
                {dayjs(currentOrder.created_at).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
            </Descriptions>

            {currentOrder.card_keys && currentOrder.card_keys.length > 0 && (
              <div>
                <h4>发货卡密</h4>
                <Table
                  rowKey="id"
                  dataSource={currentOrder.card_keys}
                  pagination={false}
                  size="small"
                  columns={[
                    { title: 'ID', dataIndex: 'id', width: 60 },
                    { title: '卡密内容', dataIndex: 'content' },
                    {
                      title: '状态', dataIndex: 'status', width: 80,
                      render: (v: number) => {
                        const s = CARD_KEY_STATUS_MAP[v] || { label: '未知', color: 'default' };
                        return <Tag color={s.color}>{s.label}</Tag>;
                      },
                    },
                  ]}
                />
              </div>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};

export default Orders;
