'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useAuth } from '@/contexts/AuthContext';
import { getUserOrders, Order, OrderStatus, ApiError } from '@/lib/api';
import { Navbar } from '@/components/Navbar';
import { Loader2, Package, User as UserIcon, Mail, Calendar, ShoppingBag, ExternalLink, LogOut } from 'lucide-react';

const statusMap: Record<OrderStatus, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
  [OrderStatus.Pending]: { label: '待支付', variant: 'secondary' },
  [OrderStatus.Paid]: { label: '已支付', variant: 'default' },
  [OrderStatus.Delivered]: { label: '已发货', variant: 'default' },
  [OrderStatus.Cancelled]: { label: '已取消', variant: 'destructive' },
  [OrderStatus.Refunded]: { label: '已退款', variant: 'outline' },
};

function ProfileContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user, isAuthenticated, logout, isLoading: authLoading } = useAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loadingOrders, setLoadingOrders] = useState(true);
  const [ordersError, setOrdersError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState(searchParams.get('tab') || 'profile');

  // 重定向未登录用户
  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/auth');
    }
  }, [authLoading, isAuthenticated, router]);

  // 获取订单列表
  useEffect(() => {
    if (isAuthenticated) {
      fetchOrders();
    }
  }, [isAuthenticated]);

  const fetchOrders = async () => {
    try {
      setLoadingOrders(true);
      setOrdersError(null);
      const result = await getUserOrders();
      setOrders(result.list);
    } catch (err) {
      if (err instanceof ApiError) {
        setOrdersError(err.message);
      } else {
        setOrdersError('加载订单失败');
      }
    } finally {
      setLoadingOrders(false);
    }
  };

  const handleLogout = async () => {
    await logout();
    router.push('/');
  };

  if (authLoading || !isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background">
      <Navbar />

      <div className="container max-w-4xl mx-auto px-4 py-8">
        {/* 页面标题 */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold">个人中心</h1>
          <p className="text-muted-foreground mt-2">管理您的账户信息和订单</p>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="grid w-full grid-cols-2 lg:w-96">
            <TabsTrigger value="profile">
              <UserIcon className="h-4 w-4 mr-2" />
              个人信息
            </TabsTrigger>
            <TabsTrigger value="orders">
              <ShoppingBag className="h-4 w-4 mr-2" />
              我的订单
            </TabsTrigger>
          </TabsList>

          {/* 个人信息 Tab */}
          <TabsContent value="profile">
            <Card>
              <CardHeader>
                <CardTitle>账户信息</CardTitle>
                <CardDescription>您的个人账户详情</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <UserIcon className="h-4 w-4" />
                      <span>用户名</span>
                    </div>
                    <p className="font-medium text-lg">{user.username}</p>
                  </div>

                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Mail className="h-4 w-4" />
                      <span>邮箱地址</span>
                    </div>
                    <p className="font-medium text-lg">{user.email}</p>
                  </div>

                  {user.nickname && (
                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <UserIcon className="h-4 w-4" />
                        <span>显示名称</span>
                      </div>
                      <p className="font-medium text-lg">{user.nickname}</p>
                    </div>
                  )}

                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Calendar className="h-4 w-4" />
                      <span>注册时间</span>
                    </div>
                    <p className="font-medium">
                      {new Date(user.created_at).toLocaleDateString('zh-CN', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </p>
                  </div>
                </div>

                <div className="pt-6 border-t">
                  <Button
                    variant="outline"
                    onClick={handleLogout}
                    className="gap-2"
                  >
                    <LogOut className="h-4 w-4" />
                    退出登录
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* 订单列表 Tab */}
          <TabsContent value="orders">
            <Card>
              <CardHeader>
                <CardTitle>我的订单</CardTitle>
                <CardDescription>查看您的购买记录和订单状态</CardDescription>
              </CardHeader>
              <CardContent>
                {loadingOrders ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                    <span className="ml-3 text-muted-foreground">加载订单中...</span>
                  </div>
                ) : ordersError ? (
                  <div className="text-center py-12">
                    <p className="text-destructive mb-4">{ordersError}</p>
                    <Button variant="outline" onClick={fetchOrders}>
                      重试
                    </Button>
                  </div>
                ) : orders.length === 0 ? (
                  <div className="text-center py-12">
                    <Package className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
                    <p className="text-lg font-medium text-muted-foreground mb-2">
                      暂无订单
                    </p>
                    <p className="text-sm text-muted-foreground mb-6">
                      您还没有购买任何商品
                    </p>
                    <Button onClick={() => router.push('/')}>
                      去逛逛
                    </Button>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {orders.map((order) => (
                      <div
                        key={order.id}
                        className="border rounded-lg p-4 hover:bg-muted/30 transition-colors"
                      >
                        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4 mb-3">
                          <div className="flex-1">
                            <div className="flex items-center gap-3 mb-2">
                              <p className="font-mono text-sm font-medium">
                                {order.order_no}
                              </p>
                              <Badge variant={statusMap[order.status]?.variant || 'secondary'}>
                                {statusMap[order.status]?.label || '未知'}
                              </Badge>
                            </div>
                            <p className="text-sm text-muted-foreground">
                              {new Date(order.created_at).toLocaleString('zh-CN')}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-2xl font-bold">
                              ¥{order.total_amount.toFixed(2)}
                            </p>
                          </div>
                        </div>

                        <div className="flex items-center justify-between pt-3 border-t">
                          <div className="text-sm text-muted-foreground">
                            <p className="font-medium text-foreground mb-1">
                              {order.product_title}
                            </p>
                            <p>数量：{order.quantity}</p>
                          </div>

                          {order.status === OrderStatus.Delivered && order.card_key && (
                            <Button
                              size="sm"
                              variant="outline"
                              className="gap-2"
                              onClick={() => {
                                // 可以打开一个模态框显示卡密
                                navigator.clipboard.writeText(order.card_key!);
                              }}
                            >
                              复制卡密
                            </Button>
                          )}

                          {order.status === OrderStatus.Paid && (
                            <Button
                              size="sm"
                              variant="outline"
                              className="gap-2"
                              onClick={() => router.push(`/profile?order=${order.order_no}`)}
                            >
                              <ExternalLink className="h-3 w-3" />
                              查看详情
                            </Button>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        <div className="mt-8 text-center">
          <Button variant="ghost" onClick={() => router.push('/')}>
            返回首页继续购物
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function ProfilePage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    }>
      <ProfileContent />
    </Suspense>
  );
}
