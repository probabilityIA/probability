import { describe, it, expect } from 'vitest';
import {
    mapLoginRequest,
    mapChangePasswordRequest,
    mapGeneratePasswordRequest,
    mapGenerateBusinessTokenRequest,
} from './index';

describe('mapLoginRequest', () => {
    it('debería mapear email y password correctamente', () => {
        const result = mapLoginRequest({ email: 'juan@test.com', password: 'pass123' });

        expect(result).toEqual({ email: 'juan@test.com', password: 'pass123' });
    });

    it('debería retornar un objeto con solo las propiedades esperadas', () => {
        const result = mapLoginRequest({ email: 'test@test.com', password: 'abc' });

        expect(Object.keys(result)).toEqual(['email', 'password']);
    });
});

describe('mapChangePasswordRequest', () => {
    it('debería mapear current_password y new_password correctamente', () => {
        const result = mapChangePasswordRequest({ current_password: 'old', new_password: 'new123' });

        expect(result).toEqual({ current_password: 'old', new_password: 'new123' });
    });
});

describe('mapGeneratePasswordRequest', () => {
    it('debería mapear user_id correctamente', () => {
        const result = mapGeneratePasswordRequest({ user_id: 5 });

        expect(result).toEqual({ user_id: 5 });
    });

    it('debería manejar user_id undefined', () => {
        const result = mapGeneratePasswordRequest({});

        expect(result).toEqual({ user_id: undefined });
    });
});

describe('mapGenerateBusinessTokenRequest', () => {
    it('debería mapear business_id correctamente', () => {
        const result = mapGenerateBusinessTokenRequest({ business_id: 3 });

        expect(result).toEqual({ business_id: 3 });
    });
});
