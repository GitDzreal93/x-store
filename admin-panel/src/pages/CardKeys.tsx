import React, { useEffect, useState } from 'react';
import {
  Card, Button, Space, Select, Row, Col, Modal, Form, Input,
  message, Tag, Progress,
} from 'antd';
import { ImportOutlined } from '@ant-design/icons';
import { getProducts } from '../api/product';
import { getCardKeyCount, importCardKeys } from '../api/order';
import type { Product } from '../types';

const { TextArea } = Input;

const CardKeys: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [stockMap, setStockMap] = useState<Record<number, number>>({});
  const [loading, setLoading] = useState(false);
  const [importOpen, setImportOpen] = useState(false);
  const [importLoading, setImportLoading] = useState(false);
  const [form] = Form.useForm();

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const res = await getProducts({ page: 1, size: 100 });
      const list = res.data.list || [];
      setProducts(list);
      // Fetch card key counts for each product
      const counts: Record<number, number> = {};
      await Promise.all(
        list.map(async (p: Product) => {
          try {
            const countRes = await getCardKeyCount(p.id);
            counts[p.id] = countRes.data.available;
          } catch {
            counts[p.id] = 0;
          }
        })
      );
      setStockMap(counts);
    } catch { /* interceptor */ } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchProducts(); }, []);

  const getStockLevel = (available: number) => {
    if (available === 0) return { color: '#d9d9d9', text: '缺货', status: 'exception' as const };
    if (available < 10) return { color: '#ff4d4f', text: '紧急', status: 'exception' as const };
    if (available < 50) return { color: '#faad14', text: '预警', status: 'normal' as const };
    return { color: '#52c41a', text: '充足', status: 'success' as const };
  };

  const handleImport = () => {
    form.resetFields();
    setImportOpen(true);
  };

  const handleSubmitImport = async () => {
    try {
      const values = await form.validateFields();
      setImportLoading(true);
      const res = await importCardKeys({
        product_id: values.product_id,
        content: values.content,
      });
      message.success(`成功导入 ${res.data.imported} 个卡密`);
      setImportOpen(false);
      fetchProducts();
    } catch { /* validation */ } finally {
      setImportLoading(false);
    }
  };

  return (
    <div>
      <Card size="small" style={{ marginBottom: 16 }}>
        <Space>
          <Button type="primary" icon={<ImportOutlined />} onClick={handleImport}>
            批量导入卡密
          </Button>
          <Button onClick={fetchProducts} loading={loading}>刷新库存</Button>
        </Space>
      </Card>

      <h3 style={{ marginBottom: 16 }}>库存水位线</h3>
      <Row gutter={[16, 16]}>
        {products.map((p) => {
          const available = stockMap[p.id] ?? 0;
          const level = getStockLevel(available);
          const maxStock = Math.max(available, p.stock, 100);
          return (
            <Col xs={24} sm={12} lg={8} xl={6} key={p.id}>
              <Card size="small" hoverable>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                  <span style={{ fontWeight: 600, fontSize: 14, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', flex: 1 }}>
                    {p.title}
                  </span>
                  <Tag color={level.color} style={{ marginLeft: 8 }}>{level.text}</Tag>
                </div>
                <Progress
                  percent={Math.round((available / maxStock) * 100)}
                  status={level.status}
                  size="small"
                  format={() => `${available}`}
                />
                <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 8, color: '#999', fontSize: 12 }}>
                  <span>可用卡密: {available}</span>
                  <span>DB库存: {p.stock}</span>
                </div>
              </Card>
            </Col>
          );
        })}
      </Row>

      <Modal
        title="批量导入卡密"
        open={importOpen}
        onCancel={() => setImportOpen(false)}
        onOk={handleSubmitImport}
        confirmLoading={importLoading}
        width={560}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="product_id" label="选择商品" rules={[{ required: true, message: '请选择商品' }]}>
            <Select
              placeholder="选择商品"
              showSearch
              optionFilterProp="label"
              options={products.map((p) => ({ value: p.id, label: p.title }))}
            />
          </Form.Item>
          <Form.Item
            name="content"
            label="卡密内容"
            rules={[{ required: true, message: '请输入卡密内容' }]}
            extra="每行一个卡密，系统会自动去重"
          >
            <TextArea rows={10} placeholder={'CHATGPT-PLUS-KEY-001\nCHATGPT-PLUS-KEY-002\nCHATGPT-PLUS-KEY-003'} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default CardKeys;
