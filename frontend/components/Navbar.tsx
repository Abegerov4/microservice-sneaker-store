"use client";

import { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { useCart } from "@/context/CartContext";
import { useAuth } from "@/context/AuthContext";
import {
  ShoppingBag,
  Search,
  X,
  Menu,
  User,
  Package,
  Settings,
  LogOut,
  ChevronDown,
} from "lucide-react";

const NAV_LINKS = [
  { label: "New Releases", href: "/" },
  { label: "Men", href: "/" },
  { label: "Women", href: "/" },
  { label: "Brands", href: "/" },
];

export default function Navbar() {
  const { count } = useCart();
  const { user, isLoggedIn, logout } = useAuth();
  const [scrolled, setScrolled] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [profileOpen, setProfileOpen] = useState(false);
  const searchRef = useRef<HTMLInputElement>(null);
  const profileRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  useEffect(() => {
    if (searchOpen) searchRef.current?.focus();
  }, [searchOpen]);

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (profileRef.current && !profileRef.current.contains(e.target as Node)) {
        setProfileOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  useEffect(() => {
    if (menuOpen) document.body.style.overflow = "hidden";
    else document.body.style.overflow = "";
    return () => { document.body.style.overflow = ""; };
  }, [menuOpen]);

  const initials = (user?.full_name || user?.email || "U")[0].toUpperCase();

  return (
    <>
      <nav
        className={`sticky top-0 z-50 transition-all duration-300 ${
          scrolled
            ? "bg-black/98 backdrop-blur-xl shadow-[0_1px_0_rgba(255,255,255,0.06)]"
            : "bg-black"
        }`}
      >
        <div className="max-w-[1400px] mx-auto px-4 sm:px-6 h-16 flex items-center justify-between gap-4">
          {/* Logo */}
          <Link
            href="/"
            className="flex-shrink-0 flex items-center gap-1 group"
            onClick={() => setMenuOpen(false)}
          >
            <span className="text-[22px] font-black tracking-tight text-white leading-none">
              SNKR
            </span>
            <span className="text-[10px] font-black tracking-[0.35em] text-[#C9A84C] ml-0.5 leading-none">
              VAULT
            </span>
          </Link>

          {/* Center nav — desktop */}
          <div className="hidden lg:flex items-center gap-8 absolute left-1/2 -translate-x-1/2">
            {NAV_LINKS.map((link) => (
              <Link
                key={link.label}
                href={link.href}
                className="text-[13px] font-medium text-gray-400 hover:text-white transition-colors tracking-wide"
              >
                {link.label}
              </Link>
            ))}
          </div>

          {/* Right controls */}
          <div className="flex items-center gap-1">
            {/* Search toggle */}
            <button
              onClick={() => setSearchOpen((v) => !v)}
              className="w-10 h-10 flex items-center justify-center rounded-xl text-gray-400 hover:text-white hover:bg-white/8 transition-all"
              aria-label="Search"
            >
              {searchOpen ? <X size={18} /> : <Search size={18} />}
            </button>

            {/* Auth */}
            {isLoggedIn ? (
              <div className="relative" ref={profileRef}>
                <button
                  onClick={() => setProfileOpen((v) => !v)}
                  className="flex items-center gap-1.5 h-10 px-2 rounded-xl hover:bg-white/8 transition-all"
                >
                  <div className="w-7 h-7 rounded-full bg-[#C9A84C] text-black text-[11px] font-black flex items-center justify-center">
                    {initials}
                  </div>
                  <ChevronDown
                    size={14}
                    className={`text-gray-500 transition-transform ${profileOpen ? "rotate-180" : ""}`}
                  />
                </button>

                {profileOpen && (
                  <div className="absolute right-0 top-full mt-2 w-52 bg-[#1a1a1a] border border-white/10 rounded-2xl shadow-2xl overflow-hidden z-50">
                    <div className="px-4 py-3 border-b border-white/8">
                      <p className="text-sm font-semibold text-white truncate">
                        {user?.full_name || "User"}
                      </p>
                      <p className="text-xs text-gray-500 truncate mt-0.5">{user?.email}</p>
                    </div>
                    <div className="py-1">
                      <Link
                        href="/profile"
                        onClick={() => setProfileOpen(false)}
                        className="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-300 hover:text-white hover:bg-white/5 transition-colors"
                      >
                        <User size={14} /> Profile
                      </Link>
                      <Link
                        href="/orders"
                        onClick={() => setProfileOpen(false)}
                        className="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-300 hover:text-white hover:bg-white/5 transition-colors"
                      >
                        <Package size={14} /> My Orders
                      </Link>
                      {user?.role === "ADMIN" && (
                        <Link
                          href="/admin"
                          onClick={() => setProfileOpen(false)}
                          className="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-300 hover:text-white hover:bg-white/5 transition-colors"
                        >
                          <Settings size={14} /> Admin Panel
                        </Link>
                      )}
                      <div className="border-t border-white/8 my-1" />
                      <button
                        onClick={() => { logout(); setProfileOpen(false); }}
                        className="flex items-center gap-3 w-full px-4 py-2.5 text-sm text-red-400 hover:text-red-300 hover:bg-white/5 transition-colors"
                      >
                        <LogOut size={14} /> Sign Out
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="hidden sm:flex items-center gap-2">
                <Link
                  href="/auth/login"
                  className="text-[13px] font-medium text-gray-400 hover:text-white transition-colors px-3 h-10 flex items-center"
                >
                  Sign In
                </Link>
                <Link
                  href="/auth/register"
                  className="text-[13px] font-semibold bg-white text-black px-4 h-9 rounded-xl flex items-center hover:bg-gray-100 transition-colors"
                >
                  Register
                </Link>
              </div>
            )}

            {/* Cart */}
            <Link
              href="/cart"
              className="relative w-10 h-10 flex items-center justify-center rounded-xl text-gray-400 hover:text-white hover:bg-white/8 transition-all"
              aria-label="Cart"
            >
              <ShoppingBag size={19} />
              {count > 0 && (
                <span className="absolute -top-0.5 -right-0.5 bg-[#C9A84C] text-black text-[10px] font-black min-w-[18px] h-[18px] rounded-full flex items-center justify-center px-1 leading-none">
                  {count > 9 ? "9+" : count}
                </span>
              )}
            </Link>

            {/* Mobile hamburger */}
            <button
              onClick={() => setMenuOpen((v) => !v)}
              className="lg:hidden w-10 h-10 flex items-center justify-center rounded-xl text-gray-400 hover:text-white hover:bg-white/8 transition-all"
              aria-label="Menu"
            >
              {menuOpen ? <X size={20} /> : <Menu size={20} />}
            </button>
          </div>
        </div>

        {/* Search bar dropdown */}
        <div
          className={`overflow-hidden transition-all duration-300 ${
            searchOpen ? "max-h-24 border-t border-white/8" : "max-h-0"
          }`}
        >
          <div className="max-w-[1400px] mx-auto px-4 sm:px-6 py-4">
            <div className="relative max-w-xl">
              <Search
                size={15}
                className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-500"
              />
              <input
                ref={searchRef}
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search sneakers, brands, colorways..."
                className="w-full bg-white/8 border border-white/10 text-white placeholder-gray-500 rounded-xl pl-10 pr-4 py-2.5 text-sm focus:outline-none focus:border-white/20 focus:bg-white/10 transition-all"
              />
            </div>
          </div>
        </div>
      </nav>

      {/* Mobile menu overlay */}
      {menuOpen && (
        <div className="fixed inset-0 z-40 bg-black flex flex-col pt-16">
          <div className="flex-1 overflow-y-auto px-6 py-6">
            <nav className="space-y-1 mb-8">
              {NAV_LINKS.map((link) => (
                <Link
                  key={link.label}
                  href={link.href}
                  onClick={() => setMenuOpen(false)}
                  className="flex items-center justify-between py-4 border-b border-white/8 text-lg font-semibold text-white hover:text-[#C9A84C] transition-colors"
                >
                  {link.label}
                  <span className="text-gray-600">→</span>
                </Link>
              ))}
            </nav>

            {isLoggedIn ? (
              <div className="space-y-1">
                <Link href="/profile" onClick={() => setMenuOpen(false)} className="flex items-center gap-3 py-3 text-gray-400 hover:text-white transition-colors">
                  <User size={16} /> Profile
                </Link>
                <Link href="/orders" onClick={() => setMenuOpen(false)} className="flex items-center gap-3 py-3 text-gray-400 hover:text-white transition-colors">
                  <Package size={16} /> My Orders
                </Link>
                <button onClick={() => { logout(); setMenuOpen(false); }} className="flex items-center gap-3 py-3 text-red-400 w-full">
                  <LogOut size={16} /> Sign Out
                </button>
              </div>
            ) : (
              <div className="flex flex-col gap-3 mt-4">
                <Link
                  href="/auth/login"
                  onClick={() => setMenuOpen(false)}
                  className="w-full text-center py-3.5 border border-white/20 rounded-xl font-semibold text-white hover:bg-white/8 transition-colors"
                >
                  Sign In
                </Link>
                <Link
                  href="/auth/register"
                  onClick={() => setMenuOpen(false)}
                  className="w-full text-center py-3.5 bg-white rounded-xl font-semibold text-black hover:bg-gray-100 transition-colors"
                >
                  Register
                </Link>
              </div>
            )}
          </div>

          <div className="px-6 py-4 border-t border-white/8">
            <p className="text-xs text-gray-600 text-center tracking-widest uppercase">
              SNKR VAULT — Premium Sneaker Marketplace
            </p>
          </div>
        </div>
      )}
    </>
  );
}
