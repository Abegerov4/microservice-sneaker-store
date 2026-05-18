"use client";

import { useState, useEffect, useRef } from "react";

const SLIDES = [
  {
    id: 0,
    tag: "NEW ARRIVAL · 2026",
    brand: "AIR JORDAN",
    title: "JORDAN 4",
    subtitle: "TORO BRAVO",
    description:
      "Fire red meets premium leather. An icon reborn for a new generation of street culture.",
    image: "/sneakers/jordan4.jpg",
    accentRgb: "180, 20, 40",
  },
  {
    id: 1,
    tag: "CLASSIC ICON",
    brand: "NIKE",
    title: "AIR FORCE 1",
    subtitle: "LOW '07",
    description:
      "The shoe that started it all. Timeless design that transcends trends and defines generations.",
    image: "/sneakers/af1.jpg",
    accentRgb: "60, 80, 160",
  },
  {
    id: 2,
    tag: "LIMITED DROP",
    brand: "ADIDAS ORIGINALS",
    title: "YEEZY BOOST",
    subtitle: "350 V2",
    description:
      "Where design vision meets Boost technology. Every step forward, engineered to perfection.",
    image: "/sneakers/yeezy350.jpg",
    accentRgb: "190, 148, 70",
  },
  {
    id: 3,
    tag: "TRENDING NOW",
    brand: "NIKE",
    title: "DUNK LOW",
    subtitle: "RETRO",
    description:
      "Born for the courts, built for the streets. The silhouette that never goes out of style.",
    image: "/sneakers/dunklow.jpg",
    accentRgb: "30, 110, 55",
  },
];

export default function HeroCarousel() {
  const [current, setCurrent] = useState(0);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const resetTimer = () => {
    if (timerRef.current) clearInterval(timerRef.current);
    timerRef.current = setInterval(() => {
      setCurrent((c) => (c + 1) % SLIDES.length);
    }, 5000);
  };

  useEffect(() => {
    resetTimer();
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const goTo = (idx: number) => {
    if (idx === current) return;
    setCurrent(idx);
    resetTimer();
  };

  const goPrev = () => goTo((current - 1 + SLIDES.length) % SLIDES.length);
  const goNext = () => goTo((current + 1) % SLIDES.length);

  return (
    <section
      className="relative overflow-hidden bg-[#080808] select-none"
      style={{ height: "88vh", minHeight: "540px", maxHeight: "900px" }}
    >
      {SLIDES.map((slide, i) => {
        const active = i === current;
        return (
          <div
            key={slide.id}
            aria-hidden={!active}
            className="absolute inset-0"
            style={{
              opacity: active ? 1 : 0,
              transition: "opacity 850ms ease-in-out",
              zIndex: active ? 2 : 1,
              pointerEvents: active ? "auto" : "none",
            }}
          >
            {/* Accent colour glow */}
            <div
              className="absolute inset-0"
              style={{
                background: `radial-gradient(ellipse 65% 80% at 72% 52%, rgba(${slide.accentRgb}, 0.22) 0%, transparent 65%)`,
              }}
            />

            {/* Text-readability gradient (left → transparent) */}
            <div
              className="absolute inset-0 z-10"
              style={{
                background:
                  "linear-gradient(90deg, #080808 28%, rgba(8,8,8,0.88) 50%, rgba(8,8,8,0.3) 72%, transparent 100%)",
              }}
            />

            {/* Sneaker image — right 60% */}
            <div className="absolute inset-y-0 right-0 z-[5]" style={{ width: "60%" }}>
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={slide.image}
                alt={`${slide.brand} ${slide.title}`}
                className="absolute inset-0 w-full h-full object-contain"
                style={{
                  padding: "5% 8%",
                  transform: active
                    ? "scale(1) translateX(0px)"
                    : "scale(1.06) translateX(28px)",
                  transition:
                    "transform 1100ms cubic-bezier(0.25, 0.46, 0.45, 0.94)",
                }}
              />
            </div>

            {/* Text block */}
            <div
              className="absolute inset-y-0 left-0 z-20 flex flex-col justify-center px-8 sm:px-12 md:px-20 lg:px-28"
              style={{
                transform: active ? "translateY(0)" : "translateY(20px)",
                transition:
                  "transform 900ms cubic-bezier(0.25, 0.46, 0.45, 0.94) 80ms",
              }}
            >
              <p className="text-[10px] sm:text-xs font-black uppercase tracking-[0.35em] text-gray-500 mb-3">
                {slide.tag}
              </p>
              <p className="text-[10px] sm:text-xs font-semibold uppercase tracking-[0.22em] text-gray-600 mb-2">
                {slide.brand}
              </p>
              <h2
                className="font-black text-white leading-[0.88] tracking-tighter mb-1"
                style={{ fontSize: "clamp(2.8rem, 7.5vw, 5.5rem)" }}
              >
                {slide.title}
              </h2>
              <h3
                className="font-black text-gray-600 leading-none tracking-tight mb-6 md:mb-8"
                style={{ fontSize: "clamp(1.1rem, 2.6vw, 2.1rem)" }}
              >
                &ldquo;{slide.subtitle}&rdquo;
              </h3>
              <p className="text-gray-500 text-sm md:text-base leading-relaxed max-w-xs md:max-w-sm mb-8 md:mb-10 hidden sm:block">
                {slide.description}
              </p>
              <div className="flex items-center gap-4 md:gap-6">
                <button
                  onClick={() =>
                    document
                      .getElementById("catalog")
                      ?.scrollIntoView({ behavior: "smooth" })
                  }
                  className="bg-white text-black text-[10px] sm:text-xs font-black uppercase tracking-widest px-6 md:px-8 py-3 md:py-4 hover:bg-gray-100 transition-colors"
                >
                  Shop Now
                </button>
                <button
                  onClick={() =>
                    document
                      .getElementById("catalog")
                      ?.scrollIntoView({ behavior: "smooth" })
                  }
                  className="text-white text-[10px] sm:text-xs font-semibold uppercase tracking-widest border-b border-white/20 hover:border-white pb-px transition-colors hidden sm:block"
                >
                  Explore All
                </button>
              </div>
            </div>
          </div>
        );
      })}

      {/* ── Arrow navigation ── */}
      <button
        onClick={goPrev}
        className="absolute left-4 md:left-6 top-1/2 -translate-y-1/2 z-30 w-11 h-11 border border-white/10 bg-black/20 hover:bg-white/10 backdrop-blur-sm flex items-center justify-center transition-all hover:border-white/30 group"
        aria-label="Previous slide"
      >
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="white"
          strokeWidth="2.5"
          className="opacity-50 group-hover:opacity-100 transition-opacity"
        >
          <polyline points="15 18 9 12 15 6" />
        </svg>
      </button>
      <button
        onClick={goNext}
        className="absolute right-4 md:right-6 top-1/2 -translate-y-1/2 z-30 w-11 h-11 border border-white/10 bg-black/20 hover:bg-white/10 backdrop-blur-sm flex items-center justify-center transition-all hover:border-white/30 group"
        aria-label="Next slide"
      >
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="white"
          strokeWidth="2.5"
          className="opacity-50 group-hover:opacity-100 transition-opacity"
        >
          <polyline points="9 18 15 12 9 6" />
        </svg>
      </button>

      {/* ── Bottom bar: dots + counter ── */}
      <div className="absolute bottom-6 md:bottom-8 inset-x-0 z-30 flex items-center justify-between px-8 sm:px-12 md:px-20 lg:px-28">
        {/* Pill dots */}
        <div className="flex items-center gap-2">
          {SLIDES.map((_, i) => (
            <button
              key={i}
              onClick={() => goTo(i)}
              aria-label={`Go to slide ${i + 1}`}
              style={{
                height: "3px",
                borderRadius: "2px",
                background:
                  i === current ? "white" : "rgba(255,255,255,0.22)",
                transition: "all 400ms ease",
                width: i === current ? "32px" : "8px",
              }}
            />
          ))}
        </div>

        {/* Slide counter */}
        <span className="font-mono text-[11px] tracking-widest text-white/25">
          <span className="text-white/70 font-semibold text-sm">
            {String(current + 1).padStart(2, "0")}
          </span>
          {" — "}
          {String(SLIDES.length).padStart(2, "0")}
        </span>
      </div>

      {/* ── Top progress lines ── */}
      <div className="absolute top-0 inset-x-0 z-30 flex gap-px">
        {SLIDES.map((_, i) => (
          <div key={i} className="flex-1 h-[2px] bg-white/8 overflow-hidden">
            <div
              className="h-full bg-white/50"
              style={{
                width: i < current ? "100%" : i === current ? "100%" : "0%",
                transition:
                  i === current
                    ? "width 5000ms linear"
                    : i < current
                    ? "none"
                    : "none",
              }}
            />
          </div>
        ))}
      </div>
    </section>
  );
}
