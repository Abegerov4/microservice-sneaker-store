"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Order, OrderItem, getOrdersByUser, getOrder, getOrderItems } from "@/lib/api";
import { useAuth } from "@/context/AuthContext";

const STATUS_STEPS = ["pending", "confirmed", "shipped", "delivered"];

const STATUS_META: Record<string, { label: string; color: string; bg: string }> = {
  pending:   { label: "Pending",   color: "text-amber-600",  bg: "bg-amber-100" },
  confirmed: { label: "Confirmed", color: "text-blue-600",   bg: "bg-blue-100" },
  shipped:   { label: "Shipped",   color: "text-violet-600", bg: "bg-violet-100" },
  delivered: { label: "Delivered", color: "text-green-600",  bg: "bg-green-100" },
  cancelled: { label: "Cancelled", color: "text-red-600",    bg: "bg-red-100" },
};

function OrderStatusBar({ status }: { status: string }) {
  if (status === "cancelled") {
    return (
      <div className="flex items-center gap-2 mt-3">
        <span className="w-2 h-2 rounded-full bg-red-500" />
        <span className="text-xs text-red-500 font-medium">Order cancelled</span>
      </div>
    );
  }
  const currentStep = STATUS_STEPS.indexOf(status);
  return (
    <div className="flex items-center gap-0 mt-4">
      {STATUS_STEPS.map((step, i) => {
        const done = i <= currentStep;
        const active = i === currentStep;
        return (
          <div key={step} className="flex items-center flex-1 last:flex-none">
            <div className={`w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-bold transition-all flex-shrink-0 ${
              done ? "bg-black text-white" : "bg-gray-100 text-gray-400"
            } ${active ? "ring-2 ring-black ring-offset-2" : ""}`}>
              {done && i < currentStep ? (
                <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="3"><polyline points="20 6 9 17 4 12"/></svg>
              ) : (
                i + 1
              )}
            </div>
            {i < STATUS_STEPS.length - 1 && (
              <div className={`flex-1 h-0.5 mx-1 transition-all ${i < currentStep ? "bg-black" : "bg-gray-100"}`} />
            )}
          </div>
        );
      })}
    </div>
  );
}

export default function OrdersPage() {
  const { isLoggedIn, loaded, userId, user } = useAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<string | null>(null);
  const [freshOrders, setFreshOrders] = useState<Record<string, Order>>({});
  const [freshItems, setFreshItems] = useState<Record<string, OrderItem[]>>({});
  const router = useRouter();

  useEffect(() => {
    if (!loaded) return;
    if (!isLoggedIn || !userId) { router.push("/auth/login"); return; }
    getOrdersByUser(userId).then(setOrders).catch(console.error).finally(() => setLoading(false));
  }, [loaded, isLoggedIn, userId, router]);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#F5F5F0]">
        <div className="max-w-2xl mx-auto px-6 py-10 space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-white rounded-3xl h-48 animate-pulse shadow-sm" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#F5F5F0]">
      <div className="max-w-2xl mx-auto px-6 py-10">
        <div className="mb-8">
          <h1 className="text-2xl font-black">My Orders</h1>
          {user && <p className="text-gray-400 text-sm mt-1">{user.email}</p>}
        </div>

        {orders.length === 0 ? (
          <div className="bg-white rounded-3xl p-16 text-center shadow-sm">
            <div className="text-6xl mb-4">📦</div>
            <h2 className="text-lg font-bold mb-2">No orders yet</h2>
            <p className="text-gray-400 text-sm mb-8">Looks like you haven&apos;t ordered anything yet</p>
            <Link href="/" className="bg-black text-white px-8 py-3 rounded-2xl font-semibold hover:bg-gray-800 transition-colors text-sm">
              Start Shopping
            </Link>
          </div>
        ) : (
          <div className="space-y-4">
            {orders.map((order) => {
              const meta = STATUS_META[order.status] ?? { label: order.status, color: "text-gray-600", bg: "bg-gray-100" };
              const isOpen = expanded === order.id;
              const displayOrder = freshOrders[order.id] ?? order;
              const displayItems = freshItems[order.id] ?? displayOrder.items;
              return (
                <div key={order.id} className="bg-white rounded-3xl shadow-sm overflow-hidden">
                  <button
                    onClick={async () => {
                      const willOpen = !isOpen;
                      setExpanded(willOpen ? order.id : null);
                      if (willOpen && !freshOrders[order.id]) {
                        try {
                          const [fresh, items] = await Promise.all([
                            getOrder(order.id),
                            getOrderItems(order.id),
                          ]);
                          setFreshOrders((prev) => ({ ...prev, [order.id]: fresh }));
                          setFreshItems((prev) => ({ ...prev, [order.id]: items }));
                        } catch {
                          // fall back to cached data
                        }
                      }
                    }}
                    className="w-full text-left p-6"
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-mono text-xs text-gray-400 mb-1">#{order.id?.slice(0, 8).toUpperCase()}</p>
                        <p className="font-bold text-gray-900">
                          {order.items?.length || 0} item{(order.items?.length || 0) !== 1 ? "s" : ""}
                          <span className="text-gray-400 font-normal ml-2">·</span>
                          <span className="font-black ml-2">${order.total_amount?.toLocaleString()}</span>
                        </p>
                        <p className="text-xs text-gray-400 mt-1">
                          {order.created_at ? new Date(order.created_at).toLocaleDateString("en-US", { day: "numeric", month: "long", year: "numeric" }) : "—"}
                        </p>
                      </div>
                      <div className="flex items-center gap-3">
                        <span className={`text-xs font-semibold px-3 py-1 rounded-full ${meta.bg} ${meta.color}`}>
                          {meta.label}
                        </span>
                        <svg className={`text-gray-400 transition-transform ${isOpen ? "rotate-180" : ""}`} width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <polyline points="6 9 12 15 18 9"/>
                        </svg>
                      </div>
                    </div>
                    <OrderStatusBar status={order.status} />
                  </button>

                  {isOpen && (
                    <div className="border-t border-gray-50 px-6 pb-6">
                      <div className="pt-4 space-y-2">
                        {displayItems?.map((item, idx) => (
                          <div key={idx} className="flex justify-between text-sm">
                            <span className="text-gray-600">
                              {item.product_name} ×{item.quantity}
                              {item.size && <span className="text-gray-400"> · EU {item.size}</span>}
                            </span>
                            <span className="font-semibold">${(item.price * item.quantity).toLocaleString()}</span>
                          </div>
                        ))}
                      </div>
                      {displayOrder.shipping_address && (
                        <div className="mt-4 pt-4 border-t border-gray-50">
                          <p className="text-xs text-gray-400 mb-1">Shipping to</p>
                          <p className="text-sm text-gray-700">{displayOrder.shipping_address}</p>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
