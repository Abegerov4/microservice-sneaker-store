"use client";

import Image from "next/image";
import Link from "next/link";
import { useState } from "react";
import { Heart, ShoppingBag } from "lucide-react";
import { Product } from "@/lib/api";
import { getProductImage } from "@/lib/imageMap";
import { useCart } from "@/context/CartContext";

export default function ProductCard({ product }: { product: Product }) {
  const { addItem } = useCart();
  const [selectedSize, setSelectedSize] = useState<string>("");
  const [added, setAdded] = useState(false);
  const [wishlisted, setWishlisted] = useState(false);
  const imgSrc = getProductImage(product.name, product.brand, product.image_url);

  const handleAdd = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    const size = selectedSize || product.sizes?.[0] || "One Size";
    addItem(product, size);
    setAdded(true);
    setTimeout(() => setAdded(false), 1600);
  };

  const isOutOfStock = product.stock <= 0;
  const isLowStock = product.stock > 0 && product.stock <= 5;

  return (
    <Link href={`/products/${product.id}`} className="block group">
      <div className="bg-white rounded-2xl overflow-hidden transition-all duration-400 group-hover:shadow-[0_12px_48px_rgba(0,0,0,0.12)] group-hover:-translate-y-1 h-full flex flex-col">
        {/* Image */}
        <div className="relative bg-[#F5F5F0] overflow-hidden">
          <div className="aspect-square relative">
            <Image
              src={imgSrc}
              alt={product.name}
              fill
              className={`object-cover transition-transform duration-700 group-hover:scale-[1.04] ${
                isOutOfStock ? "opacity-50 grayscale" : ""
              }`}
              sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 25vw"
            />
          </div>

          {/* Wishlist button */}
          <button
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              setWishlisted((v) => !v);
            }}
            className={`absolute top-3 right-3 w-8 h-8 rounded-full flex items-center justify-center transition-all duration-200 shadow-sm ${
              wishlisted
                ? "bg-red-500 opacity-100 scale-100"
                : "bg-white/90 backdrop-blur-sm opacity-0 group-hover:opacity-100 scale-90 group-hover:scale-100"
            }`}
          >
            <Heart
              size={14}
              className={wishlisted ? "text-white fill-white" : "text-gray-600"}
              fill={wishlisted ? "currentColor" : "none"}
            />
          </button>

          {/* Badges */}
          {isLowStock && !isOutOfStock && (
            <div className="absolute top-3 left-3">
              <span className="bg-orange-500 text-white text-[9px] font-black px-2.5 py-1 rounded-full uppercase tracking-wider shadow-sm">
                {product.stock} left
              </span>
            </div>
          )}

          {isOutOfStock && (
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="bg-black/75 backdrop-blur-sm text-white text-xs font-bold px-4 py-2 rounded-full">
                Sold Out
              </span>
            </div>
          )}

          {/* Quick-add panel — slides up on hover */}
          {!isOutOfStock && (
            <div
              className="absolute inset-x-0 bottom-0 translate-y-full group-hover:translate-y-0 transition-transform duration-350 ease-out"
              onClick={(e) => e.preventDefault()}
            >
              {product.sizes && product.sizes.length > 0 ? (
                <div className="bg-white/96 backdrop-blur-md p-3 border-t border-gray-100">
                  <p className="text-[9px] font-black uppercase tracking-[0.2em] text-gray-400 mb-2">
                    Select Size (EU)
                  </p>
                  <div className="flex flex-wrap gap-1 mb-2.5">
                    {product.sizes.slice(0, 7).map((s) => (
                      <button
                        key={s}
                        onClick={(e) => {
                          e.preventDefault();
                          e.stopPropagation();
                          setSelectedSize(s === selectedSize ? "" : s);
                        }}
                        className={`text-[11px] w-9 h-8 rounded-lg font-semibold transition-all ${
                          selectedSize === s
                            ? "bg-black text-white"
                            : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                        }`}
                      >
                        {s}
                      </button>
                    ))}
                  </div>
                  <button
                    onClick={handleAdd}
                    className={`w-full py-2.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center justify-center gap-2 ${
                      added
                        ? "bg-green-500 text-white"
                        : "bg-black text-white hover:bg-gray-800"
                    }`}
                  >
                    {added ? (
                      "✓ Added to Cart"
                    ) : (
                      <>
                        <ShoppingBag size={12} />
                        Add to Cart
                      </>
                    )}
                  </button>
                </div>
              ) : (
                <div className="p-2">
                  <button
                    onClick={handleAdd}
                    className={`w-full py-2.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center justify-center gap-2 shadow-lg ${
                      added
                        ? "bg-green-500 text-white"
                        : "bg-black/90 backdrop-blur-sm text-white hover:bg-black"
                    }`}
                  >
                    {added ? "✓ Added" : <><ShoppingBag size={12} /> Add to Cart</>}
                  </button>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Info */}
        <div className="p-4 flex flex-col flex-1">
          <p className="text-[9px] font-black uppercase tracking-[0.22em] text-gray-400 mb-1">
            {product.brand}
          </p>
          <h3 className="font-semibold text-gray-900 text-[13px] leading-snug line-clamp-2 flex-1 mb-3">
            {product.name}
          </h3>
          <div className="flex items-center justify-between">
            <span className="text-[15px] font-black text-gray-900">
              ${product.price?.toLocaleString()}
            </span>
            {!isOutOfStock && (
              <span className="text-[10px] text-gray-400 font-medium">
                {product.sizes?.length
                  ? `${product.sizes.length} sizes`
                  : `Stock: ${product.stock}`}
              </span>
            )}
            {isOutOfStock && (
              <span className="text-[10px] text-red-400 font-semibold">Sold Out</span>
            )}
          </div>
        </div>
      </div>
    </Link>
  );
}
