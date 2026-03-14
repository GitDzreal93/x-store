import React, { useEffect, useState } from 'react';
import {
  Table, Button, Space, Modal, Form, Input, InputNumber,
  Switch, message, Popconfirm, Card,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { getCategories, createCategory, updateCategory, deleteCategory } from '../api/category';
import type { Category } from '../types';

const Categories: React.FC = () => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<Category | null>(null);
  const [form] = Form.useForm();

  const fetchData = async () => {
    setLoading(true);
    try {
      const res = await getCategories(true);
      setCategories(res.data || []);
    } catch { /* interceptor */ } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); }, []);

  const handleCreate = () => {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({ sort: 0, status: 1 });
    setModalOpen(true);
  };

  const handleEdit = (record: Category) => {
    setEditing(record);
    form.setFieldsValue(record);
    setModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteCategory(id);
      message.success('删除成功');
      fetchData();
    } catch { /* interceptor */ }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editing) {
        await updateCategory(editing.id, values);
        message.success('更新成功');
      } else {
        await createCategory(values);
        message.success('创建成功');
      }
      setModalOpen(false);
      fetchData();
    } catch { /* validation */ }
  };

  const handleToggleStatus = async (record: Category) => {
    const newStatus = record.status === 1 ? 0 : 1;
    try {
      await updateCategory(record.id, { status: newStatus });
      message.success(newStatus === 1 ? '已启用' : '已禁用');
      fetchData();
    } catch { /* interceptor */ }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    {
      title: '图标', dataIndex: 'icon', width: 60,
      render: (v: string) => <span style={{ fontSize: 20 }}>{v}</span>,
    },
    { title: '名称', dataIndex: 'name' },
    { title: '排序', dataIndex: 'sort', width: 70 },
    {
      title: '状态', dataIndex: 'status', width: 80,
      render: (v: number, record: Category) => (
        <Switch
          checked={v === 1}
          checkedChildren="启用"
          unCheckedChildren="禁用"
          onChange={() => handleToggleStatus(record)}
          size="small"
        />
      ),
    },
    {
      title: '操作', width: 120,
      render: (_: any, record: Category) => (
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
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增分类
        </Button>
      </Card>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={categories}
        loading={loading}
        pagination={false}
        size="middle"
      />

      <Modal
        title={editing ? '编辑分类' : '新增分类'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
        width={480}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="分类名称" rules={[{ required: true, message: '请输入分类名称' }]}>
            <Input placeholder="AI 大模型" />
          </Form.Item>
          <Space size={16}>
            <Form.Item name="icon" label="图标 (Emoji)" style={{ width: 120 }}>
              <Input placeholder="🤖" maxLength={4} />
            </Form.Item>
            <Form.Item name="sort" label="排序权重" style={{ width: 120 }}>
              <InputNumber min={0} style={{ width: '100%' }} />
            </Form.Item>
          </Space>
        </Form>
      </Modal>
    </div>
  );
};

export default Categories;
