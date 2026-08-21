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

const INTEGRATION_TYPES = [
  { id: 1, code: 'Shopify', name: 'Shopify', category: 'ecommerce', image_url: 'integration-types/1765744750_shopify.png' },
  { id: 2, code: 'Whastap', name: 'WhatsApp', category: 'messaging', image_url: 'integration-types/1765744972_whatsap.png' },
  { id: 3, code: 'Mercado Libre', name: 'Mercado Libre', category: 'ecommerce', image_url: 'integration-types/1771905467_mercado-libre-logo.png' },
  { id: 4, code: 'woocommerce', name: 'WooCommerce', category: 'ecommerce', image_url: 'integration-types/1765745053_woocomerce.webp' },
  { id: 7, code: 'factus', name: 'Factus', category: 'invoicing', image_url: 'integration-types/1771736455_factus.png' },
  { id: 8, code: 'siigo', name: 'Siigo', category: 'invoicing', image_url: 'integration-types/1771738597_siigo.png' },
  { id: 9, code: 'alegra', name: 'Alegra', category: 'invoicing', image_url: 'integration-types/1771738659_alegra.png' },
  { id: 10, code: 'world_office', name: 'World Office', category: 'invoicing', image_url: 'integration-types/1771738901_worl-office.png' },
  { id: 11, code: 'helisa', name: 'Helisa', category: 'invoicing', image_url: 'integration-types/1771739081_logo-helisa.png' },
  { id: 12, code: 'envioclick', name: 'EnvioClick', category: 'shipping', image_url: 'integration-types/1771901154_envioclik.png' },
  { id: 13, code: 'enviame', name: 'Enviame', category: 'shipping', image_url: 'integration-types/1771905179_enviame.png' },
  { id: 15, code: 'mipaquete', name: 'MiPaquete', category: 'shipping', image_url: 'integration-types/1771905337_mipaquete.png' },
  { id: 16, code: 'vtex', name: 'VTEX', category: 'ecommerce', image_url: 'integration-types/1771905400_vtex.png' },
  { id: 17, code: 'tiendanube', name: 'Tiendanube', category: 'ecommerce', image_url: 'integration-types/1771905372_tiendanube.png' },
  { id: 18, code: 'magento', name: 'Magento', category: 'ecommerce', image_url: 'integration-types/1771905298_magneto.png' },
  { id: 19, code: 'amazon', name: 'Amazon', category: 'ecommerce', image_url: 'integration-types/1771905134_amazon.png' },
  { id: 20, code: 'falabella', name: 'Falabella', category: 'ecommerce', image_url: 'integration-types/1771905254_falabella.png' },
  { id: 21, code: 'exito', name: 'Exito', category: 'ecommerce', image_url: 'integration-types/1771905220_exito.png' },
  { id: 22, code: 'nequi', name: 'Nequi', category: 'payment', image_url: 'integration-types/1772147410_nequi.png' },
  { id: 23, code: 'bold_pay', name: 'Bold', category: 'payment', image_url: 'integration-types/1777354669_bold.png' },
  { id: 24, code: 'wompi', name: 'Wompi', category: 'payment', image_url: 'integration-types/1772147507_logowompi.png' },
  { id: 25, code: 'stripe', name: 'Stripe', category: 'payment', image_url: 'integration-types/1772147475_stripe.png' },
  { id: 26, code: 'payu', name: 'PayU', category: 'payment', image_url: 'integration-types/1772147440_payu.png' },
  { id: 27, code: 'epayco', name: 'ePayco', category: 'payment', image_url: 'integration-types/1772147322_epayco.png' },
  { id: 28, code: 'melipago', name: 'MercadoPago', category: 'payment', image_url: 'integration-types/1772147377_mercado-pago.png' },
  { id: 33, code: 'jumpseller', name: 'Jumpseller', category: 'ecommerce', image_url: 'integration-types/1784241131_jumpseller.png' },
  { id: 34, code: 'shipit', name: 'Shipit', category: 'shipping', image_url: 'integration-types/1785346702_shipit.png' },
  { id: 35, code: 'tiktok', name: 'TikTok', category: 'ecommerce', image_url: 'integration-types/tiktok.png' },
  { id: 42, code: 'bancolombia_qr', name: 'Bancolombia QR', category: 'payment', image_url: 'integration-types/bancolombia.png' },
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

const CHANNEL_TYPES = ['Shopify', 'woocommerce', 'Mercado Libre', 'Whastap', 'amazon', 'platform'];
const CHANNEL_LABELS = {
  'Shopify': 'Shopify',
  'woocommerce': 'WooCommerce',
  'Mercado Libre': 'Mercado Libre',
  'Whastap': 'WhatsApp',
  'amazon': 'Amazon',
  'platform': 'Manual',
};

const orders = Array.from({ length: 64 }, (_, i) => {
  const customer = customers[i % customers.length];
  const status = ORDER_STATUSES[between(0, ORDER_STATUSES.length - 1)];
  const payment = PAYMENT_STATUSES[between(0, PAYMENT_STATUSES.length - 1)];
  const fulfillment = FULFILLMENT_STATUSES[between(0, FULFILLMENT_STATUSES.length - 1)];
  const channelCode = CHANNEL_TYPES[i % CHANNEL_TYPES.length];
  const channelType = INTEGRATION_TYPES.find((t) => t.code === channelCode);
  const itemCount = between(1, 4);
  const items = Array.from({ length: itemCount }, (_, j) => {
    const product = products[(i + j) % products.length];
    const quantity = between(1, 3);
    return {
      id: `${i}-${j}`,
      product_id: product.id,
      sku: product.sku,
      name: product.name,
      quantity,
      unit_price: product.price,
      total_price: product.price * quantity,
      image_url: null,
    };
  });
  const subtotal = items.reduce((acc, it) => acc + it.unit_price * it.quantity, 0);
  const shipping = between(8000, 24000);
  const tax = Math.round(subtotal * 0.19);
  const isCod = rnd() > 0.55;
  const hasShipment = status.id >= 4;
  const isPaid = payment.code === 'paid';

  return {
    id: `ord-${String(i + 1).padStart(4, '0')}-uuid`,
    created_at: daysAgo(between(0, 60)),
    updated_at: daysAgo(between(0, 5)),
    business_id: 1,
    integration_id: (i % 6) + 1,
    integration_type: channelCode,
    integration_name: CHANNEL_LABELS[channelCode],
    integration_logo_url: channelType ? channelType.image_url : null,
    platform: CHANNEL_LABELS[channelCode],
    external_id: `${between(100000, 999999)}`,
    order_number: `ORD-${10240 + i}`,
    internal_number: `INT-${5000 + i}`,
    subtotal,
    tax,
    discount: 0,
    shipping_cost: shipping,
    total_amount: subtotal + shipping,
    currency: 'COP',
    cod_total: isCod ? subtotal + shipping : 0,
    customer_id: customer.id,
    customer_name: customer.name,
    customer_email: customer.email,
    customer_phone: customer.phone,
    customer_dni: customer.document_number,
    shipping_street: customer.address,
    shipping_city: customer.city,
    shipping_state: 'Colombia',
    shipping_country: 'CO',
    shipping_postal_code: '110111',
    payment_method_id: 1,
    is_paid: isPaid,
    paid_at: isPaid ? daysAgo(between(0, 20)) : null,
    tracking_number: hasShipment ? `GU${between(100000000, 999999999)}` : null,
    tracking_link: hasShipment ? 'https://example.com/tracking' : null,
    guide_id: hasShipment ? `${between(1000, 9999)}` : null,
    guide_link: hasShipment ? 'https://example.com/guia.pdf' : null,
    warehouse_id: 1,
    warehouse_name: 'Bodega principal',
    driver_id: null,
    driver_name: '',
    is_last_mile: false,
    weight: between(200, 4000),
    order_type_id: 1,
    order_type_name: 'Venta',
    status: status.code,
    original_status: status.code,
    status_id: status.id,
    order_status: { id: status.id, code: status.code, name: status.name, category: 'order', color: null },
    payment_status_id: payment.id,
    payment_status: { id: payment.id, code: payment.code, name: payment.name, category: 'payment', color: null },
    fulfillment_status_id: fulfillment.id,
    fulfillment_status: { id: fulfillment.id, code: fulfillment.code, name: fulfillment.name, category: 'fulfillment', color: null },
    notes: null,
    approved: true,
    user_id: 1,
    user_name: 'Sebastian Camacho',
    is_confirmed: status.id >= 2,
    invoiceable: true,
    invoice_status: isPaid ? 'issued' : null,
    invoice_url: isPaid ? 'https://example.com/factura.pdf' : null,
    order_items: items,
    items,
    items_count: itemCount,
    is_cod: isCod,
    has_shipment: hasShipment,
    has_invoice: isPaid,
    total: subtotal + shipping,
    channel: CHANNEL_LABELS[channelCode],
    status_code: status.code,
    status_name: status.name,
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

const typeByCode = (code) => INTEGRATION_TYPES.find((t) => t.code === code);

const CONNECTED = [
  { id: 1, code: 'Shopify', name: 'Tienda Shopify', status: 'connected', synced: 1842, days: 0 },
  { id: 2, code: 'woocommerce', name: 'WooCommerce principal', status: 'connected', synced: 634, days: 0 },
  { id: 3, code: 'Mercado Libre', name: 'MercadoLibre CO', status: 'error', synced: 221, days: 2 },
  { id: 4, code: 'siigo', name: 'Siigo contabilidad', status: 'connected', synced: 0, days: 1 },
  { id: 5, code: 'Whastap', name: 'WhatsApp Business', status: 'connected', synced: 0, days: 0 },
  { id: 6, code: 'envioclick', name: 'EnvioClick', status: 'connected', synced: 0, days: 0 },
  { id: 7, code: 'bold_pay', name: 'Bold pagos', status: 'disconnected', synced: 0, days: null },
  { id: 8, code: 'tiendanube', name: 'Tiendanube CO', status: 'connected', synced: 96, days: 1 },
];

const integrations = CONNECTED.map((row) => {
  const type = typeByCode(row.code);
  return {
    id: row.id,
    business_id: 1,
    name: row.name,
    code: row.code,
    category: type.category,
    category_code: type.category,
    integration_type_id: type.id,
    integration_type_name: type.name,
    integration_type_code: type.code,
    integration_type: {
      id: type.id,
      code: type.code,
      name: type.name,
      category: type.category,
      image_url: type.image_url,
    },
    image_url: type.image_url,
    is_active: row.status !== 'disconnected',
    is_default: row.id <= 2,
    status: row.status,
    last_sync_at: row.days === null ? null : daysAgo(row.days),
    orders_synced: row.synced,
    created_at: daysAgo(120),
    updated_at: daysAgo(row.days === null ? 30 : row.days),
  };
});

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
  INTEGRATION_TYPES,
  tickets,
  users,
  businesses,
  daysAgo,
};
