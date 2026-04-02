"use client";

import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { X, Clock, Mail, CreditCard, Loader2, AlertCircle, ExternalLink } from "lucide-react";
import { useState, useEffect } from "react";
import { Product } from "./ProductCard";
import { createOrder, createPayment, pollPaymentStatus, ApiError } from "@/lib/api";

interface CheckoutDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  product: Product | null;
  onPaymentSuccess?: (orderNo: string, email: string) => void;
  onPaymentError?: (error: string) => void;
}

interface PaymentChannel {
  id: number;
  name: string;
  icon: string;
}

export function CheckoutDrawer({ open, onOpenChange, product, onPaymentSuccess, onPaymentError }: CheckoutDrawerProps) {
  const [email, setEmail] = useState("");
  const [selectedPayment, setSelectedPayment] = useState<PaymentChannel | null>(null);
  const [countdown, setCountdown] = useState(15 * 60); // 15 minutes in seconds
  const [isProcessing, setIsProcessing] = useState(false);
  const [isPolling, setIsPolling] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [paymentUrl, setPaymentUrl] = useState<string | null>(null);

  // 模拟支付渠道（后续从后端 API 获取）
  const paymentMethods: PaymentChannel[] = [
    { id: 1, name: "模拟支付", icon: "💳" },
    // { id: 2, name: "微信支付", icon: "💚" },
    // { id: 3, name: "支付宝", icon: "💙" },
    // { id: 4, name: "Stripe", icon: "💳" },
  ];

  useEffect(() => {
    if (open && product) {
      setCountdown(15 * 60);
      const timer = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) {
            clearInterval(timer);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
      return () => clearInterval(timer);
    }
  }, [open, product]);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  // 处理支付流程
  const handlePayment = async () => {
    if (!product || !email || !selectedPayment) return;

    setIsProcessing(true);
    setError(null);
    setPaymentUrl(null);

    try {
      // 1. 创建订单
      const order = await createOrder({
        product_id: parseInt(product.id),
        email: email,
        quantity: 1,
      });

      // 2. 创建支付
      const payment = await createPayment({
        order_no: order.order_no,
        channel_id: selectedPayment.id,
      });

      // 3. 处理支付结果
      if (payment.payment_url) {
        // 有支付链接，显示给用户
        setPaymentUrl(payment.payment_url);

        // 开始轮询支付状态
        setIsPolling(true);
        pollPaymentStatus(payment.payment_id)
          .then((result) => {
            setIsPolling(false);
            if (result.status === "success") {
              onPaymentSuccess?.(order.order_no, email);
              onOpenChange(false);
              // 重置状态
              setEmail("");
              setSelectedPayment(null);
              setPaymentUrl(null);
            } else {
              setError("支付失败，请重试");
            }
          })
          .catch((err) => {
            setIsPolling(false);
            setError(err instanceof Error ? err.message : "支付查询超时");
          });
      } else if (payment.qr_code) {
        // 有二维码（可以显示二维码）
        // 这里简化处理，直接跳转到支付链接
        setError("请完成支付后关闭此窗口");
      } else {
        // 模拟支付直接成功
        await new Promise((resolve) => setTimeout(resolve, 1500));
        onPaymentSuccess?.(order.order_no, email);
        onOpenChange(false);
        setEmail("");
        setSelectedPayment(null);
      }
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
        onPaymentError?.(err.message);
      } else {
        setError("创建订单失败，请稍后重试");
        onPaymentError?.("创建订单失败");
      }
    } finally {
      setIsProcessing(false);
    }
  };

  if (!product) return null;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-lg overflow-y-auto" showCloseButton={true}>
        <SheetHeader>
          <SheetTitle className="flex items-center justify-between">
            <span>收银台</span>
            <div className="flex items-center gap-2 text-sm font-normal text-muted-foreground">
              <Clock className="h-4 w-4" />
              <span className={countdown < 300 ? "text-destructive" : ""}>
                支付倒计时 {formatTime(countdown)}
              </span>
            </div>
          </SheetTitle>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {/* Order Summary */}
          <div className="space-y-2">
            <h3 className="text-sm font-semibold">订单明细</h3>
            <div className="flex items-center justify-between p-3 border rounded-md bg-muted/30">
              <div className="flex-1">
                <p className="font-medium">{product?.title}</p>
                <p className="text-sm text-muted-foreground mt-1">
                  {product?.tags.join(" · ")}
                </p>
              </div>
              <div className="text-right">
                <p className="font-bold text-lg">¥ {product?.price.toFixed(2)}</p>
              </div>
            </div>
          </div>

          <Separator />

          {/* Email Input */}
          <div className="space-y-2">
            <Label htmlFor="email" className="flex items-center gap-2">
              <Mail className="h-4 w-4" />
              接收邮箱
            </Label>
            <Input
              id="email"
              type="email"
              placeholder="请输入您的邮箱地址..."
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="h-10"
              disabled={isProcessing || isPolling}
            />
            <p className="text-xs text-muted-foreground">
              卡密将发送至此邮箱，请确保邮箱地址正确
            </p>
          </div>

          <Separator />

          {/* Payment Methods */}
          <div className="space-y-3">
            <Label className="flex items-center gap-2">
              <CreditCard className="h-4 w-4" />
              支付方式
            </Label>
            <div className="grid grid-cols-1 gap-2">
              {paymentMethods.map((method) => (
                <button
                  key={method.id}
                  onClick={() => !isProcessing && !isPolling && setSelectedPayment(method)}
                  disabled={isProcessing || isPolling}
                  className={`flex items-center gap-3 p-3 border rounded-md transition-all ${
                    selectedPayment?.id === method.id
                      ? "border-primary bg-primary/5"
                      : "hover:bg-muted/50"
                  } ${(isProcessing || isPolling) ? "opacity-50 cursor-not-allowed" : ""}`}
                >
                  <span className="text-2xl">{method.icon}</span>
                  <span className="flex-1 text-left font-medium">{method.name}</span>
                  {selectedPayment?.id === method.id && (
                    <div className="h-4 w-4 rounded-full bg-primary" />
                  )}
                </button>
              ))}
            </div>
          </div>

          {/* Payment URL */}
          {paymentUrl && (
            <>
              <Separator />
              <div className="space-y-3 p-4 bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-800 rounded-md">
                <div className="flex items-center gap-2 text-sm">
                  <ExternalLink className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                  <span className="font-medium text-blue-700 dark:text-blue-300">
                    请在新窗口完成支付
                  </span>
                </div>
                <p className="text-xs text-blue-600 dark:text-blue-400">
                  支付完成后将自动检测支付状态
                </p>
                {isPolling && (
                  <div className="flex items-center gap-2 text-sm text-blue-600 dark:text-blue-400 mt-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <span>等待支付中...</span>
                  </div>
                )}
              </div>
            </>
          )}

          {/* Error Message */}
          {error && (
            <div className="flex items-start gap-2 p-3 bg-destructive/10 border border-destructive/20 rounded-md">
              <AlertCircle className="h-5 w-5 text-destructive flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <p className="text-sm font-medium text-destructive">支付失败</p>
                <p className="text-xs text-destructive/80 mt-1">{error}</p>
              </div>
            </div>
          )}

          <Separator />

          {/* Total & Checkout Button */}
          <div className="space-y-4">
            <div className="flex items-center justify-between text-lg font-semibold">
              <span>总计</span>
              <span className="text-2xl">¥ {product?.price.toFixed(2)}</span>
            </div>
            <Button
              className="w-full h-12 text-base font-semibold"
              disabled={!email || !selectedPayment || countdown === 0 || isProcessing || isPolling}
              onClick={handlePayment}
            >
              {isProcessing ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  处理中...
                </>
              ) : isPolling ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  等待支付...
                </>
              ) : (
                "🔒 立即安全支付"
              )}
            </Button>
            {paymentUrl && (
              <Button
                variant="outline"
                className="w-full"
                onClick={() => window.open(paymentUrl, "_blank")}
              >
                <ExternalLink className="h-4 w-4 mr-2" />
                重新打开支付页面
              </Button>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
