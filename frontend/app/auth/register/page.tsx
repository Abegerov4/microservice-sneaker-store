"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { createUser } from "@/lib/api";
import { useAuth } from "@/context/AuthContext";

export default function RegisterPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [showPw, setShowPw] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { setAuth } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const user = await createUser(email, password, fullName);
      setAuth(user.id, user);
      router.push("/");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[calc(100vh-64px)] bg-[#F5F5F0] flex">
      {/* Left panel */}
      <div className="hidden lg:flex flex-1 bg-black items-center justify-center p-16">
        <div className="max-w-sm">
          <div className="flex items-baseline gap-1 mb-8">
            <span className="text-2xl font-black text-white">SNKR</span>
            <span className="text-[11px] font-black text-[#C9A84C] tracking-[0.35em] ml-0.5">VAULT</span>
          </div>
          <h2 className="text-5xl font-black text-white leading-none mb-6">
            JOIN THE<br />
            <span className="text-gray-600">SNEAKER</span><br />
            CULTURE.
          </h2>
          <p className="text-gray-500 text-sm leading-relaxed">
            Create an account to start shopping, track your orders, and get early access to limited releases.
          </p>
        </div>
      </div>

      {/* Right panel */}
      <div className="flex-1 flex items-center justify-center px-6 py-12">
        <div className="w-full max-w-sm">
          <h1 className="text-2xl font-black mb-1">Create account</h1>
          <p className="text-gray-400 text-sm mb-8">Sign up for free today</p>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Full Name</label>
              <input
                type="text"
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                required
                autoComplete="name"
                className="w-full border border-gray-200 bg-white rounded-2xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black"
                placeholder="John Doe"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoComplete="email"
                className="w-full border border-gray-200 bg-white rounded-2xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-black"
                placeholder="you@example.com"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">Password</label>
              <div className="relative">
                <input
                  type={showPw ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  minLength={6}
                  autoComplete="new-password"
                  className="w-full border border-gray-200 bg-white rounded-2xl px-4 py-3 pr-11 text-sm focus:outline-none focus:ring-2 focus:ring-black"
                  placeholder="Minimum 6 characters"
                />
                <button
                  type="button"
                  onClick={() => setShowPw(!showPw)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    {showPw
                      ? <><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/><path d="M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19"/><line x1="1" y1="1" x2="23" y2="23"/></>
                      : <><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></>
                    }
                  </svg>
                </button>
              </div>
              {password.length > 0 && (
                <div className="flex gap-1 mt-2">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className={`flex-1 h-1 rounded-full transition-colors ${
                      password.length >= i * 3 ? (password.length >= 9 ? "bg-green-500" : "bg-amber-400") : "bg-gray-100"
                    }`} />
                  ))}
                </div>
              )}
            </div>

            {error && (
              <div className="bg-red-50 border border-red-100 rounded-2xl px-4 py-3 text-red-600 text-sm">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-black text-white py-3.5 rounded-2xl font-bold text-sm hover:bg-gray-800 disabled:opacity-50 transition-colors"
            >
              {loading ? "Creating account..." : "Create Account"}
            </button>
          </form>

          <p className="text-center text-sm text-gray-400 mt-8">
            Already have an account?{" "}
            <Link href="/auth/login" className="font-bold text-black hover:underline underline-offset-2">
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
