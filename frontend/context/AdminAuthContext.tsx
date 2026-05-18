"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";

interface AdminAuthContextType {
  token: string | null;
  isAdmin: boolean;
  loaded: boolean;
  setAdminAuth: (token: string) => void;
  adminLogout: () => void;
}

const AdminAuthContext = createContext<AdminAuthContextType | null>(null);

export function AdminAuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem("adminToken");
    if (stored) setToken(stored);
    setLoaded(true);
  }, []);

  const setAdminAuth = (t: string) => {
    setToken(t);
    localStorage.setItem("adminToken", t);
  };

  const adminLogout = () => {
    setToken(null);
    localStorage.removeItem("adminToken");
  };

  return (
    <AdminAuthContext.Provider value={{ token, isAdmin: !!token, loaded, setAdminAuth, adminLogout }}>
      {children}
    </AdminAuthContext.Provider>
  );
}

export function useAdminAuth() {
  const ctx = useContext(AdminAuthContext);
  if (!ctx) throw new Error("useAdminAuth must be inside AdminAuthProvider");
  return ctx;
}
