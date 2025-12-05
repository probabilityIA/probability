-- =====================================================
-- Script para insertar 20 órdenes de prueba
-- Business ID: NULL (global)
-- Integration ID: 1 (Shopify)
-- =====================================================

-- Orden 1
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-1', '#WA01001', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-1',
    1200.00, 228.00, 0.00, 0.00, 1428.00, 'USD',
    'Juan Pérez', 'juan.perez@example.com', '+57 3001234567', '1012345678',
    'Calle 72 #10-34', 'Bogotá', 'Cundinamarca', 'Colombia', '110111',
    4.6097, -74.0817, 1, true, NOW() - INTERVAL '2 hours',
    'processing', 'processing',
    '[{"sku":"PROD-001","name":"Laptop HP 15","quantity":1,"price":1200.00,"weight":2.5}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_1","utm_source":"google","utm_medium":"cpc","device_type":"desktop"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 1"}'::jsonb,
    '{"transaction_id":"TXN-001","processor":"stripe","last_4_digits":"4532"}'::jsonb,
    NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days' + INTERVAL '5 minutes'
);

-- Orden 2
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-2', '#WA01002', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-2',
    50.00, 9.50, 0.00, 15.00, 74.50, 'USD',
    'María García', 'maria.garcia@example.com', '+57 3109876543', '1023456789',
    'Carrera 15 #45-67', 'Medellín', 'Antioquia', 'Colombia', '050001',
    6.2442, -75.5812, 2, false, 'pending', 'pending',
    '[{"sku":"PROD-002","name":"Mouse Logitech","quantity":2,"price":25.00,"weight":0.2}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_2","utm_source":"google","utm_medium":"cpc","device_type":"mobile"}'::jsonb,
    '{"provider":"Servientrega","service_type":"standard","zone":"Zone 2"}'::jsonb,
    '{"transaction_id":"TXN-002","processor":"stripe","last_4_digits":"5678"}'::jsonb,
    NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days' + INTERVAL '5 minutes'
);

-- Orden 3
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at, delivered_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-3', '#WA01003', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-3',
    350.00, 66.50, 35.00, 15.00, 396.50, 'USD',
    'Carlos López', 'carlos.lopez@example.com', '+57 3201234567', '1034567890',
    'Avenida 68 #23-45', 'Cali', 'Valle del Cauca', 'Colombia', '760001',
    3.4516, -76.5320, 3, true, NOW() - INTERVAL '8 hours',
    'delivered', 'delivered',
    '[{"sku":"PROD-004","name":"Monitor Samsung 24\"","quantity":1,"price":350.00,"weight":5.0}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_3","utm_source":"google","utm_medium":"cpc","device_type":"tablet"}'::jsonb,
    '{"provider":"Servientrega","service_type":"same_day","zone":"Zone 3"}'::jsonb,
    '{"transaction_id":"TXN-003","processor":"stripe","last_4_digits":"9012"}'::jsonb,
    NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days' + INTERVAL '5 minutes',
    NOW() - INTERVAL '7 days'
);

-- Orden 4
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-4', '#WA01004', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-4',
    150.00, 28.50, 0.00, 15.00, 193.50, 'USD',
    'Ana Martínez', 'ana.martinez@example.com', '+57 3112345678', '1045678901',
    'Calle 100 #15-20', 'Barranquilla', 'Atlántico', 'Colombia', '080001',
    10.9639, -74.7964, 4, true, NOW() - INTERVAL '1 hour',
    'shipped', 'shipped',
    '[{"sku":"PROD-005","name":"Auriculares Sony","quantity":1,"price":150.00,"weight":0.3}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_4","utm_source":"google","utm_medium":"cpc","device_type":"desktop"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 4"}'::jsonb,
    '{"transaction_id":"TXN-004","processor":"stripe","last_4_digits":"3456"}'::jsonb,
    NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days' + INTERVAL '5 minutes'
);

-- Orden 5
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-5', '#WA01005', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-5',
    120.00, 22.80, 0.00, 15.00, 157.80, 'USD',
    'Luis Rodríguez', 'luis.rodriguez@example.com', '+57 3223456789', '1056789012',
    'Carrera 7 #32-16', 'Cartagena', 'Bolívar', 'Colombia', '130001',
    10.3910, -75.4794, 5, false, 'cancelled', 'cancelled',
    '[{"sku":"PROD-006","name":"Webcam Logitech","quantity":1,"price":90.00,"weight":0.4},{"sku":"PROD-007","name":"Cable HDMI","quantity":2,"price":15.00,"weight":0.1}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_5","utm_source":"google","utm_medium":"cpc","device_type":"mobile"}'::jsonb,
    '{"provider":"Servientrega","service_type":"standard","zone":"Zone 5"}'::jsonb,
    '{"transaction_id":"TXN-005","processor":"stripe","last_4_digits":"7890"}'::jsonb,
    NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day' + INTERVAL '5 minutes'
);

-- Orden 6
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-6', '#WA01006', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-6',
    80.00, 15.20, 8.00, 15.00, 102.20, 'USD',
    'Carmen Fernández', 'carmen.fernandez@example.com', '+57 3334567890', '1067890123',
    'Transversal 45 #12-34', 'Bucaramanga', 'Santander', 'Colombia', '680001',
    7.1193, -73.1227, 6, true, NOW() - INTERVAL '3 hours',
    'processing', 'processing',
    '[{"sku":"PROD-003","name":"Teclado Mecánico","quantity":1,"price":80.00,"weight":0.8}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_1","utm_source":"google","utm_medium":"cpc","device_type":"tablet"}'::jsonb,
    '{"provider":"Servientrega","service_type":"standard","zone":"Zone 1"}'::jsonb,
    '{"transaction_id":"TXN-006","processor":"stripe","last_4_digits":"1234"}'::jsonb,
    NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days' + INTERVAL '5 minutes'
);

-- Orden 7
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-7', '#WA01007', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-7',
    1200.00, 228.00, 0.00, 0.00, 1428.00, 'USD',
    'José González', 'jose.gonzalez@example.com', '+57 3445678901', '1078901234',
    'Diagonal 85 #30-15', 'Pereira', 'Risaralda', 'Colombia', '660001',
    4.8133, -75.6961, 7, true, NOW() - INTERVAL '5 hours',
    'shipped', 'shipped',
    '[{"sku":"PROD-001","name":"Laptop HP 15","quantity":1,"price":1200.00,"weight":2.5}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_2","utm_source":"google","utm_medium":"cpc","device_type":"desktop"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 2"}'::jsonb,
    '{"transaction_id":"TXN-007","processor":"stripe","last_4_digits":"5678"}'::jsonb,
    NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days' + INTERVAL '5 minutes'
);

-- Orden 8
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-8', '#WA01008', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-8',
    50.00, 9.50, 0.00, 15.00, 74.50, 'USD',
    'Laura Sánchez', 'laura.sanchez@example.com', '+57 3556789012', '1089012345',
    'Calle 26 #50-21', 'Manizales', 'Caldas', 'Colombia', '170001',
    5.0689, -75.5174, 8, false, 'pending', 'pending',
    '[{"sku":"PROD-002","name":"Mouse Logitech","quantity":2,"price":25.00,"weight":0.2}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_3","utm_source":"google","utm_medium":"cpc","device_type":"mobile"}'::jsonb,
    '{"provider":"Servientrega","service_type":"same_day","zone":"Zone 3"}'::jsonb,
    '{"transaction_id":"TXN-008","processor":"stripe","last_4_digits":"9012"}'::jsonb,
    NOW() - INTERVAL '8 hours', NOW() - INTERVAL '8 hours' + INTERVAL '5 minutes'
);

-- Orden 9
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at, delivered_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-9', '#WA01009', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-9',
    350.00, 66.50, 35.00, 15.00, 396.50, 'USD',
    'Miguel Torres', 'miguel.torres@example.com', '+57 3667890123', '1090123456',
    'Carrera 50 #80-45', 'Ibagué', 'Tolima', 'Colombia', '730001',
    4.4389, -75.2322, 1, true, NOW() - INTERVAL '12 hours',
    'delivered', 'delivered',
    '[{"sku":"PROD-004","name":"Monitor Samsung 24\"","quantity":1,"price":350.00,"weight":5.0}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_4","utm_source":"google","utm_medium":"cpc","device_type":"tablet"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 4"}'::jsonb,
    '{"transaction_id":"TXN-009","processor":"stripe","last_4_digits":"3456"}'::jsonb,
    NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days' + INTERVAL '5 minutes',
    NOW() - INTERVAL '12 days'
);

-- Orden 10
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-10', '#WA01010', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-10',
    150.00, 28.50, 0.00, 15.00, 193.50, 'USD',
    'Isabel Ramírez', 'isabel.ramirez@example.com', '+57 3778901234', '1001234567',
    'Avenida Boyacá #134-56', 'Santa Marta', 'Magdalena', 'Colombia', '470001',
    11.2408, -74.2120, 2, true, NOW() - INTERVAL '6 hours',
    'processing', 'processing',
    '[{"sku":"PROD-005","name":"Auriculares Sony","quantity":1,"price":150.00,"weight":0.3}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_5","utm_source":"google","utm_medium":"cpc","device_type":"desktop"}'::jsonb,
    '{"provider":"Servientrega","service_type":"standard","zone":"Zone 5"}'::jsonb,
    '{"transaction_id":"TXN-010","processor":"stripe","last_4_digits":"7890"}'::jsonb,
    NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days' + INTERVAL '5 minutes'
);

-- Orden 11
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-11', '#WA01011', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-11',
    120.00, 22.80, 12.00, 15.00, 145.80, 'USD',
    'Francisco Flores', 'francisco.flores@example.com', '+57 3889012345', '1012345679',
    'Calle 72 #10-34', 'Bogotá', 'Cundinamarca', 'Colombia', '110111',
    4.6097, -74.0817, 3, true, NOW() - INTERVAL '4 hours',
    'shipped', 'shipped',
    '[{"sku":"PROD-006","name":"Webcam Logitech","quantity":1,"price":90.00,"weight":0.4},{"sku":"PROD-007","name":"Cable HDMI","quantity":2,"price":15.00,"weight":0.1}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_1","utm_source":"google","utm_medium":"cpc","device_type":"mobile"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 1"}'::jsonb,
    '{"transaction_id":"TXN-011","processor":"stripe","last_4_digits":"1234"}'::jsonb,
    NOW() - INTERVAL '9 days', NOW() - INTERVAL '9 days' + INTERVAL '5 minutes'
);

-- Orden 12
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-12', '#WA01012', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-12',
    80.00, 15.20, 0.00, 15.00, 110.20, 'USD',
    'Patricia Díaz', 'patricia.diaz@example.com', '+57 3990123456', '1023456780',
    'Carrera 15 #45-67', 'Medellín', 'Antioquia', 'Colombia', '050001',
    6.2442, -75.5812, 4, false, 'pending', 'pending',
    '[{"sku":"PROD-003","name":"Teclado Mecánico","quantity":1,"price":80.00,"weight":0.8}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_2","utm_source":"google","utm_medium":"cpc","device_type":"tablet"}'::jsonb,
    '{"provider":"Servientrega","service_type":"same_day","zone":"Zone 2"}'::jsonb,
    '{"transaction_id":"TXN-012","processor":"stripe","last_4_digits":"5678"}'::jsonb,
    NOW() - INTERVAL '12 hours', NOW() - INTERVAL '12 hours' + INTERVAL '5 minutes'
);

-- Orden 13
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-13', '#WA01013', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-13',
    1200.00, 228.00, 120.00, 0.00, 1308.00, 'USD',
    'Antonio Morales', 'antonio.morales@example.com', '+57 3001234568', '1034567891',
    'Avenida 68 #23-45', 'Cali', 'Valle del Cauca', 'Colombia', '760001',
    3.4516, -76.5320, 5, true, NOW() - INTERVAL '10 hours',
    'processing', 'processing',
    '[{"sku":"PROD-001","name":"Laptop HP 15","quantity":1,"price":1200.00,"weight":2.5}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_3","utm_source":"google","utm_medium":"cpc","device_type":"desktop"}'::jsonb,
    '{"provider":"Servientrega","service_type":"standard","zone":"Zone 3"}'::jsonb,
    '{"transaction_id":"TXN-013","processor":"stripe","last_4_digits":"9012"}'::jsonb,
    NOW() - INTERVAL '11 days', NOW() - INTERVAL '11 days' + INTERVAL '5 minutes'
);

-- Orden 14
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-14', '#WA01014', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-14',
    50.00, 9.50, 0.00, 15.00, 74.50, 'USD',
    'Rosa Jiménez', 'rosa.jimenez@example.com', '+57 3112345679', '1045678902',
    'Calle 100 #15-20', 'Barranquilla', 'Atlántico', 'Colombia', '080001',
    10.9639, -74.7964, 6, false, 'cancelled', 'cancelled',
    '[{"sku":"PROD-002","name":"Mouse Logitech","quantity":2,"price":25.00,"weight":0.2}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_4","utm_source":"google","utm_medium":"cpc","device_type":"mobile"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 4"}'::jsonb,
    '{"transaction_id":"TXN-014","processor":"stripe","last_4_digits":"3456"}'::jsonb,
    NOW() - INTERVAL '18 hours', NOW() - INTERVAL '18 hours' + INTERVAL '5 minutes'
);

-- Orden 15
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at, delivered_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-15', '#WA01015', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-15',
    350.00, 66.50, 0.00, 15.00, 431.50, 'USD',
    'Manuel Ruiz', 'manuel.ruiz@example.com', '+57 3223456780', '1056789013',
    'Carrera 7 #32-16', 'Cartagena', 'Bolívar', 'Colombia', '130001',
    10.3910, -75.4794, 7, true, NOW() - INTERVAL '20 hours',
    'delivered', 'delivered',
    '[{"sku":"PROD-004","name":"Monitor Samsung 24\"","quantity":1,"price":350.00,"weight":5.0}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_5","utm_source":"google","utm_medium":"cpc","device_type":"tablet"}'::jsonb,
    '{"provider":"Servientrega","service_type":"same_day","zone":"Zone 5"}'::jsonb,
    '{"transaction_id":"TXN-015","processor":"stripe","last_4_digits":"7890"}'::jsonb,
    NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days' + INTERVAL '5 minutes',
    NOW() - INTERVAL '18 days'
);

-- Orden 16
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-16', '#WA01016', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-16',
    150.00, 28.50, 15.00, 15.00, 178.50, 'USD',
    'Teresa Álvarez', 'teresa.alvarez@example.com', '+57 3334567891', '1067890124',
    'Transversal 45 #12-34', 'Bucaramanga', 'Santander', 'Colombia', '680001',
    7.1193, -73.1227, 8, true, NOW() - INTERVAL '7 hours',
    'shipped', 'shipped',
    '[{"sku":"PROD-005","name":"Auriculares Sony","quantity":1,"price":150.00,"weight":0.3}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_1","utm_source":"google","utm_medium":"cpc","device_type":"desktop"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 1"}'::jsonb,
    '{"transaction_id":"TXN-016","processor":"stripe","last_4_digits":"1234"}'::jsonb,
    NOW() - INTERVAL '13 days', NOW() - INTERVAL '13 days' + INTERVAL '5 minutes'
);

-- Orden 17
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-17', '#WA01017', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-17',
    120.00, 22.80, 0.00, 15.00, 157.80, 'USD',
    'Pedro Romero', 'pedro.romero@example.com', '+57 3445678902', '1078901235',
    'Diagonal 85 #30-15', 'Pereira', 'Risaralda', 'Colombia', '660001',
    4.8133, -75.6961, 1, false, 'pending', 'pending',
    '[{"sku":"PROD-006","name":"Webcam Logitech","quantity":1,"price":90.00,"weight":0.4},{"sku":"PROD-007","name":"Cable HDMI","quantity":2,"price":15.00,"weight":0.1}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_2","utm_source":"google","utm_medium":"cpc","device_type":"mobile"}'::jsonb,
    '{"provider":"Servientrega","service_type":"standard","zone":"Zone 2"}'::jsonb,
    '{"transaction_id":"TXN-017","processor":"stripe","last_4_digits":"5678"}'::jsonb,
    NOW() - INTERVAL '14 hours', NOW() - INTERVAL '14 hours' + INTERVAL '5 minutes'
);

-- Orden 18
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-18', '#WA01018', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-18',
    80.00, 15.20, 8.00, 15.00, 102.20, 'USD',
    'Lucía Navarro', 'lucia.navarro@example.com', '+57 3556789013', '1089012346',
    'Calle 26 #50-21', 'Manizales', 'Caldas', 'Colombia', '170001',
    5.0689, -75.5174, 2, true, NOW() - INTERVAL '9 hours',
    'processing', 'processing',
    '[{"sku":"PROD-003","name":"Teclado Mecánico","quantity":1,"price":80.00,"weight":0.8}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_3","utm_source":"google","utm_medium":"cpc","device_type":"tablet"}'::jsonb,
    '{"provider":"Servientrega","service_type":"same_day","zone":"Zone 3"}'::jsonb,
    '{"transaction_id":"TXN-018","processor":"stripe","last_4_digits":"9012"}'::jsonb,
    NOW() - INTERVAL '16 days', NOW() - INTERVAL '16 days' + INTERVAL '5 minutes'
);

-- Orden 19
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-19', '#WA01019', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-19',
    1200.00, 228.00, 0.00, 0.00, 1428.00, 'USD',
    'Javier Castro', 'javier.castro@example.com', '+57 3667890124', '1090123457',
    'Carrera 50 #80-45', 'Ibagué', 'Tolima', 'Colombia', '730001',
    4.4389, -75.2322, 3, true, NOW() - INTERVAL '11 hours',
    'shipped', 'shipped',
    '[{"sku":"PROD-001","name":"Laptop HP 15","quantity":1,"price":1200.00,"weight":2.5}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_4","utm_source":"google","utm_medium":"cpc","device_type":"desktop"}'::jsonb,
    '{"provider":"Servientrega","service_type":"express","zone":"Zone 4"}'::jsonb,
    '{"transaction_id":"TXN-019","processor":"stripe","last_4_digits":"3456"}'::jsonb,
    NOW() - INTERVAL '17 days', NOW() - INTERVAL '17 days' + INTERVAL '5 minutes'
);

-- Orden 20
INSERT INTO orders (
    id, created_at, updated_at, business_id, integration_id, integration_type, platform,
    external_id, order_number, internal_number, subtotal, tax, discount, shipping_cost,
    total_amount, currency, customer_name, customer_email, customer_phone, customer_dni,
    shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
    shipping_lat, shipping_lng, payment_method_id, is_paid, paid_at, status, original_status,
    items, metadata, shipping_details, payment_details, occurred_at, imported_at, delivered_at
) VALUES (
    gen_random_uuid(), NOW(), NOW(), NULL, 1, 'whatsapp', 'whatsapp',
    'WA-' || EXTRACT(EPOCH FROM NOW())::bigint || '-20', '#WA01020', 'ORD-' || EXTRACT(EPOCH FROM NOW())::bigint || '-20',
    150.00, 28.50, 15.00, 15.00, 178.50, 'USD',
    'Elena Ortiz', 'elena.ortiz@example.com', '+57 3778901235', '1001234568',
    'Avenida Boyacá #134-56', 'Santa Marta', 'Magdalena', 'Colombia', '470001',
    11.2408, -74.2120, 4, true, NOW() - INTERVAL '13 hours',
    'delivered', 'delivered',
    '[{"sku":"PROD-005","name":"Auriculares Sony","quantity":1,"price":150.00,"weight":0.3}]'::jsonb,
    '{"source":"test_seed","campaign":"campaign_5","utm_source":"google","utm_medium":"cpc","device_type":"mobile"}'::jsonb,
    '{"provider":"Servientrega","service_type":"standard","zone":"Zone 5"}'::jsonb,
    '{"transaction_id":"TXN-020","processor":"stripe","last_4_digits":"7890"}'::jsonb,
    NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days' + INTERVAL '5 minutes',
    NOW() - INTERVAL '22 days'
);

-- Mensaje de confirmación
SELECT 'Se han insertado 20 órdenes de prueba exitosamente' AS resultado;
