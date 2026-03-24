import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PayGatewayUseCases } from './use-cases';
import { IPayGatewayRepository } from '../domain/ports';
import {
    PaymentGatewayType,
    PaymentGatewayTypesResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makePaymentGatewayType = (overrides: Partial<PaymentGatewayType> = {}): PaymentGatewayType => ({
    id: 1,
    name: 'Stripe',
    code: 'stripe',
    image_url: 'https://example.com/stripe.png',
    is_active: true,
    in_development: false,
    ...overrides,
});

const makeSuccessResponse = (data: PaymentGatewayType[] = [makePaymentGatewayType()]): PaymentGatewayTypesResponse => ({
    success: true,
    data,
    message: 'OK',
});

const makeErrorResponse = (): PaymentGatewayTypesResponse => ({
    success: false,
    data: [],
    message: 'Error al obtener tipos de pasarela',
});

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IPayGatewayRepository {
    return {
        listPaymentGatewayTypes: vi.fn(),
    } as unknown as IPayGatewayRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('PayGatewayUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: PayGatewayUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new PayGatewayUseCases(repo as unknown as IPayGatewayRepository);
    });

    // ---------------------------------------------------------------
    // getPaymentGatewayTypes
    // ---------------------------------------------------------------
    describe('getPaymentGatewayTypes', () => {
        it('debería retornar la lista de tipos de pasarela cuando el repositorio responde con éxito', async () => {
            const gatewayTypes = [
                makePaymentGatewayType(),
                makePaymentGatewayType({ id: 2, name: 'PayPal', code: 'paypal' }),
            ];
            vi.mocked(repo.listPaymentGatewayTypes).mockResolvedValue(makeSuccessResponse(gatewayTypes));

            const result = await useCases.getPaymentGatewayTypes();

            expect(result).toEqual(gatewayTypes);
            expect(repo.listPaymentGatewayTypes).toHaveBeenCalledOnce();
        });

        it('debería retornar un arreglo vacío cuando el repositorio responde con success=false', async () => {
            vi.mocked(repo.listPaymentGatewayTypes).mockResolvedValue(makeErrorResponse());

            const result = await useCases.getPaymentGatewayTypes();

            expect(result).toEqual([]);
            expect(repo.listPaymentGatewayTypes).toHaveBeenCalledOnce();
        });

        it('debería retornar un arreglo vacío cuando no hay tipos de pasarela', async () => {
            vi.mocked(repo.listPaymentGatewayTypes).mockResolvedValue(makeSuccessResponse([]));

            const result = await useCases.getPaymentGatewayTypes();

            expect(result).toEqual([]);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Error de red');
            vi.mocked(repo.listPaymentGatewayTypes).mockRejectedValue(expectedError);

            await expect(useCases.getPaymentGatewayTypes()).rejects.toThrow('Error de red');
        });
    });
});
