import type { TourDefinition } from '../domain/types';

export const productsTour: TourDefinition = {
    key: 'products',
    version: 2,
    title: 'Productos',
    routes: ['/products'],
    resource: 'Productos',
    autoStart: true,
    legacyStorageKey: 'products_tour_seen_v1',
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido al modulo de Productos',
            body: 'Este tutorial te guia por las funciones principales. Avanza con las flechas del teclado o cierra con la X en cualquier momento.',
        },
        {
            id: 'tabs',
            title: 'SKUs y familias',
            body: 'Aqui cambias entre las dos vistas. "SKUs / Productos" es el catalogo real: cada fila es un codigo unico con su precio, stock e integraciones.\n\n"Familias de variantes" solo agrupa SKUs del mismo producto base (color, talla, sabor); no tiene stock ni precio propio.',
            target: '[data-tour="products.tabs"]',
            placement: 'bottom',
        },
        {
            id: 'search',
            title: 'Buscar en el catalogo',
            body: 'Filtra por nombre o por codigo SKU. Los dos campos se combinan, asi que puedes acotar mucho un catalogo grande.',
            target: '[data-tour="products.search"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'integration-filter',
            title: 'Filtrar por canal',
            body: 'Muestra solo los SKUs conectados a una integracion concreta (Shopify, MercadoLibre, WooCommerce).\n\nSirve para revisar que no quede producto sin mapear: un producto sin integracion no puede recibir pedidos de ese canal.',
            target: '[data-tour="products.integration-filter"]',
            placement: 'left',
            optional: true,
        },
        {
            id: 'create',
            title: 'Crear un SKU',
            body: 'Abre el formulario de producto: codigo, nombre, precio, stock, imagen y familia.\n\nRegla practica: si cambia el precio, el stock o el empaque, es un SKU distinto.',
            target: '[data-tour="products.create"]',
            placement: 'left',
        },
    ],
};
