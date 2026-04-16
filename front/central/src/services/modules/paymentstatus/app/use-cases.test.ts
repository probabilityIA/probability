import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PaymentStatusUseCases } from './use-cases';
import { IPaymentStatusRepository } from '../domain/ports';
import {
    PaymentStatusInfo,
    PaymentStatusesResponse,
    GetPaymentStatusesParams,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makePaymentStatus = (overrides: Partial<PaymentStatusInfo> = {}): PaymentStatusInfo => ({
    id: 1,
    code: 'paid',
    name: 'Pagado',
    description: 'Pago completado',
    category: 'completed',
    color: '#00FF00',
    icon: 'check-circle',
    is_active: true,
    ...overrides,
});

const makePaymentStatusesResponse = (overrides: Partial<PaymentStatusesResponse> = {}): PaymentStatusesResponse => ({
    success: true,
    message: 'OK',
    data: [
        makePaymentStatus(),
        makePaymentStatus({ id: 2, code: 'pending', name: 'Pendiente', category: 'active', color: '#FFA500' }),
        makePaymentStatus({ id: 3, code: 'failed', name: 'Fallido', category: 'failed', color: '#FF0000' }),
    ],
    ...overrides,
});

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IPaymentStatusRepository {
    return {
        getPaymentStatuses: vi.fn(),
    } as unknown as IPaymentStatusRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('PaymentStatusUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: PaymentStatusUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new PaymentStatusUseCases(repo as unknown as IPaymentStatusRepository);
    });

    // ---------------------------------------------------------------
    // getPaymentStatuses
    // ---------------------------------------------------------------
    describe('getPaymentStatuses', () => {
        it('debería retornar la lista de estados de pago cuando el repositorio tiene éxito', async () => {
            const response = makePaymentStatusesResponse();
            vi.mocked(repo.getPaymentStatuses).mockResolvedValue(response);

            const result = await useCases.getPaymentStatuses({ is_active: true });

            expect(result).toEqual(response);
            expect(repo.getPaymentStatuses).toHaveBeenCalledOnce();
            expect(repo.getPaymentStatuses).toHaveBeenCalledWith({ is_active: true });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            const response = makePaymentStatusesResponse();
            vi.mocked(repo.getPaymentStatuses).mockResolvedValue(response);

            await useCases.getPaymentStatuses();

            expect(repo.getPaymentStatuses).toHaveBeenCalledWith(undefined);
        });

        it('debería retornar lista vacía cuando no hay estados de pago', async () => {
            const emptyResponse = makePaymentStatusesResponse({ data: [] });
            vi.mocked(repo.getPaymentStatuses).mockResolvedValue(emptyResponse);

            const result = await useCases.getPaymentStatuses();

            expect(result.data).toEqual([]);
            expect(result.success).toBe(true);
        });

        it('debería filtrar por is_active=false cuando se proporciona', async () => {
            const inactiveResponse = makePaymentStatusesResponse({
                data: [makePaymentStatus({ is_active: false })],
            });
            vi.mocked(repo.getPaymentStatuses).mockResolvedValue(inactiveResponse);

            const result = await useCases.getPaymentStatuses({ is_active: false });

            expect(result.data).toHaveLength(1);
            expect(result.data[0].is_active).toBe(false);
            expect(repo.getPaymentStatuses).toHaveBeenCalledWith({ is_active: false });
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getPaymentStatuses).mockRejectedValue(expectedError);

            await expect(useCases.getPaymentStatuses()).rejects.toThrow('Fallo de base de datos');
        });

        it('debería propagar errores de red', async () => {
            vi.mocked(repo.getPaymentStatuses).mockRejectedValue(new Error('Network error'));

            await expect(useCases.getPaymentStatuses({ is_active: true })).rejects.toThrow('Network error');
        });
    });
});
