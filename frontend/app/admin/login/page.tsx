"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Eye, EyeOff, Lock } from "lucide-react";
import { login } from "@/lib/api";
import { useAdminAuth } from "@/context/AdminAuthContext";

export default function AdminLoginPage() {
  const { setAdminAuth } = useAdminAuth();
  const router = useRouter();
  const [email, setEmail] = useState("admin@sneakerstore.com");
  const [password, setPassword] = useState("");
  const [showPw, setShowPw] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const res = await login(email, password);
      if (!res.success) {
        setError("Invalid credentials");
        return;
      }
      if (res.role !== "ADMIN") {
        setError("Access denied: admin account required");
        return;
      }
      setAdminAuth(res.token);
      router.push("/admin");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#0A0A0A] flex items-center justify-center px-4">
      <div className="w-full max-w-sm">
        {/* Logo */}
        <div className="text-center mb-10">
          <div className="flex items-baseline justify-center gap-1 mb-2">
            <span className="text-3xl font-black text-white">SNKR</span>
            <span className="text-[12px] font-black text-[#C9A84C] tracking-[0.35em] ml-0.5">
              VAULT
            </span>
          </div>
          <p className="text-[11px] font-bold text-gray-600 uppercase tracking-[0.3em]">
            Admin Portal
          </p>
        </div>

        <div className="bg-[#161616] border border-white/8 rounded-2xl p-8">
          <div className="flex items-center justify-center w-12 h-12 bg-white/5 rounded-2xl mb-6 mx-auto">
            <Lock size={20} className="text-gray-400" />
          </div>
          <h1 className="text-xl font-black text-white text-center mb-1">Sign In</h1>
          <p className="text-gray-500 text-sm text-center mb-8">
            Admin access only
          </p>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">
                Email
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-white/25 transition-colors"
                required
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">
                Password
              </label>
              <div className="relative">
                <input
                  type={showPw ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter password"
                  className="w-full bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl px-4 py-3 pr-11 text-sm focus:outline-none focus:border-white/25 transition-colors"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPw((v) => !v)}
                  className="absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 transition-colors"
                >
                  {showPw ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
            </div>

            {error && (
              <div className="bg-red-500/10 border border-red-500/20 rounded-xl px-4 py-3 text-red-400 text-sm">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-white text-black py-3 rounded-xl font-bold text-sm hover:bg-gray-100 disabled:opacity-50 transition-colors mt-2"
            >
              {loading ? "Signing in..." : "Sign In"}
            </button>
          </form>
        </div>

        <p className="text-center text-xs text-gray-700 mt-6">
          Not an admin?{" "}
          <a href="/" className="text-gray-500 hover:text-white transition-colors">
            Return to store
          </a>
        </p>
      </div>
    </div>
  );
}
