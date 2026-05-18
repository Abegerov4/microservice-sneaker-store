"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { ArrowRight, Shield, RotateCcw, Truck, Flame, Sparkles, TrendingUp } from "lucide-react";
import { Product, listProducts, getBrands, searchProducts, getProductsByBrand } from "@/lib/api";
import ProductCard from "@/components/ProductCard";
import HeroCarousel from "@/components/HeroCarousel";

const SIZES = ["36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46"];

const MARQUEE_TEXT =
  "NEW RELEASES · FREE SHIPPING ON $100+ · 100% AUTHENTICATED · 30-DAY RETURNS · PREMIUM SNEAKER MARKETPLACE · EXCLUSIVE DROPS · ";

export default function Home() {
  const [products, setProducts] = useState<Product[]>([]);
  const [brands, setBrands] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [selectedBrand, setSelectedBrand] = useState("");
  const [selectedSize, setSelectedSize] = useState("");
  const [apiProducts, setApiProducts] = useState<Product[] | null>(null);
  const [apiLoading, setApiLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [subscribed, setSubscribed] = useState(false);

  useEffect(() => {
    Promise.all([listProducts(), getBrands()])
      .then(([prods, b]) => {
        setProducts(prods);
        setBrands(b);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (!selectedBrand && !selectedSize) {
      setApiProducts(null);
      return;
    }
    setApiLoading(true);
    if (selectedSize) {
      searchProducts({ brand: selectedBrand, size: selectedSize })
        .then(setApiProducts)
        .catch(() => setApiProducts(null))
        .finally(() => setApiLoading(false));
    } else if (selectedBrand) {
      getProductsByBrand(selectedBrand)
        .then(setApiProducts)
        .catch(() => setApiProducts(null))
        .finally(() => setApiLoading(false));
    }
  }, [selectedBrand, selectedSize]);

  const source = apiProducts ?? products;
  const filtered = source.filter((p) => {
    const matchesSearch =
      !search ||
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      p.brand.toLowerCase().includes(search.toLowerCase());
    return matchesSearch;
  });

  const clearFilters = () => {
    setSelectedBrand("");
    setSelectedSize("");
    setSearch("");
  };
  const hasFilters = !!(selectedBrand || selectedSize || search);
  const newArrivals = products.slice(0, 4);

  return (
    <div className="min-h-screen bg-[#F5F5F0]">
      {/* ── Hero ── */}
      <HeroCarousel />

      {/* ── Marquee strip ── */}
      <div className="bg-black text-white py-3 overflow-hidden">
        <div
          className="flex whitespace-nowrap"
          style={{ animation: "marquee 28s linear infinite" }}
        >
          <span className="text-[11px] font-bold tracking-[0.22em] uppercase text-gray-400 flex-shrink-0">
            {MARQUEE_TEXT.repeat(4)}
          </span>
        </div>
      </div>

      {/* ── Brand filter bar ── */}
      {brands.length > 0 && (
        <div className="bg-white border-b border-gray-100 sticky top-16 z-40">
          <div className="max-w-[1400px] mx-auto px-4 sm:px-6 py-4">
            <div className="flex gap-2 flex-wrap items-center">
              <button
                onClick={() => setSelectedBrand("")}
                className={`px-4 py-2 rounded-full text-[11px] font-bold tracking-wider uppercase transition-all ${
                  !selectedBrand
                    ? "bg-black text-white shadow-sm"
                    : "bg-gray-100 text-gray-500 hover:bg-gray-200"
                }`}
              >
                All
              </button>
              {brands.map((b) => (
                <button
                  key={b}
                  onClick={() => setSelectedBrand(b === selectedBrand ? "" : b)}
                  className={`px-4 py-2 rounded-full text-[11px] font-bold tracking-wider uppercase transition-all ${
                    selectedBrand === b
                      ? "bg-black text-white shadow-sm"
                      : "bg-gray-100 text-gray-500 hover:bg-gray-200"
                  }`}
                >
                  {b}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* ── New Arrivals section ── */}
      {!hasFilters && newArrivals.length > 0 && (
        <section className="max-w-[1400px] mx-auto px-4 sm:px-6 pt-12 pb-8">
          <div className="flex items-end justify-between mb-6">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <Sparkles size={13} className="text-[#C9A84C]" />
                <span className="text-[10px] font-black uppercase tracking-[0.3em] text-[#C9A84C]">
                  Fresh Drop
                </span>
              </div>
              <h2 className="text-2xl sm:text-3xl font-black text-gray-900">New Arrivals</h2>
            </div>
            <button
              onClick={() =>
                document.getElementById("catalog")?.scrollIntoView({ behavior: "smooth" })
              }
              className="flex items-center gap-1.5 text-xs font-bold text-gray-500 hover:text-black transition-colors uppercase tracking-wider"
            >
              View All <ArrowRight size={13} />
            </button>
          </div>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {newArrivals.map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
        </section>
      )}

      {/* ── Exclusive banner ── */}
      {!hasFilters && (
        <section className="max-w-[1400px] mx-auto px-4 sm:px-6 py-6">
          <div className="bg-black rounded-3xl overflow-hidden relative">
            <div
              className="absolute inset-0 opacity-40"
              style={{
                backgroundImage:
                  "radial-gradient(circle at 75% 50%, rgba(201,168,76,0.5) 0%, transparent 55%)",
              }}
            />
            <div className="relative z-10 px-8 sm:px-14 py-10 sm:py-14 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-6">
              <div>
                <div className="flex items-center gap-2 mb-3">
                  <Flame size={14} className="text-[#C9A84C]" />
                  <span className="text-[10px] font-black uppercase tracking-[0.35em] text-[#C9A84C]">
                    Exclusive
                  </span>
                </div>
                <h2 className="text-3xl sm:text-4xl font-black text-white leading-tight mb-3">
                  Every Pair<br />Authenticated.
                </h2>
                <p className="text-gray-400 text-sm leading-relaxed max-w-sm">
                  We source the most coveted releases. Each pair verified before delivery.
                </p>
              </div>
              <button
                onClick={() =>
                  document.getElementById("catalog")?.scrollIntoView({ behavior: "smooth" })
                }
                className="bg-white text-black px-7 py-3.5 rounded-xl font-bold text-sm hover:bg-gray-100 transition-colors whitespace-nowrap flex-shrink-0"
              >
                Shop Now
              </button>
            </div>
          </div>
        </section>
      )}

      {/* ── Main Catalog ── */}
      <section id="catalog" className="max-w-[1400px] mx-auto px-4 sm:px-6 py-10">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-8">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <TrendingUp size={13} className="text-gray-400" />
              <span className="text-[10px] font-black uppercase tracking-[0.3em] text-gray-400">
                Catalog
              </span>
            </div>
            <h2 className="text-2xl sm:text-3xl font-black text-gray-900">
              {selectedBrand || "All Sneakers"}
            </h2>
          </div>

          <div className="flex flex-wrap gap-2.5 items-center">
            <div className="relative">
              <svg
                className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
                width="13"
                height="13"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2.5"
              >
                <circle cx="11" cy="11" r="8" />
                <path d="m21 21-4.35-4.35" />
              </svg>
              <input
                type="text"
                placeholder="Search..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="bg-white border border-gray-200 rounded-xl pl-9 pr-4 py-2.5 text-[13px] focus:outline-none focus:ring-2 focus:ring-black w-44 sm:w-48 transition-all"
              />
            </div>

            <select
              value={selectedSize}
              onChange={(e) => setSelectedSize(e.target.value)}
              className="bg-white border border-gray-200 rounded-xl px-3 py-2.5 text-[13px] focus:outline-none focus:ring-2 focus:ring-black"
            >
              <option value="">All Sizes</option>
              {SIZES.map((s) => (
                <option key={s} value={s}>
                  EU {s}
                </option>
              ))}
            </select>

            {hasFilters && (
              <button
                onClick={clearFilters}
                className="flex items-center gap-1.5 text-[12px] font-bold text-gray-400 hover:text-black transition-colors px-3 py-2.5 bg-white border border-gray-200 rounded-xl"
              >
                <svg
                  width="10"
                  height="10"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="3"
                >
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
                Clear
              </button>
            )}

            <span className="text-[12px] text-gray-400 font-medium hidden sm:block">
              {loading || apiLoading ? "Loading..." : `${filtered.length} results`}
            </span>
          </div>
        </div>

        {/* Grid */}
        {loading || apiLoading ? (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {[...Array(8)].map((_, i) => (
              <div key={i} className="bg-white rounded-2xl overflow-hidden">
                <div className="aspect-square bg-gray-100 animate-pulse" />
                <div className="p-4 space-y-2">
                  <div className="bg-gray-100 animate-pulse h-2.5 w-1/4 rounded-full" />
                  <div className="bg-gray-100 animate-pulse h-4 w-3/4 rounded-full" />
                  <div className="bg-gray-100 animate-pulse h-5 w-1/3 rounded-full mt-3" />
                </div>
              </div>
            ))}
          </div>
        ) : filtered.length === 0 ? (
          <div className="text-center py-28 bg-white rounded-3xl">
            <div className="text-5xl mb-5">🔍</div>
            <h3 className="text-xl font-black text-gray-900 mb-2">No results</h3>
            <p className="text-gray-400 mb-6 text-sm">Try adjusting your search or filters</p>
            <button
              onClick={clearFilters}
              className="bg-black text-white px-7 py-3 rounded-xl font-bold text-sm hover:bg-gray-800 transition-colors"
            >
              Clear Filters
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {filtered.map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
        )}
      </section>

      {/* ── Trust badges ── */}
      <section className="max-w-[1400px] mx-auto px-4 sm:px-6 pb-10">
        <div className="bg-white rounded-3xl grid grid-cols-1 sm:grid-cols-3 divide-y sm:divide-y-0 sm:divide-x divide-gray-100">
          {[
            {
              icon: <Shield size={22} className="text-gray-700" />,
              title: "100% Authentic",
              desc: "Every sneaker verified before shipment",
            },
            {
              icon: <Truck size={22} className="text-gray-700" />,
              title: "Free Shipping",
              desc: "On all orders over $100",
            },
            {
              icon: <RotateCcw size={22} className="text-gray-700" />,
              title: "30-Day Returns",
              desc: "Hassle-free return policy",
            },
          ].map((item) => (
            <div key={item.title} className="flex items-center gap-4 px-8 py-6">
              <div className="w-10 h-10 bg-gray-50 rounded-xl flex items-center justify-center flex-shrink-0">
                {item.icon}
              </div>
              <div>
                <p className="text-sm font-black text-gray-900">{item.title}</p>
                <p className="text-xs text-gray-400 mt-0.5">{item.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* ── Newsletter ── */}
      <section className="max-w-[1400px] mx-auto px-4 sm:px-6 pb-8">
        <div
          className="relative overflow-hidden rounded-3xl bg-[#0A0A0A] px-8 sm:px-14 py-12 sm:py-16 text-center"
          style={{
            backgroundImage:
              "radial-gradient(circle at 50% 0%, rgba(201,168,76,0.18) 0%, transparent 60%)",
          }}
        >
          <span className="text-[10px] font-black uppercase tracking-[0.4em] text-[#C9A84C] block mb-3">
            Stay ahead
          </span>
          <h2 className="text-3xl sm:text-4xl font-black text-white mb-3">
            Drop Alerts & Exclusive Access
          </h2>
          <p className="text-gray-500 mb-8 max-w-md mx-auto text-sm leading-relaxed">
            Be first to know about new releases, restocks, and limited drops.
          </p>
          {subscribed ? (
            <div className="inline-flex items-center gap-2 bg-white/10 text-white px-6 py-3 rounded-xl text-sm font-semibold">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="#C9A84C"
                strokeWidth="2.5"
              >
                <polyline points="20 6 9 17 4 12" />
              </svg>
              You&apos;re on the list!
            </div>
          ) : (
            <form
              onSubmit={(e) => {
                e.preventDefault();
                if (email) setSubscribed(true);
              }}
              className="flex max-w-md mx-auto gap-2"
            >
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter your email"
                required
                className="flex-1 bg-white/8 border border-white/10 text-white placeholder-gray-600 rounded-xl px-5 py-3.5 text-sm focus:outline-none focus:border-white/25 transition-all"
              />
              <button
                type="submit"
                className="bg-[#C9A84C] text-black px-6 py-3.5 rounded-xl font-bold text-sm hover:bg-[#d4b55c] transition-colors flex-shrink-0"
              >
                Subscribe
              </button>
            </form>
          )}
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="bg-black text-white mt-6">
        <div className="max-w-[1400px] mx-auto px-4 sm:px-6 pt-14 pb-8">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12">
            <div className="col-span-2 md:col-span-1">
              <div className="flex items-baseline gap-1 mb-4">
                <span className="text-2xl font-black text-white">SNKR</span>
                <span className="text-[11px] font-black text-[#C9A84C] tracking-[0.35em] ml-0.5">
                  VAULT
                </span>
              </div>
              <p className="text-gray-500 text-sm leading-relaxed max-w-[200px]">
                Premium sneakers, authenticated and delivered worldwide.
              </p>
            </div>

            {[
              {
                title: "Shop",
                links: [
                  { label: "New Releases", href: "/" },
                  { label: "All Sneakers", href: "/" },
                  { label: "Men", href: "/" },
                  { label: "Women", href: "/" },
                ],
              },
              {
                title: "Account",
                links: [
                  { label: "Sign In", href: "/auth/login" },
                  { label: "Register", href: "/auth/register" },
                  { label: "My Orders", href: "/orders" },
                  { label: "Profile", href: "/profile" },
                ],
              },
              {
                title: "Info",
                links: [
                  { label: "Authenticity", href: "/" },
                  { label: "Shipping Policy", href: "/" },
                  { label: "Returns", href: "/" },
                  { label: "Contact", href: "/" },
                ],
              },
            ].map((col) => (
              <div key={col.title}>
                <h4 className="text-[10px] font-black uppercase tracking-[0.3em] text-gray-600 mb-4">
                  {col.title}
                </h4>
                <ul className="space-y-2.5">
                  {col.links.map((link) => (
                    <li key={link.label}>
                      <Link
                        href={link.href}
                        className="text-[13px] text-gray-400 hover:text-white transition-colors"
                      >
                        {link.label}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>

          <div className="border-t border-white/8 pt-6 flex flex-col sm:flex-row items-center justify-between gap-3">
            <p className="text-xs text-gray-600">© 2026 SNKR VAULT. All rights reserved.</p>
            <p className="text-xs text-gray-600">Premium Sneaker Marketplace</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
