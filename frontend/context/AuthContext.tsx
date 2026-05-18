"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { User } from "@/lib/api";

interface AuthContextType {
  user: User | null;
  userId: string | null;
  isLoggedIn: boolean;
  loaded: boolean;
  setAuth: (userId: string, user: User) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const storedId = localStorage.getItem("userId");
    const storedUser = localStorage.getItem("user");
    if (storedId && storedUser) {
      setUserId(storedId);
      setUser(JSON.parse(storedUser));
    }
    setLoaded(true);
  }, []);

  const setAuth = (id: string, u: User) => {
    setUserId(id);
    setUser(u);
    localStorage.setItem("userId", id);
    localStorage.setItem("user", JSON.stringify(u));
  };

  const logout = () => {
    setUserId(null);
    setUser(null);
    localStorage.removeItem("userId");
    localStorage.removeItem("user");
  };

  return (
    <AuthContext.Provider value={{ user, userId, isLoggedIn: !!userId, loaded, setAuth, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be inside AuthProvider");
  return ctx;
}
