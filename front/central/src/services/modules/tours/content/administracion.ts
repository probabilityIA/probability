import type { TourDefinition } from '../domain/types';

export const iamTour: TourDefinition = {
    key: 'iam',
    version: 1,
    title: 'Usuarios y permisos',
    routes: ['/users', '/roles', '/permissions', '/resources', '/businesses'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Usuarios y permisos',
            body: 'El control de quien entra y que puede hacer. Se arma en cuatro capas: recursos, permisos, roles y usuarios.',
        },
        {
            id: 'businesses',
            superAdminOnly: true,
            title: 'Empresas',
            body: 'Cada negocio es un inquilino aislado: sus órdenes, productos y clientes no se cruzan con los de otro.',
            target: 'a[href="/businesses"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'users',
            title: 'Usuarios',
            body: 'Las personas que entran a la plataforma. Un usuario pertenece a un negocio y recibe uno o varios roles.',
            target: 'a[href="/users"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'roles',
            title: 'Roles',
            body: 'Un paquete de permisos con nombre: Administrador, Bodeguero, Vendedor. Es lo que se le asigna a la persona.',
            target: 'a[href="/roles"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'permissions',
            title: 'Permisos',
            body: 'La pieza atómica: recurso más acción. Por ejemplo Órdenes + Crear, o Productos + Eliminar.',
            target: 'a[href="/permissions"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'resources',
            title: 'Recursos',
            body: 'Los módulos sobre los que se pueden dar permisos. Un recurso nuevo aparece aquí antes de poder incluirlo en un rol.',
            target: 'a[href="/resources"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'orden',
            title: 'El orden correcto',
            body: 'Crea el recurso, luego sus permisos, luego el rol que los agrupa, y al final asignale el rol al usuario. Al revés no funciona.',
        },
    ],
};

export const ticketsTour: TourDefinition = {
    key: 'tickets',
    version: 1,
    title: 'Soporte',
    routes: ['/tickets'],
    autoStart: true,
    steps: [
        {
            id: 'welcome',
            title: 'Tickets de soporte',
            body: 'Si algo falla en la operación, se reporta aquí y queda con historial y responsable, en vez de perderse en un chat.',
        },
        {
            id: 'estados',
            title: 'Estados',
            body: 'Cada ticket avanza por estados y guarda quien lo movió y cuando. Eso permite medir cuanto se demora el soporte en responder.',
            target: 'table',
            placement: 'top',
            optional: true,
        },
    ],
};

export const announcementsTour: TourDefinition = {
    key: 'announcements',
    version: 1,
    title: 'Anuncios',
    routes: ['/announcements'],
    autoStart: true,
    superAdminOnly: true,
    steps: [
        {
            id: 'welcome',
            title: 'Anuncios',
            body: 'Mensajes que se le muestran a los negocios dentro de la plataforma: mantenimientos, funciones nuevas o avisos comerciales.',
        },
        {
            id: 'targets',
            title: 'A quien le llega',
            body: 'Puedes dirigir un anuncio a todos o solo a ciertos negocios. Sin destinatarios, el anuncio no se muestra a nadie.',
        },
        {
            id: 'formatos',
            title: 'Cinta o modal',
            body: 'La cinta es discreta y va arriba; el modal interrumpe. Usa el modal solo para lo que de verdad no se puede pasar por alto.',
        },
    ],
};

export const commercialTour: TourDefinition = {
    key: 'commercial',
    version: 1,
    title: 'Comercial',
    routes: ['/commercial', '/marketing-leads', '/siigo-referrals'],
    autoStart: true,
    superAdminOnly: true,
    steps: [
        {
            id: 'welcome',
            title: 'Comercial',
            body: 'El seguimiento de prospectos y referidos antes de que se conviertan en negocios activos.',
        },
        {
            id: 'leads',
            title: 'Leads de encuestas',
            body: 'Los contactos que dejaron sus datos en el sitio publico o en una encuesta.',
            target: 'a[href="/marketing-leads"]',
            placement: 'bottom',
            optional: true,
        },
        {
            id: 'siigo',
            title: 'Referidos Siigo',
            body: 'Los que llegaron por la alianza con Siigo, con su rango de órdenes declarado.',
            target: 'a[href="/siigo-referrals"]',
            placement: 'bottom',
            optional: true,
        },
    ],
};
