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
    avatar_url: 'avatars/1766115229_planeta-kaiosama.webp',
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

add('POST', '/auth/recovery-channels', ({ body }) => ok([
  { channel: 'email', masked: (body.email || 'demo@probability.co').replace(/^(.{2}).*(@.*)$/, '$1****$2'), available: true },
  { channel: 'whatsapp', masked: '+57 300 *** 4567', available: true },
]));

add('POST', '/auth/forgot-password', ({ body }) => ok(
  { channel: body.channel || 'email' },
  'Codigo enviado',
));

add('POST', '/auth/verify-otp', ({ body }) => {
  if ((body.code || '') !== '123456') {
    return { status: 400, payload: { success: false, error: 'Codigo invalido o vencido' } };
  }
  return ok({ token: 'mock-reset-token' }, 'Codigo verificado');
});

add('POST', '/auth/reset-password', () => ok(null, 'Contrasena restablecida'));

add('GET', '/auth/verify', () => ok({ valid: true }));

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
  const status = url.searchParams.get('status');
  if (status && status !== 'all') rows = rows.filter((o) => o.status === status);
  const orderNumber = url.searchParams.get('order_number');
  if (orderNumber) {
    const q = orderNumber.toLowerCase();
    rows = rows.filter((o) =>
      o.order_number.toLowerCase().includes(q) ||
      o.customer_name.toLowerCase().includes(q) ||
      (o.customer_email || '').toLowerCase().includes(q) ||
      (o.tracking_number || '').toLowerCase().includes(q));
  }
  const customerEmail = url.searchParams.get('customer_email');
  if (customerEmail) rows = rows.filter((o) => o.customer_email === customerEmail);
  const integrationType = url.searchParams.get('integration_type');
  if (integrationType) rows = rows.filter((o) => o.integration_type === integrationType);
  const isCod = url.searchParams.get('is_cod');
  if (isCod === 'true') rows = rows.filter((o) => o.is_cod);
  const isPaid = url.searchParams.get('is_paid');
  if (isPaid === 'true') rows = rows.filter((o) => o.is_paid);
  if (isPaid === 'false') rows = rows.filter((o) => !o.is_paid);
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/orders/:id', ({ params }) => ok(d.orders.find((o) => o.id === params.id)));
add('GET', '/orders/:id/raw', ({ params }) => ok({ raw: d.orders.find((o) => o.id === params.id) }));
add('PUT', '/orders/:id', ({ params, body }) => ok({ ...d.orders.find((o) => o.id === params.id), ...body }));
add('POST', '/orders/:id/cancel', () => ok(null, 'Orden cancelada'));

add('GET', '/mobile/orders/:id/full', ({ params }) => {
  const order = d.orders.find((o) => o.id === params.id);
  if (!order) return { status: 404, payload: { success: false, error: 'orden no encontrada' } };
  const shipment = d.shipments.find((s) => s.order_id === order.id) || null;
  const invoice = d.invoices.find((i) => i.order_id === order.id) || null;
  return ok({
    order: {
      id: order.id,
      order_number: order.order_number,
      internal_number: order.internal_number,
      integration_id: order.integration_id,
      integration_type: order.integration_type,
      platform: order.platform,
      status: order.status,
      status_id: order.status_id,
      is_paid: order.is_paid,
      is_cod: order.is_cod,
      cod_total: order.cod_total,
      subtotal: order.subtotal,
      tax: order.tax,
      discount: order.discount,
      shipping_cost: order.shipping_cost,
      total_amount: order.total_amount,
      currency: order.currency,
      customer_name: order.customer_name,
      customer_email: order.customer_email,
      customer_phone: order.customer_phone,
      shipping_street: order.shipping_street,
      shipping_city: order.shipping_city,
      warehouse_name: order.warehouse_name,
      user_name: order.user_name,
      created_at: order.created_at,
      updated_at: order.updated_at,
    },
    items: order.order_items.map((it, idx) => ({
      id: idx + 1,
      product_id: String(it.product_id),
      sku: it.sku,
      name: it.name,
      quantity: it.quantity,
      unit_price: it.unit_price,
      total_price: it.total_price,
    })),
    shipment: shipment && {
      id: shipment.id,
      tracking_number: shipment.tracking_number,
      tracking_url: shipment.tracking_url,
      carrier: shipment.carrier,
      guide_url: shipment.guide_url,
      status: shipment.status,
      carrier_status: shipment.carrier_status,
      destination_city: shipment.destination_city,
      insurance_cost: shipment.insurance_cost,
      total_cost: shipment.total_cost,
      carrier_cost: shipment.carrier_cost,
      applied_margin: shipment.applied_margin,
      cod_carrier_fee: shipment.cod_carrier_fee,
      cod_probability_margin: shipment.cod_probability_margin,
      created_at: shipment.created_at,
    },
    invoice: invoice && {
      id: invoice.id,
      invoice_number: invoice.invoice_number,
      status: invoice.status,
      total_amount: invoice.total,
      currency: 'COP',
      invoice_url: 'https://example.com/factura.pdf',
      cufe: invoice.cufe,
      issued_at: invoice.issued_at,
      created_at: invoice.created_at,
    },
  });
});

add('GET', '/order-statuses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.ORDER_STATUSES, url) }));
add('GET', '/payment-statuses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.PAYMENT_STATUSES, url) }));
add('GET', '/fulfillment-statuses', ({ url, paginate }) => ({ status: 200, payload: paginate(d.FULFILLMENT_STATUSES, url) }));

add('GET', '/products', ({ url, paginate }) => {
  let rows = d.products;
  const name = url.searchParams.get('name');
  if (name) {
    const q = name.toLowerCase();
    rows = rows.filter((p) => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q));
  }
  const sku = url.searchParams.get('sku');
  if (sku) rows = rows.filter((p) => p.sku.toLowerCase().includes(sku.toLowerCase()));
  const status = url.searchParams.get('status');
  if (status === 'in_stock') rows = rows.filter((p) => p.stock >= 20);
  if (status === 'low_stock') rows = rows.filter((p) => p.stock > 0 && p.stock < 20);
  if (status === 'out_of_stock') rows = rows.filter((p) => p.stock === 0);
  if (status === 'inactive') rows = rows.filter((p) => !p.is_active);
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/products/:id', ({ params }) => ok(d.products.find((p) => p.id === params.id)));
add('GET', '/products/:id/integrations', ({ params }) => {
  const product = d.products.find((p) => p.id === params.id);
  if (!product) return ok([]);
  return ok(product.channels.map((code, index) => {
    const type = d.INTEGRATION_TYPES.find((t) => t.code === code);
    return {
      id: index + 1,
      product_id: product.id,
      integration_id: type ? type.id : index + 1,
      integration_type: code,
      integration_name: type ? type.name : code,
      external_product_id: `${1000 + index}${product.sku}`,
      created_at: product.created_at,
      updated_at: product.updated_at,
    };
  }));
});

add('GET', '/customers', ({ url, paginate }) => {
  let rows = d.customers;
  const search = url.searchParams.get('search') || url.searchParams.get('name');
  if (search) {
    const q = search.toLowerCase();
    rows = rows.filter((c) =>
      c.name.toLowerCase().includes(q) ||
      (c.email || '').toLowerCase().includes(q) ||
      (c.phone || '').includes(q) ||
      (c.dni || '').includes(q));
  }
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/customers/:id', ({ params }) => ok(
  d.customers.find((c) => c.id === Number(params.id)),
));

add('GET', '/shipments', ({ url, paginate }) => {
  let rows = d.shipments;
  const status = url.searchParams.get('status');
  if (status && status !== 'all') rows = rows.filter((s) => s.status === status);
  const tracking = url.searchParams.get('tracking_number');
  if (tracking) {
    const q = tracking.toLowerCase();
    rows = rows.filter((s) =>
      (s.tracking_number || '').toLowerCase().includes(q) ||
      (s.order_number || '').toLowerCase().includes(q) ||
      (s.client_name || '').toLowerCase().includes(q));
  }
  const carrier = url.searchParams.get('carrier');
  if (carrier) rows = rows.filter((s) => s.carrier === carrier);
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/shipments/origin-addresses', () => ok(d.warehouses.map((w) => ({
  id: w.id, business_id: 1, alias: w.name, company: 'Probability Demo',
  first_name: 'Sebastian', last_name: 'Camacho', email: 'bodega@probability.co',
  phone: '3001234567', street: w.address, suburb: 'Centro', city_dane_code: '11001',
  city: w.city, state: 'Colombia', postal_code: '110111', is_default: w.is_default,
}))));
add('GET', '/shipments/track/:tracking', ({ params }) => {
  const shipment = d.shipments.find((s) => s.tracking_number === params.tracking);
  if (!shipment) return { status: 404, payload: { success: false, error: 'guia no encontrada' } };
  return ok({
    tracking_number: shipment.tracking_number,
    carrier: shipment.carrier,
    status: shipment.status,
    events: [
      { date: shipment.created_at, status: 'created', description: 'Guia generada' },
      { date: d.daysAgo(2), status: 'in_transit', description: 'En camino al destino' },
      ...(shipment.status === 'delivered'
        ? [{ date: shipment.delivered_at, status: 'delivered', description: 'Entregada al destinatario' }]
        : []),
    ],
  });
});
add('GET', '/shipments/:id', ({ params }) => ok(d.shipments.find((s) => s.id === Number(params.id))));
add('POST', '/shipments/:id/cancel', () => ok(null, 'Guia cancelada'));

add('GET', '/invoices', ({ url, paginate }) => ({ status: 200, payload: paginate(d.invoices, url) }));
add('GET', '/invoices/:id', ({ params }) => ok(d.invoices.find((i) => i.id === Number(params.id))));

add('GET', '/pay/wallet/balance', () => ok({ balance: 2480500, currency: 'COP', updated_at: d.daysAgo(0) }));
add('GET', '/pay/wallet/history', ({ url, paginate }) => ({ status: 200, payload: paginate(d.walletMovements, url) }));

add('GET', '/inventory/movements', ({ url, paginate }) => {
  let rows = d.inventoryMovements;
  const warehouseId = url.searchParams.get('warehouse_id');
  if (warehouseId) rows = rows.filter((m) => m.warehouse_id === Number(warehouseId));
  const typeCode = url.searchParams.get('movement_type_code');
  if (typeCode) rows = rows.filter((m) => m.movement_type_code === typeCode);
  const productId = url.searchParams.get('product_id');
  if (productId) rows = rows.filter((m) => m.product_id === productId);
  rows = [...rows].sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/inventory/movement-types', ({ url, paginate }) => ({
  status: 200,
  payload: paginate(d.MOVEMENT_TYPES.map((t) => ({ ...t, is_active: true })), url),
}));
add('GET', '/inventory/warehouse/:id', ({ url, paginate, params }) => {
  let rows = d.inventoryLevels.filter((l) => l.warehouse_id === Number(params.id));
  const lowStock = url.searchParams.get('low_stock');
  if (lowStock === 'true') rows = rows.filter((l) => l.available_qty <= (l.reorder_point || 0));
  return { status: 200, payload: paginate(rows, url) };
});
add('GET', '/inventory/product/:id', ({ params }) => ok(
  d.inventoryLevels.filter((l) => l.product_id === params.id),
));
add('POST', '/inventory/adjust', ({ body }) => {
  const level = d.inventoryLevels.find(
    (l) => l.product_id === body.product_id && l.warehouse_id === Number(body.warehouse_id),
  );
  const previous = level ? level.quantity : 0;
  const next = previous + Number(body.quantity || 0);
  if (level) {
    level.quantity = next;
    level.available_qty = next - level.reserved_qty;
  }
  return ok({
    id: d.inventoryMovements.length + 1,
    product_id: body.product_id,
    warehouse_id: Number(body.warehouse_id),
    movement_type_code: 'adjustment',
    movement_type_name: 'Ajuste por conteo',
    quantity: Number(body.quantity || 0),
    previous_qty: previous,
    new_qty: next,
    reason: body.reason,
    notes: body.notes || null,
    created_at: new Date().toISOString(),
  }, 'Ajuste registrado');
});
add('POST', '/inventory/transfer', () => ok(null, 'Traslado registrado'));

add('GET', '/warehouses', ({ url, paginate }) => ({
  status: 200,
  payload: paginate(d.warehouses, url),
}));
add('GET', '/warehouses/:id', ({ params }) => ok(
  d.warehouses.find((w) => w.id === Number(params.id)),
));
add('GET', '/warehouses/:id/locations', ({ params }) => ok(
  d.warehouseLocations.filter((l) => l.warehouse_id === Number(params.id)),
));

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
add('GET', '/integration-types', ({ url, paginate }) => ({
  status: 200,
  payload: paginate(d.INTEGRATION_TYPES.map((t) => ({ ...t, is_active: true })), url),
}));
add('GET', '/integration-types/active', () => ok(d.INTEGRATION_TYPES.map((t) => ({ ...t, is_active: true }))));
add('GET', '/integration-types/:id', ({ params }) => ok(
  d.INTEGRATION_TYPES.find((t) => t.id === Number(params.id)),
));

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
