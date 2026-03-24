import { describe, it, expect } from 'vitest';
import {
    mapLoginResponse,
    mapUserRolesPermissionsResponse,
    mapChangePasswordResponse,
    mapGeneratePasswordResponse,
    mapGenerateBusinessTokenResponse,
} from './index';

// -----------------------------------------------------------------
// Fixtures
// -----------------------------------------------------------------

const loginFixture = {
    success: true,
    data: {
        user: { id: 1, name: 'Juan', email: 'juan@test.com', phone: '300', avatar_url: '', is_active: true },
        token: 'jwt-token',
        require_password_change: false,
        businesses: [
            {
                id: 1, name: 'Mi Negocio', code: 'mi-neg', business_type_id: 1,
                business_type: { id: 1, name: 'Restaurante', code: 'restaurant', description: '', icon: '' },
                timezone: 'America/Bogota', address: '', description: '', logo_url: '',
                primary_color: '#000', secondary_color: '#fff', tertiary_color: '#ccc', quaternary_color: '#eee',
                navbar_image_url: '', custom_domain: '', is_active: true,
                enable_delivery: false, enable_pickup: false, enable_reservations: false,
            },
        ],
        scope: 'platform',
        is_super_admin: false,
    },
};

const rolesPermissionsFixture = {
    success: true,
    data: {
        is_super: false,
        business_id: 1,
        business_name: 'Mi Negocio',
        business_type_id: 1,
        business_type_name: 'Restaurante',
        role: { id: 1, name: 'Admin', description: 'Administrador' },
        resources: [
            { resource: 'orders', actions: ['read', 'write'], active: true },
        ],
        subscription_status: 'active',
    },
};

// -----------------------------------------------------------------
// Tests
// -----------------------------------------------------------------

describe('mapLoginResponse', () => {
    it('debería retornar la respuesta completa como LoginSuccessResponse', () => {
        const result = mapLoginResponse(loginFixture);

        expect(result.success).toBe(true);
        expect(result.data.token).toBe('jwt-token');
        expect(result.data.user.email).toBe('juan@test.com');
    });

    it('debería preservar los campos de businesses', () => {
        const result = mapLoginResponse(loginFixture);

        expect(result.data.businesses).toHaveLength(1);
        expect(result.data.businesses[0].name).toBe('Mi Negocio');
    });

    it('debería manejar respuesta con businesses vacío', () => {
        const empty = { ...loginFixture, data: { ...loginFixture.data, businesses: [] } };
        const result = mapLoginResponse(empty);

        expect(result.data.businesses).toEqual([]);
    });
});

describe('mapUserRolesPermissionsResponse', () => {
    it('debería retornar la respuesta como UserRolesPermissionsSuccessResponse', () => {
        const result = mapUserRolesPermissionsResponse(rolesPermissionsFixture);

        expect(result.success).toBe(true);
        expect(result.data.role.name).toBe('Admin');
    });

    it('debería preservar los campos de resources', () => {
        const result = mapUserRolesPermissionsResponse(rolesPermissionsFixture);

        expect(result.data.resources).toHaveLength(1);
        expect(result.data.resources[0].actions).toEqual(['read', 'write']);
    });

    it('debería manejar resources vacío', () => {
        const empty = { ...rolesPermissionsFixture, data: { ...rolesPermissionsFixture.data, resources: [] } };
        const result = mapUserRolesPermissionsResponse(empty);

        expect(result.data.resources).toEqual([]);
    });
});

describe('mapChangePasswordResponse', () => {
    it('debería retornar success y message', () => {
        const result = mapChangePasswordResponse({ success: true, message: 'Contraseña actualizada' });

        expect(result.success).toBe(true);
        expect(result.message).toBe('Contraseña actualizada');
    });
});

describe('mapGeneratePasswordResponse', () => {
    it('debería retornar success, email, password y message', () => {
        const input = { success: true, email: 'test@test.com', password: 'newPass', message: 'OK' };
        const result = mapGeneratePasswordResponse(input);

        expect(result.success).toBe(true);
        expect(result.email).toBe('test@test.com');
        expect(result.password).toBe('newPass');
    });
});

describe('mapGenerateBusinessTokenResponse', () => {
    it('debería retornar success y data.token', () => {
        const input = { success: true, data: { token: 'biz-token' }, message: 'OK' };
        const result = mapGenerateBusinessTokenResponse(input);

        expect(result.success).toBe(true);
        expect(result.data.token).toBe('biz-token');
    });
});
