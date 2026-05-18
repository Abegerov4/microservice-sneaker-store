"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { updateUser, changePassword, countOrdersByUser } from "@/lib/api";

type Tab = "profile" | "security";

export default function ProfilePage() {
  const { user, userId, isLoggedIn, loaded, setAuth } = useAuth();
  const router = useRouter();
  const [tab, setTab] = useState<Tab>("profile");

  const [fullName, setFullName] = useState("");
  const [phone, setPhone] = useState("");
  const [saving, setSaving] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState("");

  const [orderCount, setOrderCount] = useState<number | null>(null);

  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [pwSaving, setPwSaving] = useState(false);
  const [pwSuccess, setPwSuccess] = useState(false);
  const [pwError, setPwError] = useState("");

  useEffect(() => {
    if (!loaded) return;
    if (!isLoggedIn) { router.push("/auth/login"); return; }
    setFullName(user?.full_name || "");
    setPhone(user?.phone || "");
    if (userId) {
      countOrdersByUser(userId).then(setOrderCount).catch(() => {});
    }
  }, [loaded, isLoggedIn, user]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleSave = async () => {
    if (!userId || !user) return;
    setSaving(true); setError(""); setSuccess(false);
    try {
      const updated = await updateUser(userId, { full_name: fullName, phone });
      setAuth(userId, updated);
      setSuccess(true);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to save");
    } finally {
      setSaving(false);
    }
  };

  const handleChangePassword = async () => {
    if (!userId) return;
    setPwError(""); setPwSuccess(false);
    if (newPassword.length < 6) { setPwError("New password must be at least 6 characters"); return; }
    if (newPassword !== confirmPassword) { setPwError("Passwords do not match"); return; }
    setPwSaving(true);
    try {
      await changePassword(userId, oldPassword, newPassword);
      setPwSuccess(true);
      setOldPassword(""); setNewPassword(""); setConfirmPassword("");
    } catch (e: unknown) {
      setPwError(e instanceof Error ? e.message : "Failed to change password");
    } finally {
      setPwSaving(false);
    }
  };

  if (!isLoggedIn) return null;

  return (
    <div className="min-h-screen bg-[#F5F5F0]">
      <div className="max-w-xl mx-auto px-6 py-10">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <div className="w-14 h-14 rounded-full bg-black text-white flex items-center justify-center text-xl font-black">
            {(user?.full_name || user?.email || "U")[0].toUpperCase()}
          </div>
          <div>
            <h1 className="text-xl font-black text-gray-900">{user?.full_name || "My Account"}</h1>
            <p className="text-sm text-gray-400">{user?.email}</p>
            {orderCount !== null && (
              <p className="text-xs text-gray-400 mt-1">
                {orderCount} order{orderCount !== 1 ? "s" : ""} placed
              </p>
            )}
          </div>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 bg-white rounded-2xl p-1 mb-6 shadow-sm">
          {(["profile", "security"] as Tab[]).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`flex-1 py-2.5 rounded-xl text-sm font-semibold transition-all capitalize ${
                tab === t ? "bg-black text-white" : "text-gray-500 hover:text-black"
              }`}
            >
              {t === "profile" ? "Edit Profile" : "Security"}
            </button>
          ))}
        </div>

        {tab === "profile" && (
          <div className="bg-white rounded-3xl shadow-sm p-6 space-y-5">
            {success && (
              <div className="bg-green-50 border border-green-100 rounded-2xl px-4 py-3 text-green-700 text-sm flex items-center gap-2">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                Profile updated successfully.
              </div>
            )}
            {error && (
              <div className="bg-red-50 border border-red-100 rounded-2xl px-4 py-3 text-red-600 text-sm">{error}</div>
            )}

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Full Name</label>
              <input
                type="text"
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                placeholder="Your full name"
                className="w-full border border-gray-200 rounded-2xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Phone</label>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="+7 777 123 45 67"
                className="w-full border border-gray-200 rounded-2xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Email</label>
              <input
                type="email"
                value={user?.email || ""}
                disabled
                className="w-full border border-gray-100 rounded-2xl px-4 py-3 text-sm bg-gray-50 text-gray-400 cursor-not-allowed"
              />
              <p className="text-xs text-gray-400 mt-1.5">Email address cannot be changed</p>
            </div>

            <button
              onClick={handleSave}
              disabled={saving}
              className="w-full bg-black text-white rounded-2xl py-3 text-sm font-bold hover:bg-gray-800 disabled:opacity-50 transition-colors"
            >
              {saving ? "Saving..." : "Save Changes"}
            </button>
          </div>
        )}

        {tab === "security" && (
          <div className="bg-white rounded-3xl shadow-sm p-6 space-y-5">
            {pwSuccess && (
              <div className="bg-green-50 border border-green-100 rounded-2xl px-4 py-3 text-green-700 text-sm flex items-center gap-2">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                Password changed successfully.
              </div>
            )}
            {pwError && (
              <div className="bg-red-50 border border-red-100 rounded-2xl px-4 py-3 text-red-600 text-sm">{pwError}</div>
            )}

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Current Password</label>
              <input
                type="password"
                value={oldPassword}
                onChange={(e) => setOldPassword(e.target.value)}
                placeholder="Enter current password"
                className="w-full border border-gray-200 rounded-2xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">New Password</label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="Minimum 6 characters"
                className="w-full border border-gray-200 rounded-2xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black"
              />
              {newPassword.length > 0 && (
                <div className="flex gap-1 mt-2">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className={`flex-1 h-1 rounded-full transition-colors ${
                      newPassword.length >= i * 3 ? (newPassword.length >= 9 ? "bg-green-500" : "bg-amber-400") : "bg-gray-100"
                    }`} />
                  ))}
                </div>
              )}
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Confirm New Password</label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Repeat new password"
                className={`w-full border rounded-2xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black ${
                  confirmPassword && confirmPassword !== newPassword ? "border-red-300" : "border-gray-200"
                }`}
              />
              {confirmPassword && confirmPassword !== newPassword && (
                <p className="text-xs text-red-500 mt-1.5">Passwords do not match</p>
              )}
            </div>

            <button
              onClick={handleChangePassword}
              disabled={pwSaving || !oldPassword || !newPassword || !confirmPassword}
              className="w-full bg-black text-white rounded-2xl py-3 text-sm font-bold hover:bg-gray-800 disabled:opacity-50 transition-colors"
            >
              {pwSaving ? "Changing..." : "Change Password"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
