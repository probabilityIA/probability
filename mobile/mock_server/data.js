function seedRandom(seed) {
  let s = seed;
  return () => {
    s = (s * 1103515245 + 12345) % 2147483648;
    return s / 2147483648;
  };
}

const rnd = seedRandom(20260820);
const pick = (arr) => arr[Math.floor(rnd() * arr.length)];
const between = (a, b) => a + Math.floor(rnd() * (b - a + 1));
const daysAgo = (d) => new Date(Date.now() - d * 86400000).toISOString();

const CITIES = ['Bogota', 'Medellin', 'Cali', 'Barranquilla', 'Bucaramanga', 'Pereira', 'Cartagena'];
const CHANNELS = ['Shopify', 'WooCommerce', 'MercadoLibre', 'WhatsApp', 'Amazon', 'Manual'];
const CARRIERS = ['Interrapidisimo', 'Coordinadora', 'Servientrega', 'Envia', 'TCC'];
const FIRST = ['Ana', 'Carlos', 'Laura', 'Diego', 'Sofia', 'Andres', 'Valentina', 'Julian', 'Camila', 'Mateo'];
const LAST = ['Gomez', 'Rodriguez', 'Martinez', 'Lopez', 'Perez', 'Ramirez', 'Torres', 'Vargas', 'Castro', 'Rojas'];
const PRODUCTS = [
  'Audifonos inalambricos', 'Reloj deportivo', 'Cargador rapido 65W', 'Mouse ergonomico',
  'Teclado mecanico', 'Camara web HD', 'Parlante bluetooth', 'Base para portatil',
  'Cable USB-C 2m', 'Power bank 20000mAh', 'Silla gamer', 'Monitor 24 pulgadas',
  'Disco SSD 1TB', 'Memoria RAM 16GB', 'Hub USB 7 puertos', 'Lampara de escritorio',
];

const ORDER_STATUSES = [
  { id: 1, code: 'pending', name: 'Pendiente' },
  { id: 2, code: 'confirmed', name: 'Confirmada' },
  { id: 3, code: 'processing', name: 'En proceso' },
  { id: 4, code: 'shipped', name: 'Enviada' },
  { id: 5, code: 'delivered', name: 'Entregada' },
  { id: 6, code: 'cancelled', name: 'Cancelada' },
];

const PAYMENT_STATUSES = [
  { id: 1, code: 'pending', name: 'Pendiente' },
  { id: 2, code: 'paid', name: 'Pagada' },
  { id: 3, code: 'refunded', name: 'Reembolsada' },
  { id: 4, code: 'failed', name: 'Fallida' },
];

const FULFILLMENT_STATUSES = [
  { id: 1, code: 'unfulfilled', name: 'Sin preparar' },
  { id: 2, code: 'partial', name: 'Parcial' },
  { id: 3, code: 'fulfilled', name: 'Preparada' },
];

const customers = Array.from({ length: 48 }, (_, i) => {
  const name = `${pick(FIRST)} ${pick(LAST)}`;
  return {
    id: i + 1,
    business_id: 1,
    name,
    email: `${name.toLowerCase().replace(/\s/g, '.')}${i}@correo.com`,
    phone: `30${between(0, 9)}${between(1000000, 9999999)}`,
    document_type: 'CC',
    document_number: `${between(10000000, 99999999)}`,
    city: pick(CITIES),
    address: `Calle ${between(1, 180)} # ${between(1, 90)}-${between(1, 90)}`,
    total_orders: between(1, 24),
    total_spent: between(80000, 4200000),
    is_active: rnd() > 0.1,
    created_at: daysAgo(between(5, 400)),
  };
});

const products = Array.from({ length: 42 }, (_, i) => {
  const name = `${PRODUCTS[i % PRODUCTS.length]} ${i > 15 ? `v${Math.floor(i / 16) + 1}` : ''}`.trim();
  const price = between(18000, 890000);
  return {
    id: i + 1,
    business_id: 1,
    sku: `SKU-${String(1000 + i)}`,
    name,
    description: `${name} de alta calidad, garantia de 12 meses.`,
    price,
    cost: Math.round(price * 0.62),
    stock: between(0, 320),
    track_inventory: true,
    is_active: rnd() > 0.08,
    category: pick(['Tecnologia', 'Accesorios', 'Oficina', 'Audio']),
    image_url: null,
    weight: between(100, 4500),
    created_at: daysAgo(between(10, 500)),
    variants_count: between(0, 4),
  };
});

const warehouses = [
  { id: 1, business_id: 1, name: 'Bodega principal', code: 'BOG-01', city: 'Bogota', address: 'Calle 13 # 68-50', is_active: true, is_default: true, occupancy: 68 },
  { id: 2, business_id: 1, name: 'Bodega norte', code: 'BOG-02', city: 'Bogota', address: 'Autopista Norte km 19', is_active: true, is_default: false, occupancy: 41 },
  { id: 3, business_id: 1, name: 'Bodega Medellin', code: 'MDE-01', city: 'Medellin', address: 'Cra 48 # 20-114', is_active: true, is_default: false, occupancy: 83 },
  { id: 4, business_id: 1, name: 'Punto de venta Cali', code: 'CLO-01', city: 'Cali', address: 'Av 6N # 25-30', is_active: false, is_default: false, occupancy: 12 },
];

const orders = Array.from({ length: 64 }, (_, i) => {
  const customer = customers[i % customers.length];
  const status = ORDER_STATUSES[between(0, ORDER_STATUSES.length - 1)];
  const payment = PAYMENT_STATUSES[between(0, PAYMENT_STATUSES.length - 1)];
  const fulfillment = FULFILLMENT_STATUSES[between(0, FULFILLMENT_STATUSES.length - 1)];
  const itemCount = between(1, 4);
  const items = Array.from({ length: itemCount }, (_, j) => {
    const product = products[(i + j) % products.length];
    const quantity = between(1, 3);
    return {
      id: i * 10 + j + 1,
      product_id: product.id,
      sku: product.sku,
      name: product.name,
      quantity,
      unit_price: product.price,
      total: product.price * quantity,
      image_url: null,
    };
  });
  const subtotal = items.reduce((acc, it) => acc + it.total, 0);
  const shipping = between(8000, 24000);
  const isCod = rnd() > 0.55;
  return {
    id: i + 1,
    business_id: 1,
    order_number: `ORD-${10240 + i}`,
    external_id: `${between(100000, 999999)}`,
    channel: pick(CHANNELS),
    customer_id: customer.id,
    customer_name: customer.name,
    customer_email: customer.email,
    customer_phone: customer.phone,
    shipping_city: customer.city,
    shipping_address: customer.address,
    status_id: status.id,
    status_code: status.code,
    status_name: status.name,
    payment_status_code: payment.code,
    payment_status_name: payment.name,
    fulfillment_status_code: fulfillment.code,
    fulfillment_status_name: fulfillment.name,
    items,
    items_count: itemCount,
    subtotal,
    discount: 0,
    shipping_cost: shipping,
    tax: Math.round(subtotal * 0.19),
    total: subtotal + shipping,
    is_cod: isCod,
    cod_total: isCod ? subtotal + shipping : 0,
    currency: 'COP',
    has_shipment: status.id >= 4,
    has_invoice: payment.code === 'paid',
    created_at: daysAgo(between(0, 60)),
    updated_at: daysAgo(between(0, 5)),
  };
});

const shipments = orders
  .filter((o) => o.has_shipment)
  .map((o, i) => ({
    id: i + 1,
    business_id: 1,
    order_id: o.id,
    order_number: o.order_number,
    guide_number: `GU${between(100000000, 999999999)}`,
    carrier: pick(CARRIERS),
    status: pick(['created', 'in_transit', 'delivered', 'returned']),
    destination_city: o.shipping_city,
    recipient: o.customer_name,
    total_cost: o.shipping_cost,
    carrier_cost: Math.round(o.shipping_cost * 0.85),
    applied_margin: Math.round(o.shipping_cost * 0.15),
    is_cod: o.is_cod,
    cod_total: o.cod_total,
    cod_carrier_fee: o.is_cod ? Math.round(o.cod_total * 0.03) : 0,
    guide_url: 'https://example.com/guia.pdf',
    created_at: o.created_at,
  }));

const invoices = orders
  .filter((o) => o.has_invoice)
  .map((o, i) => ({
    id: i + 1,
    business_id: 1,
    order_id: o.id,
    order_number: o.order_number,
    invoice_number: `FE-${2400 + i}`,
    cufe: `${between(10000000, 99999999)}abcdef`,
    provider: pick(['Siigo', 'Factus', 'Alegra']),
    status: pick(['issued', 'pending', 'failed', 'cancelled']),
    customer_name: o.customer_name,
    subtotal: o.subtotal,
    tax: o.tax,
    total: o.total,
    issued_at: o.created_at,
    created_at: o.created_at,
  }));

const walletMovements = Array.from({ length: 30 }, (_, i) => {
  const isCredit = rnd() > 0.55;
  const amount = between(15000, 480000);
  return {
    id: i + 1,
    business_id: 1,
    type: isCredit ? 'credit' : 'debit',
    concept: isCredit ? pick(['Recarga Bold', 'Ajuste manual', 'Recaudo contra entrega']) : pick(['Cobro de guia', 'Comision transportadora', 'Factura electronica']),
    amount,
    balance_after: between(120000, 3200000),
    reference: `MOV-${5000 + i}`,
    created_at: daysAgo(between(0, 45)),
  };
});

const inventoryMovements = Array.from({ length: 36 }, (_, i) => {
  const product = products[i % products.length];
  return {
    id: i + 1,
    product_id: product.id,
    product_name: product.name,
    sku: product.sku,
    warehouse_id: warehouses[i % 3].id,
    warehouse_name: warehouses[i % 3].name,
    type: pick(['entry', 'exit', 'adjustment', 'transfer']),
    quantity: between(1, 60),
    stock_after: between(0, 300),
    reason: pick(['Compra a proveedor', 'Venta', 'Ajuste por conteo', 'Traslado entre bodegas', 'Devolucion']),
    user_name: pick(FIRST),
    created_at: daysAgo(between(0, 30)),
  };
});

const drivers = Array.from({ length: 8 }, (_, i) => ({
  id: i + 1,
  business_id: 1,
  name: `${pick(FIRST)} ${pick(LAST)}`,
  document_number: `${between(10000000, 99999999)}`,
  phone: `31${between(0, 9)}${between(1000000, 9999999)}`,
  license_number: `LIC-${between(10000, 99999)}`,
  is_active: rnd() > 0.15,
  is_available: rnd() > 0.4,
  vehicle_id: (i % 5) + 1,
  created_at: daysAgo(between(30, 400)),
}));

const vehicles = Array.from({ length: 5 }, (_, i) => ({
  id: i + 1,
  business_id: 1,
  plate: `${['ABC', 'XYZ', 'JKL', 'PQR', 'MNO'][i]}${between(100, 999)}`,
  brand: pick(['Chevrolet', 'Renault', 'Nissan', 'Yamaha']),
  model: `${between(2016, 2025)}`,
  type: pick(['van', 'motorcycle', 'truck']),
  capacity_kg: between(150, 3200),
  is_active: true,
  is_available: rnd() > 0.3,
}));

const deliveryRoutes = Array.from({ length: 10 }, (_, i) => ({
  id: i + 1,
  business_id: 1,
  name: `Ruta ${pick(CITIES)} ${i + 1}`,
  code: `RT-${300 + i}`,
  status: pick(['draft', 'in_progress', 'completed']),
  driver_id: (i % drivers.length) + 1,
  driver_name: drivers[i % drivers.length].name,
  vehicle_id: (i % vehicles.length) + 1,
  vehicle_plate: vehicles[i % vehicles.length].plate,
  stops_count: between(3, 14),
  completed_stops: between(0, 3),
  total_distance_km: between(8, 120),
  scheduled_date: daysAgo(between(-3, 10)),
  created_at: daysAgo(between(0, 20)),
}));

const integrations = [
  { id: 1, business_id: 1, name: 'Tienda Shopify', type: 'shopify', category: 'ecommerce', is_active: true, is_default: true, status: 'connected', last_sync_at: daysAgo(0), orders_synced: 1842 },
  { id: 2, business_id: 1, name: 'WooCommerce principal', type: 'woocommerce', category: 'ecommerce', is_active: true, is_default: false, status: 'connected', last_sync_at: daysAgo(0), orders_synced: 634 },
  { id: 3, business_id: 1, name: 'MercadoLibre CO', type: 'mercadolibre', category: 'ecommerce', is_active: true, is_default: false, status: 'error', last_sync_at: daysAgo(2), orders_synced: 221 },
  { id: 4, business_id: 1, name: 'Siigo contabilidad', type: 'siigo', category: 'invoicing', is_active: true, is_default: true, status: 'connected', last_sync_at: daysAgo(1), orders_synced: 0 },
  { id: 5, business_id: 1, name: 'WhatsApp Business', type: 'whatsapp', category: 'messaging', is_active: true, is_default: true, status: 'connected', last_sync_at: daysAgo(0), orders_synced: 0 },
  { id: 6, business_id: 1, name: 'EnvioClick', type: 'envioclick', category: 'transport', is_active: true, is_default: true, status: 'connected', last_sync_at: daysAgo(0), orders_synced: 0 },
  { id: 7, business_id: 1, name: 'Bold pagos', type: 'bold', category: 'pay', is_active: false, is_default: false, status: 'disconnected', last_sync_at: null, orders_synced: 0 },
];

const tickets = Array.from({ length: 14 }, (_, i) => ({
  id: i + 1,
  business_id: 1,
  code: `TK-${900 + i}`,
  subject: pick(['Guia no genera', 'Producto sin stock', 'Error al facturar', 'Cliente no recibio pedido', 'Duda con la billetera']),
  status: pick(['open', 'in_progress', 'closed']),
  priority: pick(['low', 'medium', 'high']),
  requester: pick(FIRST),
  assigned_to: pick(FIRST),
  comments_count: between(0, 9),
  created_at: daysAgo(between(0, 40)),
  updated_at: daysAgo(between(0, 6)),
}));

const users = Array.from({ length: 12 }, (_, i) => ({
  id: i + 1,
  business_id: 1,
  name: `${pick(FIRST)} ${pick(LAST)}`,
  email: `usuario${i + 1}@probability.co`,
  phone: `30${between(0, 9)}${between(1000000, 9999999)}`,
  role: pick(['Administrador', 'Operador', 'Bodega', 'Contabilidad']),
  is_active: rnd() > 0.15,
  last_login_at: daysAgo(between(0, 20)),
  created_at: daysAgo(between(20, 500)),
}));

const businesses = [
  { id: 1, name: 'Probability Demo', logo_url: null, primary_color: '#5E17EB', secondary_color: '#14F5A0', accent_color: '#F472B6', is_active: true, nit: '901234567-1', city: 'Bogota' },
  { id: 2, name: 'Tienda Norte SAS', logo_url: null, primary_color: '#0F172A', secondary_color: '#06B6D4', accent_color: '#BE185D', is_active: true, nit: '900876543-2', city: 'Medellin' },
];

module.exports = {
  ORDER_STATUSES,
  PAYMENT_STATUSES,
  FULFILLMENT_STATUSES,
  CARRIERS,
  CHANNELS,
  customers,
  products,
  warehouses,
  orders,
  shipments,
  invoices,
  walletMovements,
  inventoryMovements,
  drivers,
  vehicles,
  deliveryRoutes,
  integrations,
  tickets,
  users,
  businesses,
  daysAgo,
};
