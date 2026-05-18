import type { Metadata } from "next";
import { Geist } from "next/font/google";
import "./globals.css";
import { CartProvider } from "@/context/CartContext";
import { AuthProvider } from "@/context/AuthContext";
import Navbar from "@/components/Navbar";
import AIChat from "@/components/AIChat";

const geist = Geist({ subsets: ["latin"], variable: "--font-geist" });

export const metadata: Metadata = {
  title: "SNKR VAULT — Premium Sneaker Marketplace",
  description: "Shop the most coveted sneakers. 100% authenticated, free shipping on $100+.",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ru" className={`${geist.variable} h-full`}>
      <body className="min-h-full bg-[#F5F5F0] antialiased">
        <AuthProvider>
          <CartProvider>
            <Navbar />
            <main>{children}</main>
            <AIChat />
          </CartProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
