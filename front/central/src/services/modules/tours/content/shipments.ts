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
            title: 'Bienvenido al modulo de Envios',
            body: 'Aqui vive el seguimiento de cada guia despues de generarla: estado en la transportadora, novedades y entregas.',
        },
        {
            id: 'orders',
            title: 'De donde salen los envios',
            body: 'Un envio siempre nace de una orden. Si no ves una guia aqui, es porque la orden todavia no la genero.',
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
            body: 'Historial de tarifas consultadas. Cada cotizacion guarda que transportadoras respondieron y a que precio, para poder auditar por que se eligio una.',
            target: 'a[href="/shipments/quotes"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'filters',
            title: 'Filtros',
            body: 'Acota por transportadora, estado del envio, ciudad de destino o rango de fechas.',
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
            body: 'Aqui se concilia la plata del pago contra entrega: lo que el mensajero le cobra al cliente, la comision de la transportadora y lo que efectivamente te queda a ti.',
        },
        {
            id: 'concepto',
            title: 'Los tres numeros que importan',
            body: 'Valor a recaudar: producto + envio + comision, lo que paga el cliente.\n\nNeto: lo que te queda a ti, sin la comision.\n\nComision del carrier: lo que se queda la transportadora por manejar el efectivo.',
        },
        {
            id: 'cortes',
            title: 'Cortes de pago',
            body: 'La transportadora no te gira orden por orden: agrupa varias en un corte. Hasta que el corte no se confirma, la orden queda como pendiente de liquidacion.',
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
            title: 'Cotizaciones de envio',
            body: 'Cada vez que se pide tarifa para una orden queda registrada aqui, con todas las transportadoras que respondieron y su precio.',
        },
        {
            id: 'formula',
            title: 'Como se arma el precio',
            body: 'El costo de la guia es flete + seguro minimo + seguro extra si asegura, mas el margen de contra entrega cuando aplica.\n\nEl total suma ademas la comision de la transportadora si el envio es contra entrega.',
        },
        {
            id: 'tabla',
            title: 'Historial',
            body: 'Sirve para auditar: si un cliente reclama que le cobraron distinto, aqui esta que se le mostro y cuando.',
            target: 'table',
            placement: 'top',
            optional: true,
        },
    ],
};

export const shippingMarginsTour: TourDefinition = {
    key: 'shipping-margins',
    version: 1,
    title: 'Margenes de envio',
    routes: ['/shipping-margins'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Margenes de envio',
            body: 'Configuracion de plataforma: cuanto se le suma a la tarifa de la transportadora antes de mostrarsela al negocio.',
        },
        {
            id: 'cod',
            title: 'Margen de contra entrega',
            body: 'Monto fijo por negocio que se agrega sobre la comision del carrier en los envios contra entrega. Se aplica solo si la tarifa soporta contra entrega.',
        },
        {
            id: 'cuidado',
            title: 'Cuidado al tocarlo',
            body: 'Un cambio aqui afecta el precio que ve el negocio en el cotizador y el valor que se le declara a la transportadora. Verifica siempre con una cotizacion real despues de cambiarlo.',
        },
    ],
};
