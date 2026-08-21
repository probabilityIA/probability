const d = require('./data');

const ok = (data, message) => ({ status: 200, payload: { success: true, data, message } });
const list = (items) => ({ items });

const routes = [];
const add = (method, pattern, fn) => routes.push({ method, pattern, fn });

function matchPath(pattern, path) {
  const p = pattern.split('/').filter(Boolean);
  const s = path.split('/').filter(Boolean);
  if (p.length !== s.length) return null;
  const params = {};
  for (let i = 0; i < p.length; i++) {
    if (p[i].startsWith(':')) params[p[i].slice(1)] = s[i];
    else if (p[i] !== s[i]) return null;
  }
  return params;
}

function resolve(method, path) {
  for (const r of routes) {
    if (r.method !== method) continue;
    const params = matchPath(r.pattern, path);
    if (params) {
      const bound = (ctx) => r.fn({ ...ctx, params });
      bound.params = params;
      return bound;
    }
  }
  return null;
}

add('POST', '/auth/login', ({ body }) => ok({
  user: {
    id: 1,
    name: 'Sebastian Camacho',
    email: body.email || 'demo@probability.co',
    phone: '3001234567',
    avatar_url: null,
    is_active: true,
    last_login_at: d.daysAgo(0),
  },
  token: 'mock-jwt-token',
  require_password_change: false,
  is_super_admin: false,
  scope: 'business',
  businesses: d.businesses,
}));

add('POST', '/auth/change-password', () => ok(null, 'Contrasena actualizada'));
add('POST', '/auth/generate-password', () => ({ status: 200, payload: { success: true, message: 'ok', password: 'Temporal123*' } }));

add('GET', '/auth/roles-permissions', () => ({
  status: 200,
  payload: {
    is_super: false,
    business_id: 1,
    business_name: 'Probability Demo',
    business_type_id: 1,
    business_type_name: 'Ecommerce',
    role: 'Administrador',
    subscription_status: 'active',
    resources: [
      'orders', 'products', 'customers', 'shipments', 'inventory', 'warehouses',
      'invoicing', 'wallet', 'routes', 'drivers', 'vehicles', 'integrations',
      'storefront', 'notifications', 'users', 'roles', 'permissions', 'tickets',
      'dashboard', 'pay',
    ].map((resource) => ({ resource, actions: ['read', 'create', 'update', 'delete'] })),
  },
}));

add('GET', '/businesses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.businesses, url) }));
add('GET', '/businesses/simple', () => ok(d.businesses.map((b) => ({ id: b.id, name: b.name, logo_url: b.logo_url }))));
add('GET', '/businesses/:id', ({ params }) => ok(d.businesses.find((b) => b.id === Number(params.id)) || d.businesses[0]));

add('GET', '/dashboard/stats', () => {
  const byKey = (rows, key) => {
    const map = {};
    rows.forEach((r) => { map[r[key]] = (map[r[key]] || 0) + 1; });
    return Object.entries(map).map(([k, v]) => [k, v]);
  };
  return ok({
    total_orders: d.orders.length,
    orders_by_integration_type: byKey(d.orders, 'channel').map(([k, v]) => ({ integration_type: k, count: v })),
    top_customers: d.customers.slice(0, 5).map((c) => ({
      customer_name: c.name, customer_email: c.email, order_count: c.total_orders,
    })),
    orders_by_location: byKey(d.orders, 'shipping_city').map(([k, v]) => ({ city: k, state: 'Colombia', order_count: v })),
    top_drivers: d.drivers.slice(0, 5).map((x, i) => ({ driver_name: x.name, driver_id: x.id, order_count: 40 - i * 6 })),
    drivers_by_location: d.drivers.slice(0, 5).map((x, i) => ({ driver_name: x.name, city: d.customers[i].city, state: 'Colombia', order_count: 20 - i * 3 })),
    top_products: d.products.slice(0, 5).map((p, i) => ({
      product_name: p.name, product_id: String(p.id), sku: p.sku,
      order_count: 34 - i * 5, total_sold: 120 - i * 17,
    })),
    products_by_category: byKey(d.products, 'category').map(([k, v]) => ({ category: k, count: v })),
    products_by_brand: [
      { brand: 'Genius', count: 12 }, { brand: 'Logitech', count: 9 },
      { brand: 'Xiaomi', count: 8 }, { brand: 'Sin marca', count: 13 },
    ],
    shipments_by_status: byKey(d.shipments, 'status').map(([k, v]) => ({ status: k, count: v })),
    shipments_by_carrier: byKey(d.shipments, 'carrier').map(([k, v]) => ({ carrier: k, count: v })),
    shipments_by_warehouse: d.warehouses.map((w, i) => ({ warehouse_name: w.name, warehouse_id: w.id, count: 18 - i * 4 })),
    orders_by_business: d.businesses.map((b, i) => ({ business_id: b.id, business_name: b.name, order_count: 48 - i * 20 })),
    total_sales: d.orders.reduce((a, o) => a + o.total, 0),
    total_customers: d.customers.length,
    total_products: d.products.length,
    wallet_balance: 2480500,
    sales_by_day: Array.from({ length: 14 }, (_, i) => ({
      date: d.daysAgo(13 - i).slice(0, 10), total: 400000 + i * 63000, orders: 3 + (i % 7),
    })),
  });
});

add('GET', '/orders', ({ url, paginate }) => {
  let rows = d.orders;
  const status = url.searchParams.get('status_code') || url.searchParams.get('status');
  if (status && status !== 'all') rows = rows.filter((o) => o.status_code === status);
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/orders/:id', ({ params }) => ok(d.orders.find((o) => o.id === Number(params.id))));
add('GET', '/orders/:id/raw', ({ params }) => ok({ raw: d.orders.find((o) => o.id === Number(params.id)) }));
add('PUT', '/orders/:id', ({ params, body }) => ok({ ...d.orders.find((o) => o.id === Number(params.id)), ...body }));
add('POST', '/orders/:id/cancel', () => ok(null, 'Orden cancelada'));

add('GET', '/order-statuses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.ORDER_STATUSES, url) }));
add('GET', '/payment-statuses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.PAYMENT_STATUSES, url) }));
add('GET', '/fulfillment-statuses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.FULFILLMENT_STATUSES, url) }));

add('GET', '/products', ({ url, paginate }) => ({ status: 200, payload: paginate(d.products, url) }));
add('GET', '/products/:id', ({ params }) => ok(d.products.find((p) => p.id === Number(params.id))));

add('GET', '/customers', ({ url, paginate }) => ({ status: 200, payload: paginate(d.customers, url) }));
add('GET', '/customers/:id', ({ params }) => ok(d.customers.find((c) => c.id === Number(params.id))));

add('GET', '/shipments', ({ url, paginate }) => ({ status: 200, payload: paginate(d.shipments, url) }));
add('GET', '/shipments/:id', ({ params }) => ok(d.shipments.find((s) => s.id === Number(params.id))));
add('POST', '/shipments/:id/cancel', () => ok(null, 'Guia cancelada'));
add('GET', '/shipments/origin-addresses', () => ok(d.warehouses.map((w) => ({ id: w.id, name: w.name, address: w.address, city: w.city }))));

add('GET', '/invoices', ({ url, paginate }) => ({ status: 200, payload: paginate(d.invoices, url) }));
add('GET', '/invoices/:id', ({ params }) => ok(d.invoices.find((i) => i.id === Number(params.id))));

add('GET', '/pay/wallet/balance', () => ok({ balance: 2480500, currency: 'COP', updated_at: d.daysAgo(0) }));
add('GET', '/pay/wallet/history', ({ url, paginate }) => ({ status: 200, payload: paginate(d.walletMovements, url) }));

add('GET', '/inventory/movements', ({ url, paginate }) => ({ status: 200, payload: paginate(d.inventoryMovements, url) }));
add('GET', '/inventory/movement-types', () => ok([
  { code: 'entry', name: 'Entrada' },
  { code: 'exit', name: 'Salida' },
  { code: 'adjustment', name: 'Ajuste' },
  { code: 'transfer', name: 'Traslado' },
]));
add('GET', '/inventory/warehouse/:id', ({ url, paginate, params }) => {
  const rows = d.products.map((p) => ({
    product_id: p.id, sku: p.sku, name: p.name, stock: p.stock,
    warehouse_id: Number(params.id), reserved: 0, available: p.stock,
  }));
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/inventory/product/:id', ({ params }) => ok(
  d.warehouses.map((w) => ({ warehouse_id: w.id, warehouse_name: w.name, stock: 10 + w.id * 7, product_id: Number(params.id) })),
));

add('GET', '/warehouses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.warehouses, url) }));
add('GET', '/warehouses/:id', ({ params }) => ok(d.warehouses.find((w) => w.id === Number(params.id))));

add('GET', '/routes', ({ url, paginate }) => ({ status: 200, payload: paginate(d.deliveryRoutes, url) }));
add('GET', '/routes/:id', ({ params }) => ok(d.deliveryRoutes.find((r) => r.id === Number(params.id))));
add('GET', '/routes/available-drivers', () => ok(d.drivers.filter((x) => x.is_available)));
add('GET', '/routes/available-vehicles', () => ok(d.vehicles.filter((x) => x.is_available)));
add('GET', '/routes/assignable-orders', ({ url, paginate }) => ({ status: 200, payload: paginate(d.orders.slice(0, 18), url) }));

add('GET', '/drivers', ({ url, paginate }) => ({ status: 200, payload: paginate(d.drivers, url) }));
add('GET', '/vehicles', ({ url, paginate }) => ({ status: 200, payload: paginate(d.vehicles, url) }));

add('GET', '/integrations', ({ url, paginate }) => ({ status: 200, payload: paginate(d.integrations, url) }));
add('GET', '/integrations/:id', ({ params }) => ok(d.integrations.find((i) => i.id === Number(params.id))));
add('GET', '/integration-categories', () => ok([
  { code: 'ecommerce', name: 'Ecommerce' },
  { code: 'invoicing', name: 'Facturacion' },
  { code: 'messaging', name: 'Mensajeria' },
  { code: 'transport', name: 'Transporte' },
  { code: 'pay', name: 'Pagos' },
]));
add('GET', '/integration-types', ({ url, paginate }) => ({ status: 200, payload: paginate([
  { id: 1, code: 'shopify', name: 'Shopify', category: 'ecommerce', is_active: true },
  { id: 2, code: 'woocommerce', name: 'WooCommerce', category: 'ecommerce', is_active: true },
  { id: 3, code: 'mercadolibre', name: 'MercadoLibre', category: 'ecommerce', is_active: true },
  { id: 4, code: 'siigo', name: 'Siigo', category: 'invoicing', is_active: true },
  { id: 5, code: 'factus', name: 'Factus', category: 'invoicing', is_active: true },
  { id: 6, code: 'whatsapp', name: 'WhatsApp', category: 'messaging', is_active: true },
  { id: 7, code: 'envioclick', name: 'EnvioClick', category: 'transport', is_active: true },
  { id: 8, code: 'bold', name: 'Bold', category: 'pay', is_active: true },
], url) }));

add('GET', '/users', ({ url, paginate }) => ({ status: 200, payload: paginate(d.users, url) }));
add('GET', '/roles', ({ url, paginate }) => ({ status: 200, payload: paginate([
  { id: 1, name: 'Administrador', level: 1, scope: 'business', users_count: 3 },
  { id: 2, name: 'Operador', level: 2, scope: 'business', users_count: 5 },
  { id: 3, name: 'Bodega', level: 3, scope: 'business', users_count: 2 },
  { id: 4, name: 'Contabilidad', level: 3, scope: 'business', users_count: 2 },
], url) }));
add('GET', '/permissions', ({ url, paginate }) => ({ status: 200, payload: paginate(
  ['orders', 'products', 'customers', 'shipments', 'inventory'].flatMap((r, i) =>
    ['read', 'create', 'update', 'delete'].map((a, j) => ({ id: i * 4 + j + 1, resource: r, action: a, name: `${r}.${a}` }))),
  url,
) }));
add('GET', '/resources', ({ url, paginate }) => ({ status: 200, payload: paginate(
  ['orders', 'products', 'customers', 'shipments', 'inventory', 'warehouses', 'invoicing', 'wallet']
    .map((r, i) => ({ id: i + 1, code: r, name: r })),
  url,
) }));

add('GET', '/notification-configs', ({ url, paginate }) => ({ status: 200, payload: paginate([
  { id: 1, event_code: 'order.created', event_name: 'Orden creada', channel: 'whatsapp', is_active: true, template: 'orden_creada' },
  { id: 2, event_code: 'order.shipped', event_name: 'Orden enviada', channel: 'whatsapp', is_active: true, template: 'orden_enviada' },
  { id: 3, event_code: 'order.delivered', event_name: 'Orden entregada', channel: 'whatsapp', is_active: false, template: 'orden_entregada' },
  { id: 4, event_code: 'wallet.low_balance', event_name: 'Saldo bajo', channel: 'email', is_active: true, template: 'saldo_bajo' },
], url) }));

add('GET', '/tickets', ({ url, paginate }) => ({ status: 200, payload: paginate(d.tickets, url) }));
add('GET', '/tickets/:id', ({ params }) => ok(d.tickets.find((t) => t.id === Number(params.id))));

module.exports = { resolve, routes, ok, list };
