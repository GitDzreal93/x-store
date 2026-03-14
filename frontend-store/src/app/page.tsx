"use client";

import { Navbar } from "@/components/Navbar";
import { Sidebar } from "@/components/Sidebar";
import { ProductCard, Product } from "@/components/ProductCard";
import { CheckoutDrawer } from "@/components/CheckoutDrawer";
import { ProductDetailModal } from "@/components/ProductDetailModal";
import { OrderSuccess } from "@/components/OrderSuccess";
import { useState, useCallback, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Search, Grid, List } from "lucide-react";

// Mock data - 后续会从 API 获取
const mockProducts: Product[] = [
  {
    id: "1",
    title: "ChatGPT Plus 独享账号合集",
    tags: ["🏷️ 自动发货", "📦 独享资源", "🔒 质保30天"],
    price: 120.0,
    sales: 12400,
    likes: 224,
    updatedAt: "2026.03.14",
    isNew: true,
  },
  {
    id: "2",
    title: "Netflix 顶级 4K 车票 (1个月)",
    tags: ["🏷️ 自动发货", "👥 拼车账号", "🎬 4K HDR"],
    price: 25.0,
    sales: 85100,
    likes: 189,
    updatedAt: "2026.03.13",
    isHot: true,
  },
  {
    id: "3",
    title: "Github Copilot 一年授权激活码",
    tags: ["🏷️ 激活码", "💻 开发工具", "⚡ 极速交付"],
    price: 150.0,
    sales: 8200,
    likes: 156,
    updatedAt: "2026.03.10",
  },
  {
    id: "4",
    title: "Spotify Premium 家庭版 (6个月)",
    tags: ["🏷️ 自动发货", "👥 家庭共享", "🎵 无损音质"],
    price: 45.0,
    sales: 12300,
    likes: 98,
    updatedAt: "2026.03.12",
  },
  {
    id: "5",
    title: "Adobe Creative Cloud 全系列 (1年)",
    tags: ["🏷️ 激活码", "🎨 设计工具", "💎 官方正品"],
    price: 299.0,
    sales: 5600,
    likes: 201,
    updatedAt: "2026.03.11",
  },
  {
    id: "6",
    title: "Disney+ 4K 会员账号 (3个月)",
    tags: ["🏷️ 自动发货", "🎬 4K HDR", "👨‍👩‍👧‍👦 家庭共享"],
    price: 35.0,
    sales: 9800,
    likes: 142,
    updatedAt: "2026.03.09",
  },
];

export default function Home() {
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isOrderSuccessOpen, setIsOrderSuccessOpen] = useState(false);
  const [checkoutEmail, setCheckoutEmail] = useState("");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");

  // 点击卡片 -> 打开商品详情弹窗
  const handleViewDetail = useCallback((product: Product) => {
    setSelectedProduct(product);
    setIsDetailOpen(true);
  }, []);

  // 从详情弹窗或卡片直接 -> 打开收银台
  const handlePurchase = useCallback((product: Product) => {
    setSelectedProduct(product);
    setIsDetailOpen(false);
    setIsDrawerOpen(true);
  }, []);

  // 收银台支付成功 -> 打开履约成功页
  const handlePaymentSuccess = useCallback((email: string) => {
    setCheckoutEmail(email);
    setIsDrawerOpen(false);
    setIsOrderSuccessOpen(true);
  }, []);

  // 关闭成功页 -> 回到首页
  const handleCloseSuccess = useCallback(() => {
    setIsOrderSuccessOpen(false);
    setSelectedProduct(null);
    setCheckoutEmail("");
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="flex">
        <Sidebar />
        <main className="md:ml-64 flex-1 p-4 md:p-6">
          {/* Toolbar: Search + Sort + View Toggle */}
          <div className="mb-6 flex items-center justify-between gap-4">
            <div className="flex-1 max-w-2xl">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  type="search"
                  placeholder="输入关键字搜索商品..."
                  className="pl-9 h-10 w-full"
                />
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" className="h-9">
                综合排序
              </Button>
              <div className="flex items-center border rounded-md">
                <Button
                  variant={viewMode === "grid" ? "default" : "ghost"}
                  size="sm"
                  className="h-9 rounded-r-none"
                  onClick={() => setViewMode("grid")}
                >
                  <Grid className="h-4 w-4" />
                </Button>
                <Button
                  variant={viewMode === "list" ? "default" : "ghost"}
                  size="sm"
                  className="h-9 rounded-l-none"
                  onClick={() => setViewMode("list")}
                >
                  <List className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>

          {/* Product Grid */}
          <div
            className={
              viewMode === "grid"
                ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
                : "space-y-4"
            }
          >
            {mockProducts.map((product, index) => (
              <ProductCard
                key={product.id}
                product={product}
                index={index}
                onPurchase={handlePurchase}
                onViewDetail={handleViewDetail}
              />
            ))}
          </div>

          {/* Product Detail Modal */}
          <ProductDetailModal
            open={isDetailOpen}
            onOpenChange={setIsDetailOpen}
            product={selectedProduct}
            onPurchase={handlePurchase}
          />

          {/* Checkout Drawer */}
          <CheckoutDrawer
            open={isDrawerOpen}
            onOpenChange={setIsDrawerOpen}
            product={selectedProduct}
            onPaymentSuccess={handlePaymentSuccess}
          />

          {/* Order Success */}
          <OrderSuccess
            open={isOrderSuccessOpen}
            onClose={handleCloseSuccess}
            product={selectedProduct}
            email={checkoutEmail}
          />
        </main>
      </div>
    </div>
  );
}
