import type { TourDefinition } from '../domain/types';

export const inventoryTour: TourDefinition = {
    key: 'inventory',
    version: 1,
    title: 'Inventario',
    routes: ['/inventory'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bienvenido al modulo de Inventario',
            body: 'Este es el stock real por bodega. Todo lo que ves aqui alimenta lo que se puede vender y lo que se puede despachar.',
        },
        {
            id: 'products',
            title: 'Productos',
            body: 'El catalogo: que vendes, a que precio y con que codigo. El inventario cuenta unidades de esos SKUs.',
            target: 'a[href="/products"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'warehouses',
            title: 'Bodegas',
            body: 'Donde esta guardada la mercancia. Un mismo SKU puede tener stock en varias bodegas, y cada bodega tiene sus ubicaciones internas.',
            target: 'a[href="/warehouses"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'movements',
            title: 'Movimientos',
            body: 'Toda entrada y salida queda registrada: recibos, despachos, ajustes y traslados. Si el stock no cuadra, la respuesta esta aqui.',
            target: 'a[href="/inventory/movements"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'kardex',
            title: 'Kardex',
            body: 'La historia completa de un SKU: saldo inicial, cada movimiento y saldo final. Es el reporte contable del inventario.',
            target: 'a[href="/inventory/kardex"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'operations',
            title: 'Operaciones',
            body: 'Las tareas de bodega del dia: recibos por procesar, picking y packing pendientes.',
            target: 'a[href="/inventory/operations"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'audit',
            title: 'Auditoria',
            body: 'Conteos fisicos: cuentas lo que hay en el estante y el sistema te muestra la diferencia contra lo que deberia haber.',
            target: 'a[href="/inventory/audit"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'lpn',
            title: 'LPN y Scan',
            body: 'LPN identifica una estiba o caja completa con un solo codigo. Scan es la vista para el operario con lector de codigo de barras.',
            target: 'a[href="/inventory/lpn"]',
            placement: 'bottom',
            optional: true,
        },
    ],
};

export const warehousesTour: TourDefinition = {
    key: 'warehouses',
    version: 1,
    title: 'Bodegas',
    routes: ['/warehouses'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Bodegas',
            body: 'Cada bodega es un lugar fisico con su direccion y su contacto. El stock siempre vive en una bodega, nunca "en el aire".',
        },
        {
            id: 'contacto',
            title: 'Los datos de contacto importan',
            body: 'La direccion, el telefono y el correo de la bodega se usan como origen al generar una guia.\n\nUna bodega sin contacto valido hace fallar la generacion de guia aunque la cotizacion haya funcionado.',
        },
        {
            id: 'ubicaciones',
            title: 'Ubicaciones internas',
            body: 'Entra a una bodega para definir su estructura: pasillos, estantes y posiciones. Eso es lo que permite el picking guiado y el LPN.',
            target: 'table',
            placement: 'top',
            optional: true,
        },
    ],
};
