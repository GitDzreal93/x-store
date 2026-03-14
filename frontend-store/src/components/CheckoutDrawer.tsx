"use client";

import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { X, Clock, Mail, CreditCard } from "lucide-react";
import { useState, useEffect } from "react";
import { Product } from "./ProductCard";

interface CheckoutDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  product: Product | null;
  onPaymentSuccess?: (email: string) => void;
}

export function CheckoutDrawer({ open, onOpenChange, product, onPaymentSuccess }: CheckoutDrawerProps) {
  const [email, setEmail] = useState("");
  const [selectedPayment, setSelectedPayment] = useState<string | null>(null);
  const [countdown, setCountdown] = useState(15 * 60); // 15 minutes in seconds
  const [isProcessing, setIsProcessing] = useState(false);

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

  const paymentMethods = [
    { id: "wechat", name: "微信支付", icon: "💚" },
    { id: "alipay", name: "支付宝", icon: "💙" },
    { id: "stripe", name: "Stripe", icon: "💳" },
  ];

  if (!product) return null;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-lg overflow-y-auto" showCloseButton={true}>
        <SheetHeader>
          <SheetTitle className="flex items-center justify-between">
            <span>收银台</span>
            <div className="flex items-center gap-2 text-sm font-normal text-muted-foreground">
              <Clock className="h-4 w-4" />
              <span>支付倒计时 {formatTime(countdown)}</span>
            </div>
          </SheetTitle>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {/* Order Summary */}
          <div className="space-y-2">
            <h3 className="text-sm font-semibold">订单明细</h3>
            <div className="flex items-center justify-between p-3 border rounded-md bg-muted/30">
              <div className="flex-1">
                <p className="font-medium">{product.title}</p>
                <p className="text-sm text-muted-foreground mt-1">
                  {product.tags.join(" · ")}
                </p>
              </div>
              <div className="text-right">
                <p className="font-bold text-lg">¥ {product.price.toFixed(2)}</p>
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
                  onClick={() => setSelectedPayment(method.id)}
                  className={`flex items-center gap-3 p-3 border rounded-md transition-all ${
                    selectedPayment === method.id
                      ? "border-primary bg-primary/5"
                      : "hover:bg-muted/50"
                  }`}
                >
                  <span className="text-2xl">{method.icon}</span>
                  <span className="flex-1 text-left font-medium">{method.name}</span>
                  {selectedPayment === method.id && (
                    <div className="h-4 w-4 rounded-full bg-primary" />
                  )}
                </button>
              ))}
            </div>
          </div>

          <Separator />

          {/* Total & Checkout Button */}
          <div className="space-y-4">
            <div className="flex items-center justify-between text-lg font-semibold">
              <span>总计</span>
              <span className="text-2xl">¥ {product.price.toFixed(2)}</span>
            </div>
            <Button
              className="w-full h-12 text-base font-semibold"
              disabled={!email || !selectedPayment || countdown === 0 || isProcessing}
              onClick={() => {
                setIsProcessing(true);
                // Simulate payment processing
                setTimeout(() => {
                  setIsProcessing(false);
                  if (onPaymentSuccess) {
                    onPaymentSuccess(email);
                  }
                  setEmail("");
                  setSelectedPayment(null);
                }, 1500);
              }}
            >
              {isProcessing ? "处理中..." : "🔒 立即安全支付"}
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
