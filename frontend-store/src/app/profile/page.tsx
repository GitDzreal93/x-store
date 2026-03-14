'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface Order {
  id: number;
  order_no: string;
  product_id: number;
  amount: number;
  status: number;
  created_at: string;
  product?: { title: string };
}

const statusMap: Record<number, { label: string; color: string }> = {
  0: { label: '待支付', color: 'bg-yellow-500' },
  1: { label: '已支付', color: 'bg-blue-500' },
  2: { label: '已发货', color: 'bg-green-500' },
  3: { label: '已完成', color: 'bg-gray-500' },
  4: { label: '已取消', color: 'bg-red-500' },
};

export default function ProfilePage() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('user');
    
    if (!token || !userData) {
      router.push('/auth');
      return;
    }

    setUser(JSON.parse(userData));
    fetchOrders(token);
  }, [router]);

  const fetchOrders = async (token: string) => {
    try {
      const res = await fetch('/api/user/orders', {
        headers: { 'Authorization': `Bearer ${token}` },
      });
      const data = await res.json();
      if (data.code === 0) {
        setOrders(data.data.list || []);
      }
    } catch (error) {
      console.error('Failed to fetch orders:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    router.push('/');
  };

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-4xl mx-auto space-y-6">
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <div>
                <CardTitle>个人中心</CardTitle>
                <CardDescription>欢迎回来，{user.nickname || user.username}</CardDescription>
              </div>
              <Button variant="outline" onClick={handleLogout}>退出登录</Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-500">用户名</p>
                <p className="font-medium">{user.username}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">邮箱</p>
                <p className="font-medium">{user.email}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>我的订单</CardTitle>
            <CardDescription>查看您的购买记录</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <p className="text-center text-gray-500">加载中...</p>
            ) : orders.length === 0 ? (
              <p className="text-center text-gray-500">暂无订单</p>
            ) : (
              <div className="space-y-4">
                {orders.map((order) => (
                  <div key={order.id} className="border rounded-lg p-4 hover:bg-gray-50">
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <p className="font-medium">订单号：{order.order_no}</p>
                        <p className="text-sm text-gray-500">
                          {new Date(order.created_at).toLocaleString('zh-CN')}
                        </p>
                      </div>
                      <Badge className={statusMap[order.status]?.color || 'bg-gray-500'}>
                        {statusMap[order.status]?.label || '未知'}
                      </Badge>
                    </div>
                    <div className="flex justify-between items-center">
                      <p className="text-sm text-gray-600">
                        商品 ID: {order.product_id}
                      </p>
                      <p className="font-semibold text-lg">¥{order.amount.toFixed(2)}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <div className="text-center">
          <Button variant="outline" onClick={() => router.push('/')}>
            返回首页
          </Button>
        </div>
      </div>
    </div>
  );
}
