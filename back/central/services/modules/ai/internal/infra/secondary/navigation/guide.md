## home
Tablero con los indicadores del negocio: «Órdenes Totales», «Órdenes del Día», «Órdenes Mensuales», pronóstico de órdenes, días de mayor demanda, órdenes por transportadora, mapa de órdenes, top de departamentos, productos y clientes. Se puede filtrar por período.

## orders
Lista de órdenes de todos los canales (Shopify, WooCommerce, MercadoLibre, WhatsApp y las creadas a mano). Filtros: nombre o correo del cliente, estado, plataforma, contra entrega, estado de pago, de fulfillment y de factura, y rango de fechas.
Pestañas de esta sección: «Órdenes», «Envíos», «Recaudo contra entrega» y «Cotizaciones».
- Crear una orden a mano: botón «Nueva orden» (arriba), llenar cliente, dirección, envío, contra entrega, pago, productos y bodega.
- Cargar muchas órdenes: «Carga masiva de órdenes», con archivo CSV o Excel.
- Generar una guía de envío: en la fila de la orden, botón «Ver recomendación IA» (ícono de robot). Se abre «Generar Guía de Envío» con 4 pasos: Origen y Destino, Cotización (elegir transportadora), Detalles y Pago (se paga con el saldo de la billetera). Al terminar se puede descargar el PDF de la guía. Si la orden ya tiene guía, el botón aparece desactivado.
- Contra entrega: se marca en la orden, en la sección de contra entrega del formulario al crearla («Nueva orden») o al editarla («Editar orden»). La guía toma ese dato; no se activa sola.
- Desde el detalle de la orden también está «Cotizar y Generar Guía».
- Varias guías a la vez: «Generación masiva de guías», eligiendo la bodega de origen.
- Otras acciones por fila: «Ver orden», «Editar orden», «Cambiar estado», «Ver guía de envío», «Eliminar orden». También hay «Descargar órdenes en Excel» y «Cotizador Expres».

## products
Catálogo de productos con pestañas «SKUs / Productos» y «Familias de variantes». Botón «Crear» para un producto nuevo, carga y descarga de peso y dimensiones en Excel, y «Grupos de clientes y precios por catalogo». Filtros por integración, estado, categoría, marca, precio, stock y peso. En cada producto: «Ver detalle», «Editar», «Integraciones» y «Eliminar».

## shipments
Lista de envíos con su guía y estado (Pendiente, Recolectado, En Tránsito, En Reparto, Novedad, Devuelto, Cancelado). Se busca por cliente, número de orden o número de guía. El detalle muestra el «Historial de rastreo» y permite cancelar la guía. Aquí también está «Configuración de envíos».

## shipments.cod
Informe de recaudo contra entrega: cuánto recaudaron las transportadoras y cuánto le corresponde al negocio. Sub-pestañas «Resumen», «Ordenes» y «Cortes de pago», con períodos rápidos «Esta semana», «Semana pasada» y «Este mes».

## shipments.quotes
Cotizaciones de envío hechas desde Shopify, WooCommerce o el panel. Cada cotización tiene «Crear orden», que permite crear solo la orden o «Crear orden + guía».

## customers
Clientes del negocio con botón «Nuevo cliente». El detalle de cada cliente muestra total de órdenes, entregadas, score de entrega, total gastado, ticket promedio, productos comprados y direcciones.

## delivery
Última milla con repartidores propios. Pestañas «Rutas», «Conductores», «Vehículos» y «Geozonas». En Rutas está «Nueva ruta»: se agregan paradas, se ordenan por la ruta más corta y cada parada se marca como entregada o fallida.

## delivery.drivers
Conductores propios; botón «Nuevo conductor».

## delivery.vehicles
Vehículos propios (moto, carro, van, camión); botón «Nuevo vehiculo».

## delivery.geozones
Zonas de cobertura dibujadas en el mapa; botón «Nueva geozona».

## inventory
Stock por bodega. Botones «Añadir stock», «Ajustar stock» y «Transferir stock». La sección tiene pestañas para Productos, Bodegas, Stock, Movimientos, Trazabilidad, Kardex, Operaciones, Slotting ABC, Auditoría, LPN, Scan y Sync Logs (algunas dependen del rol).

## inventory.movements
Historial de entradas y salidas de inventario, por producto, bodega o familia.

## warehouses
Bodegas del negocio con botón «Nueva bodega». Cada bodega guarda dirección, contacto y datos para la transportadora (la dirección de origen de las guías sale de aquí). Dentro de una bodega se ve la jerarquía física (zonas, pasillos, racks, niveles, posiciones) y un plano 2D.

## invoicing
Facturas electrónicas. Filtros por estado (Emitida, Pendiente, Cancelada, Fallida), contra entrega y fechas. Botones «Crear Facturas» para facturar varias órdenes, «Configuraciones de facturación», «Comparar con proveedor» y «Sincronizar facturas anuladas».

## invoicing.configs
Configuración de la facturación electrónica: con qué proveedor se factura cada integración y si la facturación es automática. Botón «Nueva Configuración de Facturación».

## wallet
Billetera del negocio: saldo e «Historial de Transacciones». Cada guía generada se descuenta de aquí. Para recargar se abre «¿Cómo quieres pagar?» y se indica el monto a recargar.

## wallet.finanzas
Resumen financiero del negocio y utilidad de los envíos por mes.

## notifications
Mensajes automáticos a los clientes. Tiene «Reglas de Notificación» y las pestañas «Conversaciones» (chats de WhatsApp), «Auditoria», «Campañas», «Plantillas WhatsApp» y «Flujos».

## storefront
Tienda B2B para que los clientes del negocio hagan pedidos: catálogo con buscador y filtros, «Agregar al carrito» y envío del pedido desde el carrito. También «Mis Pedidos» y «Clientes».

## integrations
Conexión con otros sistemas. Botón «Crear Integración»: se elige la categoría, luego el proveedor y se llenan los datos (Shopify, MercadoLibre, Jumpseller y Tiendanube se conectan iniciando sesión en su cuenta; WhatsApp con «Conectar WhatsApp»).
Categorías: tiendas en línea (Shopify, WooCommerce, MercadoLibre, Tiendanube, Jumpseller, VTEX, Magento, Amazon, Falabella, Éxito, TikTok), facturación electrónica (Siigo, Alegra, Factus, Helisa, Softpymes, World Office), mensajería (WhatsApp), pagos (Bold, ePayco, Nequi, PayU, Stripe, Wompi) y transportadoras (EnvioClick, MiPaquete, Enviame, Shipit). Algunas aparecen como «Próximamente».
En cada integración: «Probar conexión», «Sincronizar órdenes» (con rango de fechas), «Editar integración» y «Eliminar integración».

## users
Equipo del negocio: «Crear usuario», «Asignar rol», «Editar usuario» y «Eliminar usuario».

## roles
Roles y lo que cada uno puede hacer: «Crear rol», «Asignar permisos», «Editar rol», «Eliminar rol».

## permissions
Listado de permisos del sistema.

## businesses
Datos del negocio: nombre, datos fiscales y colores de la marca.

## website_config
Editor «Disena tu tienda» del sitio web del negocio: plantilla, tema de colores y secciones que se pueden ordenar u ocultar.

## subscription
Estado del plan de Probability (Activo, Vencido, Suspendido, Pendiente) y sus pagos.

## profile
Datos de contacto de la cuenta, cambio de contraseña en «Seguridad de Cuenta» y ajustes de los tutoriales.

## tickets
Tickets de soporte y desarrollo de la plataforma, con botón «Nuevo ticket».

## accounting
Contabilidad de la plataforma: reporte, «Movimientos», «Facturas» y «Configuración».

## announcements
Anuncios que se muestran a los negocios y sus estadísticas.

## commercial
Prospectos y seguimiento comercial.

## shipping_margins
Márgenes por transportadora: guías generadas, cobrado al cliente, costo real y ganancia.
