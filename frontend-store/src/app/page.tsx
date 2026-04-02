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
import { Search, Grid, List, Loader2, AlertCircle } from "lucide-react";
import { getProducts, getCategories, transformProduct, Category } from "@/lib/api";
import { ApiError } from "@/lib/api-client";

export default function Home() {
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isOrderSuccessOpen, setIsOrderSuccessOpen] = useState(false);
  const [checkoutEmail, setCheckoutEmail] = useState("");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");

  // 数据状态
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null);
  const [searchKeyword, setSearchKeyword] = useState("");

  // 加载状态
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 分页状态
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  // 加载商品列表
  const loadProducts = useCallback(async (pageNum: number = 1) => {
    try {
      setIsLoading(true);
      setError(null);

      const response = await getProducts({
        page: pageNum,
        size: 20,
        category_id: selectedCategory ?? undefined,
        keyword: searchKeyword || undefined,
      });

      const transformedProducts = response.list.map(transformProduct);

      if (pageNum === 1) {
        setProducts(transformedProducts);
      } else {
        setProducts((prev) => [...prev, ...transformedProducts]);
      }

      setTotal(response.total);
      setHasMore(transformedProducts.length === response.size && products.length + transformedProducts.length < response.total);
      setPage(pageNum);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("加载商品失败，请稍后重试");
      }
    } finally {
      setIsLoading(false);
    }
  }, [selectedCategory, searchKeyword]);

  // 加载分类列表
  useEffect(() => {
    const loadCategories = async () => {
      try {
        const data = await getCategories();
        setCategories(data);
      } catch (err) {
        console.error("加载分类失败:", err);
      }
    };
    loadCategories();
  }, []);

  // 初始加载商品
  useEffect(() => {
    loadProducts(1);
  }, [loadProducts]);

  // 搜索防抖
  useEffect(() => {
    const timer = setTimeout(() => {
      loadProducts(1);
    }, 500);

    return () => clearTimeout(timer);
  }, [searchKeyword]);

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
  const [orderNo, setOrderNo] = useState<string>();
  const handlePaymentSuccess = useCallback((orderNo: string, email: string) => {
    setOrderNo(orderNo);
    setCheckoutEmail(email);
    setIsDrawerOpen(false);
    setIsOrderSuccessOpen(true);
  }, []);

  // 关闭成功页 -> 回到首页
  const handleCloseSuccess = useCallback(() => {
    setIsOrderSuccessOpen(false);
    setSelectedProduct(null);
    setCheckoutEmail("");
    setOrderNo(undefined);
  }, []);

  // 处理分类选择
  const handleCategoryChange = (categoryId: number | null) => {
    setSelectedCategory(categoryId);
    setPage(1);
  };

  // 加载更多
  const handleLoadMore = () => {
    if (!isLoading && hasMore) {
      loadProducts(page + 1);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="flex">
        <Sidebar />
        <main className="md:ml-64 flex-1 p-4 md:p-6">
          {/* Toolbar: Search + Sort + View Toggle */}
          <div className="mb-6 space-y-4">
            {/* 分类筛选 */}
            {categories.length > 0 && (
              <div className="flex items-center gap-2 flex-wrap">
                <Button
                  variant={selectedCategory === null ? "default" : "outline"}
                  size="sm"
                  onClick={() => handleCategoryChange(null)}
                  className="h-8"
                >
                  全部
                </Button>
                {categories.map((category) => (
                  <Button
                    key={category.id}
                    variant={selectedCategory === category.id ? "default" : "outline"}
                    size="sm"
                    onClick={() => handleCategoryChange(category.id)}
                    className="h-8"
                  >
                    {category.icon} {category.name}
                  </Button>
                ))}
              </div>
            )}

            {/* 搜索和视图切换 */}
            <div className="flex items-center justify-between gap-4">
              <div className="flex-1 max-w-2xl">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="search"
                    placeholder="输入关键字搜索商品..."
                    className="pl-9 h-10 w-full"
                    value={searchKeyword}
                    onChange={(e) => setSearchKeyword(e.target.value)}
                  />
                </div>
              </div>
              <div className="flex items-center gap-2">
                <div className="text-sm text-muted-foreground">
                  共 {total} 件商品
                </div>
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
          </div>

          {/* 错误提示 */}
          {error && (
            <div className="flex items-center gap-2 p-4 bg-destructive/10 text-destructive rounded-md mb-6">
              <AlertCircle className="h-5 w-5" />
              <div>
                <p className="font-medium">加载失败</p>
                <p className="text-sm">{error}</p>
              </div>
              <Button
                variant="outline"
                size="sm"
                className="ml-auto"
                onClick={() => loadProducts(1)}
              >
                重试
              </Button>
            </div>
          )}

          {/* 商品列表 */}
          {isLoading && products.length === 0 ? (
            <div className="flex items-center justify-center py-20">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
              <span className="ml-3 text-muted-foreground">加载中...</span>
            </div>
          ) : products.length === 0 ? (
            <div className="flex items-center justify-center py-20 text-center">
              <div>
                <p className="text-lg font-medium text-muted-foreground">暂无商品</p>
                <p className="text-sm text-muted-foreground mt-1">
                  {selectedCategory || searchKeyword
                    ? "试试其他分类或搜索关键词"
                    : "敬请期待"}
                </p>
              </div>
            </div>
          ) : (
            <>
              <div
                className={
                  viewMode === "grid"
                    ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
                    : "space-y-4"
                }
              >
                {products.map((product, index) => (
                  <ProductCard
                    key={product.id}
                    product={product}
                    index={index}
                    onPurchase={handlePurchase}
                    onViewDetail={handleViewDetail}
                  />
                ))}
              </div>

              {/* 加载更多 */}
              {hasMore && (
                <div className="flex justify-center mt-8">
                  <Button
                    variant="outline"
                    size="lg"
                    onClick={handleLoadMore}
                    disabled={isLoading}
                    className="min-w-32"
                  >
                    {isLoading ? (
                      <>
                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                        加载中...
                      </>
                    ) : (
                      "加载更多"
                    )}
                  </Button>
                </div>
              )}

              {!hasMore && products.length > 0 && (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  已加载全部商品
                </div>
              )}
            </>
          )}

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
            onPaymentError={(error) => {
              console.error("Payment error:", error);
              // 可以在这里添加 toast 提示
            }}
          />

          {/* Order Success */}
          <OrderSuccess
            open={isOrderSuccessOpen}
            onClose={handleCloseSuccess}
            product={selectedProduct}
            email={checkoutEmail}
            orderNo={orderNo}
          />
        </main>
      </div>
    </div>
  );
}
