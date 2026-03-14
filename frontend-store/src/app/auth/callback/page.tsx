'use client';

import { Suspense, useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

function CallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [errorMsg, setErrorMsg] = useState('');

  useEffect(() => {
    const token = searchParams.get('token');
    const error = searchParams.get('error');

    if (error) {
      setStatus('error');
      setErrorMsg(decodeURIComponent(error));
      return;
    }

    if (token) {
      localStorage.setItem('token', token);

      // 用 token 获取用户信息
      fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8082'}/api/user/profile`, {
        headers: { Authorization: `Bearer ${token}` },
      })
        .then(res => res.json())
        .then(data => {
          if (data.code === 0 && data.data) {
            localStorage.setItem('user', JSON.stringify(data.data));
          }
          setStatus('success');
          setTimeout(() => router.push('/profile'), 1500);
        })
        .catch(() => {
          setStatus('success');
          setTimeout(() => router.push('/'), 1500);
        });
    } else {
      setStatus('error');
      setErrorMsg('未收到授权信息');
    }
  }, [searchParams, router]);

  return (
    <Card className="w-full max-w-sm text-center">
      <CardHeader>
        <CardTitle className="text-xl">
          {status === 'loading' && '\uD83D\uDD04 正在登录...'}
          {status === 'success' && '\u2705 登录成功'}
          {status === 'error' && '\u274C 登录失败'}
        </CardTitle>
      </CardHeader>
      <CardContent>
        {status === 'loading' && (
          <p className="text-muted-foreground">正在处理第三方授权，请稍候...</p>
        )}
        {status === 'success' && (
          <p className="text-muted-foreground">即将跳转到个人中心...</p>
        )}
        {status === 'error' && (
          <div className="space-y-4">
            <p className="text-red-500 text-sm">{errorMsg}</p>
            <Button onClick={() => router.push('/auth')} className="w-full">
              返回登录
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

export default function OAuthCallbackPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-500 to-pink-500 p-4">
      <Suspense fallback={
        <Card className="w-full max-w-sm text-center">
          <CardHeader><CardTitle className="text-xl">正在加载...</CardTitle></CardHeader>
        </Card>
      }>
        <CallbackContent />
      </Suspense>
    </div>
  );
}
