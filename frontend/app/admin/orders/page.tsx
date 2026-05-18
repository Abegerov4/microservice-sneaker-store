"use client";

import React, { useEffect, useState } from "react";

function parseGoDate(s: string): Date {
  if (!s) return new Date(NaN);
  const parts = s.split(" ");
  if (parts.length >= 3) return new Date(`${parts[0]}T${parts[1]}${parts[2]}`);
  return new Date(s);
}
import { ShoppingCart, X, ChevronDown, ChevronUp, Calendar } from "lucide-react";
import { useAdminAuth } from "@/context/AdminAuthContext";
import {
  Order,
  OrderItem,
  adminListOrders,
  adminUpdateOrderStatus,
  adminCancelOrder,
  adminGetOrdersByStatus,
  adminGetOrderItems,
  getOrdersByDateRange,
} from "@/lib/api";

const STATUS_OPTIONS = ["pending", "confirmed", "shipped", "delivered", "cancelled"];

const STATUS_STYLES: Record<string, string> = {
  pending: "bg-yellow-500/15 text-yellow-400",
  confirmed: "bg-blue-500/15 text-blue-400",
  shipped: "bg-purple-500/15 text-purple-400",
  delivered: "bg-green-500/15 text-green-400",
  cancelled: "bg-red-500/15 text-red-400",
};

function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={`px-2.5 py-1 rounded-full text-[11px] font-bold capitalize ${
        STATUS_STYLES[status] || "bg-white/10 text-gray-400"
      }`}
    >
      {status}
    </span>
  );
}

export default function AdminOrdersPage() {
  const { token } = useAdminAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [expanded, setExpanded] = useState<string | null>(null);
  const [expandedItems, setExpandedItems] = useState<Record<string, OrderItem[]>>({});
  const [itemsLoading, setItemsLoading] = useState<string | null>(null);
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const [dateActive, setDateActive] = useState(false);

  const loadAll = async () => {
    if (!token) return;
    setStatusFilter("");
    setDateActive(false);
    try {
      const os = await adminListOrders(token);
      setOrders(os);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to load orders");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAll();
  }, [token]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleStatusFilter = async (status: string) => {
    if (!token) return;
    setStatusFilter(status);
    setDateActive(false);
    setLoading(true);
    try {
      const os = await adminGetOrdersByStatus(token, status);
      setOrders(os);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to filter orders");
    } finally {
      setLoading(false);
    }
  };

  const handleDateFilter = async () => {
    if (!dateFrom || !dateTo) return;
    setStatusFilter("");
    setDateActive(true);
    setLoading(true);
    try {
      const os = await getOrdersByDateRange(dateFrom, dateTo);
      setOrders(os);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to filter by date");
    } finally {
      setLoading(false);
    }
  };

  const handleExpand = async (id: string) => {
    if (expanded === id) {
      setExpanded(null);
      return;
    }
    setExpanded(id);
    if (!expandedItems[id] && token) {
      setItemsLoading(id);
      try {
        const items = await adminGetOrderItems(token, id);
        setExpandedItems((prev) => ({ ...prev, [id]: items }));
      } catch {
        // fall back to order.items
      } finally {
        setItemsLoading(null);
      }
    }
  };

  const handleStatusChange = async (id: string, status: string) => {
    if (!token) return;
    try {
      await adminUpdateOrderStatus(token, id, status);
      loadAll();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to update status");
    }
  };

  const handleCancel = async (id: string) => {
    if (!token || !confirm("Cancel this order?")) return;
    try {
      await adminCancelOrder(token, id);
      loadAll();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to cancel");
    }
  };

  const totalRevenue = orders
    .filter((o) => o.status !== "cancelled")
    .reduce((sum, o) => sum + (o.total_amount || 0), 0);

  return (
    <div>
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <ShoppingCart size={14} className="text-[#C9A84C]" />
            <span className="text-[10px] font-black uppercase tracking-[0.3em] text-gray-600">
              Management
            </span>
          </div>
          <h1 className="text-2xl font-black text-white">Orders</h1>
          <p className="text-gray-600 text-sm mt-1">
            {orders.length} total · Revenue: ${totalRevenue.toLocaleString()}
          </p>
        </div>
      </div>

      {error && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-4 mb-4 text-red-400 text-sm flex items-center justify-between">
          <span>{error}</span>
          <button onClick={() => setError("")}>
            <X size={16} />
          </button>
        </div>
      )}

      {/* Status filter pills */}
      <div className="flex gap-2 mb-4 flex-wrap">
        <button
          onClick={loadAll}
          className={`px-3 py-1.5 rounded-xl text-[12px] font-bold transition-colors ${
            !statusFilter && !dateActive
              ? "bg-white text-black"
              : "border border-white/10 text-gray-500 hover:text-white hover:border-white/20"
          }`}
        >
          All
        </button>
        {STATUS_OPTIONS.map((s) => (
          <button
            key={s}
            onClick={() => handleStatusFilter(s)}
            className={`px-3 py-1.5 rounded-xl text-[12px] font-bold capitalize transition-colors ${
              statusFilter === s
                ? "bg-white text-black"
                : "border border-white/10 text-gray-500 hover:text-white hover:border-white/20"
            }`}
          >
            {s}
          </button>
        ))}
      </div>

      {/* Date range filter */}
      <div className={`flex flex-wrap gap-2 items-center mb-5 p-3 rounded-xl border transition-colors ${dateActive ? "border-[#C9A84C]/30 bg-[#C9A84C]/5" : "border-white/6 bg-white/3"}`}>
        <Calendar size={13} className="text-gray-600" />
        <span className="text-[11px] text-gray-500 font-semibold uppercase tracking-wider">Date Range</span>
        <input
          type="date"
          value={dateFrom}
          onChange={(e) => setDateFrom(e.target.value)}
          className="bg-white/5 border border-white/10 text-gray-300 rounded-lg px-3 py-1.5 text-xs focus:outline-none focus:border-white/20"
        />
        <span className="text-gray-600 text-xs">to</span>
        <input
          type="date"
          value={dateTo}
          onChange={(e) => setDateTo(e.target.value)}
          className="bg-white/5 border border-white/10 text-gray-300 rounded-lg px-3 py-1.5 text-xs focus:outline-none focus:border-white/20"
        />
        <button
          onClick={handleDateFilter}
          disabled={!dateFrom || !dateTo}
          className="bg-white text-black px-3 py-1.5 rounded-lg text-xs font-bold hover:bg-gray-100 disabled:opacity-40 transition-colors"
        >
          Apply
        </button>
        {dateActive && (
          <button
            onClick={() => { setDateFrom(""); setDateTo(""); loadAll(); }}
            className="text-gray-500 hover:text-white transition-colors text-xs"
          >
            Clear
          </button>
        )}
      </div>

      <div className="bg-[#161616] border border-white/6 rounded-2xl overflow-hidden">
        {loading ? (
          <div className="p-12 text-center text-gray-600 text-sm">Loading orders...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/6">
                  {["Order", "User", "Status", "Total", "Date", "Actions"].map((h) => (
                    <th
                      key={h}
                      className={`px-4 py-3 text-[10px] font-bold text-gray-600 uppercase tracking-wider ${
                        ["Total", "Actions"].includes(h) ? "text-right" : "text-left"
                      }`}
                    >
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-white/4">
                {orders.map((o) => (
                  <React.Fragment key={o.id}>
                    <tr
                      className="hover:bg-white/3 transition-colors cursor-pointer"
                      onClick={() => handleExpand(o.id)}
                    >
                      <td className="px-4 py-3.5">
                        <p className="font-mono text-[11px] text-gray-500">
                          #{o.id.slice(0, 8).toUpperCase()}
                        </p>
                      </td>
                      <td className="px-4 py-3.5">
                        <p className="font-mono text-[11px] text-gray-600">
                          {o.user_id.slice(0, 8)}...
                        </p>
                      </td>
                      <td className="px-4 py-3.5">
                        <StatusBadge status={o.status} />
                      </td>
                      <td className="px-4 py-3.5 text-right font-semibold text-white">
                        ${o.total_amount?.toLocaleString()}
                      </td>
                      <td className="px-4 py-3.5 text-gray-500 text-xs">
                        {parseGoDate(o.created_at).toLocaleDateString("en-US", {
                          month: "short",
                          day: "numeric",
                          year: "numeric",
                        })}
                      </td>
                      <td
                        className="px-4 py-3.5 text-right"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <div className="flex justify-end items-center gap-2">
                          <select
                            value={o.status}
                            onChange={(e) => handleStatusChange(o.id, e.target.value)}
                            disabled={o.status === "cancelled"}
                            className="bg-white/5 border border-white/10 text-gray-400 rounded-lg px-2 py-1 text-xs focus:outline-none disabled:opacity-40"
                          >
                            {STATUS_OPTIONS.map((s) => (
                              <option key={s} value={s} className="bg-[#1a1a1a]">
                                {s}
                              </option>
                            ))}
                          </select>
                          {o.status !== "cancelled" && o.status !== "delivered" && (
                            <button
                              onClick={() => handleCancel(o.id)}
                              className="text-xs text-red-400 hover:text-red-300 font-semibold transition-colors"
                            >
                              Cancel
                            </button>
                          )}
                          {expanded === o.id ? (
                            <ChevronUp size={14} className="text-gray-600" />
                          ) : (
                            <ChevronDown size={14} className="text-gray-600" />
                          )}
                        </div>
                      </td>
                    </tr>
                    {expanded === o.id && (
                      <tr className="bg-white/3">
                        <td colSpan={6} className="px-8 py-4">
                          <div className="text-xs text-gray-500 space-y-2">
                            <p>
                              <span className="text-gray-600 font-semibold">Shipping: </span>
                              {o.shipping_address}
                            </p>
                            <div className="space-y-1">
                              {itemsLoading === o.id ? (
                                <p className="text-gray-600">Loading items...</p>
                              ) : (expandedItems[o.id] || o.items || []).map((item, i) => (
                                <div key={i} className="flex justify-between text-gray-500">
                                  <span>
                                    {item.product_name} × {item.quantity}
                                    {item.size && ` · EU ${item.size}`}
                                  </span>
                                  <span className="text-white font-medium">
                                    ${(item.price * item.quantity).toLocaleString()}
                                  </span>
                                </div>
                              ))}
                            </div>
                          </div>
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                ))}
                {orders.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-4 py-12 text-center text-gray-600 text-sm">
                      No orders found
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
