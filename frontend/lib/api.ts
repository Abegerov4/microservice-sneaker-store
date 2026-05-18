const API_BASE = "/api/v1";
const ADMIN_BASE = "/api/v1/admin";

export interface Product {
  id: string;
  name: string;
  brand: string;
  description: string;
  price: number;
  sizes: string[];
  stock: number;
  image_url: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: string;
  email: string;
  full_name: string;
  phone: string;
  active: boolean;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  product_id: string;
  product_name: string;
  quantity: number;
  price: number;
  size: string;
}

export interface Order {
  id: string;
  user_id: string;
  items: OrderItem[];
  total_amount: number;
  status: string;
  shipping_address: string;
  created_at: string;
  updated_at: string;
}

export interface OrderItemInput {
  product_id: string;
  quantity: number;
  size: string;
}

export async function searchProducts(params: { brand?: string; size?: string }): Promise<Product[]> {
  const qs = new URLSearchParams();
  if (params.brand) qs.set("brand", params.brand);
  if (params.size) qs.set("size", params.size);
  const res = await fetch(`${API_BASE}/products/search?${qs}`);
  if (!res.ok) throw new Error("Failed to search products");
  const data = await res.json();
  return data.products || [];
}

export async function getProductsByBrand(brand: string): Promise<Product[]> {
  const res = await fetch(`${API_BASE}/products/by-brand/${encodeURIComponent(brand)}`);
  if (!res.ok) throw new Error("Failed to load products by brand");
  const data = await res.json();
  return data.products || [];
}

export async function listProducts(): Promise<Product[]> {
  const res = await fetch(`${API_BASE}/products`);
  if (!res.ok) throw new Error("Failed to load products");
  const data = await res.json();
  return data.products || [];
}

export async function getProduct(id: string): Promise<Product> {
  const res = await fetch(`${API_BASE}/products/${id}`);
  if (!res.ok) throw new Error("Product not found");
  return res.json();
}

export async function getBrands(): Promise<string[]> {
  const res = await fetch(`${API_BASE}/products/brands`);
  if (!res.ok) return [];
  const data = await res.json();
  return data.brands || [];
}

export async function createUser(email: string, password: string, full_name: string): Promise<User> {
  const res = await fetch(`${API_BASE}/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, full_name }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Registration failed");
  }
  return res.json();
}

export async function login(
  email: string,
  password: string
): Promise<{ success: boolean; user_id: string; role: string; token: string; message: string }> {
  const res = await fetch(`${API_BASE}/users/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Login failed");
  }
  return res.json();
}

export async function getUser(id: string): Promise<User> {
  const res = await fetch(`${API_BASE}/users/${id}`);
  if (!res.ok) throw new Error("User not found");
  return res.json();
}

export async function changePassword(id: string, old_password: string, new_password: string): Promise<void> {
  const res = await fetch(`${API_BASE}/users/${id}/password`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ old_password, new_password }),
  });
  if (!res.ok) {
    const d = await res.json().catch(() => ({}));
    throw new Error(d.error || "Failed to change password");
  }
}

export async function updateUser(id: string, data: { full_name?: string; phone?: string }): Promise<User> {
  const res = await fetch(`${API_BASE}/users/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const d = await res.json().catch(() => ({}));
    throw new Error(d.error || "Failed to update profile");
  }
  return res.json();
}

export async function createOrder(
  user_id: string,
  items: OrderItemInput[],
  shipping_address: string
): Promise<Order> {
  const res = await fetch(`${API_BASE}/orders`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_id, items, shipping_address }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Failed to create order");
  }
  return res.json();
}

export async function getOrdersByUser(user_id: string): Promise<Order[]> {
  const res = await fetch(`${API_BASE}/orders/user/${user_id}`);
  if (!res.ok) return [];
  const data = await res.json();
  return data.orders || [];
}

export async function getOrder(id: string): Promise<Order> {
  const res = await fetch(`${API_BASE}/orders/${id}`);
  if (!res.ok) throw new Error("Order not found");
  return res.json();
}

export async function getOrderItems(id: string): Promise<OrderItem[]> {
  const res = await fetch(`${API_BASE}/orders/${id}/items`);
  if (!res.ok) throw new Error("Failed to load order items");
  const data = await res.json();
  return data.items || [];
}

export async function countOrdersByUser(userId: string): Promise<number> {
  const res = await fetch(`${API_BASE}/orders/user/${userId}/count`);
  if (!res.ok) return 0;
  const data = await res.json();
  return data.count || 0;
}

export async function getOrdersByDateRange(from: string, to: string): Promise<Order[]> {
  const fromRFC = `${from}T00:00:00Z`;
  const toRFC = `${to}T23:59:59Z`;
  const qs = new URLSearchParams({ from: fromRFC, to: toRFC });
  const res = await fetch(`${API_BASE}/orders/by-date?${qs}`);
  if (!res.ok) throw new Error("Failed to load orders by date range");
  const data = await res.json();
  return data.orders || [];
}

export async function getUserByEmail(email: string): Promise<User> {
  const res = await fetch(`${API_BASE}/users/by-email?email=${encodeURIComponent(email)}`);
  if (!res.ok) throw new Error("User not found");
  return res.json();
}

// ── Admin API ──────────────────────────────────────────────────────────────────

function adminHeaders(token: string) {
  return { "Content-Type": "application/json", Authorization: `Bearer ${token}` };
}

export interface DashStats {
  total_products: number;
  total_orders: number;
  total_users: number;
  total_revenue: number;
  pending_orders: number;
  active_users: number;
}

export async function adminGetStats(token: string): Promise<DashStats> {
  const res = await fetch(`${ADMIN_BASE}/stats`, { headers: adminHeaders(token) });
  if (!res.ok) throw new Error("Failed to load stats");
  return res.json();
}

// Products
export async function adminListProducts(token: string): Promise<Product[]> {
  const res = await fetch(`${ADMIN_BASE}/products`, { headers: adminHeaders(token) });
  if (!res.ok) throw new Error("Failed to load products");
  const data = await res.json();
  return data.products || [];
}

export async function adminCreateProduct(token: string, product: Partial<Product>): Promise<Product> {
  const res = await fetch(`${ADMIN_BASE}/products`, {
    method: "POST",
    headers: adminHeaders(token),
    body: JSON.stringify(product),
  });
  if (!res.ok) {
    const d = await res.json().catch(() => ({}));
    throw new Error(d.error || "Failed to create product");
  }
  return res.json();
}

export async function adminUpdateProduct(token: string, id: string, product: Partial<Product>): Promise<Product> {
  const res = await fetch(`${ADMIN_BASE}/products/${id}`, {
    method: "PUT",
    headers: adminHeaders(token),
    body: JSON.stringify(product),
  });
  if (!res.ok) {
    const d = await res.json().catch(() => ({}));
    throw new Error(d.error || "Failed to update product");
  }
  return res.json();
}

export async function adminDeleteProduct(token: string, id: string): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/products/${id}`, {
    method: "DELETE",
    headers: adminHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to delete product");
}

export async function adminBulkDeleteProducts(token: string, ids: string[]): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/products/bulk-delete`, {
    method: "POST",
    headers: adminHeaders(token),
    body: JSON.stringify({ ids }),
  });
  if (!res.ok) throw new Error("Failed to bulk delete products");
}

export async function adminGetLowStockProducts(token: string): Promise<Product[]> {
  const res = await fetch(`${ADMIN_BASE}/products/low-stock`, { headers: adminHeaders(token) });
  if (!res.ok) throw new Error("Failed to load low stock products");
  const data = await res.json();
  return data.products || [];
}

export async function adminUpdateStock(token: string, id: string, delta: number): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/products/${id}/stock`, {
    method: "PATCH",
    headers: adminHeaders(token),
    body: JSON.stringify({ delta }),
  });
  if (!res.ok) throw new Error("Failed to update stock");
}

// Orders
export async function adminListOrders(token: string): Promise<Order[]> {
  const res = await fetch(`${ADMIN_BASE}/orders`, { headers: adminHeaders(token) });
  if (!res.ok) throw new Error("Failed to load orders");
  const data = await res.json();
  return data.orders || [];
}

export async function adminUpdateOrderStatus(token: string, id: string, status: string): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/orders/${id}/status`, {
    method: "PATCH",
    headers: adminHeaders(token),
    body: JSON.stringify({ status }),
  });
  if (!res.ok) throw new Error("Failed to update order status");
}

export async function adminCancelOrder(token: string, id: string): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/orders/${id}/cancel`, {
    method: "DELETE",
    headers: adminHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to cancel order");
}

export async function adminGetOrdersByStatus(token: string, status: string): Promise<Order[]> {
  const res = await fetch(`${ADMIN_BASE}/orders/by-status/${encodeURIComponent(status)}`, {
    headers: adminHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to load orders by status");
  const data = await res.json();
  return data.orders || [];
}

export async function adminGetOrderItems(token: string, id: string): Promise<OrderItem[]> {
  const res = await fetch(`${ADMIN_BASE}/orders/${id}/items`, {
    headers: adminHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to load order items");
  const data = await res.json();
  return data.items || [];
}

// Users
export async function adminListUsers(token: string): Promise<User[]> {
  const res = await fetch(`${ADMIN_BASE}/users`, { headers: adminHeaders(token) });
  if (!res.ok) throw new Error("Failed to load users");
  const data = await res.json();
  return data.users || [];
}

export async function adminSearchUsers(token: string, query: string): Promise<User[]> {
  const res = await fetch(`${ADMIN_BASE}/users/search?q=${encodeURIComponent(query)}`, {
    headers: adminHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to search users");
  const data = await res.json();
  return data.users || [];
}

export async function adminUpdateUserStatus(token: string, id: string, active: boolean): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/users/${id}/status`, {
    method: "PATCH",
    headers: adminHeaders(token),
    body: JSON.stringify({ active }),
  });
  if (!res.ok) throw new Error("Failed to update user status");
}

export async function adminResetPassword(token: string, id: string, new_password: string): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/users/${id}/password/reset`, {
    method: "PATCH",
    headers: adminHeaders(token),
    body: JSON.stringify({ new_password }),
  });
  if (!res.ok) throw new Error("Failed to reset password");
}

export async function adminDeleteUser(token: string, id: string): Promise<void> {
  const res = await fetch(`${ADMIN_BASE}/users/${id}`, {
    method: "DELETE",
    headers: adminHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to delete user");
}

// ── AI ────────────────────────────────────────────────────────────────────────
const AI_BASE = `${API_BASE}/ai`;

export interface AIMessage {
  role: "user" | "assistant";
  content: string;
}

export interface SneakerRecommendation {
  product_id: string;
  name: string;
  brand: string;
  price: number;
  image_url: string;
  reason: string;
  match_score: number;
}

export interface TrendingSneaker {
  product_id: string;
  name: string;
  brand: string;
  price: number;
  image_url: string;
  trend_reason: string;
  trend_score: number;
}

export async function aiChat(
  message: string,
  sessionId?: string,
  userId?: string
): Promise<{ reply: string; session_id: string }> {
  const res = await fetch(`${AI_BASE}/chat`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ message, session_id: sessionId || "", user_id: userId || "" }),
  });
  if (!res.ok) throw new Error("AI chat failed");
  return res.json();
}

export async function aiRecommend(
  preferences: string,
  budget?: number,
  size?: string,
  userId?: string
): Promise<{ recommendations: SneakerRecommendation[] }> {
  const res = await fetch(`${AI_BASE}/recommend`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      preferences,
      budget: budget || 0,
      size: size || "",
      user_id: userId || "",
    }),
  });
  if (!res.ok) throw new Error("AI recommend failed");
  return res.json();
}

export async function aiTrending(limit = 8): Promise<{ sneakers: TrendingSneaker[] }> {
  const res = await fetch(`${AI_BASE}/trending?limit=${limit}`);
  if (!res.ok) throw new Error("AI trending failed");
  return res.json();
}

export async function aiSearchByStyle(
  style_description: string,
  size = "",
  max_price = 0
): Promise<{ results: SneakerRecommendation[]; ai_summary: string }> {
  const res = await fetch(`${AI_BASE}/search-by-style`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ style_description, size, max_price }),
  });
  if (!res.ok) throw new Error("AI style search failed");
  return res.json();
}
