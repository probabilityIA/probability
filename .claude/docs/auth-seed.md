# Auth Seed - Datos basicos (Scopes / Resources / Actions / Roles / Permissions / Business Types)

Snapshot tomado de la DB actual (probability) para recrear el sistema de autorizacion en otra DB.

## Modelo

```
business_type (1) ── role (N) ─┬─ role_permissions (N) ── permission (N) ── resource
                               │                                          └─ action
                               └─ scope                                    └─ scope
user ── user_roles ── role
user ── user_businesses ── business
```

- `scope`: `platform` (super admin global) vs `business` (acotado al negocio).
- `role.level`: 1 = mas alto. `is_system` = no editable desde UI.
- `permission` = combinacion (resource, action, scope, business_type).
- Super Admin (role_id=1, scope=platform) tiene acceso total via codigo, no requiere filas en `role_permissions`.

## 1. business_type

| id | code | name | description |
|----|------|------|-------------|
| 1 | Probability | Probability | Probability |

```sql
INSERT INTO business_type (id, name, code, description, icon, is_active, created_at, updated_at)
VALUES (1, 'Probability', 'Probability', 'Probability', '', true, NOW(), NOW());
```

## 2. scope

| id | code | name | is_system |
|----|------|------|-----------|
| 1 | platform | Platform | true |
| 2 | business | Business | false |

```sql
INSERT INTO scope (id, name, code, description, is_system, created_at, updated_at) VALUES
(1, 'Platform', 'platform', 'Scope for platform-wide permissions', true, NOW(), NOW()),
(2, 'Business', 'business', NULL, false, NOW(), NOW());
```

## 3. action

| id | name | description |
|----|------|-------------|
| 1 | Create | Create new records |
| 2 | Read | Read/view information |
| 3 | Update | Modify existing records |
| 4 | Delete | Delete records |
| 5 | Manage | Full control (includes all actions) |
| 6 | Approve | Approve requests or documents |
| 7 | Reject | Reject requests or documents |
| 8 | Assign | Assign resources or tasks |
| 9 | Schedule | Schedule events or tasks |
| 10 | Report | Generate reports |
| 11 | Configure | Configure system parameters |
| 12 | Audit | Audit system actions |
| 13 | Migrate | Execute data migrations |

```sql
INSERT INTO action (id, name, description, created_at, updated_at) VALUES
(1,'Create','Create new records',NOW(),NOW()),
(2,'Read','Read/view information',NOW(),NOW()),
(3,'Update','Modify existing records',NOW(),NOW()),
(4,'Delete','Delete records',NOW(),NOW()),
(5,'Manage','Full control (includes all actions)',NOW(),NOW()),
(6,'Approve','Approve requests or documents',NOW(),NOW()),
(7,'Reject','Reject requests or documents',NOW(),NOW()),
(8,'Assign','Assign resources or tasks',NOW(),NOW()),
(9,'Schedule','Schedule events or tasks',NOW(),NOW()),
(10,'Report','Generate reports',NOW(),NOW()),
(11,'Configure','Configure system parameters',NOW(),NOW()),
(12,'Audit','Audit system actions',NOW(),NOW()),
(13,'Migrate','Execute data migrations',NOW(),NOW());
```

## 4. resource

| id | name | description |
|----|------|-------------|
| 1 | Usuarios | Usuarios |
| 2 | Permisos | Permisos |
| 3 | Roles | Roles |
| 4 | Recursos | Recursos |
| 5 | Empresas | Empresas |
| 6 | Ordenes | Ordenes |
| 7 | Productos | Productos |
| 8 | Integraciones | Integraciones |
| 9 | Envios | Envios |
| 10 | Facturacion | Facturacion |
| 11 | Integraciones-Platform | Integraciones-Platform |
| 12 | Integraciones-E-commerce | Integraciones-E-commerce |
| 13 | Integraciones-Facturacion-Electronica | Integraciones-Facturacion-Electronica |
| 14 | Integraciones-Mensajeria | Integraciones-Mensajeria |
| 15 | Integraciones-Pagos | Integraciones-Pagos |
| 16 | Integraciones-Logistica | Integraciones-Logistica |
| 17 | Integraciones-Tipos-de-integracion | Integraciones-Tipos-de-integracion |
| 18 | Notificaciones | Configuracion de notificaciones por integracion |
| 19 | Clientes | Gestion de clientes |
| 20 | Ultima Milla | Gestion de ultima milla (delivery) |
| 21 | Billetera | Gestion de billetera y transacciones |
| 22 | Inventario | Gestion de inventario |
| 23 | Bodegas | Gestion de bodegas y ubicaciones |
| 24 | Storefront | Modulo de tienda para clientes finales |
| 25 | Inventario-Stock | Vista de stock por bodega/SKU |
| 26 | Inventario-Movimientos | Historial y registro de movimientos |
| 27 | Inventario-Trazabilidad | Trazabilidad por lote/serial |
| 28 | Inventario-Kardex | Kardex contable por SKU |
| 29 | Inventario-Operaciones | Recepciones, despachos y conteos |
| 30 | Inventario-Slotting | Analitica de slotting ABC |
| 31 | Inventario-Auditoria | Auditoria de inventario |
| 32 | Inventario-LPN | Gestion de License Plate Numbers |
| 33 | Inventario-Scan | App movil de captura por escaneo |
| 34 | Inventario-Sync-Logs | Logs de sincronizacion de inventario |

```sql
INSERT INTO resource (id, name, description, created_at, updated_at) VALUES
(1,'Usuarios','Usuarios',NOW(),NOW()),
(2,'Permisos','Permisos',NOW(),NOW()),
(3,'Roles','Roles',NOW(),NOW()),
(4,'Recursos','Recursos',NOW(),NOW()),
(5,'Empresas','Empresas',NOW(),NOW()),
(6,'Ordenes','Ordenes',NOW(),NOW()),
(7,'Productos','Productos',NOW(),NOW()),
(8,'Integraciones','Integraciones',NOW(),NOW()),
(9,'Envios','Envios',NOW(),NOW()),
(10,'Facturacion','Facturacion',NOW(),NOW()),
(11,'Integraciones-Platform','Integraciones-Platform',NOW(),NOW()),
(12,'Integraciones-E-commerce','Integraciones-E-commerce',NOW(),NOW()),
(13,'Integraciones-Facturacion-Electronica','Integraciones-Facturacion-Electronica',NOW(),NOW()),
(14,'Integraciones-Mensajeria','Integraciones-Mensajeria',NOW(),NOW()),
(15,'Integraciones-Pagos','Integraciones-Pagos',NOW(),NOW()),
(16,'Integraciones-Logistica','Integraciones-Logistica',NOW(),NOW()),
(17,'Integraciones-Tipos-de-integracion','Integraciones-Tipos-de-integracion',NOW(),NOW()),
(18,'Notificaciones','Configuracion de notificaciones por integracion',NOW(),NOW()),
(19,'Clientes','Gestion de clientes',NOW(),NOW()),
(20,'Ultima Milla','Gestion de ultima milla (delivery)',NOW(),NOW()),
(21,'Billetera','Gestion de billetera y transacciones',NOW(),NOW()),
(22,'Inventario','Gestion de inventario',NOW(),NOW()),
(23,'Bodegas','Gestion de bodegas y ubicaciones',NOW(),NOW()),
(24,'Storefront','Modulo de tienda para clientes finales',NOW(),NOW()),
(25,'Inventario-Stock','Vista de stock por bodega/SKU',NOW(),NOW()),
(26,'Inventario-Movimientos','Historial y registro de movimientos',NOW(),NOW()),
(27,'Inventario-Trazabilidad','Trazabilidad por lote/serial',NOW(),NOW()),
(28,'Inventario-Kardex','Kardex contable por SKU',NOW(),NOW()),
(29,'Inventario-Operaciones','Recepciones, despachos y conteos',NOW(),NOW()),
(30,'Inventario-Slotting','Analitica de slotting ABC',NOW(),NOW()),
(31,'Inventario-Auditoria','Auditoria de inventario',NOW(),NOW()),
(32,'Inventario-LPN','Gestion de License Plate Numbers',NOW(),NOW()),
(33,'Inventario-Scan','App movil de captura por escaneo',NOW(),NOW()),
(34,'Inventario-Sync-Logs','Logs de sincronizacion de inventario',NOW(),NOW());
```

## 5. role

| id | name | level | is_system | scope_id | business_type_id |
|----|------|-------|-----------|----------|------------------|
| 1 | Super Admin | 1 | true | 1 (platform) | NULL |
| 2 | Operador | 2 | true | 1 (platform) | NULL |
| 4 | Administrador | 1 | false | 2 (business) | 1 |
| 5 | cliente_final | 5 | true | 2 (business) | NULL |

```sql
INSERT INTO role (id, name, description, level, is_system, scope_id, business_type_id, created_at, updated_at) VALUES
(1,'Super Admin','Super Administrator with full access',1,true,1,NULL,NOW(),NOW()),
(2,'Operador','Operador',2,true,1,NULL,NOW(),NOW()),
(4,'Administrador','Administrador',1,false,2,1,NOW(),NOW()),
(5,'cliente_final','Rol para clientes finales del storefront',5,true,2,NULL,NOW(),NOW());
```

## 6. permission

Patron: por cada `resource` se crean 4 permisos (Create/Read/Update/Delete) en `scope_id=2` (business). `business_type_id=1` para los recursos de negocio originales; `NULL` para los recursos transversales (Notificaciones, Storefront e Inventario-*).

Total: 127 permisos (IDs 1..127, con algunos huecos historicos).

```sql
-- Ordenes (resource 6)
INSERT INTO permission (id,name,description,resource_id,action_id,scope_id,business_type_id,created_at,updated_at) VALUES
(1,'Ver Ordenes','Ver Ordenes',6,2,2,1,NOW(),NOW()),
(2,'Editar Ordenes','Editar Ordenes',6,3,2,1,NOW(),NOW()),
-- Integraciones (resource 8)
(3,'Crear Integraciones','',8,1,2,1,NOW(),NOW()),
(4,'Ver Integraciones','',8,2,2,1,NOW(),NOW()),
-- Productos (resource 7)
(5,'Ver Producto','Ver Productos',7,2,2,1,NOW(),NOW()),
-- Envios (resource 9)
(6,'Ver Envios','ver envios',9,2,2,1,NOW(),NOW()),
(7,'Crear Envios','',9,1,2,1,NOW(),NOW()),
(8,'Eliminar Envios','',9,4,2,1,NOW(),NOW()),
(9,'Actualizar Envios','',9,3,2,1,NOW(),NOW()),
-- Empresas (resource 5)
(10,'Create Empresas','',5,1,2,1,NOW(),NOW()),
(11,'Delete Empresas','',5,4,2,1,NOW(),NOW()),
(12,'Read Empresas','',5,2,2,1,NOW(),NOW()),
(13,'Update Empresas','',5,3,2,1,NOW(),NOW()),
-- Usuarios (resource 1)
(14,'Create Usuarios','',1,1,2,1,NOW(),NOW()),
(15,'Delete Usuarios','',1,4,2,1,NOW(),NOW()),
(16,'Read Usuarios','',1,2,2,1,NOW(),NOW()),
(17,'Update Usuarios','',1,3,2,1,NOW(),NOW()),
-- Facturacion (resource 10)
(18,'Create Facturacion','',10,1,2,1,NOW(),NOW()),
(19,'Delete Facturacion','',10,4,2,1,NOW(),NOW()),
(20,'Read Facturacion','',10,2,2,1,NOW(),NOW()),
(21,'Update Facturacion','',10,3,2,1,NOW(),NOW()),
-- Integraciones-Tipos-de-integracion (resource 17)
(22,'Create Integraciones-Tipos-de-integracion','',17,1,2,1,NOW(),NOW()),
(23,'Delete Integraciones-Tipos-de-integracion','',17,4,2,1,NOW(),NOW()),
(24,'Read Integraciones-Tipos-de-integracion','',17,2,2,1,NOW(),NOW()),
(25,'Update Integraciones-Tipos-de-integracion','',17,3,2,1,NOW(),NOW()),
-- Integraciones-Logistica (resource 16)
(26,'Create Integraciones-Logistica','',16,1,2,1,NOW(),NOW()),
(27,'Delete Integraciones-Logistica','',16,4,2,1,NOW(),NOW()),
(28,'Read Integraciones-Logistica','',16,2,2,1,NOW(),NOW()),
(29,'Update Integraciones-Logistica','',16,3,2,1,NOW(),NOW()),
-- Integraciones-Pagos (resource 15)
(30,'Create Integraciones-Pagos','',15,1,2,1,NOW(),NOW()),
(31,'Delete Integraciones-Pagos','',15,4,2,1,NOW(),NOW()),
(32,'Read Integraciones-Pagos','',15,2,2,1,NOW(),NOW()),
(33,'Update Integraciones-Pagos','',15,3,2,1,NOW(),NOW()),
-- Integraciones-Mensajeria (resource 14)
(34,'Create Integraciones-Mensajeria','',14,1,2,1,NOW(),NOW()),
(35,'Delete Integraciones-Mensajeria','',14,4,2,1,NOW(),NOW()),
(36,'Read Integraciones-Mensajeria','',14,2,2,1,NOW(),NOW()),
(37,'Update Integraciones-Mensajeria','',14,3,2,1,NOW(),NOW()),
-- Integraciones-Facturacion-Electronica (resource 13)
(38,'Create Integraciones-Facturacion-Electronica','',13,1,2,1,NOW(),NOW()),
(39,'Delete Integraciones-Facturacion-Electronica','',13,4,2,1,NOW(),NOW()),
(40,'Read Integraciones-Facturacion-Electronica','',13,2,2,1,NOW(),NOW()),
(41,'Update Integraciones-Facturacion-Electronica','',13,3,2,1,NOW(),NOW()),
-- Integraciones-E-commerce (resource 12)
(42,'Create Integraciones-E-commerce','',12,1,2,1,NOW(),NOW()),
(43,'Delete Integraciones-E-commerce','',12,4,2,1,NOW(),NOW()),
(44,'Read Integraciones-E-commerce','',12,2,2,1,NOW(),NOW()),
(45,'Update Integraciones-E-commerce','',12,3,2,1,NOW(),NOW()),
-- Integraciones-Platform (resource 11)
(46,'Create Integraciones-Platform','',11,1,2,1,NOW(),NOW()),
(47,'Delete Integraciones-Platform','',11,4,2,1,NOW(),NOW()),
(48,'Read Integraciones-Platform','',11,2,2,1,NOW(),NOW()),
(49,'Update Integraciones-Platform','',11,3,2,1,NOW(),NOW()),
-- Notificaciones (resource 18, business_type_id NULL)
(50,'Create Notificaciones','Crear configuraciones de notificacion',18,1,2,NULL,NOW(),NOW()),
(51,'Read Notificaciones','Ver configuraciones de notificacion y auditoria de mensajes',18,2,2,NULL,NOW(),NOW()),
(52,'Update Notificaciones','Editar reglas de notificacion',18,3,2,NULL,NOW(),NOW()),
(53,'Delete Notificaciones','Eliminar configuraciones de notificacion',18,4,2,NULL,NOW(),NOW()),
-- Clientes (resource 19)
(54,'Create Clientes','Create Clientes',19,1,2,1,NOW(),NOW()),
(55,'Read Clientes','Read Clientes',19,2,2,1,NOW(),NOW()),
(56,'Update Clientes','Update Clientes',19,3,2,1,NOW(),NOW()),
(57,'Delete Clientes','Delete Clientes',19,4,2,1,NOW(),NOW()),
-- Ultima Milla (resource 20)
(58,'Create Ultima Milla','Create Ultima Milla',20,1,2,1,NOW(),NOW()),
(59,'Read Ultima Milla','Read Ultima Milla',20,2,2,1,NOW(),NOW()),
(60,'Update Ultima Milla','Update Ultima Milla',20,3,2,1,NOW(),NOW()),
(61,'Delete Ultima Milla','Delete Ultima Milla',20,4,2,1,NOW(),NOW()),
-- Billetera (resource 21)
(62,'Create Billetera','Create Billetera',21,1,2,1,NOW(),NOW()),
(63,'Read Billetera','Read Billetera',21,2,2,1,NOW(),NOW()),
(64,'Update Billetera','Update Billetera',21,3,2,1,NOW(),NOW()),
(65,'Delete Billetera','Delete Billetera',21,4,2,1,NOW(),NOW()),
-- Inventario (resource 22)
(66,'Create Inventario','Create Inventario',22,1,2,1,NOW(),NOW()),
(67,'Read Inventario','Read Inventario',22,2,2,1,NOW(),NOW()),
(68,'Update Inventario','Update Inventario',22,3,2,1,NOW(),NOW()),
(69,'Delete Inventario','Delete Inventario',22,4,2,1,NOW(),NOW()),
-- Bodegas (resource 23)
(70,'Create Bodegas','Create Bodegas',23,1,2,1,NOW(),NOW()),
(71,'Read Bodegas','Read Bodegas',23,2,2,1,NOW(),NOW()),
(72,'Update Bodegas','Update Bodegas',23,3,2,1,NOW(),NOW()),
(73,'Delete Bodegas','Delete Bodegas',23,4,2,1,NOW(),NOW()),
-- Permisos (resource 2)
(74,'Create Permisos','Create Permisos',2,1,2,1,NOW(),NOW()),
(75,'Read Permisos','Read Permisos',2,2,2,1,NOW(),NOW()),
(76,'Update Permisos','Update Permisos',2,3,2,1,NOW(),NOW()),
(77,'Delete Permisos','Delete Permisos',2,4,2,1,NOW(),NOW()),
-- Roles (resource 3)
(78,'Create Roles','Create Roles',3,1,2,1,NOW(),NOW()),
(79,'Read Roles','Read Roles',3,2,2,1,NOW(),NOW()),
(80,'Update Roles','Update Roles',3,3,2,1,NOW(),NOW()),
(81,'Delete Roles','Delete Roles',3,4,2,1,NOW(),NOW()),
-- Recursos (resource 4)
(82,'Create Recursos','Create Recursos',4,1,2,1,NOW(),NOW()),
(83,'Read Recursos','Read Recursos',4,2,2,1,NOW(),NOW()),
(84,'Update Recursos','Update Recursos',4,3,2,1,NOW(),NOW()),
(85,'Delete Recursos','Delete Recursos',4,4,2,1,NOW(),NOW()),
-- Storefront (resource 24, business_type_id NULL)
(86,'Read Storefront','Read Storefront',24,2,2,NULL,NOW(),NOW()),
(87,'Create Storefront','Create Storefront',24,1,2,NULL,NOW(),NOW()),
-- Inventario-Stock (25)
(88,'Create Inventario-Stock','Crear en Inventario-Stock',25,1,2,NULL,NOW(),NOW()),
(89,'Read Inventario-Stock','Ver Inventario-Stock',25,2,2,NULL,NOW(),NOW()),
(90,'Update Inventario-Stock','Editar Inventario-Stock',25,3,2,NULL,NOW(),NOW()),
(91,'Delete Inventario-Stock','Eliminar Inventario-Stock',25,4,2,NULL,NOW(),NOW()),
-- Inventario-Movimientos (26)
(92,'Create Inventario-Movimientos','Crear en Inventario-Movimientos',26,1,2,NULL,NOW(),NOW()),
(93,'Read Inventario-Movimientos','Ver Inventario-Movimientos',26,2,2,NULL,NOW(),NOW()),
(94,'Update Inventario-Movimientos','Editar Inventario-Movimientos',26,3,2,NULL,NOW(),NOW()),
(95,'Delete Inventario-Movimientos','Eliminar Inventario-Movimientos',26,4,2,NULL,NOW(),NOW()),
-- Inventario-Trazabilidad (27)
(96,'Create Inventario-Trazabilidad','Crear en Inventario-Trazabilidad',27,1,2,NULL,NOW(),NOW()),
(97,'Read Inventario-Trazabilidad','Ver Inventario-Trazabilidad',27,2,2,NULL,NOW(),NOW()),
(98,'Update Inventario-Trazabilidad','Editar Inventario-Trazabilidad',27,3,2,NULL,NOW(),NOW()),
(99,'Delete Inventario-Trazabilidad','Eliminar Inventario-Trazabilidad',27,4,2,NULL,NOW(),NOW()),
-- Inventario-Kardex (28)
(100,'Create Inventario-Kardex','Crear en Inventario-Kardex',28,1,2,NULL,NOW(),NOW()),
(101,'Read Inventario-Kardex','Ver Inventario-Kardex',28,2,2,NULL,NOW(),NOW()),
(102,'Update Inventario-Kardex','Editar Inventario-Kardex',28,3,2,NULL,NOW(),NOW()),
(103,'Delete Inventario-Kardex','Eliminar Inventario-Kardex',28,4,2,NULL,NOW(),NOW()),
-- Inventario-Operaciones (29)
(104,'Create Inventario-Operaciones','Crear en Inventario-Operaciones',29,1,2,NULL,NOW(),NOW()),
(105,'Read Inventario-Operaciones','Ver Inventario-Operaciones',29,2,2,NULL,NOW(),NOW()),
(106,'Update Inventario-Operaciones','Editar Inventario-Operaciones',29,3,2,NULL,NOW(),NOW()),
(107,'Delete Inventario-Operaciones','Eliminar Inventario-Operaciones',29,4,2,NULL,NOW(),NOW()),
-- Inventario-Slotting (30)
(108,'Create Inventario-Slotting','Crear en Inventario-Slotting',30,1,2,NULL,NOW(),NOW()),
(109,'Read Inventario-Slotting','Ver Inventario-Slotting',30,2,2,NULL,NOW(),NOW()),
(110,'Update Inventario-Slotting','Editar Inventario-Slotting',30,3,2,NULL,NOW(),NOW()),
(111,'Delete Inventario-Slotting','Eliminar Inventario-Slotting',30,4,2,NULL,NOW(),NOW()),
-- Inventario-Auditoria (31)
(112,'Create Inventario-Auditoria','Crear en Inventario-Auditoria',31,1,2,NULL,NOW(),NOW()),
(113,'Read Inventario-Auditoria','Ver Inventario-Auditoria',31,2,2,NULL,NOW(),NOW()),
(114,'Update Inventario-Auditoria','Editar Inventario-Auditoria',31,3,2,NULL,NOW(),NOW()),
(115,'Delete Inventario-Auditoria','Eliminar Inventario-Auditoria',31,4,2,NULL,NOW(),NOW()),
-- Inventario-LPN (32)
(116,'Create Inventario-LPN','Crear en Inventario-LPN',32,1,2,NULL,NOW(),NOW()),
(117,'Read Inventario-LPN','Ver Inventario-LPN',32,2,2,NULL,NOW(),NOW()),
(118,'Update Inventario-LPN','Editar Inventario-LPN',32,3,2,NULL,NOW(),NOW()),
(119,'Delete Inventario-LPN','Eliminar Inventario-LPN',32,4,2,NULL,NOW(),NOW()),
-- Inventario-Scan (33)
(120,'Create Inventario-Scan','Crear en Inventario-Scan',33,1,2,NULL,NOW(),NOW()),
(121,'Read Inventario-Scan','Ver Inventario-Scan',33,2,2,NULL,NOW(),NOW()),
(122,'Update Inventario-Scan','Editar Inventario-Scan',33,3,2,NULL,NOW(),NOW()),
(123,'Delete Inventario-Scan','Eliminar Inventario-Scan',33,4,2,NULL,NOW(),NOW()),
-- Inventario-Sync-Logs (34)
(124,'Create Inventario-Sync-Logs','Crear en Inventario-Sync-Logs',34,1,2,NULL,NOW(),NOW()),
(125,'Read Inventario-Sync-Logs','Ver Inventario-Sync-Logs',34,2,2,NULL,NOW(),NOW()),
(126,'Update Inventario-Sync-Logs','Editar Inventario-Sync-Logs',34,3,2,NULL,NOW(),NOW()),
(127,'Delete Inventario-Sync-Logs','Eliminar Inventario-Sync-Logs',34,4,2,NULL,NOW(),NOW());
```

## 7. role_permissions

- **role 1 (Super Admin)**: sin filas. Bypass por codigo.
- **role 2 (Operador)**: sin filas (asignar segun necesidad).
- **role 4 (Administrador)**: permisos 1..21, 38..85 (69 permisos: todo lo de negocio basico).
- **role 5 (cliente_final)**: permisos 86, 87 (solo Storefront Read/Create).

```sql
-- Administrador (role 4): 1-21 + 38-85
INSERT INTO role_permissions (role_id, permission_id)
SELECT 4, p FROM generate_series(1,21) p
UNION ALL SELECT 4, p FROM generate_series(38,85) p;

-- cliente_final (role 5)
INSERT INTO role_permissions (role_id, permission_id) VALUES (5,86),(5,87);
```

## 8. Secuencias

Tras los inserts con IDs explicitos, alinear las secuencias:

```sql
SELECT setval('scope_id_seq',         (SELECT MAX(id) FROM scope));
SELECT setval('action_id_seq',        (SELECT MAX(id) FROM action));
SELECT setval('resource_id_seq',      (SELECT MAX(id) FROM resource));
SELECT setval('business_type_id_seq', (SELECT MAX(id) FROM business_type));
SELECT setval('role_id_seq',          (SELECT MAX(id) FROM role));
SELECT setval('permission_id_seq',    (SELECT MAX(id) FROM permission));
```

## 9. Tablas relacionadas (no se siembran aqui)

- `user`, `user_roles`, `user_businesses`: dependen de los usuarios reales. Crear via API (`POST /api/v1/auth/register` + asignacion de roles) en lugar de insertar a mano.
- `business`: crear via API; cada empresa apunta a `business_type_id=1`.

## 10. Resumen rapido

- 1 business_type, 2 scopes, 13 actions, 34 resources.
- 4 roles activos (Super Admin platform-level, Administrador business-level, Operador platform-level, cliente_final business-level).
- 127 permisos (CRUD por recurso), de los cuales 69 estan asignados al Administrador y 2 al cliente_final.
- Super Admin no necesita permisos en DB; se autoriza por `scope=platform` + business_id=0 en JWT.
