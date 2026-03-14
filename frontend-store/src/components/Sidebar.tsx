"use client";

import { Search, Package, Bot, Gamepad2, Gift, Sparkles } from "lucide-react";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

interface Category {
  id: string;
  name: string;
  icon: React.ReactNode;
  count?: number;
}

const categories: Category[] = [
  { id: "streaming", name: "流媒体账号", icon: <Package className="h-4 w-4" />, count: 1250 },
  { id: "ai", name: "AI 大模型", icon: <Bot className="h-4 w-4" />, count: 342 },
  { id: "gaming", name: "游戏充值", icon: <Gamepad2 className="h-4 w-4" />, count: 89 },
  { id: "giftcard", name: "礼品卡", icon: <Gift className="h-4 w-4" />, count: 156 },
];

const attributes = [
  { id: "auto", name: "自动发货", icon: <Sparkles className="h-4 w-4" /> },
  { id: "manual", name: "人工代充", icon: <Package className="h-4 w-4" /> },
];

export function Sidebar() {
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");

  return (
    <aside className="hidden md:block fixed left-0 top-16 h-[calc(100vh-4rem)] w-64 border-r bg-background overflow-y-auto">
      <div className="p-4 space-y-6">
        {/* Sidebar Search */}
        <div className="space-y-2">
          <Input
            type="search"
            placeholder="输入关键词搜索您想要的模型"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="h-9"
          />
          <p className="text-xs text-muted-foreground">(共 1,837 个商品)</p>
        </div>

        {/* 热门分类 */}
        <div className="space-y-2">
          <h3 className="text-sm font-semibold text-foreground">热门分类</h3>
          <div className="space-y-1">
            {categories.map((category) => (
              <button
                key={category.id}
                onClick={() => setSelectedCategory(category.id)}
                className={cn(
                  "w-full flex items-center justify-between px-3 py-2 rounded-md text-sm transition-colors",
                  selectedCategory === category.id
                    ? "bg-accent text-accent-foreground"
                    : "text-muted-foreground hover:bg-accent/50 hover:text-foreground"
                )}
              >
                <div className="flex items-center gap-2">
                  {category.icon}
                  <span>{category.name}</span>
                </div>
                {category.count && (
                  <span className="text-xs text-muted-foreground">
                    {category.count}
                  </span>
                )}
              </button>
            ))}
          </div>
        </div>

        {/* 商品属性 */}
        <div className="space-y-2">
          <h3 className="text-sm font-semibold text-foreground">商品属性</h3>
          <div className="space-y-1">
            {attributes.map((attr) => (
              <button
                key={attr.id}
                className="w-full flex items-center gap-2 px-3 py-2 rounded-md text-sm text-muted-foreground hover:bg-accent/50 hover:text-foreground transition-colors"
              >
                {attr.icon}
                <span>{attr.name}</span>
              </button>
            ))}
          </div>
        </div>

        {/* 其他 */}
        <div className="space-y-2">
          <h3 className="text-sm font-semibold text-foreground">其他</h3>
          <button className="w-full flex items-center gap-2 px-3 py-2 rounded-md text-sm text-muted-foreground hover:bg-accent/50 hover:text-foreground transition-colors">
            <Gift className="h-4 w-4" />
            <span>礼品卡</span>
          </button>
        </div>
      </div>
    </aside>
  );
}
