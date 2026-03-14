import React, { useEffect, useState } from 'react';
import { Table, Switch, Button, Modal, Form, Input, message, Tag } from 'antd';
import { EditOutlined, EyeOutlined, EyeInvisibleOutlined } from '@ant-design/icons';
import request from '../utils/request';

interface OAuthProvider {
  id: number;
  provider: string;
  name: string;
  enabled: boolean;
  client_id: string;
  client_secret: string;
  redirect_url: string;
  sort: number;
  created_at: string;
  updated_at: string;
}

const OAuthProviders: React.FC = () => {
  const [providers, setProviders] = useState<OAuthProvider[]>([]);
  const [loading, setLoading] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [currentProvider, setCurrentProvider] = useState<OAuthProvider | null>(null);
  const [form] = Form.useForm();

  const fetchProviders = async () => {
    setLoading(true);
    try {
      const res = await request.get<any, { code: number; data: OAuthProvider[] }>('/admin/oauth-providers');
      if (res.code === 0) {
        setProviders(res.data || []);
      }
    } catch (error) {
      message.error('获取配置失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProviders();
  }, []);

  const handleToggle = async (record: OAuthProvider) => {
    try {
      const res = await request.post<any, { code: number; data: OAuthProvider }>(
        `/admin/oauth-providers/${record.id}/toggle`
      );
      if (res.code === 0) {
        message.success(res.data.enabled ? '已启用' : '已禁用');
        fetchProviders();
      }
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleEdit = async (record: OAuthProvider) => {
    try {
      const res = await request.get<any, { code: number; data: OAuthProvider }>(
        `/admin/oauth-providers/${record.id}`
      );
      if (res.code === 0) {
        setCurrentProvider(res.data);
        form.setFieldsValue({
          client_id: res.data.client_id,
          client_secret: res.data.client_secret,
          redirect_url: res.data.redirect_url,
        });
        setEditModalVisible(true);
      }
    } catch (error) {
      message.error('获取详情失败');
    }
  };

  const handleUpdate = async () => {
    try {
      const values = await form.validateFields();
      if (!currentProvider) return;

      const res = await request.put<any, { code: number; message: string }>(
        `/admin/oauth-providers/${currentProvider.id}`,
        values
      );
      if (res.code === 0) {
        message.success('更新成功');
        setEditModalVisible(false);
        fetchProviders();
      }
    } catch (error) {
      message.error('更新失败');
    }
  };

  const columns = [
    {
      title: '提供商',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: OAuthProvider) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <span style={{ fontSize: 20 }}>
            {record.provider === 'github' ? '🐙' : '🔵'}
          </span>
          <span style={{ fontWeight: 600 }}>{text}</span>
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 100,
      render: (enabled: boolean, record: OAuthProvider) => (
        <Switch
          checked={enabled}
          onChange={() => handleToggle(record)}
          checkedChildren="启用"
          unCheckedChildren="禁用"
        />
      ),
    },
    {
      title: 'Client ID',
      dataIndex: 'client_id',
      key: 'client_id',
      render: (text: string) => {
        if (!text || text.length < 10) return text;
        return (
          <code style={{ fontSize: 12, color: '#666' }}>
            {text.substring(0, 8)}...{text.substring(text.length - 8)}
          </code>
        );
      },
    },
    {
      title: 'Client Secret',
      dataIndex: 'client_secret',
      key: 'client_secret',
      render: (text: string) => (
        <Tag color="default">{text === '******' ? '已配置' : '未配置'}</Tag>
      ),
    },
    {
      title: '回调地址',
      dataIndex: 'redirect_url',
      key: 'redirect_url',
      ellipsis: true,
      render: (text: string) => (
        <span style={{ fontSize: 12, color: '#999' }}>{text}</span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_: any, record: OAuthProvider) => (
        <Button
          type="link"
          icon={<EditOutlined />}
          onClick={() => handleEdit(record)}
        >
          编辑
        </Button>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h3 style={{ margin: 0 }}>第三方登录配置</h3>
        <Button onClick={fetchProviders}>刷新</Button>
      </div>

      <Table
        columns={columns}
        dataSource={providers}
        rowKey="id"
        loading={loading}
        pagination={false}
      />

      <Modal
        title={`编辑 ${currentProvider?.name} 配置`}
        open={editModalVisible}
        onOk={handleUpdate}
        onCancel={() => setEditModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label="Client ID"
            name="client_id"
            rules={[{ required: true, message: '请输入 Client ID' }]}
          >
            <Input placeholder="OAuth 应用的 Client ID" />
          </Form.Item>

          <Form.Item
            label="Client Secret"
            name="client_secret"
            rules={[{ required: true, message: '请输入 Client Secret' }]}
          >
            <Input.Password
              placeholder="OAuth 应用的 Client Secret"
              iconRender={(visible) =>
                visible ? <EyeOutlined /> : <EyeInvisibleOutlined />
              }
            />
          </Form.Item>

          <Form.Item
            label="回调地址"
            name="redirect_url"
            rules={[{ required: true, message: '请输入回调地址' }]}
          >
            <Input placeholder="http://localhost:8082/api/oauth/github/callback" />
          </Form.Item>

          <div style={{ padding: 12, background: '#f5f5f5', borderRadius: 4, fontSize: 12, color: '#666' }}>
            <p style={{ margin: 0, marginBottom: 8 }}><strong>配置说明：</strong></p>
            <ul style={{ margin: 0, paddingLeft: 20 }}>
              <li>GitHub: 到 <a href="https://github.com/settings/developers" target="_blank" rel="noreferrer">github.com/settings/developers</a> 创建 OAuth App</li>
              <li>Google: 到 <a href="https://console.cloud.google.com" target="_blank" rel="noreferrer">console.cloud.google.com</a> 创建 OAuth 凭据</li>
              <li>回调地址必须与 OAuth 应用配置中的一致</li>
            </ul>
          </div>
        </Form>
      </Modal>
    </div>
  );
};

export default OAuthProviders;
