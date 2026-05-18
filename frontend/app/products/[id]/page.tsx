"use client";

import { use, useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { ArrowLeft, Heart, Shield, Truck, RotateCcw, ShoppingBag } from "lucide-react";
import { Product, getProduct } from "@/lib/api";
import { getProductImage } from "@/lib/imageMap";
import { useCart } from "@/context/CartContext";

export default function ProductPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedSize, setSelectedSize] = useState("");
  const [added, setAdded] = useState(false);
  const [sizeError, setSizeError] = useState(false);
  const [wishlisted, setWishlisted] = useState(false);
  const { addItem } = useCart();

  useEffect(() => {
    getProduct(id)
      .then(setProduct)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id]);

  const handleAdd = () => {
    if (!product) return;
    if (product.sizes?.length > 0 && !selectedSize) {
      setSizeError(true);
      return;
    }
    addItem(product, selectedSize || "One Size");
    setAdded(true);
    setSizeError(false);
    setTimeout(() => setAdded(false), 2000);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-[#F5F5F0]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 py-10">
          <div className="h-4 w-28 bg-gray-200 animate-pulse rounded-full mb-10" />
          <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
            <div className="bg-white rounded-3xl aspect-square animate-pulse" />
            <div className="space-y-5 py-4">
              <div className="bg-gray-200 h-3 w-20 rounded-full animate-pulse" />
              <div className="bg-gray-200 h-10 w-3/4 rounded-xl animate-pulse" />
              <div className="bg-gray-200 h-12 w-32 rounded-xl animate-pulse" />
              <div className="flex gap-2 mt-4">
                {[...Array(6)].map((_, i) => (
                  <div key={i} className="bg-gray-200 h-12 w-14 rounded-xl animate-pulse" />
                ))}
              </div>
              <div className="bg-gray-200 h-14 rounded-2xl animate-pulse mt-4" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="min-h-screen bg-[#F5F5F0] flex items-center justify-center">
        <div className="text-center bg-white rounded-3xl p-16 shadow-sm mx-4">
          <div className="text-6xl mb-4">😕</div>
          <h2 className="text-xl font-black mb-2">Product not found</h2>
          <p className="text-gray-400 text-sm mb-6">
            This sneaker may have been removed or sold out
          </p>
          <Link
            href="/"
            className="bg-black text-white px-6 py-3 rounded-xl font-bold text-sm hover:bg-gray-800 transition-colors"
          >
            Back to Catalog
          </Link>
        </div>
      </div>
    );
  }

  const imgSrc = getProductImage(product.name, product.brand, product.image_url);
  const isOutOfStock = product.stock <= 0;
  const isLowStock = product.stock > 0 && product.stock <= 5;

  return (
    <div className="min-h-screen bg-[#F5F5F0]">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 py-8">
        {/* Back link */}
        <Link
          href="/"
          className="inline-flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-black transition-colors mb-8 group"
        >
          <ArrowLeft
            size={15}
            className="group-hover:-translate-x-0.5 transition-transform"
          />
          Back to catalog
        </Link>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 lg:gap-14">
          {/* ── Image ── */}
          <div className="relative">
            <div className="relative bg-white rounded-3xl overflow-hidden shadow-sm aspect-square">
              <Image
                src={imgSrc}
                alt={product.name}
                fill
                className="object-cover"
                sizes="(max-width: 768px) 100vw, 50vw"
                priority
              />
              {isLowStock && (
                <div className="absolute top-4 left-4">
                  <span className="bg-orange-500 text-white text-[10px] font-black px-3 py-1.5 rounded-full uppercase tracking-wider">
                    Only {product.stock} left
                  </span>
                </div>
              )}
              {isOutOfStock && (
                <div className="absolute inset-0 bg-black/40 backdrop-blur-[1px] flex items-center justify-center">
                  <span className="bg-white text-black font-black px-8 py-3 rounded-full text-sm">
                    Sold Out
                  </span>
                </div>
              )}
            </div>

            {/* Wishlist float */}
            <button
              onClick={() => setWishlisted((v) => !v)}
              className={`absolute top-4 right-4 w-10 h-10 rounded-full flex items-center justify-center shadow-md transition-all ${
                wishlisted ? "bg-red-500" : "bg-white/90 backdrop-blur-sm hover:bg-white"
              }`}
            >
              <Heart
                size={16}
                className={wishlisted ? "text-white fill-white" : "text-gray-600"}
                fill={wishlisted ? "currentColor" : "none"}
              />
            </button>
          </div>

          {/* ── Details ── */}
          <div className="flex flex-col py-2">
            <p className="text-[10px] font-black uppercase tracking-[0.25em] text-gray-400 mb-2">
              {product.brand}
            </p>
            <h1 className="text-3xl lg:text-4xl font-black text-gray-900 leading-tight mb-3">
              {product.name}
            </h1>

            {product.description && (
              <p className="text-gray-500 leading-relaxed mb-5 text-sm">{product.description}</p>
            )}

            <div className="text-4xl font-black text-gray-900 mb-7">
              ${product.price?.toLocaleString()}
            </div>

            {/* Stock indicator */}
            <div
              className={`flex items-center gap-2 mb-6 text-sm font-semibold ${
                isOutOfStock
                  ? "text-red-500"
                  : isLowStock
                  ? "text-orange-500"
                  : "text-green-600"
              }`}
            >
              <div
                className={`w-2 h-2 rounded-full ${
                  isOutOfStock
                    ? "bg-red-500"
                    : isLowStock
                    ? "bg-orange-500 animate-pulse"
                    : "bg-green-500"
                }`}
              />
              {isOutOfStock
                ? "Out of stock"
                : isLowStock
                ? `Only ${product.stock} pairs left`
                : `In stock · ${product.stock} pairs available`}
            </div>

            {/* Size selector */}
            {product.sizes?.length > 0 && (
              <div className="mb-6">
                <div className="flex items-center justify-between mb-3">
                  <p className="text-sm font-black text-gray-900">Size (EU)</p>
                  {sizeError && (
                    <p className="text-xs text-red-500 font-semibold">Please select a size</p>
                  )}
                </div>
                <div className="grid grid-cols-5 sm:grid-cols-6 gap-2">
                  {product.sizes.map((s) => (
                    <button
                      key={s}
                      onClick={() => {
                        setSelectedSize(s);
                        setSizeError(false);
                      }}
                      className={`h-12 rounded-xl font-bold text-sm transition-all ${
                        selectedSize === s
                          ? "bg-black text-white shadow-lg scale-[1.05]"
                          : sizeError
                          ? "border-2 border-red-300 text-gray-700 bg-white hover:border-gray-400"
                          : "bg-white border border-gray-200 text-gray-700 hover:border-black"
                      }`}
                    >
                      {s}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Add to Cart */}
            <button
              onClick={handleAdd}
              disabled={isOutOfStock}
              className={`w-full py-4 rounded-2xl font-black text-sm tracking-wide transition-all flex items-center justify-center gap-2 ${
                isOutOfStock
                  ? "bg-gray-100 text-gray-400 cursor-not-allowed"
                  : added
                  ? "bg-green-500 text-white scale-[0.99]"
                  : "bg-black text-white hover:bg-gray-800 active:scale-[0.99]"
              }`}
            >
              {added ? (
                <>
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="3"
                  >
                    <polyline points="20 6 9 17 4 12" />
                  </svg>
                  Added to Cart!
                </>
              ) : (
                <>
                  <ShoppingBag size={16} />
                  {isOutOfStock ? "Sold Out" : "Add to Cart"}
                </>
              )}
            </button>

            {/* Trust badges */}
            <div className="mt-8 grid grid-cols-3 gap-3 border-t border-gray-200 pt-8">
              {[
                {
                  icon: <Shield size={18} className="text-gray-600" />,
                  label: "Authentic",
                  sub: "100% verified",
                },
                {
                  icon: <Truck size={18} className="text-gray-600" />,
                  label: "Free shipping",
                  sub: "Orders over $100",
                },
                {
                  icon: <RotateCcw size={18} className="text-gray-600" />,
                  label: "30-day returns",
                  sub: "Hassle-free",
                },
              ].map((f) => (
                <div key={f.label} className="text-center bg-white rounded-2xl p-3">
                  <div className="flex justify-center mb-1.5">{f.icon}</div>
                  <p className="text-[11px] font-black text-gray-700">{f.label}</p>
                  <p className="text-[10px] text-gray-400 mt-0.5">{f.sub}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
