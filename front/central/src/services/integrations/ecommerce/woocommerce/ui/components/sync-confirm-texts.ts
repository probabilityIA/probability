export const CONFIRM_TEXTS = {
    stock: {
        title: 'Sincronizar stock en WooCommerce',
        description:
            'Se actualizará la cantidad disponible en WooCommerce de los productos que ya están vinculados a este canal. No se creará ningún producto nuevo en la tienda.',
        countLabel: 'productos vinculados a los que se les actualizará el stock',
        warning:
            'Los productos de Probability que no existan en la tienda se omiten y se reportan al final. Para crearlos hay que usar la acción "Crear en WooCommerce".',
        confirmText: 'Actualizar stock',
    },
    createInWoo: {
        title: 'Crear productos en tu tienda WooCommerce',
        description:
            'Se crearán productos NUEVOS y publicados en tu tienda WooCommerce, uno por cada producto de Probability que no exista allá. Quedarán visibles para tus clientes.',
        countLabel: 'productos nuevos que se crearán en la tienda',
        warning:
            'Esta acción no se puede deshacer desde Probability: para quitarlos habría que borrarlos a mano en WooCommerce. Revisa el listado antes de continuar.',
        confirmText: 'Sí, crear en la tienda',
    },
    createInProbability: {
        title: 'Crear productos en Probability',
        description:
            'Se crearán en tu catálogo de Probability los productos que existen en WooCommerce y aún no están aquí. No se modifica nada en la tienda.',
        countLabel: 'productos que se crearán en Probability',
        confirmText: 'Crear en Probability',
    },
    associate: {
        title: 'Asociar productos a este canal',
        description:
            'Se vincularán por SKU los productos que ya existen en los dos lados, para que WooCommerce los reconozca como propios. No se crea ni se modifica ningún producto, y no se toca el stock.',
        countLabel: 'productos que se vincularán al canal',
        confirmText: 'Asociar',
    },
} as const;
