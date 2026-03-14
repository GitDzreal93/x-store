"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
  Download,
  Heart,
  Clock,
  ShieldCheck,
  Mail,
  Zap,
  Package,
} from "lucide-react";
import { Product } from "./ProductCard";
import { motion, AnimatePresence } from "framer-motion";

interface ProductDetailModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  product: Product | null;
  onPurchase: (product: Product) => void;
}

export function ProductDetailModal({
  open,
  onOpenChange,
  product,
  onPurchase,
}: ProductDetailModalProps) {
  if (!product) return null;

  const handlePurchase = () => {
    onOpenChange(false);
    onPurchase(product);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className="sm:max-w-2xl backdrop-blur-xl bg-background/80 border-border/50 shadow-2xl"
        showCloseButton={true}
      >
        <DialogHeader>
          <DialogTitle className="text-xl font-bold leading-tight">
            {product.title}
          </DialogTitle>
          <DialogDescription className="sr-only">
            商品详情信息
          </DialogDescription>
        </DialogHeader>

        <div className="grid grid-cols-1 sm:grid-cols-5 gap-6 mt-2">
          {/* Left: Product Details */}
          <div className="sm:col-span-3 space-y-4">
            {/* Tags */}
            <div className="flex flex-wrap gap-1.5">
              {product.tags.map((tag, index) => (
                <Badge
                  key={index}
                  variant="outline"
                  className="text-xs px-2 py-0.5 font-normal"
                >
                  {tag}
                </Badge>
              ))}
              {product.isNew && (
                <Badge variant="secondary" className="text-xs px-2 py-0.5">
                  NEW
                </Badge>
              )}
              {product.isHot && (
                <Badge variant="destructive" className="text-xs px-2 py-0.5">
                  HOT
                </Badge>
              )}
            </div>

            {/* Description (Markdown-like rich text) */}
            <div className="prose prose-sm dark:prose-invert max-w-none text-sm text-muted-foreground space-y-3">
              <p>
                本商品为正规渠道采购，支持售后质保。购买后系统自动发货至您的邮箱，请放心选购。
              </p>
              <div className="space-y-2">
                <h4 className="text-sm font-semibold text-foreground">
                  商品说明
                </h4>
                <ul className="list-disc list-inside space-y-1 text-sm">
                  <li>购买后即时自动发货到邮箱</li>
                  <li>质保期内如遇问题可免费换新</li>
                  <li>严禁用于非法用途，违者后果自负</li>
                </ul>
              </div>
            </div>

            {/* Stats Bar */}
            <div className="flex items-center gap-4 text-sm text-muted-foreground pt-2">
              <div className="flex items-center gap-1.5">
                <Download className="h-3.5 w-3.5" />
                <span>
                  {product.sales >= 1000
                    ? `${(product.sales / 1000).toFixed(1)}k`
                    : product.sales}{" "}
                  销量
                </span>
              </div>
              {product.likes !== undefined && (
                <div className="flex items-center gap-1.5">
                  <Heart className="h-3.5 w-3.5" />
                  <span>{product.likes} 好评</span>
                </div>
              )}
              <div className="flex items-center gap-1.5">
                <Clock className="h-3.5 w-3.5" />
                <span>{product.updatedAt}</span>
              </div>
            </div>
          </div>

          {/* Right: Specs & Purchase */}
          <div className="sm:col-span-2 space-y-4">
            {/* Price */}
            <div className="p-4 rounded-lg bg-muted/40 border space-y-3">
              <div className="text-center">
                <span className="text-3xl font-bold text-foreground">
                  ¥ {product.price.toFixed(2)}
                </span>
              </div>

              <Separator />

              {/* Trust Indicators */}
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <ShieldCheck className="h-3.5 w-3.5 text-green-500" />
                  <span>官方正品保障</span>
                </div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Zap className="h-3.5 w-3.5 text-yellow-500" />
                  <span>付款后极速自动发货</span>
                </div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Mail className="h-3.5 w-3.5 text-blue-500" />
                  <span>卡密同步发送至邮箱</span>
                </div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Package className="h-3.5 w-3.5 text-purple-500" />
                  <span>质保期内免费售后</span>
                </div>
              </div>
            </div>

            {/* Purchase Button */}
            <Button
              className="w-full h-11 text-base font-semibold"
              onClick={handlePurchase}
            >
              立即购买
            </Button>
            <p className="text-center text-xs text-muted-foreground">
              点击后将跳转至收银台完成支付
            </p>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
