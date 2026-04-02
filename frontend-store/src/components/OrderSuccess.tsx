"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Check,
  Copy,
  Mail,
  ArrowLeft,
  Eye,
  EyeOff,
  ExternalLink,
} from "lucide-react";
import { Product } from "./ProductCard";

interface OrderSuccessProps {
  open: boolean;
  onClose: () => void;
  product: Product | null;
  email: string;
  orderNo?: string;
}

// Mock card key for demo (生产环境应从后端 API 获取)
const MOCK_CARD_KEY = "XSTORE-GPT4-A8K2-M9XP-7LQR";

export function OrderSuccess({
  open,
  onClose,
  product,
  email,
  orderNo,
}: OrderSuccessProps) {
  const [isRevealed, setIsRevealed] = useState(false);
  const [isCopied, setIsCopied] = useState(false);
  const [scratchProgress, setScratchProgress] = useState(0);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(MOCK_CARD_KEY);
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    } catch {
      // Fallback
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    }
  };

  const handleReveal = () => {
    setIsRevealed(true);
  };

  const handleScratch = () => {
    // Simulate scratch progress
    setScratchProgress((prev) => {
      const next = Math.min(prev + 25, 100);
      if (next >= 100) {
        setTimeout(() => setIsRevealed(true), 300);
      }
      return next;
    });
  };

  if (!open || !product) return null;

  return (
    <AnimatePresence>
      {open && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
        >
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            transition={{ type: "spring", damping: 25, stiffness: 300 }}
            className="w-full max-w-md"
          >
            <Card className="overflow-hidden border-0 shadow-2xl">
              {/* Success Header */}
              <div className="bg-gradient-to-br from-green-500 to-emerald-600 p-6 text-white text-center">
                {/* Animated Checkmark */}
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{
                    type: "spring",
                    damping: 10,
                    stiffness: 200,
                    delay: 0.2,
                  }}
                  className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-white/20 backdrop-blur-sm"
                >
                  <motion.div
                    initial={{ pathLength: 0, opacity: 0 }}
                    animate={{ pathLength: 1, opacity: 1 }}
                    transition={{ delay: 0.5, duration: 0.4 }}
                  >
                    <Check className="h-8 w-8 text-white" strokeWidth={3} />
                  </motion.div>
                </motion.div>

                <motion.h2
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.4 }}
                  className="text-xl font-bold"
                >
                  支付成功
                </motion.h2>
                <motion.p
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.5 }}
                  className="text-sm text-white/80 mt-1"
                >
                  您的数字商品已准备就绪
                </motion.p>
              </div>

              {/* Order Info */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.6 }}
                className="p-6 space-y-5"
              >
                {/* Order Details */}
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">订单号</span>
                  <span className="font-mono font-medium">{orderId}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">商品名称</span>
                  <span className="font-medium truncate ml-4 text-right">
                    {product.title}
                  </span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">支付金额</span>
                  <span className="font-bold text-foreground">
                    ¥ {product.price.toFixed(2)}
                  </span>
                </div>

                {/* Card Key Reveal Area */}
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 0.8 }}
                  className="relative overflow-hidden rounded-lg border bg-muted/30"
                >
                  <div className="p-4 space-y-3">
                    <div className="flex items-center justify-between">
                      <h3 className="text-sm font-semibold flex items-center gap-2">
                        <Badge
                          variant="secondary"
                          className="text-xs bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300"
                        >
                          卡密
                        </Badge>
                        数字资产
                      </h3>
                      <button
                        onClick={() =>
                          isRevealed
                            ? setIsRevealed(false)
                            : handleReveal()
                        }
                        className="text-muted-foreground hover:text-foreground transition-colors"
                      >
                        {isRevealed ? (
                          <EyeOff className="h-4 w-4" />
                        ) : (
                          <Eye className="h-4 w-4" />
                        )}
                      </button>
                    </div>

                    {/* Scratch Card / Flip Card Area */}
                    <div className="relative min-h-[60px] flex items-center justify-center">
                      {!isRevealed ? (
                        // Scratch Card Surface
                        <div
                          className="w-full cursor-pointer select-none"
                          onClick={handleScratch}
                        >
                          <div className="relative overflow-hidden rounded-md">
                            {/* Scratch overlay */}
                            <motion.div
                              className="absolute inset-0 bg-gradient-to-r from-slate-300 to-slate-400 dark:from-slate-600 dark:to-slate-700 rounded-md flex items-center justify-center z-10"
                              animate={{
                                opacity: 1 - scratchProgress / 100,
                              }}
                              transition={{ duration: 0.3 }}
                            >
                              <span className="text-sm font-medium text-white/90">
                                {scratchProgress === 0
                                  ? "🎁 点击刮开查看卡密"
                                  : `刮开进度 ${scratchProgress}%`}
                              </span>
                            </motion.div>
                            {/* Hidden content below */}
                            <div className="p-3 bg-muted/50 rounded-md">
                              <p className="font-mono text-sm text-center tracking-wider text-muted-foreground">
                                ****-****-****-****
                              </p>
                            </div>
                          </div>
                        </div>
                      ) : (
                        // Revealed Card Key with flip animation
                        <motion.div
                          initial={{ rotateY: 90, opacity: 0 }}
                          animate={{ rotateY: 0, opacity: 1 }}
                          transition={{
                            type: "spring",
                            damping: 20,
                            stiffness: 200,
                          }}
                          className="w-full"
                        >
                          <div className="p-3 bg-green-50 dark:bg-green-950/30 border border-green-200 dark:border-green-800 rounded-md">
                            <p className="font-mono text-sm text-center tracking-wider font-semibold text-green-700 dark:text-green-300">
                              {MOCK_CARD_KEY}
                            </p>
                          </div>
                        </motion.div>
                      )}
                    </div>

                    {/* Copy Button */}
                    {isRevealed && (
                      <motion.div
                        initial={{ opacity: 0, y: 5 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.2 }}
                      >
                        <Button
                          variant="outline"
                          size="sm"
                          className="w-full"
                          onClick={handleCopy}
                        >
                          {isCopied ? (
                            <>
                              <Check className="h-3.5 w-3.5 mr-2 text-green-500" />
                              已复制到剪贴板
                            </>
                          ) : (
                            <>
                              <Copy className="h-3.5 w-3.5 mr-2" />
                              点击复制卡密
                            </>
                          )}
                        </Button>
                      </motion.div>
                    )}
                  </div>
                </motion.div>

                {/* Email Notice */}
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 1 }}
                  className="flex items-center gap-2 p-3 rounded-md bg-blue-50 dark:bg-blue-950/30 text-sm"
                >
                  <Mail className="h-4 w-4 text-blue-500 flex-shrink-0" />
                  <span className="text-blue-700 dark:text-blue-300">
                    已同步发送至您的邮箱{" "}
                    <span className="font-medium">{email || "xxx@gmail.com"}</span>
                  </span>
                </motion.div>

                {/* Actions */}
                <div className="flex gap-3 pt-2">
                  <Button
                    variant="outline"
                    className="flex-1"
                    onClick={onClose}
                  >
                    <ArrowLeft className="h-4 w-4 mr-2" />
                    返回首页
                  </Button>
                  <Button
                    variant="default"
                    className="flex-1"
                    onClick={() => {
                      // 导航到订单详情页或个人中心订单页
                      if (orderNo) {
                        window.location.href = `/profile?order=${orderNo}`;
                      } else {
                        window.location.href = "/profile";
                      }
                    }}
                  >
                    <ExternalLink className="h-4 w-4 mr-2" />
                    查看订单
                  </Button>
                </div>
              </motion.div>
            </Card>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
