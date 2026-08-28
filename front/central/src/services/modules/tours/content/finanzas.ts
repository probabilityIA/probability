import type { TourDefinition } from '../domain/types';

export const walletTour: TourDefinition = {
    key: 'wallet',
    version: 1,
    title: 'Billetera',
    routes: ['/wallet'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido a tu Billetera',
            body: 'Es el saldo con el que se pagan las guias. Cada envio que generas descuenta de aqui.',
        },
        {
            id: 'saldos',
            title: 'Saldos',
            body: 'Cuanto tienes disponible ahora mismo y el detalle de cada movimiento: recargas, cobros de guias y devoluciones.',
            target: 'a[href="/wallet/saldos"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'finanzas',
            title: 'Finanzas',
            body: 'La vista de analisis: cuanto gastas en envios, en que transportadoras y como evoluciona mes a mes.',
            target: 'a[href="/wallet/finanzas"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'aviso',
            title: 'Si el saldo llega a cero',
            body: 'No se pueden generar mas guias hasta recargar. Vale la pena dejar configurada la alerta de saldo bajo para no quedar tirado a mitad de un dia de despachos.',
        },
    ],
};

export const accountingTour: TourDefinition = {
    key: 'accounting',
    version: 1,
    title: 'Contabilidad',
    routes: ['/accounting'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Contabilidad',
            body: 'La vista de plataforma sobre el dinero: que se facturo, que se cobro y que quedo pendiente.',
        },
        {
            id: 'movimientos',
            title: 'Movimientos',
            body: 'El detalle linea por linea de cada ingreso y egreso registrado.',
            target: 'a[href="/accounting/movimientos"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'facturas',
            title: 'Facturas',
            body: 'Las facturas emitidas a los negocios por el uso de la plataforma.',
            target: 'a[href="/accounting/facturas"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'config',
            title: 'Configuracion',
            body: 'Conceptos, tarifas y parametros con los que se calculan los cobros.',
            target: 'a[href="/accounting/configuracion"]',
            placement: 'bottom',
            optional: true,
        },
    ],
};

export const invoicingTour: TourDefinition = {
    key: 'invoicing',
    version: 1,
    title: 'Facturacion',
    routes: ['/invoicing'],
    resource: 'Facturacion',
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido a Facturacion',
            body: 'Desde aqui se emiten las facturas electronicas de tus ventas hacia la DIAN, a traves de tu proveedor tecnologico.',
        },
        {
            id: 'facturas',
            title: 'Facturas',
            body: 'El listado de lo emitido, con su estado en la DIAN. Una factura rechazada se ve aqui con el motivo del rechazo.',
            target: 'a[href="/invoicing/invoices"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'flujo',
            title: 'De donde sale una factura',
            body: 'Se genera a partir de una orden: cliente, items, impuestos y totales salen de ahi.\n\nSi un dato del cliente esta incompleto (NIT, tipo de documento, correo), la DIAN la rechaza.',
        },
        {
            id: 'config',
            title: 'Antes de facturar',
            body: 'Hay que configurar el proveedor y la resolucion de facturacion. Sin eso, el modulo no puede emitir nada.',
        },
    ],
};

export const subscriptionTour: TourDefinition = {
    key: 'subscription',
    version: 1,
    title: 'Suscripcion',
    routes: ['/subscription'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Tu suscripcion',
            body: 'Que plan tienes, que modulos incluye y hasta cuando esta vigente.',
        },
        {
            id: 'cuotas',
            title: 'Cuotas y excedentes',
            body: 'Los planes traen una cuota incluida de envios, ordenes o facturas. Al pasarte se cobra un valor por unidad adicional, que se ve reflejado aqui.',
        },
    ],
};
