"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import {
  LayoutDashboard,
  Package,
  ShoppingCart,
  Users,
  LogOut,
  Store,
  Menu,
  X,
} from "lucide-react";
import { AdminAuthProvider, useAdminAuth } from "@/context/AdminAuthContext";

const NAV_ITEMS = [
  { href: "/admin", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/products", label: "Products", icon: Package },
  { href: "/admin/orders", label: "Orders", icon: ShoppingCart },
  { href: "/admin/users", label: "Users", icon: Users },
];

function AdminLayoutInner({ children }: { children: React.ReactNode }) {
  const { isAdmin, loaded, adminLogout } = useAdminAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  useEffect(() => {
    if (loaded && !isAdmin && pathname !== "/admin/login") {
      router.push("/admin/login");
    }
  }, [isAdmin, loaded, pathname, router]);

  if (pathname === "/admin/login") return <>{children}</>;
  if (!loaded) return null;
  if (!isAdmin) return null;

  const Sidebar = ({ mobile = false }: { mobile?: boolean }) => (
    <aside
      className={`${
        mobile
          ? "fixed inset-y-0 left-0 z-50 w-64"
          : "fixed inset-y-0 left-0 z-50 w-60 hidden lg:flex"
      } bg-[#0D0D0D] border-r border-white/6 flex flex-col`}
    >
      {/* Logo */}
      <div className="h-16 flex items-center px-6 border-b border-white/6">
        <div className="flex items-baseline gap-1">
          <span className="text-lg font-black text-white">SNKR</span>
          <span className="text-[9px] font-black text-[#C9A84C] tracking-[0.35em] ml-0.5">
            VAULT
          </span>
        </div>
        <span className="ml-3 text-[10px] font-bold text-gray-600 uppercase tracking-wider">
          Admin
        </span>
        {mobile && (
          <button
            onClick={() => setSidebarOpen(false)}
            className="ml-auto text-gray-500 hover:text-white"
          >
            <X size={18} />
          </button>
        )}
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-4 space-y-0.5">
        {NAV_ITEMS.map((item) => {
          const active = pathname === item.href;
          const Icon = item.icon;
          return (
            <Link
              key={item.href}
              href={item.href}
              onClick={() => setSidebarOpen(false)}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all ${
                active
                  ? "bg-white/10 text-white"
                  : "text-gray-500 hover:text-gray-200 hover:bg-white/5"
              }`}
            >
              <Icon size={16} className={active ? "text-[#C9A84C]" : "text-gray-600"} />
              {item.label}
              {active && (
                <div className="ml-auto w-1.5 h-1.5 rounded-full bg-[#C9A84C]" />
              )}
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="px-3 py-4 border-t border-white/6 space-y-0.5">
        <Link
          href="/"
          className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-gray-500 hover:text-gray-200 hover:bg-white/5 transition-all"
        >
          <Store size={16} className="text-gray-600" />
          View Store
        </Link>
        <button
          onClick={() => {
            adminLogout();
            router.push("/admin/login");
          }}
          className="flex items-center gap-3 w-full px-3 py-2.5 rounded-xl text-sm font-medium text-gray-500 hover:text-red-400 hover:bg-white/5 transition-all"
        >
          <LogOut size={16} className="text-gray-600" />
          Sign Out
        </button>
      </div>
    </aside>
  );

  return (
    <div className="flex min-h-screen bg-[#0A0A0A]">
      {/* Desktop sidebar */}
      <Sidebar />

      {/* Mobile sidebar */}
      {sidebarOpen && (
        <>
          <div
            className="fixed inset-0 bg-black/60 z-40 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          />
          <Sidebar mobile />
        </>
      )}

      {/* Main content */}
      <div className="flex-1 lg:ml-60 flex flex-col min-h-screen">
        {/* Mobile topbar */}
        <div className="lg:hidden h-14 bg-[#0D0D0D] border-b border-white/6 flex items-center px-4 gap-3">
          <button
            onClick={() => setSidebarOpen(true)}
            className="text-gray-400 hover:text-white"
          >
            <Menu size={20} />
          </button>
          <div className="flex items-baseline gap-1">
            <span className="text-base font-black text-white">SNKR</span>
            <span className="text-[9px] font-black text-[#C9A84C] tracking-[0.3em] ml-0.5">VAULT</span>
          </div>
          <span className="text-[10px] text-gray-600 font-bold uppercase tracking-wider ml-1">
            Admin
          </span>
        </div>

        <main className="flex-1 p-6 lg:p-8">{children}</main>
      </div>
    </div>
  );
}

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <AdminAuthProvider>
      <AdminLayoutInner>{children}</AdminLayoutInner>
    </AdminAuthProvider>
  );
}
