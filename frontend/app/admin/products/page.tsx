"use client";

import { useEffect, useState } from "react";
import { Plus, Search, X, Package } from "lucide-react";
import { useAdminAuth } from "@/context/AdminAuthContext";
import {
  Product,
  adminListProducts,
  adminCreateProduct,
  adminUpdateProduct,
  adminDeleteProduct,
  adminUpdateStock,
  adminBulkDeleteProducts,
  adminGetLowStockProducts,
} from "@/lib/api";

const EMPTY_FORM = {
  name: "",
  brand: "",
  description: "",
  price: 0,
  sizes: "",
  stock: 0,
  image_url: "",
};
type ProductForm = typeof EMPTY_FORM;

function DarkModal({
  title,
  onClose,
  children,
}: {
  title: string;
  onClose: () => void;
  children: React.ReactNode;
}) {
  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-[#1a1a1a] border border-white/10 rounded-2xl shadow-2xl w-full max-w-lg">
        <div className="flex items-center justify-between px-6 py-4 border-b border-white/8">
          <h2 className="font-bold text-white text-sm">{title}</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-white transition-colors"
          >
            <X size={18} />
          </button>
        </div>
        <div className="p-6">{children}</div>
      </div>
    </div>
  );
}

function FormFields({
  form,
  setForm,
}: {
  form: ProductForm;
  setForm: (f: ProductForm) => void;
}) {
  const fields = [
    { key: "name", label: "Name", type: "text" },
    { key: "brand", label: "Brand", type: "text" },
    { key: "description", label: "Description", type: "text" },
    { key: "price", label: "Price ($)", type: "number" },
    { key: "stock", label: "Stock", type: "number" },
    { key: "image_url", label: "Image URL", type: "text" },
  ] as const;

  return (
    <div className="space-y-3">
      {fields.map(({ key, label, type }) => (
        <div key={key}>
          <label className="block text-[11px] font-bold text-gray-500 uppercase tracking-wider mb-1.5">
            {label}
          </label>
          <input
            type={type}
            value={((form as Record<string, unknown>)[key] ?? "") as string}
            onChange={(e) =>
              setForm({
                ...form,
                [key]: type === "number" ? Number(e.target.value) : e.target.value,
              })
            }
            className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-white/25 transition-colors"
          />
        </div>
      ))}
      <div>
        <label className="block text-[11px] font-bold text-gray-500 uppercase tracking-wider mb-1.5">
          Sizes (comma-separated)
        </label>
        <input
          type="text"
          value={form.sizes}
          onChange={(e) => setForm({ ...form, sizes: e.target.value })}
          placeholder="40,41,42,43,44"
          className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-white/25 transition-colors"
        />
      </div>
    </div>
  );
}

export default function AdminProductsPage() {
  const { token } = useAdminAuth();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [brandFilter, setBrandFilter] = useState("");
  const [showLowStock, setShowLowStock] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);
  const [editProduct, setEditProduct] = useState<Product | null>(null);
  const [stockProduct, setStockProduct] = useState<Product | null>(null);
  const [stockDelta, setStockDelta] = useState<string>("0");
  const [form, setForm] = useState<ProductForm>(EMPTY_FORM);
  const [saving, setSaving] = useState(false);
  const [selected, setSelected] = useState<Set<string>>(new Set());

  const load = async () => {
    if (!token) return;
    setShowLowStock(false);
    try {
      const ps = await adminListProducts(token);
      setProducts(ps);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to load");
    } finally {
      setLoading(false);
    }
    return;
  };

  const loadLowStock = async () => {
    if (!token) return;
    setShowLowStock(true);
    setLoading(true);
    try {
      const ps = await adminGetLowStockProducts(token);
      setProducts(ps);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to load low stock");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, [token]); // eslint-disable-line react-hooks/exhaustive-deps

  const brands = Array.from(new Set(products.map((p) => p.brand))).sort();

  const filtered = products.filter((p) => {
    const matchSearch =
      !search ||
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      p.brand.toLowerCase().includes(search.toLowerCase());
    const matchBrand = !brandFilter || p.brand === brandFilter;
    return matchSearch && matchBrand;
  });

  const openCreate = () => {
    setForm(EMPTY_FORM);
    setCreateOpen(true);
  };
  const openEdit = (p: Product) => {
    setForm({
      name: p.name,
      brand: p.brand,
      description: p.description,
      price: p.price,
      sizes: p.sizes.join(","),
      stock: p.stock,
      image_url: p.image_url,
    });
    setEditProduct(p);
  };

  const handleCreate = async () => {
    if (!token) return;
    setSaving(true);
    try {
      await adminCreateProduct(token, {
        ...form,
        sizes: form.sizes
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean),
        price: Number(form.price),
        stock: Number(form.stock),
      });
      setCreateOpen(false);
      load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to create");
    } finally {
      setSaving(false);
    }
  };

  const handleUpdate = async () => {
    if (!token || !editProduct) return;
    setSaving(true);
    try {
      await adminUpdateProduct(token, editProduct.id, {
        ...form,
        sizes: form.sizes
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean),
        price: Number(form.price),
        stock: Number(form.stock),
      });
      setEditProduct(null);
      load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to update");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!token || !confirm("Delete this product?")) return;
    try {
      await adminDeleteProduct(token, id);
      load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to delete");
    }
  };

  const handleBulkDelete = async () => {
    if (!token || selected.size === 0 || !confirm(`Delete ${selected.size} products?`)) return;
    try {
      await adminBulkDeleteProducts(token, Array.from(selected));
    } catch {
      for (const id of selected) {
        await adminDeleteProduct(token, id).catch(() => {});
      }
    }
    setSelected(new Set());
    load();
  };

  const handleStockUpdate = async () => {
    if (!token || !stockProduct) return;
    const delta = parseInt(stockDelta, 10);
    if (isNaN(delta) || delta === 0) return;
    try {
      await adminUpdateStock(token, stockProduct.id, delta);
      setStockProduct(null);
      await load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to update stock");
    }
  };

  const toggleSelect = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  return (
    <div>
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Package size={14} className="text-[#C9A84C]" />
            <span className="text-[10px] font-black uppercase tracking-[0.3em] text-gray-600">
              Inventory
            </span>
          </div>
          <h1 className="text-2xl font-black text-white">Products</h1>
          <p className="text-gray-600 text-sm mt-1">{products.length} {showLowStock ? "low-stock" : "total"} items</p>
        </div>
        <div className="flex gap-3">
          <button
            onClick={showLowStock ? load : loadLowStock}
            className={`px-4 py-2 rounded-xl text-sm font-bold transition-colors border ${
              showLowStock
                ? "bg-amber-500/20 text-amber-400 border-amber-500/20 hover:bg-amber-500/30"
                : "border-white/10 text-gray-500 hover:text-amber-400 hover:border-amber-500/20"
            }`}
          >
            {showLowStock ? "⚠ Low Stock" : "Low Stock"}
          </button>
          {selected.size > 0 && (
            <button
              onClick={handleBulkDelete}
              className="bg-red-500/20 text-red-400 border border-red-500/20 px-4 py-2 rounded-xl text-sm font-bold hover:bg-red-500/30 transition-colors"
            >
              Delete {selected.size}
            </button>
          )}
          <button
            onClick={openCreate}
            className="flex items-center gap-2 bg-white text-black px-4 py-2 rounded-xl text-sm font-bold hover:bg-gray-100 transition-colors"
          >
            <Plus size={15} /> Add Product
          </button>
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

      {/* Filters */}
      <div className="flex gap-3 mb-5">
        <div className="relative flex-1 max-w-xs">
          <Search size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-600" />
          <input
            type="text"
            placeholder="Search products..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl pl-9 pr-4 py-2.5 text-sm focus:outline-none focus:border-white/20 transition-colors"
          />
        </div>
        <select
          value={brandFilter}
          onChange={(e) => setBrandFilter(e.target.value)}
          className="bg-white/5 border border-white/10 text-gray-400 rounded-xl px-4 py-2.5 text-sm focus:outline-none"
        >
          <option value="">All Brands</option>
          {brands.map((b) => (
            <option key={b} value={b}>
              {b}
            </option>
          ))}
        </select>
        {(search || brandFilter) && (
          <button
            onClick={() => {
              setSearch("");
              setBrandFilter("");
            }}
            className="text-sm text-gray-500 hover:text-white transition-colors px-3"
          >
            Clear
          </button>
        )}
      </div>

      {/* Table */}
      <div className="bg-[#161616] border border-white/6 rounded-2xl overflow-hidden">
        {loading ? (
          <div className="p-12 text-center text-gray-600 text-sm">Loading products...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/6">
                  <th className="px-4 py-3 text-left w-10">
                    <input
                      type="checkbox"
                      className="accent-[#C9A84C]"
                      checked={selected.size === filtered.length && filtered.length > 0}
                      onChange={() => {
                        if (selected.size === filtered.length)
                          setSelected(new Set());
                        else setSelected(new Set(filtered.map((p) => p.id)));
                      }}
                    />
                  </th>
                  <th className="px-4 py-3 text-left text-[10px] font-bold text-gray-600 uppercase tracking-wider">
                    Product
                  </th>
                  <th className="px-4 py-3 text-left text-[10px] font-bold text-gray-600 uppercase tracking-wider">
                    Brand
                  </th>
                  <th className="px-4 py-3 text-right text-[10px] font-bold text-gray-600 uppercase tracking-wider">
                    Price
                  </th>
                  <th className="px-4 py-3 text-right text-[10px] font-bold text-gray-600 uppercase tracking-wider">
                    Stock
                  </th>
                  <th className="px-4 py-3 text-right text-[10px] font-bold text-gray-600 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/4">
                {filtered.map((p) => (
                  <tr
                    key={p.id}
                    className={`transition-colors ${
                      selected.has(p.id)
                        ? "bg-white/5"
                        : "hover:bg-white/3"
                    }`}
                  >
                    <td className="px-4 py-3.5">
                      <input
                        type="checkbox"
                        className="accent-[#C9A84C]"
                        checked={selected.has(p.id)}
                        onChange={() => toggleSelect(p.id)}
                      />
                    </td>
                    <td className="px-4 py-3.5">
                      <p className="font-semibold text-white truncate max-w-48">{p.name}</p>
                      <p className="text-xs text-gray-600 mt-0.5">{p.sizes.join(", ")}</p>
                    </td>
                    <td className="px-4 py-3.5 text-gray-400">{p.brand}</td>
                    <td className="px-4 py-3.5 text-right font-semibold text-white">
                      ${p.price.toLocaleString()}
                    </td>
                    <td className="px-4 py-3.5 text-right">
                      <span
                        className={`px-2.5 py-1 rounded-full text-[11px] font-bold ${
                          p.stock <= 5
                            ? "bg-red-500/15 text-red-400"
                            : "bg-green-500/15 text-green-400"
                        }`}
                      >
                        {p.stock}
                      </span>
                    </td>
                    <td className="px-4 py-3.5 text-right">
                      <div className="flex justify-end gap-3">
                        <button
                          onClick={() => {
                            setStockProduct(p);
                            setStockDelta("0");
                          }}
                          className="text-xs text-blue-400 hover:text-blue-300 font-semibold transition-colors"
                        >
                          Stock
                        </button>
                        <button
                          onClick={() => openEdit(p)}
                          className="text-xs text-gray-400 hover:text-white font-semibold transition-colors"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDelete(p.id)}
                          className="text-xs text-red-400 hover:text-red-300 font-semibold transition-colors"
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
                {filtered.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-4 py-12 text-center text-gray-600 text-sm">
                      No products found
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Create modal */}
      {createOpen && (
        <DarkModal title="Add Product" onClose={() => setCreateOpen(false)}>
          <FormFields form={form} setForm={setForm} />
          <div className="flex gap-3 mt-6">
            <button
              onClick={() => setCreateOpen(false)}
              className="flex-1 border border-white/10 text-gray-400 rounded-xl py-2.5 text-sm font-medium hover:bg-white/5 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleCreate}
              disabled={saving}
              className="flex-1 bg-white text-black rounded-xl py-2.5 text-sm font-bold hover:bg-gray-100 disabled:opacity-50 transition-colors"
            >
              {saving ? "Creating..." : "Create Product"}
            </button>
          </div>
        </DarkModal>
      )}

      {/* Edit modal */}
      {editProduct && (
        <DarkModal title="Edit Product" onClose={() => setEditProduct(null)}>
          <FormFields form={form} setForm={setForm} />
          <div className="flex gap-3 mt-6">
            <button
              onClick={() => setEditProduct(null)}
              className="flex-1 border border-white/10 text-gray-400 rounded-xl py-2.5 text-sm font-medium hover:bg-white/5 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleUpdate}
              disabled={saving}
              className="flex-1 bg-white text-black rounded-xl py-2.5 text-sm font-bold hover:bg-gray-100 disabled:opacity-50 transition-colors"
            >
              {saving ? "Saving..." : "Save Changes"}
            </button>
          </div>
        </DarkModal>
      )}

      {/* Stock modal */}
      {stockProduct && (
        <DarkModal
          title={`Update Stock`}
          onClose={() => setStockProduct(null)}
        >
          <p className="text-sm text-gray-400 mb-4">
            <span className="text-white font-semibold">{stockProduct.name}</span>
            <br />
            <span className="text-gray-500 text-xs">
              Current stock:{" "}
              <span className="text-white font-bold">{stockProduct.stock ?? 0}</span>
            </span>
          </p>
          <div>
            <label className="block text-[11px] font-bold text-gray-500 uppercase tracking-wider mb-1.5">
              Adjustment (+ to add, - to remove)
            </label>
            <input
              type="text"
              inputMode="numeric"
              value={stockDelta}
              onChange={(e) => setStockDelta(e.target.value)}
              placeholder="e.g. -10 or 5"
              className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-white/25 transition-colors"
            />
            <p className="text-xs text-gray-600 mt-2">
              New stock:{" "}
              <span className="text-white font-semibold">
                {(stockProduct.stock ?? 0) + (parseInt(stockDelta, 10) || 0)}
              </span>
            </p>
          </div>
          <div className="flex gap-3 mt-6">
            <button
              onClick={() => setStockProduct(null)}
              className="flex-1 border border-white/10 text-gray-400 rounded-xl py-2.5 text-sm font-medium hover:bg-white/5 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleStockUpdate}
              className="flex-1 bg-white text-black rounded-xl py-2.5 text-sm font-bold hover:bg-gray-100 transition-colors"
            >
              Update Stock
            </button>
          </div>
        </DarkModal>
      )}
    </div>
  );
}
