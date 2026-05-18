"use client";

import { useState, useRef, useEffect } from "react";
import { Sparkles, X, Send, ChevronDown, Bot, User } from "lucide-react";
import { aiChat, SneakerRecommendation, TrendingSneaker, aiRecommend, aiTrending, aiSearchByStyle } from "@/lib/api";
import Link from "next/link";

interface Message {
  role: "user" | "assistant";
  content: string;
  recommendations?: SneakerRecommendation[];
  trending?: TrendingSneaker[];
  styleResults?: SneakerRecommendation[];
}

function RecommendationCard({ rec }: { rec: SneakerRecommendation }) {
  return (
    <Link
      href={`/products/${rec.product_id}`}
      className="flex gap-3 p-2.5 bg-white/5 hover:bg-white/10 rounded-xl border border-white/8 transition-all group"
    >
      {rec.image_url ? (
        <img
          src={rec.image_url}
          alt={rec.name}
          className="w-14 h-14 object-cover rounded-lg flex-shrink-0 bg-white/5"
        />
      ) : (
        <div className="w-14 h-14 rounded-lg bg-white/5 flex-shrink-0" />
      )}
      <div className="min-w-0 flex-1">
        <p className="text-white text-xs font-bold truncate group-hover:text-[#C9A84C] transition-colors">
          {rec.name}
        </p>
        <p className="text-gray-500 text-[10px]">{rec.brand}</p>
        <p className="text-[#C9A84C] text-xs font-bold mt-1">${rec.price.toLocaleString()}</p>
        {rec.reason && (
          <p className="text-gray-600 text-[10px] mt-0.5 line-clamp-2">{rec.reason}</p>
        )}
      </div>
    </Link>
  );
}

function TrendingCard({ sneaker }: { sneaker: TrendingSneaker }) {
  return (
    <Link
      href={`/products/${sneaker.product_id}`}
      className="flex gap-3 p-2.5 bg-white/5 hover:bg-white/10 rounded-xl border border-white/8 transition-all group"
    >
      {sneaker.image_url ? (
        <img
          src={sneaker.image_url}
          alt={sneaker.name}
          className="w-14 h-14 object-cover rounded-lg flex-shrink-0 bg-white/5"
        />
      ) : (
        <div className="w-14 h-14 rounded-lg bg-white/5 flex-shrink-0" />
      )}
      <div className="min-w-0 flex-1">
        <p className="text-white text-xs font-bold truncate group-hover:text-[#C9A84C] transition-colors">
          {sneaker.name}
        </p>
        <p className="text-gray-500 text-[10px]">{sneaker.brand}</p>
        <p className="text-[#C9A84C] text-xs font-bold mt-1">${sneaker.price.toLocaleString()}</p>
        {sneaker.trend_reason && (
          <p className="text-gray-600 text-[10px] mt-0.5 line-clamp-2">{sneaker.trend_reason}</p>
        )}
      </div>
    </Link>
  );
}

const QUICK_PROMPTS = [
  "What's trending right now?",
  "Best Nike under $200",
  "Comfortable running shoes",
  "Limited edition drops",
];

export default function AIChat() {
  const [open, setOpen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([
    {
      role: "assistant",
      content:
        "Hey! I'm VAULT AI — your personal sneaker advisor. Ask me anything about kicks, and I'll help you find the perfect pair. What are you looking for? 🔥",
    },
  ]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [sessionId, setSessionId] = useState("");
  const bottomRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (open) {
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  }, [open, messages]);

  const send = async (text?: string) => {
    const msg = (text || input).trim();
    if (!msg || loading) return;

    setInput("");
    setMessages((prev) => [...prev, { role: "user", content: msg }]);
    setLoading(true);

    try {
      const { reply, session_id } = await aiChat(msg, sessionId);
      if (!sessionId) setSessionId(session_id);

      let recs: SneakerRecommendation[] | undefined;
      let trendingList: TrendingSneaker[] | undefined;
      let styleResults: SneakerRecommendation[] | undefined;
      const lower = msg.toLowerCase();

      if (lower.includes("recommend") || lower.includes("suggest") || lower.includes("find me")) {
        try {
          const { recommendations } = await aiRecommend(msg, 0, "", "");
          if (recommendations?.length) recs = recommendations.slice(0, 3);
        } catch {
          // optional
        }
      } else if (lower.includes("trending") || lower.includes("popular") || lower.includes("hot")) {
        try {
          const { sneakers } = await aiTrending(4);
          if (sneakers?.length) trendingList = sneakers;
        } catch {
          // optional
        }
      } else if (lower.includes("style") || lower.includes("look") || lower.includes("aesthetic") || lower.includes("vibe")) {
        try {
          const { results } = await aiSearchByStyle(msg);
          if (results?.length) styleResults = results.slice(0, 3);
        } catch {
          // optional
        }
      }

      setMessages((prev) => [...prev, { role: "assistant", content: reply, recommendations: recs, trending: trendingList, styleResults }]);
    } catch {
      setMessages((prev) => [
        ...prev,
        {
          role: "assistant",
          content: "Sorry, I'm having trouble connecting right now. Please try again in a moment.",
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      {/* Floating button */}
      <button
        onClick={() => setOpen((v) => !v)}
        className={`fixed bottom-6 right-6 z-50 w-14 h-14 rounded-full shadow-2xl flex items-center justify-center transition-all duration-300 ${
          open
            ? "bg-white text-black scale-90"
            : "bg-[#0A0A0A] text-white hover:scale-110 border border-white/10"
        }`}
        aria-label="Open AI chat"
      >
        {open ? <ChevronDown size={22} /> : <Sparkles size={22} className="text-[#C9A84C]" />}
      </button>

      {/* Chat panel */}
      <div
        className={`fixed bottom-24 right-6 z-50 w-[360px] max-h-[580px] flex flex-col bg-[#0D0D0D] border border-white/10 rounded-2xl shadow-2xl transition-all duration-300 overflow-hidden ${
          open ? "opacity-100 translate-y-0 pointer-events-auto" : "opacity-0 translate-y-4 pointer-events-none"
        }`}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-white/8 flex-shrink-0">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-full bg-[#C9A84C]/15 flex items-center justify-center">
              <Sparkles size={14} className="text-[#C9A84C]" />
            </div>
            <div>
              <p className="text-white text-sm font-bold leading-none">VAULT AI</p>
              <p className="text-gray-600 text-[10px] mt-0.5">Sneaker advisor</p>
            </div>
          </div>
          <button
            onClick={() => setOpen(false)}
            className="text-gray-600 hover:text-white transition-colors"
          >
            <X size={16} />
          </button>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto px-4 py-3 space-y-3 min-h-0">
          {messages.map((m, i) => (
            <div
              key={i}
              className={`flex gap-2 ${m.role === "user" ? "flex-row-reverse" : "flex-row"}`}
            >
              <div
                className={`w-7 h-7 rounded-full flex-shrink-0 flex items-center justify-center ${
                  m.role === "user"
                    ? "bg-white/10"
                    : "bg-[#C9A84C]/15"
                }`}
              >
                {m.role === "user" ? (
                  <User size={13} className="text-gray-400" />
                ) : (
                  <Bot size={13} className="text-[#C9A84C]" />
                )}
              </div>
              <div className={`flex-1 ${m.role === "user" ? "flex flex-col items-end" : ""}`}>
                <div
                  className={`px-3 py-2 rounded-xl text-sm leading-relaxed max-w-[85%] ${
                    m.role === "user"
                      ? "bg-white text-black font-medium"
                      : "bg-white/6 text-gray-300"
                  }`}
                >
                  {m.content}
                </div>
                {m.recommendations && m.recommendations.length > 0 && (
                  <div className="mt-2 w-full space-y-2">
                    <p className="text-[10px] text-gray-600 font-bold uppercase tracking-wider">
                      Recommended for you
                    </p>
                    {m.recommendations.map((r) => (
                      <RecommendationCard key={r.product_id} rec={r} />
                    ))}
                  </div>
                )}
                {m.trending && m.trending.length > 0 && (
                  <div className="mt-2 w-full space-y-2">
                    <p className="text-[10px] text-gray-600 font-bold uppercase tracking-wider">
                      Trending now
                    </p>
                    {m.trending.map((s) => (
                      <TrendingCard key={s.product_id} sneaker={s} />
                    ))}
                  </div>
                )}
                {m.styleResults && m.styleResults.length > 0 && (
                  <div className="mt-2 w-full space-y-2">
                    <p className="text-[10px] text-gray-600 font-bold uppercase tracking-wider">
                      Style matches
                    </p>
                    {m.styleResults.map((r) => (
                      <RecommendationCard key={r.product_id} rec={r} />
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}

          {loading && (
            <div className="flex gap-2">
              <div className="w-7 h-7 rounded-full bg-[#C9A84C]/15 flex items-center justify-center flex-shrink-0">
                <Bot size={13} className="text-[#C9A84C]" />
              </div>
              <div className="bg-white/6 rounded-xl px-3 py-2">
                <div className="flex gap-1 items-center h-5">
                  {[0, 1, 2].map((i) => (
                    <span
                      key={i}
                      className="w-1.5 h-1.5 bg-gray-500 rounded-full animate-bounce"
                      style={{ animationDelay: `${i * 0.15}s` }}
                    />
                  ))}
                </div>
              </div>
            </div>
          )}

          <div ref={bottomRef} />
        </div>

        {/* Quick prompts (show only when only the initial message) */}
        {messages.length === 1 && (
          <div className="px-4 pb-2 flex flex-wrap gap-1.5 flex-shrink-0">
            {QUICK_PROMPTS.map((p) => (
              <button
                key={p}
                onClick={() => send(p)}
                className="px-2.5 py-1 rounded-full text-[11px] font-medium border border-white/10 text-gray-400 hover:border-[#C9A84C]/40 hover:text-[#C9A84C] transition-colors"
              >
                {p}
              </button>
            ))}
          </div>
        )}

        {/* Input */}
        <div className="flex gap-2 px-3 pb-3 pt-2 border-t border-white/8 flex-shrink-0">
          <input
            ref={inputRef}
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && !e.shiftKey && send()}
            placeholder="Ask about any sneaker..."
            className="flex-1 bg-white/5 border border-white/10 text-white placeholder-gray-600 rounded-xl px-3 py-2 text-sm focus:outline-none focus:border-white/20 transition-colors"
          />
          <button
            onClick={() => send()}
            disabled={!input.trim() || loading}
            className="w-9 h-9 rounded-xl bg-white flex items-center justify-center text-black disabled:opacity-30 hover:bg-gray-100 transition-colors flex-shrink-0"
          >
            <Send size={14} />
          </button>
        </div>
      </div>
    </>
  );
}
