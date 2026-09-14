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
    version: 2,
    title: 'Notificaciones',
    routes: ['/notification-config'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido a Notificaciones',
            body: 'Todo lo que le dices a tus clientes por WhatsApp en un solo lugar: avisos automáticos del pedido, conversaciones, plantillas, flujos de respuesta y campañas.',
        },
        {
            id: 'rules',
            title: 'Reglas',
            body: 'Qué evento del pedido dispara qué mensaje: confirmación, guía generada, en reparto, entregado. Aquí también programas envíos por segmento de clientes.',
            target: '[data-tour="notif-rules"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'conversations',
            title: 'Conversaciones',
            body: 'Tu bandeja estilo WhatsApp. Los chats con mensajes sin leer suben arriba con su contador, y en rojo aparece quien pidió dejar de recibir. Cada chat muestra el nombre del cliente, su pedido y su campaña.',
            target: '[data-tour="notif-tab-conversations"]',
            route: '/notification-config?tab=conversations',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'conversations-reply',
            title: 'Responder al cliente',
            body: 'Contesta con texto, imágenes o archivos. Puedes escribir libremente hasta 24 horas después del último mensaje del cliente; pasado ese tiempo Meta solo deja enviar plantillas aprobadas. Los chats se borran solos al año.',
            target: '[data-tour="notif-conversations-list"]',
            waitMs: 8000,
            placement: 'right',
            optional: true,
        },
        {
            id: 'templates',
            title: 'Plantillas WhatsApp',
            body: 'Los mensajes que Meta revisa y aprueba. Son los únicos con los que puedes escribirle primero a un cliente.',
            target: '[data-tour="notif-tab-templates"]',
            route: '/notification-config?tab=templates',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'templates-new',
            title: 'Crear y enviar a revisión',
            body: 'Crea la plantilla, guárdala como borrador y envíala a Meta. Te avisamos antes de guardar lo que Meta rechaza, como una variable al final del texto o emojis en el encabezado. Si una sigue en revisión, "Consultar estado en Meta" la actualiza, y además lo revisamos solos cada 10 minutos.',
            target: '[data-tour="notif-new-template"]',
            waitMs: 8000,
            placement: 'left',
            optional: true,
        },
        {
            id: 'flows',
            title: 'Flujos',
            body: 'Qué responde automáticamente cada botón. Por ejemplo: el cliente toca "Opción A" y le llega al instante la plantilla del camino A.',
            target: '[data-tour="notif-tab-flows"]',
            route: '/notification-config?tab=flows',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'flows-new',
            title: 'Armar un flujo',
            body: 'Parte de una plantilla inicial y crea las respuestas ahí mismo con "+ Crear respuesta". Para usar el flujo en una campaña, todas sus respuestas tienen que estar aprobadas por Meta.',
            target: '[data-tour="notif-new-flow"]',
            waitMs: 8000,
            placement: 'left',
            optional: true,
        },
        {
            id: 'campaigns',
            title: 'Campañas',
            body: 'Mensajes de marketing a muchos clientes a la vez. Salen desde el número propio de tu negocio y nunca les llegan a quienes pidieron dejar de recibir.',
            target: '[data-tour="notif-tab-campaigns"]',
            route: '/notification-config?tab=campaigns',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'campaigns-new',
            title: 'Crear una campaña',
            body: 'Elige la plantilla o el flujo, a quién le llega (todos, con filtros o a mano), el horario y el tope por día. Puedes enviarla todos los días, cada cierto número de días o en fechas del calendario, y pausarla o cancelarla cuando quieras.',
            target: '[data-tour="notif-new-campaign"]',
            waitMs: 8000,
            placement: 'left',
            optional: true,
        },
        {
            id: 'audit',
            title: 'Auditoría',
            body: 'El registro de cada mensaje: enviado, entregado, leído o fallido, con el motivo cuando algo sale mal.',
            target: '[data-tour="notif-tab-audit"]',
            route: '/notification-config?tab=audit',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'done',
            title: 'Listo',
            body: 'Ya conoces todo el módulo. Puedes repetir este recorrido cuando quieras desde el botón Tutorial.',
            route: '/notification-config?tab=conversations',
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
