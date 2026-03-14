"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Package, Download, Heart, Clock } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";
import { motion } from "framer-motion";

export interface Product {
  id: string;
  title: string;
  tags: string[];
  price: number;
  sales: number;
  likes?: number;
  updatedAt: string;
  isNew?: boolean;
  isHot?: boolean;
}

interface ProductCardProps {
  product: Product;
  onPurchase: (product: Product) => void;
  onViewDetail?: (product: Product) => void;
  index?: number;
}

export function ProductCard({ product, onPurchase, onViewDetail, index = 0 }: ProductCardProps) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, delay: index * 0.06, ease: "easeOut" }}
      whileHover={{ y: -2 }}
    >
    <Card
      className={cn(
        "group relative overflow-hidden border transition-all duration-200 cursor-pointer",
        "hover:border-primary/50 hover:shadow-md",
        isHovered && "border-primary/30"
      )}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onClick={() => onViewDetail?.(product)}
    >
      <div className="p-4 space-y-3">
        {/* Header: Title + Badges */}
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-semibold text-base leading-tight flex-1">
            {product.title}
          </h3>
          <div className="flex gap-1 flex-shrink-0">
            {product.isNew && (
              <Badge variant="secondary" className="text-xs px-2 py-0">
                NEW
              </Badge>
            )}
            {product.isHot && (
              <Badge variant="destructive" className="text-xs px-2 py-0">
                HOT
              </Badge>
            )}
          </div>
        </div>

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
        </div>

        {/* Footer: Price + Stats + Action */}
        <div className="flex items-center justify-between pt-2 border-t">
          <div className="flex items-center gap-4 text-sm">
            <span className="font-bold text-lg text-foreground">
              ¥ {product.price.toFixed(2)}
            </span>
            <div className="flex items-center gap-4 text-muted-foreground">
              <div className="flex items-center gap-1">
                <Download className="h-3.5 w-3.5" />
                <span className="text-xs">
                  {product.sales >= 1000
                    ? `${(product.sales / 1000).toFixed(1)}k`
                    : product.sales}
                </span>
              </div>
              {product.likes !== undefined && (
                <div className="flex items-center gap-1">
                  <Heart className="h-3.5 w-3.5" />
                  <span className="text-xs">{product.likes}</span>
                </div>
              )}
              <div className="flex items-center gap-1">
                <Clock className="h-3.5 w-3.5" />
                <span className="text-xs">{product.updatedAt}</span>
              </div>
            </div>
          </div>

          {/* Purchase Button - Appears on Hover */}
          <Button
            size="sm"
            className={cn(
              "transition-all duration-200",
              isHovered
                ? "opacity-100 translate-x-0"
                : "opacity-0 translate-x-2 pointer-events-none"
            )}
            onClick={(e) => {
              e.stopPropagation();
              onPurchase(product);
            }}
          >
            购买
          </Button>
        </div>
      </div>
    </Card>
    </motion.div>
  );
}
