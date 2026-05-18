const MAP: Array<[string[], string]> = [
  // Nike Air Force 1
  [["air force 1", "af1", "force 1"], "/sneakers/af1.jpg"],
  // Nike Dunk
  [["dunk low", "dunk high", "dunk sb"], "/sneakers/dunklow.jpg"],
  // Air Jordan 1
  [["jordan 1", "air jordan 1"], "/sneakers/jordan1.jpg"],
  // Air Jordan 4
  [["jordan 4", "air jordan 4"], "/sneakers/jordan4.jpg"],
  // Jordan Why Not
  [["why not"], "/sneakers/jordanwhy.jpg"],
  // Air Max 90
  [["air max 90"], "/sneakers/airmax90.jpg"],
  // Air Max 270
  [["air max 270"], "/sneakers/airmax270.jpg"],
  // Air Max Plus / TN
  [["air max plus", "air max tn", " tn "], "/sneakers/airmaxplus.jpg"],
  // Yeezy
  [["yeezy 350", "yeezy boost 350", "yeezy"], "/sneakers/yeezy350.jpg"],
  // Adidas Samba
  [["samba"], "/sneakers/samba.jpg"],
  // Adidas NMD
  [["nmd"], "/sneakers/nmd.jpg"],
  // Stan Smith
  [["stan smith"], "/sneakers/stansmith.jpg"],
  // Ultraboost
  [["ultraboost", "ultra boost"], "/sneakers/ultraboost.jpg"],
  // New Balance 990
  [["990", "990v5", "990v6"], "/sneakers/nb990.jpg"],
  // New Balance 550
  [["550"], "/sneakers/nb550.jpg"],
  // New Balance 574
  [["574", "new balance"], "/sneakers/nb574.jpg"],
  // Vans Old Skool
  [["old skool"], "/sneakers/vansoldskool.jpg"],
  // Vans Authentic
  [["vans authentic", "vans"], "/sneakers/vansauthentic.jpg"],
  // Under Armour
  [["under armour", "hovr", "ua "], "/sneakers/underarmour.jpg"],
];

export function getProductImage(name: string, brand: string, imageUrl?: string): string {
  if (imageUrl) return imageUrl;

  const text = `${name} ${brand}`.toLowerCase();

  for (const [keywords, path] of MAP) {
    if (keywords.some((kw) => text.includes(kw))) {
      return path;
    }
  }

  // Generic fallback by brand
  if (text.includes("nike")) return "/sneakers/af1.jpg";
  if (text.includes("adidas")) return "/sneakers/samba.jpg";
  if (text.includes("new balance")) return "/sneakers/nb574.jpg";
  if (text.includes("vans")) return "/sneakers/vansoldskool.jpg";

  const fallbacks = ["/sneakers/af1.jpg", "/sneakers/dunklow.jpg", "/sneakers/nb574.jpg", "/sneakers/samba.jpg"];
  const hash = (name + brand).split("").reduce((a, c) => a + c.charCodeAt(0), 0);
  return fallbacks[hash % fallbacks.length];
}
