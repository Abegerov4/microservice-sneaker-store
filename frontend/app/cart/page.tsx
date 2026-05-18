"use client";

import Image from "next/image";
import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCart } from "@/context/CartContext";
import { useAuth } from "@/context/AuthContext";
import { getProductImage } from "@/lib/imageMap";
import { createOrder } from "@/lib/api";

export default function CartPage() {
  const { items, removeItem, updateQuantity, clearCart, total } = useCart();
  const { isLoggedIn, userId } = useAuth();
  const [address, setAddress] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const router = useRouter();

  const handleCheckout = async () => {
    if (!isLoggedIn || !userId) { router.push("/auth/login"); return; }
    if (!address.trim()) { setError("Please enter a shipping address"); return; }
    setError("");
    setLoading(true);
    try {
      await createOrder(userId, items.map((i) => ({ product_id: i.product.id, quantity: i.quantity, size: i.size })), address);
      clearCart();
      setSuccess(true);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to place order");
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="min-h-screen bg-[#F5F5F0] flex items-center justify-center px-4">
        <div className="bg-white rounded-3xl p-12 max-w-md w-full text-center shadow-sm">
          <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#22c55e" strokeWidth="2.5">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
          </div>
          <h1 className="text-2xl font-black mb-2">Order Placed!</h1>
          <p className="text-gray-400 mb-8 text-sm leading-relaxed">Your order has been confirmed. We&apos;ll send you updates as it ships.</p>
          <div className="flex gap-3">
            <Link href="/orders" className="flex-1 bg-black text-white py-3 rounded-2xl font-semibold text-sm hover:bg-gray-800 transition-colors text-center">
              View Orders
            </Link>
            <Link href="/" className="flex-1 border border-gray-200 py-3 rounded-2xl font-semibold text-sm hover:bg-gray-50 transition-colors text-center">
              Keep Shopping
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="min-h-screen bg-[#F5F5F0] flex items-center justify-center px-4">
        <div className="text-center">
          <div className="text-8xl mb-6">🛍️</div>
          <h1 className="text-2xl font-black mb-2">Your cart is empty</h1>
          <p className="text-gray-400 mb-8">Find something you love in our catalog</p>
          <Link href="/" className="bg-black text-white px-8 py-4 rounded-2xl font-semibold hover:bg-gray-800 transition-colors">
            Shop Now
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#F5F5F0]">
      <div className="max-w-5xl mx-auto px-6 py-10">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-2xl font-black">Cart <span className="text-gray-300">({items.length})</span></h1>
          <button onClick={clearCart} className="text-sm text-gray-400 hover:text-red-500 transition-colors">
            Clear all
          </button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Items */}
          <div className="lg:col-span-2 space-y-3">
            {items.map((item) => (
              <div key={`${item.product.id}-${item.size}`} className="bg-white rounded-2xl p-4 flex gap-4 shadow-sm">
                <Link href={`/products/${item.product.id}`} className="relative w-24 h-24 bg-gray-50 rounded-xl overflow-hidden flex-shrink-0 hover:opacity-90 transition-opacity">
                  <Image
                    src={getProductImage(item.product.name, item.product.brand, item.product.image_url)}
                    alt={item.product.name}
                    fill
                    className="object-cover"
                    sizes="96px"
                  />
                </Link>
                <div className="flex-1 min-w-0">
                  <p className="text-[10px] font-black uppercase tracking-widest text-gray-400">{item.product.brand}</p>
                  <p className="font-semibold text-gray-900 text-sm truncate mt-0.5">{item.product.name}</p>
                  <p className="text-xs text-gray-400 mt-0.5">EU {item.size}</p>
                  <div className="flex items-center gap-2 mt-3">
                    <button
                      onClick={() => updateQuantity(item.product.id, item.size, item.quantity - 1)}
                      className="w-7 h-7 rounded-full border border-gray-200 flex items-center justify-center text-gray-500 hover:border-black hover:text-black transition-colors text-sm font-bold"
                    >
                      −
                    </button>
                    <span className="w-5 text-center text-sm font-semibold">{item.quantity}</span>
                    <button
                      onClick={() => updateQuantity(item.product.id, item.size, item.quantity + 1)}
                      className="w-7 h-7 rounded-full border border-gray-200 flex items-center justify-center text-gray-500 hover:border-black hover:text-black transition-colors text-sm font-bold"
                    >
                      +
                    </button>
                  </div>
                </div>
                <div className="flex flex-col items-end justify-between flex-shrink-0">
                  <p className="font-black text-gray-900">${(item.product.price * item.quantity).toLocaleString()}</p>
                  <button
                    onClick={() => removeItem(item.product.id, item.size)}
                    className="text-gray-300 hover:text-red-500 transition-colors"
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <polyline points="3 6 5 6 21 6"/><path d="m19 6-.867 12.142A2 2 0 0116.138 20H7.862a2 2 0 01-1.995-1.858L5 6m5 0V4a1 1 0 011-1h2a1 1 0 011 1v2"/>
                    </svg>
                  </button>
                </div>
              </div>
            ))}
          </div>

          {/* Summary */}
          <div className="space-y-4">
            <div className="bg-white rounded-2xl p-6 shadow-sm">
              <h2 className="font-black text-lg mb-5">Summary</h2>

              <div className="space-y-2 mb-4">
                {items.map((i) => (
                  <div key={`sum-${i.product.id}-${i.size}`} className="flex justify-between text-sm">
                    <span className="text-gray-500 truncate mr-2">{i.product.name} ×{i.quantity}</span>
                    <span className="font-medium flex-shrink-0">${(i.product.price * i.quantity).toLocaleString()}</span>
                  </div>
                ))}
              </div>

              <div className="border-t border-gray-100 my-4" />

              <div className="flex justify-between font-black text-xl mb-2">
                <span>Total</span>
                <span>${total.toLocaleString()}</span>
              </div>
              <p className="text-xs text-gray-400 mb-6">Taxes and shipping calculated at checkout</p>

              <div className="mb-4">
                <label className="block text-sm font-semibold text-gray-700 mb-2">Shipping Address</label>
                <textarea
                  value={address}
                  onChange={(e) => setAddress(e.target.value)}
                  placeholder="City, street, building, apt."
                  rows={3}
                  className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black resize-none"
                />
              </div>

              {error && (
                <div className="bg-red-50 border border-red-100 rounded-xl px-4 py-3 text-red-600 text-sm mb-4">
                  {error}
                </div>
              )}

              {!isLoggedIn && (
                <p className="text-sm text-gray-400 mb-4">
                  <Link href="/auth/login" className="font-semibold text-black underline underline-offset-2">Sign in</Link> to place an order
                </p>
              )}

              <button
                onClick={handleCheckout}
                disabled={loading}
                className="w-full bg-black text-white py-4 rounded-2xl font-bold hover:bg-gray-800 disabled:opacity-50 transition-colors"
              >
                {loading ? "Placing Order..." : "Place Order"}
              </button>
            </div>

            <Link href="/" className="block text-center text-sm text-gray-400 hover:text-black transition-colors">
              ← Continue Shopping
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
