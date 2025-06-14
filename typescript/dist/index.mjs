var D = Object.defineProperty;
var U = (e, t, r) => t in e ? D(e, t, { enumerable: !0, configurable: !0, writable: !0, value: r }) : e[t] = r;
var l = (e, t, r) => U(e, typeof t != "symbol" ? t + "" : t, r);
import { Hono as O } from "hono";
import { cors as C } from "hono/cors";
function A() {
  const { process: e, Deno: t } = globalThis;
  return !(typeof (t == null ? void 0 : t.noColor) == "boolean" ? t.noColor : e !== void 0 ? "NO_COLOR" in (e == null ? void 0 : e.env) : !1);
}
var j = (e) => {
  const [t, r] = [",", "."];
  return e.map((o) => o.replace(/(\d)(?=(\d\d\d)+(?!\d))/g, "$1" + t)).join(r);
}, N = (e) => {
  const t = Date.now() - e;
  return j([t < 1e3 ? t + "ms" : Math.round(t / 1e3) + "s"]);
}, T = (e) => {
  if (A())
    switch (e / 100 | 0) {
      case 5:
        return `\x1B[31m${e}\x1B[0m`;
      case 4:
        return `\x1B[33m${e}\x1B[0m`;
      case 3:
        return `\x1B[36m${e}\x1B[0m`;
      case 2:
        return `\x1B[32m${e}\x1B[0m`;
    }
  return `${e}`;
};
function q(e, t, r, s, o = 0, a) {
  const d = t === "<--" ? `${t} ${r} ${s}` : `${t} ${r} ${s} ${T(o)} ${a}`;
  e(d);
}
var P = (e = console.log) => async function(r, s) {
  const { method: o, url: a } = r.req, d = a.slice(a.indexOf("/", 8));
  q(e, "<--", o, d);
  const i = Date.now();
  await s(), q(e, "-->", o, d, r.res.status, N(i));
};
class x {
  constructor() {
    l(this, "products", /* @__PURE__ */ new Map());
    l(this, "categories", /* @__PURE__ */ new Map());
    l(this, "users", /* @__PURE__ */ new Map());
    l(this, "carts", /* @__PURE__ */ new Map());
    l(this, "orders", /* @__PURE__ */ new Map());
    this.initializeMockData();
  }
  initializeMockData() {
    this.categories.set("1", {
      id: "1",
      name: "Electronics",
      parentId: void 0,
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.categories.set("2", {
      id: "2",
      name: "Laptops",
      parentId: "1",
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.categories.set("3", {
      id: "3",
      name: "Smartphones",
      parentId: "1",
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.categories.set("4", {
      id: "4",
      name: "Clothing",
      parentId: void 0,
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.products.set("1", {
      id: "1",
      name: 'MacBook Pro 16"',
      description: "Apple MacBook Pro with M3 chip",
      price: 2499.99,
      stock: 10,
      categoryId: "2",
      imageUrls: ["https://example.com/macbook.jpg"],
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.products.set("2", {
      id: "2",
      name: "iPhone 15 Pro",
      description: "Latest iPhone with titanium design",
      price: 999.99,
      stock: 25,
      categoryId: "3",
      imageUrls: ["https://example.com/iphone.jpg"],
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.products.set("3", {
      id: "3",
      name: "T-Shirt",
      description: "Comfortable cotton t-shirt",
      price: 29.99,
      stock: 100,
      categoryId: "4",
      imageUrls: ["https://example.com/tshirt.jpg"],
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.users.set("1", {
      id: "1",
      email: "user1@example.com",
      name: "Test User 1",
      address: {
        street: "123 Test St",
        city: "Test City",
        state: "TC",
        postalCode: "12345",
        country: "USA"
      },
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.users.set("2", {
      id: "2",
      email: "user2@example.com",
      name: "Test User 2",
      address: {
        street: "456 Demo Ave",
        city: "Demo City",
        state: "DC",
        postalCode: "67890",
        country: "USA"
      },
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.carts.set("1", {
      id: "cart-1",
      userId: "1",
      items: [],
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    }), this.carts.set("2", {
      id: "cart-2",
      userId: "2",
      items: [],
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    });
  }
  // Products
  getProducts() {
    return Array.from(this.products.values());
  }
  getProduct(t) {
    return this.products.get(t);
  }
  createProduct(t) {
    return this.products.set(t.id, t), t;
  }
  updateProduct(t, r) {
    return this.products.set(t, r), r;
  }
  deleteProduct(t) {
    const r = this.products.get(t);
    return this.products.delete(t), r;
  }
  // Categories
  getCategories() {
    return Array.from(this.categories.values());
  }
  getCategory(t) {
    return this.categories.get(t);
  }
  createCategory(t) {
    return this.categories.set(t.id, t), t;
  }
  updateCategory(t, r) {
    return this.categories.set(t, r), r;
  }
  deleteCategory(t) {
    const r = this.categories.get(t);
    return this.categories.delete(t), r;
  }
  // Users
  getUsers() {
    return Array.from(this.users.values());
  }
  getUser(t) {
    return this.users.get(t);
  }
  createUser(t) {
    return this.users.set(t.id, t), t;
  }
  updateUser(t, r) {
    return this.users.set(t, r), r;
  }
  deleteUser(t) {
    const r = this.users.get(t);
    return this.users.delete(t), r;
  }
  // Carts
  getCartByUserId(t) {
    return this.carts.get(t) || {
      id: `cart-${t}`,
      userId: t,
      items: [],
      createdAt: (/* @__PURE__ */ new Date()).toISOString(),
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    };
  }
  updateCart(t, r) {
    return this.carts.set(t, r), r;
  }
  // Orders
  getOrders() {
    return Array.from(this.orders.values());
  }
  getOrder(t) {
    return this.orders.get(t);
  }
  getOrdersByUserId(t) {
    return Array.from(this.orders.values()).filter((r) => r.userId === t);
  }
  createOrder(t) {
    return this.orders.set(t.id, t), t;
  }
  updateOrder(t, r) {
    return this.orders.set(t, r), r;
  }
}
const n = new x(), y = new O();
y.get("/", (e) => {
  const t = parseInt(e.req.query("limit") || "10"), r = parseInt(e.req.query("offset") || "0"), s = e.req.query("search"), o = e.req.query("categoryId"), a = e.req.query("minPrice") ? parseFloat(e.req.query("minPrice")) : void 0, d = e.req.query("maxPrice") ? parseFloat(e.req.query("maxPrice")) : void 0;
  let i = n.getProducts();
  s && (i = i.filter(
    (u) => u.name.toLowerCase().includes(s.toLowerCase()) || u.description.toLowerCase().includes(s.toLowerCase())
  )), o && (i = i.filter((u) => u.categoryId === o)), a !== void 0 && (i = i.filter((u) => u.price >= a)), d !== void 0 && (i = i.filter((u) => u.price <= d));
  const c = {
    items: i.slice(r, r + t),
    total: i.length,
    limit: t,
    offset: r
  };
  return e.json(c);
});
y.get("/:productId", (e) => {
  const t = e.req.param("productId"), r = n.getProduct(t);
  return r ? e.json(r) : e.json({ error: { code: "NOT_FOUND", message: "Product not found" } }, 404);
});
y.post("/", async (e) => {
  const t = await e.req.json(), r = {
    id: Date.now().toString(),
    ...t,
    imageUrls: t.imageUrls || [],
    createdAt: (/* @__PURE__ */ new Date()).toISOString(),
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }, s = n.createProduct(r);
  return e.json(s, 201);
});
y.put("/:productId", async (e) => {
  const t = e.req.param("productId"), r = await e.req.json(), s = n.getProduct(t);
  if (!s)
    return e.json({ error: { code: "NOT_FOUND", message: "Product not found" } }, 404);
  const o = {
    ...s,
    ...r,
    id: t,
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }, a = n.updateProduct(t, o);
  return e.json(a);
});
y.delete("/:productId", (e) => {
  const t = e.req.param("productId");
  return n.deleteProduct(t) ? e.body(null, 204) : e.json({ error: { code: "NOT_FOUND", message: "Product not found" } }, 404);
});
const m = new O();
m.get("/", (e) => {
  const t = n.getCategories();
  return e.json(t);
});
m.get("/tree", (e) => {
  const t = n.getCategories(), r = /* @__PURE__ */ new Map(), s = [];
  return t.forEach((o) => {
    r.set(o.id, { ...o, children: [] });
  }), r.forEach((o) => {
    if (!o.parentId)
      s.push(o);
    else {
      const a = r.get(o.parentId);
      a && a.children.push(o);
    }
  }), e.json(s);
});
m.get("/:categoryId", (e) => {
  const t = e.req.param("categoryId"), r = n.getCategory(t);
  return r ? e.json(r) : e.json({ error: { code: "NOT_FOUND", message: "Category not found" } }, 404);
});
m.post("/", async (e) => {
  const t = await e.req.json(), r = {
    id: Date.now().toString(),
    ...t,
    createdAt: (/* @__PURE__ */ new Date()).toISOString(),
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }, s = n.createCategory(r);
  return e.json(s, 201);
});
m.put("/:categoryId", async (e) => {
  const t = e.req.param("categoryId"), r = await e.req.json(), s = n.getCategory(t);
  if (!s)
    return e.json({ error: { code: "NOT_FOUND", message: "Category not found" } }, 404);
  const o = {
    ...s,
    ...r,
    id: t
  }, a = n.updateCategory(t, o);
  return e.json(a);
});
m.delete("/:categoryId", (e) => {
  const t = e.req.param("categoryId");
  return n.deleteCategory(t) ? e.body(null, 204) : e.json({ error: { code: "NOT_FOUND", message: "Category not found" } }, 404);
});
const S = new O();
S.get("/", (e) => {
  const t = parseInt(e.req.query("limit") || "20"), r = parseInt(e.req.query("offset") || "0"), s = n.getUsers(), a = {
    items: s.slice(r, r + t),
    total: s.length,
    limit: t,
    offset: r
  };
  return e.json(a);
});
S.get("/:userId", (e) => {
  const t = e.req.param("userId"), r = n.getUser(t);
  return r ? e.json(r) : e.json({ error: { code: "NOT_FOUND", message: "User not found" } }, 404);
});
S.post("/", async (e) => {
  const t = await e.req.json(), r = {
    id: Date.now().toString(),
    ...t,
    createdAt: (/* @__PURE__ */ new Date()).toISOString(),
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }, s = n.createUser(r);
  return n.updateCart(s.id, {
    id: `cart-${s.id}`,
    userId: s.id,
    items: [],
    createdAt: (/* @__PURE__ */ new Date()).toISOString(),
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }), e.json(s, 201);
});
S.put("/:userId", async (e) => {
  const t = e.req.param("userId"), r = await e.req.json(), s = n.getUser(t);
  if (!s)
    return e.json({ error: { code: "NOT_FOUND", message: "User not found" } }, 404);
  const o = {
    ...s,
    ...r,
    id: t,
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }, a = n.updateUser(t, o);
  return e.json(a);
});
S.delete("/:userId", (e) => {
  const t = e.req.param("userId");
  return n.deleteUser(t) ? e.body(null, 204) : e.json({ error: { code: "NOT_FOUND", message: "User not found" } }, 404);
});
const f = new O();
f.get("/users/:userId", (e) => {
  const t = e.req.param("userId"), r = n.getCartByUserId(t);
  return e.json(r);
});
f.post("/users/:userId/items", async (e) => {
  const t = e.req.param("userId"), r = await e.req.json(), s = n.getCartByUserId(t), o = n.getProduct(r.productId);
  if (!o)
    return e.json({ error: { code: "NOT_FOUND", message: "Product not found" } }, 404);
  if (o.stock < r.quantity)
    return e.json({ error: { code: "INSUFFICIENT_STOCK", message: "Insufficient stock" } }, 400);
  const a = s.items.findIndex((i) => i.productId === r.productId);
  a >= 0 ? s.items[a].quantity += r.quantity : s.items.push({
    productId: r.productId,
    quantity: r.quantity
  }), s.updatedAt = (/* @__PURE__ */ new Date()).toISOString();
  const d = n.updateCart(t, s);
  return e.json(d);
});
f.patch("/users/:userId/items/:productId", async (e) => {
  const t = e.req.param("userId"), r = e.req.param("productId"), s = await e.req.json(), o = n.getCartByUserId(t), a = n.getProduct(r);
  if (!a)
    return e.json({ error: { code: "NOT_FOUND", message: "Product not found" } }, 404);
  const d = o.items.findIndex((g) => g.productId === r);
  if (d < 0)
    return e.json({ error: { code: "NOT_FOUND", message: "Item not found in cart" } }, 404);
  if (a.stock < s.quantity)
    return e.json({ error: { code: "INSUFFICIENT_STOCK", message: "Insufficient stock" } }, 400);
  o.items[d].quantity = s.quantity, o.updatedAt = (/* @__PURE__ */ new Date()).toISOString();
  const i = n.updateCart(t, o);
  return e.json(i);
});
f.delete("/users/:userId/items/:productId", (e) => {
  const t = e.req.param("userId"), r = e.req.param("productId"), s = n.getCartByUserId(t), o = s.items.findIndex((a) => a.productId === r);
  return o < 0 ? e.json({ error: { code: "NOT_FOUND", message: "Item not found in cart" } }, 404) : (s.items.splice(o, 1), s.updatedAt = (/* @__PURE__ */ new Date()).toISOString(), n.updateCart(t, s), e.body(null, 204));
});
f.delete("/users/:userId/items", (e) => {
  const t = e.req.param("userId"), r = n.getCartByUserId(t);
  return r.items = [], r.updatedAt = (/* @__PURE__ */ new Date()).toISOString(), n.updateCart(t, r), e.body(null, 204);
});
const w = new O();
w.get("/", (e) => {
  const t = parseInt(e.req.query("limit") || "10"), r = parseInt(e.req.query("offset") || "0"), s = e.req.query("userId"), o = e.req.query("status");
  let a = s ? n.getOrdersByUserId(s) : n.getOrders();
  o && (a = a.filter((g) => g.status === o)), a.sort((g, c) => new Date(c.createdAt).getTime() - new Date(g.createdAt).getTime());
  const i = {
    items: a.slice(r, r + t),
    total: a.length,
    limit: t,
    offset: r
  };
  return e.json(i);
});
w.get("/:orderId", (e) => {
  const t = e.req.param("orderId"), r = n.getOrder(t);
  return r ? e.json(r) : e.json({ error: { code: "NOT_FOUND", message: "Order not found" } }, 404);
});
w.get("/users/:userId", (e) => {
  const t = e.req.param("userId"), r = parseInt(e.req.query("limit") || "10"), s = parseInt(e.req.query("offset") || "0"), o = e.req.query("status");
  if (!n.getUser(t))
    return e.json({ error: { code: "NOT_FOUND", message: "User not found" } }, 404);
  let d = n.getOrders().filter((c) => c.userId === t);
  o && (d = d.filter((c) => c.status === o));
  const g = {
    items: d.slice(s, s + r),
    total: d.length,
    limit: r,
    offset: s
  };
  return e.json(g);
});
w.post("/users/:userId", async (e) => {
  const t = e.req.param("userId"), r = await e.req.json();
  if (!n.getUser(t))
    return e.json({ error: { code: "NOT_FOUND", message: "User not found" } }, 404);
  const o = n.getCartByUserId(t);
  if (o.items.length === 0)
    return e.json({ error: { code: "EMPTY_CART", message: "Cart is empty" } }, 400);
  let a = 0;
  const d = [];
  for (const u of o.items) {
    const p = n.getProduct(u.productId);
    if (!p)
      return e.json({ error: { code: "NOT_FOUND", message: `Product ${u.productId} not found` } }, 404);
    if (p.stock < u.quantity)
      return e.json({ error: { code: "INSUFFICIENT_STOCK", message: `Insufficient stock for product ${p.name}` } }, 400);
    const h = p.price;
    a += h * u.quantity, d.push({
      productId: u.productId,
      quantity: u.quantity,
      price: h,
      productName: p.name
    }), n.updateProduct(p.id, {
      ...p,
      stock: p.stock - u.quantity,
      updatedAt: (/* @__PURE__ */ new Date()).toISOString()
    });
  }
  const i = {
    id: Date.now().toString(),
    userId: t,
    items: d,
    totalAmount: a,
    status: "pending",
    shippingAddress: r.shippingAddress,
    createdAt: (/* @__PURE__ */ new Date()).toISOString(),
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }, g = n.createOrder(i), c = n.getCartByUserId(t);
  return c.items = [], c.updatedAt = (/* @__PURE__ */ new Date()).toISOString(), n.updateCart(t, c), e.json(g, 201);
});
w.patch("/status/:orderId", async (e) => {
  const t = e.req.param("orderId"), r = await e.req.json(), s = n.getOrder(t);
  if (!s)
    return e.json({ error: { code: "NOT_FOUND", message: "Order not found" } }, 404);
  if (!{
    pending: ["processing", "cancelled"],
    processing: ["shipped", "cancelled"],
    shipped: ["delivered", "cancelled"],
    delivered: [],
    cancelled: []
  }[s.status].includes(r.status))
    return e.json({
      error: {
        code: "INVALID_STATUS_TRANSITION",
        message: `Cannot transition from ${s.status} to ${r.status}`
      }
    }, 400);
  const a = {
    ...s,
    status: r.status,
    updatedAt: (/* @__PURE__ */ new Date()).toISOString()
  }, d = n.updateOrder(t, a);
  return e.json(d);
});
const I = new O();
I.use("*", P());
I.use("*", C());
I.get("/health", (e) => e.json({ status: "ok" }));
I.route("/products", y);
I.route("/categories", m);
I.route("/users", S);
I.route("/carts", f);
I.route("/orders", w);
export {
  I as default
};
