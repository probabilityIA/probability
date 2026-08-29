import type { TourDefinition } from '../domain/types';

export const shipmentsTour: TourDefinition = {
    key: 'shipments',
    version: 1,
    title: 'Envios',
    routes: ['/shipments'],
    resource: 'Envios',
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido al módulo de Envíos',
            body: 'Aquí vive el seguimiento de cada guía después de generarla: estado en la transportadora, novedades y entregas.',
        },
        {
            id: 'orders',
            title: 'De donde salen los envíos',
            body: 'Un envío siempre nace de una orden. Si no ves una guía aquí, es porque la orden todavía no la genero.',
            target: 'a[href="/orders"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'cod',
            title: 'Recaudo contra entrega',
            body: 'Si vendes con pago contra entrega, el dinero que recoge el mensajero se concilia en esta pestana: cuanto se recaudo, cuanto liquido la transportadora y que queda pendiente.',
            target: 'a[href="/shipments/cod"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'quotes',
            title: 'Cotizaciones',
            body: 'Historial de tarifas consultadas. Cada cotización guarda que transportadoras respondieron y a que precio, para poder auditar por que se eligio una.',
            target: 'a[href="/shipments/quotes"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'filters',
            title: 'Filtros',
            body: 'Acota por transportadora, estado del envío, ciudad de destino o rango de fechas.',
            target: '#orders-filters-slot',
            placement: 'bottom',
            optional: true,
        },
    ],
};

export const shipmentsCodTour: TourDefinition = {
    key: 'shipments-cod',
    version: 1,
    title: 'Recaudo contra entrega',
    routes: ['/shipments/cod'],
    resource: 'Envios',
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Recaudo contra entrega',
            body: 'Aquí se concilia la plata del pago contra entrega: lo que el mensajero le cobra al cliente, la comisión de la transportadora y lo que efectivamente te queda a ti.',
        },
        {
            id: 'concepto',
            title: 'Los tres números que importan',
            body: 'Valor a recaudar: producto + envío + comisión, lo que paga el cliente.\n\nNeto: lo que te queda a ti, sin la comisión.\n\nComision del carrier: lo que se queda la transportadora por manejar el efectivo.',
        },
        {
            id: 'cortes',
            title: 'Cortes de pago',
            body: 'La transportadora no te gira orden por orden: agrupa varias en un corte. Hasta que el corte no se confirma, la orden queda como pendiente de liquidación.',
            target: 'table',
            placement: 'top',
            optional: true,
        },
    ],
};

export const shipmentsQuotesTour: TourDefinition = {
    key: 'shipments-quotes',
    version: 1,
    title: 'Cotizaciones',
    routes: ['/shipments/quotes'],
    resource: 'Envios',
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Cotizaciones de envío',
            body: 'Cada vez que se pide tarifa para una orden queda registrada aquí, con todas las transportadoras que respondieron y su precio.',
        },
        {
            id: 'formula',
            title: 'Como se arma el precio',
            body: 'El costo de la guía es flete + seguro mínimo + seguro extra si asegura, más el margen de contra entrega cuando aplica.\n\nEl total suma además la comisión de la transportadora si el envío es contra entrega.',
        },
        {
            id: 'tabla',
            title: 'Historial',
            body: 'Sirve para auditar: si un cliente reclama que le cobraron distinto, aquí esta que se le mostro y cuando.',
            target: 'table',
            placement: 'top',
            optional: true,
        },
    ],
};

export const shippingMarginsTour: TourDefinition = {
    key: 'shipping-margins',
    version: 1,
    title: 'Márgenes de envío',
    routes: ['/shipping-margins'],
    autoStart: true,
    superAdminOnly: true,
    steps: [
        {
            id: 'welcome',
            title: 'Márgenes de envío',
            body: 'Configuración de plataforma: cuanto se le suma a la tarifa de la transportadora antes de mostrarsela al negocio.',
        },
        {
            id: 'cod',
            title: 'Margen de contra entrega',
            body: 'Monto fijo por negocio que se agrega sobre la comisión del carrier en los envíos contra entrega. Se aplica solo si la tarifa soporta contra entrega.',
        },
        {
            id: 'cuidado',
            title: 'Cuidado al tocarlo',
            body: 'Un cambio aquí afecta el precio que ve el negocio en el cotizador y el valor que se le declara a la transportadora. Verifica siempre con una cotización real después de cambiarlo.',
        },
    ],
};
