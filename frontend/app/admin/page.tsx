"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import {
  DollarSign,
  ShoppingCart,
  Package,
  Users,
  Clock,
  UserCheck,
  ArrowRight,
  TrendingUp,
} from "lucide-react";
import { useAdminAuth } from "@/context/AdminAuthContext";
import { adminGetStats, DashStats } from "@/lib/api";

interface StatCardProps {
  label: string;
  value: string | number;
  sub?: string;
  icon: React.ReactNode;
  accent?: string;
}

function StatCard({ label, value, sub, icon, accent = "text-gray-400" }: StatCardProps) {
  return (
    <div className="bg-[#161616] border border-white/6 rounded-2xl p-6">
      <div className="flex items-start justify-between mb-4">
        <p className="text-xs font-bold text-gray-500 uppercase tracking-wider">{label}</p>
        <div className={`w-9 h-9 rounded-xl bg-white/5 flex items-center justify-center ${accent}`}>
          {icon}
        </div>
      </div>
      <p className="text-3xl font-black text-white mb-1">{value}</p>
      {sub && <p className="text-xs text-gray-600 font-medium">{sub}</p>}
    </div>
  );
}

function StatCardSkeleton() {
  return (
    <div className="bg-[#161616] border border-white/6 rounded-2xl p-6 animate-pulse">
      <div className="flex items-start justify-between mb-4">
        <div className="h-3 bg-white/8 rounded w-24" />
        <div className="w-9 h-9 bg-white/8 rounded-xl" />
      </div>
      <div className="h-9 bg-white/8 rounded w-28 mb-2" />
      <div className="h-3 bg-white/5 rounded w-20" />
    </div>
  );
}

export default function AdminDashboard() {
  const { token } = useAdminAuth();
  const [stats, setStats] = useState<DashStats | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!token) return;
    adminGetStats(token)
      .then(setStats)
      .catch((e: Error) => setError(e.message));
  }, [token]);

  return (
    <div>
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-2 mb-1">
          <TrendingUp size={14} className="text-[#C9A84C]" />
          <span className="text-[10px] font-black uppercase tracking-[0.3em] text-gray-600">
            Overview
          </span>
        </div>
        <h1 className="text-2xl font-black text-white">Dashboard</h1>
        <p className="text-gray-600 text-sm mt-1">Store performance at a glance</p>
      </div>

      {error && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-4 mb-6 text-red-400 text-sm">
          {error}
        </div>
      )}

      {/* Stats grid */}
      {!stats ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          {[...Array(6)].map((_, i) => (
            <StatCardSkeleton key={i} />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          <StatCard
            label="Total Revenue"
            value={`$${stats.total_revenue.toLocaleString()}`}
            icon={<DollarSign size={16} />}
            accent="text-[#C9A84C]"
          />
          <StatCard
            label="Total Orders"
            value={stats.total_orders}
            sub={`${stats.pending_orders} pending`}
            icon={<ShoppingCart size={16} />}
            accent="text-blue-400"
          />
          <StatCard
            label="Products"
            value={stats.total_products}
            icon={<Package size={16} />}
            accent="text-purple-400"
          />
          <StatCard
            label="Registered Users"
            value={stats.total_users}
            sub={`${stats.active_users} active`}
            icon={<Users size={16} />}
            accent="text-green-400"
          />
          <StatCard
            label="Pending Orders"
            value={stats.pending_orders}
            sub="Awaiting processing"
            icon={<Clock size={16} />}
            accent="text-orange-400"
          />
          <StatCard
            label="Active Users"
            value={stats.active_users}
            sub={`${stats.total_users - stats.active_users} inactive`}
            icon={<UserCheck size={16} />}
            accent="text-cyan-400"
          />
        </div>
      )}

      {/* Quick actions */}
      <div className="bg-[#161616] border border-white/6 rounded-2xl p-6">
        <h2 className="text-sm font-black text-white mb-5 uppercase tracking-wider">
          Quick Actions
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
          {[
            { href: "/admin/products", label: "Manage Products", icon: <Package size={16} /> },
            { href: "/admin/orders", label: "View Orders", icon: <ShoppingCart size={16} /> },
            { href: "/admin/users", label: "Manage Users", icon: <Users size={16} /> },
            { href: "/", label: "View Storefront", icon: <ArrowRight size={16} /> },
          ].map((a) => (
            <Link
              key={a.href}
              href={a.href}
              className="flex items-center gap-3 px-4 py-3.5 bg-white/5 border border-white/8 rounded-xl text-sm font-medium text-gray-400 hover:text-white hover:bg-white/8 hover:border-white/15 transition-all group"
            >
              <span className="text-gray-600 group-hover:text-[#C9A84C] transition-colors">
                {a.icon}
              </span>
              {a.label}
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
