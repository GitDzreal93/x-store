import React, { useEffect, useState } from 'react';
import {
  Table, Button, Space, Tag, Modal, Form, Input, InputNumber, Select,
  Switch, message, Popconfirm, Card,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { getProducts, createProduct, updateProduct, deleteProduct } from '../api/product';
import { getCategories } from '../api/category';
import type { Product, Category, CreateProductReq } from '../types';

const { TextArea } = Input;

const Products: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [size] = useState(20);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [form] = Form.useForm();
  const [keyword, setKeyword] = useState('');
  const [filterCategoryId, setFilterCategoryId] = useState<number | undefined>();

  const fetchProducts = async (p = page) => {
    setLoading(true);
    try {
      const res = await getProducts({
        page: p, size,
        keyword: keyword || undefined,
        category_id: filterCategoryId,
      });
      setProducts(res.data.list || []);
      setTotal(res.data.total);
    } catch { /* interceptor */ } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const res = await getCategories(true);
      setCategories(res.data || []);
    } catch { /* interceptor */ }
  };

  useEffect(() => { fetchCategories(); }, []);
  useEffect(() => { fetchProducts(page); }, [page, keyword, filterCategoryId]);

  const handleCreate = () => {
    setEditingProduct(null);
    form.resetFields();
    form.setFieldsValue({ delivery_type: 'auto', status: 1, sort: 0, is_new: false, is_hot: false });
    setModalOpen(true);
  };

  const handleEdit = (record: Product) => {
    setEditingProduct(record);
    form.setFieldsValue({
      ...record,
      description: record.detail?.description || '',
      notice: record.detail?.notice || '',
    });
    setModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteProduct(id);
      message.success('删除成功');
      fetchProducts();
    } catch { /* interceptor */ }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      // tags: convert array to JSON string if needed
      if (Array.isArray(values.tags)) {
        values.tags = JSON.stringify(values.tags);
      }
      if (editingProduct) {
        await updateProduct(editingProduct.id, values);
        message.success('更新成功');
      } else {
        await createProduct(values as CreateProductReq);
        message.success('创建成功');
      }
      setModalOpen(false);
      fetchProducts();
    } catch { /* validation error */ }
  };

  const handleToggleStatus = async (record: Product) => {
    const newStatus = record.status === 1 ? 0 : 1;
    try {
      await updateProduct(record.id, { status: newStatus });
      message.success(newStatus === 1 ? '已上架' : '已下架');
      fetchProducts();
    } catch { /* interceptor */ }
  };

  const columns = [
    {
      title: 'ID', dataIndex: 'id', width: 60,
    },
    {
      title: '商品名称', dataIndex: 'title', ellipsis: true,
    },
    {
      title: '分类', dataIndex: 'category_id', width: 100,
      render: (cid: number) => {
        const cat = categories.find(c => c.id === cid);
        return cat ? <Tag>{cat.icon} {cat.name}</Tag> : '-';
      },
    },
    {
      title: '价格', dataIndex: 'price', width: 90,
      render: (v: number) => <span style={{ fontWeight: 600 }}>¥{v.toFixed(2)}</span>,
    },
    {
      title: '库存', dataIndex: 'stock', width: 70,
      render: (v: number) => (
        <span style={{ color: v === 0 ? '#ff4d4f' : v < 10 ? '#faad14' : undefined }}>{v}</span>
      ),
    },
    {
      title: '销量', dataIndex: 'sales', width: 70,
    },
    {
      title: '状态', dataIndex: 'status', width: 80,
      render: (v: number, record: Product) => (
        <Switch
          checked={v === 1}
          checkedChildren="在售"
          unCheckedChildren="下架"
          onChange={() => handleToggleStatus(record)}
          size="small"
        />
      ),
    },
    {
      title: '标签', width: 120,
      render: (_: any, record: Product) => (
        <Space size={4} wrap>
          {record.is_new && <Tag color="blue">新品</Tag>}
          {record.is_hot && <Tag color="red">热卖</Tag>}
        </Space>
      ),
    },
    {
      title: '操作', width: 120, fixed: 'right' as const,
      render: (_: any, record: Product) => (
        <Space size={4}>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card size="small" style={{ marginBottom: 16 }}>
        <Space wrap>
          <Input.Search
            placeholder="搜索商品名称"
            allowClear
            style={{ width: 240 }}
            onSearch={(v) => { setKeyword(v); setPage(1); }}
          />
          <Select
            placeholder="按分类筛选"
            allowClear
            style={{ width: 160 }}
            onChange={(v) => { setFilterCategoryId(v); setPage(1); }}
            options={categories.map(c => ({ value: c.id, label: `${c.icon} ${c.name}` }))}
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            新增商品
          </Button>
        </Space>
      </Card>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={products}
        loading={loading}
        pagination={{
          current: page, pageSize: size, total,
          onChange: (p) => setPage(p),
          showTotal: (t) => `共 ${t} 条`,
          showSizeChanger: false,
        }}
        scroll={{ x: 900 }}
        size="middle"
      />

      <Modal
        title={editingProduct ? '编辑商品' : '新增商品'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
        width={680}
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ maxHeight: '60vh', overflow: 'auto' }}>
          <Form.Item name="title" label="商品名称" rules={[{ required: true, message: '请输入商品名称' }]}>
            <Input placeholder="ChatGPT Plus 独享账号" />
          </Form.Item>
          <Space style={{ width: '100%' }} size={16}>
            <Form.Item name="category_id" label="分类" rules={[{ required: true, message: '请选择分类' }]} style={{ width: 200 }}>
              <Select
                placeholder="选择分类"
                options={categories.filter(c => c.status === 1).map(c => ({ value: c.id, label: `${c.icon} ${c.name}` }))}
              />
            </Form.Item>
            <Form.Item name="price" label="售价" rules={[{ required: true, message: '请输入售价' }]} style={{ width: 160 }}>
              <InputNumber min={0.01} step={0.01} prefix="¥" style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="original_price" label="原价（划线价）" style={{ width: 160 }}>
              <InputNumber min={0} step={0.01} prefix="¥" style={{ width: '100%' }} />
            </Form.Item>
          </Space>
          <Form.Item name="cover" label="封面图 URL">
            <Input placeholder="https://..." />
          </Form.Item>
          <Space style={{ width: '100%' }} size={16}>
            <Form.Item name="delivery_type" label="发货方式" style={{ width: 160 }}>
              <Select options={[
                { value: 'auto', label: '自动发货' },
                { value: 'manual', label: '人工代充' },
              ]} />
            </Form.Item>
            <Form.Item name="sort" label="排序权重" style={{ width: 120 }}>
              <InputNumber min={0} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="is_new" label="新品" valuePropName="checked" style={{ width: 80 }}>
              <Switch />
            </Form.Item>
            <Form.Item name="is_hot" label="热卖" valuePropName="checked" style={{ width: 80 }}>
              <Switch />
            </Form.Item>
          </Space>
          <Form.Item name="description" label="商品描述">
            <TextArea rows={4} placeholder="支持 Markdown 格式" />
          </Form.Item>
          <Form.Item name="notice" label="购买须知">
            <TextArea rows={2} placeholder="购买须知..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Products;
