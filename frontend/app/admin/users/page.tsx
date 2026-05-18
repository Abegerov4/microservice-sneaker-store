"use client";

import { useEffect, useState } from "react";
import { Users, Search, X, Eye, EyeOff } from "lucide-react";
import { useAdminAuth } from "@/context/AdminAuthContext";
import {
  User,
  adminListUsers,
  adminSearchUsers,
  adminUpdateUserStatus,
  adminResetPassword,
  adminDeleteUser,
  getUserByEmail,
  countOrdersByUser,
} from "@/lib/api";

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
      <div className="bg-[#1a1a1a] border border-white/10 rounded-2xl shadow-2xl w-full max-w-sm">
        <div className="flex items-center justify-between px-6 py-4 border-b border-white/8">
          <h2 className="font-bold text-white text-sm">{title}</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-white transition-colors">
            <X size={18} />
          </button>
        </div>
        <div className="p-6">{children}</div>
      </div>
    </div>
  );
}

export default function AdminUsersPage() {
  const { token } = useAdminAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [searching, setSearching] = useState(false);
  const [orderCounts, setOrderCounts] = useState<Record<string, number>>({});
  const [loadingCount, setLoadingCount] = useState<string | null>(null);
  const [resetUser, setResetUser] = useState<User | null>(null);
  const [newPassword, setNewPassword] = useState("");
  const [showPw, setShowPw] = useState(false);
  const [saving, setSaving] = useState(false);

  const load = async () => {
    if (!token) return;
    try {
      const us = await adminListUsers(token);
      setUsers(us);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, [token]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleSearch = async (q: string) => {
    setSearch(q);
    if (!token) return;
    if (!q.trim()) {
      load();
      return;
    }
    setSearching(true);
    try {
      if (q.includes("@")) {
        const u = await getUserByEmail(q.trim());
        setUsers([u]);
      } else {
        const us = await adminSearchUsers(token, q);
        setUsers(us);
      }
    } catch {
      setUsers([]);
    } finally {
      setSearching(false);
    }
  };

  const handleLoadOrderCount = async (userId: string) => {
    setLoadingCount(userId);
    try {
      const count = await countOrdersByUser(userId);
      setOrderCounts((prev) => ({ ...prev, [userId]: count }));
    } catch {
      setOrderCounts((prev) => ({ ...prev, [userId]: 0 }));
    } finally {
      setLoadingCount(null);
    }
  };

  const handleToggleStatus = async (user: User) => {
    if (!token) return;
    try {
      await adminUpdateUserStatus(token, user.id, !user.active);
      load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to update status");
    }
  };

  const handleResetPassword = async () => {
    if (!token || !resetUser || !newPassword) return;
    if (newPassword.length < 6) {
      setError("Password must be at least 6 characters");
      return;
    }
    setSaving(true);
    try {
      await adminResetPassword(token, resetUser.id, newPassword);
      setResetUser(null);
      setNewPassword("");
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to reset password");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (user: User) => {
    if (!token || !confirm(`Delete user ${user.email}? This cannot be undone.`)) return;
    try {
      await adminDeleteUser(token, user.id);
      load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to delete user");
    }
  };

  const activeCount = users.filter((u) => u.active).length;

  return (
    <div>
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Users size={14} className="text-[#C9A84C]" />
            <span className="text-[10px] font-black uppercase tracking-[0.3em] text-gray-600">
              Management
            </span>
          </div>
          <h1 className="text-2xl font-black text-white">Users</h1>
          <p className="text-gray-600 text-sm mt-1">
            {users.length} total · {activeCount} active
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

      {/* Search */}
      <div className="flex gap-3 mb-5">
        <div className="relative flex-1 max-w-xs">
          <Search size={13} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-600" />
          <input
            type="text"
            placeholder="Name, or exact email (user@example.com)..."
            value={search}
            onChange={(e) => handleSearch(e.target.value)}
            className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl pl-9 pr-4 py-2.5 text-sm focus:outline-none focus:border-white/20 transition-colors"
          />
          {searching && (
            <span className="absolute right-3 top-2.5 text-gray-600 text-xs">...</span>
          )}
        </div>
        {search && (
          <button
            onClick={() => {
              setSearch("");
              load();
            }}
            className="text-sm text-gray-500 hover:text-white transition-colors"
          >
            Clear
          </button>
        )}
      </div>

      <div className="bg-[#161616] border border-white/6 rounded-2xl overflow-hidden">
        {loading ? (
          <div className="p-12 text-center text-gray-600 text-sm">Loading users...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/6">
                  {["User", "Role", "Status", "Joined", "Orders", "Actions"].map((h) => (
                    <th
                      key={h}
                      className={`px-4 py-3 text-[10px] font-bold text-gray-600 uppercase tracking-wider ${
                        h === "Actions" ? "text-right" : "text-left"
                      }`}
                    >
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-white/4">
                {users.map((user) => (
                  <tr key={user.id} className="hover:bg-white/3 transition-colors">
                    <td className="px-4 py-3.5">
                      <p className="font-semibold text-white">{user.full_name || "—"}</p>
                      <p className="text-xs text-gray-500 mt-0.5">{user.email}</p>
                    </td>
                    <td className="px-4 py-3.5">
                      <span
                        className={`px-2.5 py-1 rounded-full text-[11px] font-bold ${
                          user.role === "ADMIN"
                            ? "bg-[#C9A84C]/15 text-[#C9A84C]"
                            : "bg-white/8 text-gray-500"
                        }`}
                      >
                        {user.role || "USER"}
                      </span>
                    </td>
                    <td className="px-4 py-3.5">
                      <span
                        className={`px-2.5 py-1 rounded-full text-[11px] font-bold ${
                          user.active
                            ? "bg-green-500/15 text-green-400"
                            : "bg-red-500/15 text-red-400"
                        }`}
                      >
                        {user.active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td className="px-4 py-3.5 text-xs text-gray-500">
                      {new Date(user.created_at).toLocaleDateString("en-US", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                      })}
                    </td>
                    <td className="px-4 py-3.5 text-left">
                      {orderCounts[user.id] !== undefined ? (
                        <span className="text-[11px] font-bold text-gray-400 bg-white/8 px-2 py-1 rounded-full">
                          {orderCounts[user.id]}
                        </span>
                      ) : (
                        <button
                          onClick={() => handleLoadOrderCount(user.id)}
                          disabled={loadingCount === user.id}
                          className="text-[11px] text-gray-600 hover:text-[#C9A84C] font-medium transition-colors"
                        >
                          {loadingCount === user.id ? "..." : "Count"}
                        </button>
                      )}
                    </td>
                    <td className="px-4 py-3.5 text-right">
                      <div className="flex justify-end gap-3">
                        <button
                          onClick={() => handleToggleStatus(user)}
                          className={`text-xs font-semibold transition-colors ${
                            user.active
                              ? "text-yellow-500 hover:text-yellow-400"
                              : "text-green-500 hover:text-green-400"
                          }`}
                        >
                          {user.active ? "Deactivate" : "Activate"}
                        </button>
                        <button
                          onClick={() => {
                            setResetUser(user);
                            setNewPassword("");
                          }}
                          className="text-xs text-blue-400 hover:text-blue-300 font-semibold transition-colors"
                        >
                          Reset PW
                        </button>
                        {user.role !== "ADMIN" && (
                          <button
                            onClick={() => handleDelete(user)}
                            className="text-xs text-red-400 hover:text-red-300 font-semibold transition-colors"
                          >
                            Delete
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
                {users.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-4 py-12 text-center text-gray-600 text-sm">
                      No users found
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Reset password modal */}
      {resetUser && (
        <DarkModal
          title="Reset Password"
          onClose={() => setResetUser(null)}
        >
          <p className="text-sm text-gray-400 mb-5">
            Set a new password for{" "}
            <span className="text-white font-semibold">{resetUser.email}</span>
          </p>
          <div>
            <label className="block text-[11px] font-bold text-gray-500 uppercase tracking-wider mb-1.5">
              New Password
            </label>
            <div className="relative">
              <input
                type={showPw ? "text" : "password"}
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="Minimum 6 characters"
                className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl px-4 py-2.5 pr-10 text-sm focus:outline-none focus:border-white/25 transition-colors"
              />
              <button
                type="button"
                onClick={() => setShowPw((v) => !v)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 transition-colors"
              >
                {showPw ? <EyeOff size={15} /> : <Eye size={15} />}
              </button>
            </div>
          </div>
          <div className="flex gap-3 mt-6">
            <button
              onClick={() => setResetUser(null)}
              className="flex-1 border border-white/10 text-gray-400 rounded-xl py-2.5 text-sm font-medium hover:bg-white/5 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleResetPassword}
              disabled={saving || !newPassword}
              className="flex-1 bg-white text-black rounded-xl py-2.5 text-sm font-bold hover:bg-gray-100 disabled:opacity-50 transition-colors"
            >
              {saving ? "Resetting..." : "Reset Password"}
            </button>
          </div>
        </DarkModal>
      )}
    </div>
  );
}
