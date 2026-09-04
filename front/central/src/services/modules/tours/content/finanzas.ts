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
            body: 'Es el saldo con el que se pagan las guías. Cada envío que generas descuenta de aquí.',
        },
        {
            id: 'saldos',
            title: 'Saldos',
            body: 'Cuanto tienes disponible ahora mismo y el detalle de cada movimiento: recargas, cobros de guías y devoluciones.',
            target: 'a[href="/wallet/saldos"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'finanzas',
            superAdminOnly: true,
            title: 'Finanzas',
            body: 'La vista de análisis: cuanto gastas en envíos, en que transportadoras y como evoluciona mes a mes.',
            target: 'a[href="/wallet/finanzas"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'aviso',
            title: 'Si el saldo llega a cero',
            body: 'No se pueden generar más guías hasta recargar. Vale la pena dejar configurada la alerta de saldo bajo para no quedar tirado a mitad de un día de despachos.',
        },
    ],
};

export const accountingTour: TourDefinition = {
    key: 'accounting',
    version: 1,
    title: 'Contabilidad',
    routes: ['/accounting'],
    autoStart: true,
    superAdminOnly: true,
    steps: [
        {
            id: 'welcome',
            title: 'Contabilidad',
            body: 'La vista de plataforma sobre el dinero: que se facturo, que se cobro y que quedo pendiente.',
        },
        {
            id: 'movimientos',
            title: 'Movimientos',
            body: 'El detalle línea por línea de cada ingreso y egreso registrado.',
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
            title: 'Configuración',
            body: 'Conceptos, tarifas y parámetros con los que se calculan los cobros.',
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
            title: 'Bienvenido a Facturación',
            body: 'Desde aquí se emiten las facturas electrónicas de tus ventas hacia la DIAN, a través de tu proveedor tecnológico.',
        },
        {
            id: 'facturas',
            title: 'Facturas',
            body: 'El listado de lo emitido, con su estado en la DIAN. Una factura rechazada se ve aquí con el motivo del rechazo.',
            target: 'a[href="/invoicing/invoices"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'flujo',
            title: 'De donde sale una factura',
            body: 'Se genera a partir de una orden: cliente, items, impuestos y totales salen de ahí.\n\nSi un dato del cliente esta incompleto (NIT, tipo de documento, correo), la DIAN la rechaza.',
        },
        {
            id: 'config',
            title: 'Antes de facturar',
            body: 'Hay que configurar el proveedor y la resolución de facturación. Sin eso, el módulo no puede emitir nada.',
        },
    ],
};

export const subscriptionTour: TourDefinition = {
    key: 'subscription',
    version: 2,
    title: 'Suscripción',
    routes: ['/subscription'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Tu suscripción',
            body: 'Que plan tienes, que módulos incluye y hasta cuando esta vigente.',
        },
        {
            id: 'periodo-cobro',
            title: 'Periodo de cobro',
            body: 'El periodo facturado dura un mes. Al vencer, tienes unos días de gracia para pagar sin que se suspenda la cuenta: aquí te decimos exactamente entre que fechas puedes pagar.',
            target: '#subscription-payment-window',
            placement: 'top',
        },
        {
            id: 'pagar-extender',
            title: 'Pagar / Extender',
            body: 'Este botón solo se activa cuando empieza el periodo de cobro (el día que vence el ciclo actual). No se puede adelantar el pago antes de esa fecha.',
            target: '#subscription-pay-button',
            placement: 'left',
        },
        {
            id: 'autopago',
            title: 'Pago automático',
            body: 'Actívalo para que la suscripción se pague sola desde tu billetera el día que vence, sin que tengas que acordarte de hacerlo manual.',
            target: '#subscription-auto-payment-toggle',
            placement: 'top',
        },
        {
            id: 'cuotas',
            title: 'Cuotas y excedentes',
            body: 'Los planes traen una cuota incluida de envíos, órdenes o facturas. Al pasarte se cobra un valor por unidad adicional, que se ve reflejado aquí.',
        },
    ],
};
