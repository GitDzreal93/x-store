"use client";

import { Search, ShoppingBag, User, Moon, Sun } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { MobileSidebar } from "./MobileSidebar";

export function Navbar() {
  const [isDark, setIsDark] = useState(false);

  const toggleTheme = () => {
    setIsDark(!isDark);
    document.documentElement.classList.toggle("dark");
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between px-4">
        {/* Logo & Brand */}
        <div className="flex items-center gap-4">
          <MobileSidebar />
          <div className="flex items-center gap-2">
            <span className="text-xl font-bold">x-store</span>
          </div>
          <nav className="hidden md:flex items-center gap-6 text-sm">
            <a href="#" className="text-muted-foreground hover:text-foreground transition-colors">
              首页
            </a>
            <a href="#" className="text-muted-foreground hover:text-foreground transition-colors">
              API中心
            </a>
            <a href="#" className="text-muted-foreground hover:text-foreground transition-colors">
              使用文档
            </a>
          </nav>
        </div>

        {/* Right Actions */}
        <div className="flex items-center gap-4">
          {/* Global Search - Desktop */}
          <div className="hidden lg:flex items-center gap-2 relative w-64">
            <Search className="absolute left-3 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="搜索您感兴趣的内容..."
              className="pl-9 h-9 bg-muted/50"
            />
            <kbd className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
              <span className="text-xs">⌘</span>K
            </kbd>
          </div>
          
          {/* Mobile Search Button */}
          <Button variant="ghost" size="icon" className="lg:hidden">
            <Search className="h-4 w-4" />
          </Button>

          {/* Order Query */}
          <Button variant="ghost" size="sm" className="hidden md:flex">
            <ShoppingBag className="h-4 w-4 mr-2" />
            查单
          </Button>

          {/* Theme Toggle */}
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleTheme}
            className="h-9 w-9"
          >
            {isDark ? (
              <Sun className="h-4 w-4" />
            ) : (
              <Moon className="h-4 w-4" />
            )}
          </Button>

          {/* User/Login */}
          <Button variant="ghost" size="sm">
            <User className="h-4 w-4 mr-2" />
            登录
          </Button>
        </div>
      </div>
    </header>
  );
}
