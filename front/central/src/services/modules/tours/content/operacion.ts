import type { TourDefinition } from '../domain/types';

export const customersTour: TourDefinition = {
    key: 'customers',
    version: 1,
    title: 'Clientes',
    routes: ['/customers'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido al módulo de Clientes',
            body: 'La base de clientes se arma sola con las órdenes: cada vez que llega un pedido, el cliente se crea o se actualiza.',
        },
        {
            id: 'historial',
            title: 'Historial de compras',
            body: 'Entra a un cliente para ver todo lo que ha comprado, sus direcciones y sus datos de contacto. Sirve para atención y para ventas repetidas.',
            target: 'table',
            placement: 'top',
            optional: true,
        },
        {
            id: 'duplicados',
            title: 'Cuidado con los duplicados',
            body: 'Un mismo cliente que compra por dos canales distintos con correos diferentes entra como dos registros. Revisa por teléfono o documento antes de crear uno a mano.',
        },
    ],
};

export const integrationsTour: TourDefinition = {
    key: 'integrations',
    version: 1,
    title: 'Integraciones',
    routes: ['/integrations'],
    resource: 'Integraciones',
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido a Integraciones',
            body: 'Aquí conectas Probability con tus canales de venta y tus transportadoras. Es lo primero que hay que configurar en una cuenta nueva.',
        },
        {
            id: 'canales',
            title: 'Canales de venta',
            body: 'Shopify, MercadoLibre, WooCommerce, Amazon, WhatsApp. Al conectar uno, sus pedidos empiezan a entrar solos al módulo de Órdenes.',
        },
        {
            id: 'productos',
            title: 'Mapeo de productos',
            body: 'Cada producto del canal se asocia a un SKU tuyo. Sin ese mapeo el pedido llega, pero el producto queda sin identificar y no descuenta stock.',
        },
        {
            id: 'transportadoras',
            title: 'Transportadoras',
            body: 'Las credenciales con las que se cotiza y se generan guías. Si una guía falla, el detalle del error queda en los logs de sincronización de la integración, no en el mensaje de pantalla.',
        },
    ],
};

export const deliveryTour: TourDefinition = {
    key: 'delivery',
    version: 1,
    title: 'Domicilios',
    routes: ['/delivery'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido a Domicilios',
            body: 'Este módulo es para tu flota propia: lo que entregas tu mismo en vez de mandarlo por transportadora.',
        },
        {
            id: 'routes',
            title: 'Rutas',
            body: 'Agrupa varias entregas en un recorrido y se lo asignas a un conductor. Es la vista del día a día.',
            target: 'a[href="/delivery/routes"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'drivers',
            title: 'Conductores',
            body: 'Quien reparte. Cada conductor entra a la app móvil con su usuario y va marcando las entregas.',
            target: 'a[href="/delivery/drivers"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'vehicles',
            title: 'Vehículos',
            body: 'La flota y su capacidad. La capacidad limita cuántas paradas caben en una ruta.',
            target: 'a[href="/delivery/vehicles"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'geozones',
            title: 'Geozonas',
            body: 'Las zonas que cubres dibujadas en el mapa. Definen que direcciones puedes atender con flota propia y cuales tienen que ir por transportadora.',
            target: 'a[href="/delivery/geozones"]',
            placement: 'bottom',
            optional: true,
        },
    ],
};

export const notificationsTour: TourDefinition = {
    key: 'notifications',
    version: 1,
    title: 'Notificaciones',
    routes: ['/notification-config', '/notification-channels', '/notification-event-types'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido a Notificaciones',
            body: 'Controla que se le avisa al cliente sobre su pedido, por que canal y con que mensaje.',
        },
        {
            id: 'configs',
            title: 'Configuraciones',
            body: 'La combinación de evento más canal más plantilla. Por ejemplo: "orden despachada" por WhatsApp con la plantilla de guía.',
            target: 'a[href="/notification-config"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'channels',
            title: 'Canales',
            body: 'Por donde sale el mensaje: WhatsApp o correo. Cada canal necesita su integración configurada; sin ella el mensaje se descarta en silencio.',
            target: 'a[href="/notification-channels"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'events',
            title: 'Tipos de evento',
            body: 'Los momentos que disparan un aviso: orden creada, pago confirmado, guía generada, en reparto, entregada, novedad.',
            target: 'a[href="/notification-event-types"]',
            placement: 'bottom',
            optional: true,
        },
    ],
};

export const storefrontTour: TourDefinition = {
    key: 'storefront',
    version: 1,
    title: 'Mi sitio web',
    routes: ['/website-config', '/storefront'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Tu tienda publica',
            body: 'Desde aquí configuras el sitio que ven tus clientes: catalogo, colores, logo, dominio y contenido.',
        },
        {
            id: 'catalogo',
            title: 'Que se publica',
            body: 'El catalogo sale de tus productos activos. Un SKU inactivo o sin imagen no se muestra en la tienda.',
        },
        {
            id: 'pedidos',
            title: 'Los pedidos entran como órdenes',
            body: 'Una compra en tu tienda crea una orden normal, igual que si viniera de Shopify. Se factura, se despacha y se notifica con el mismo flujo.',
        },
    ],
};
