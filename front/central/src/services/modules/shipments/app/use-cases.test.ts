import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ShipmentUseCases } from './use-cases';
import { IShipmentRepository } from '../domain/ports';
import {
    Shipment,
    PaginatedResponse,
    EnvioClickQuoteRequest,
    EnvioClickQuoteResponse,
    EnvioClickGenerateResponse,
    EnvioClickTrackingResponse,
    EnvioClickCancelResponse,
    CreateShipmentRequest,
    OriginAddress,
    CreateOriginAddressRequest,
    UpdateOriginAddressRequest,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeShipment = (overrides: Partial<Shipment> = {}): Shipment => ({
    id: 1,
    created_at: '2026-03-01T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
    status: 'pending',
    is_last_mile: false,
    is_test: false,
    ...overrides,
});

const makeOriginAddress = (overrides: Partial<OriginAddress> = {}): OriginAddress => ({
    id: 1,
    business_id: 1,
    alias: 'Bodega Principal',
    company: 'Mi Empresa',
    first_name: 'Juan',
    last_name: 'Perez',
    email: 'juan@test.com',
    phone: '3001234567',
    street: 'Calle 10 #5-20',
    city_dane_code: '11001',
    city: 'Bogota',
    state: 'Cundinamarca',
    is_default: true,
    created_at: '2026-03-01T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
    ...overrides,
});

const paginatedShipments: PaginatedResponse<Shipment> = {
    success: true,
    message: 'OK',
    data: [makeShipment()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IShipmentRepository {
    return {
        getShipments: vi.fn(),
        quoteShipment: vi.fn(),
        generateGuide: vi.fn(),
        trackShipment: vi.fn(),
        cancelShipment: vi.fn(),
        createShipment: vi.fn(),
        getOriginAddresses: vi.fn(),
        createOriginAddress: vi.fn(),
        updateOriginAddress: vi.fn(),
        deleteOriginAddress: vi.fn(),
    } as unknown as IShipmentRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('ShipmentUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: ShipmentUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new ShipmentUseCases(repo as unknown as IShipmentRepository);
    });

    // ---------------------------------------------------------------
    // getShipments
    // ---------------------------------------------------------------
    describe('getShipments', () => {
        it('debería retornar la lista paginada de envíos cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getShipments).mockResolvedValue(paginatedShipments);

            const result = await useCases.getShipments({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedShipments);
            expect(repo.getShipments).toHaveBeenCalledOnce();
            expect(repo.getShipments).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getShipments).mockResolvedValue(paginatedShipments);

            await useCases.getShipments();

            expect(repo.getShipments).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getShipments).mockRejectedValue(expectedError);

            await expect(useCases.getShipments()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // quoteShipment
    // ---------------------------------------------------------------
    describe('quoteShipment', () => {
        const quoteReq: EnvioClickQuoteRequest = {
            description: 'Paquete de prueba',
            contentValue: 50000,
            includeGuideCost: true,
            codPaymentMethod: 'cash',
            packages: [{ weight: 2, height: 10, width: 10, length: 10 }],
            origin: { address: 'Calle 10', daneCode: '11001' },
            destination: { address: 'Calle 20', daneCode: '05001' },
        };

        it('debería retornar cotización de envío', async () => {
            const quoteResponse: EnvioClickQuoteResponse = {
                success: true,
                message: 'OK',
                correlation_id: 'corr-123',
                data: { rates: [{ idRate: 1, idProduct: 1, product: 'Express', idCarrier: 1, carrier: 'Servientrega', flete: 15000, deliveryDays: 2, quotationType: 'standard' }] },
            };
            vi.mocked(repo.quoteShipment).mockResolvedValue(quoteResponse);

            const result = await useCases.quoteShipment(quoteReq);

            expect(result).toEqual(quoteResponse);
            expect(repo.quoteShipment).toHaveBeenCalledWith(quoteReq);
        });

        it('debería propagar el error cuando la cotización falla', async () => {
            vi.mocked(repo.quoteShipment).mockRejectedValue(new Error('Servicio de cotización no disponible'));

            await expect(useCases.quoteShipment(quoteReq)).rejects.toThrow('Servicio de cotización no disponible');
        });
    });

    // ---------------------------------------------------------------
    // generateGuide
    // ---------------------------------------------------------------
    describe('generateGuide', () => {
        const guideReq: EnvioClickQuoteRequest = {
            idRate: 1,
            carrier: 'Servientrega',
            description: 'Paquete',
            contentValue: 50000,
            includeGuideCost: true,
            codPaymentMethod: 'cash',
            packages: [{ weight: 2, height: 10, width: 10, length: 10 }],
            origin: { address: 'Calle 10', daneCode: '11001' },
            destination: { address: 'Calle 20', daneCode: '05001' },
        };

        it('debería generar una guía y retornar la respuesta', async () => {
            const generateResponse: EnvioClickGenerateResponse = {
                success: true,
                message: 'Guía generada',
                correlation_id: 'corr-456',
                shipment_id: 1,
            };
            vi.mocked(repo.generateGuide).mockResolvedValue(generateResponse);

            const result = await useCases.generateGuide(guideReq);

            expect(result).toEqual(generateResponse);
            expect(repo.generateGuide).toHaveBeenCalledWith(guideReq);
        });

        it('debería propagar el error cuando la generación de guía falla', async () => {
            vi.mocked(repo.generateGuide).mockRejectedValue(new Error('Error al generar guía'));

            await expect(useCases.generateGuide(guideReq)).rejects.toThrow('Error al generar guía');
        });
    });

    // ---------------------------------------------------------------
    // trackShipment
    // ---------------------------------------------------------------
    describe('trackShipment', () => {
        it('debería retornar el tracking de un envío', async () => {
            const trackingResponse: EnvioClickTrackingResponse = {
                success: true,
                message: 'OK',
                data: {
                    trackingNumber: 'TRACK-001',
                    carrier: 'Servientrega',
                    status: 'in_transit',
                    history: [{ date: '2026-03-01', status: 'shipped', description: 'Paquete enviado', location: 'Bogota' }],
                },
            };
            vi.mocked(repo.trackShipment).mockResolvedValue(trackingResponse);

            const result = await useCases.trackShipment('TRACK-001');

            expect(result).toEqual(trackingResponse);
            expect(repo.trackShipment).toHaveBeenCalledWith('TRACK-001');
        });

        it('debería propagar el error cuando el tracking falla', async () => {
            vi.mocked(repo.trackShipment).mockRejectedValue(new Error('Tracking no encontrado'));

            await expect(useCases.trackShipment('INVALID')).rejects.toThrow('Tracking no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // cancelShipment
    // ---------------------------------------------------------------
    describe('cancelShipment', () => {
        it('debería cancelar un envío y retornar confirmación', async () => {
            const cancelResponse: EnvioClickCancelResponse = {
                success: true,
                message: 'OK',
                data: { status: 'cancelled', message: 'Envío cancelado' },
            };
            vi.mocked(repo.cancelShipment).mockResolvedValue(cancelResponse);

            const result = await useCases.cancelShipment('1');

            expect(result).toEqual(cancelResponse);
            expect(repo.cancelShipment).toHaveBeenCalledWith('1');
        });

        it('debería propagar el error cuando la cancelación falla', async () => {
            vi.mocked(repo.cancelShipment).mockRejectedValue(new Error('Envío ya entregado'));

            await expect(useCases.cancelShipment('1')).rejects.toThrow('Envío ya entregado');
        });
    });

    // ---------------------------------------------------------------
    // createShipment
    // ---------------------------------------------------------------
    describe('createShipment', () => {
        const createReq: CreateShipmentRequest = {
            order_id: 'order-1',
            client_name: 'Juan Perez',
            destination_address: 'Calle 20 #10-30',
            carrier: 'Servientrega',
            status: 'pending',
        };

        it('debería crear un envío y retornar la respuesta', async () => {
            const createResponse = { success: true, message: 'Envío creado', data: makeShipment() };
            vi.mocked(repo.createShipment).mockResolvedValue(createResponse);

            const result = await useCases.createShipment(createReq);

            expect(result).toEqual(createResponse);
            expect(repo.createShipment).toHaveBeenCalledOnce();
            expect(repo.createShipment).toHaveBeenCalledWith(createReq);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createShipment).mockRejectedValue(new Error('Orden no encontrada'));

            await expect(useCases.createShipment(createReq)).rejects.toThrow('Orden no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // getOriginAddresses
    // ---------------------------------------------------------------
    describe('getOriginAddresses', () => {
        it('debería retornar las direcciones de origen', async () => {
            const addresses = [makeOriginAddress()];
            vi.mocked(repo.getOriginAddresses).mockResolvedValue(addresses);

            const result = await useCases.getOriginAddresses();

            expect(result).toEqual(addresses);
            expect(repo.getOriginAddresses).toHaveBeenCalledWith(undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.getOriginAddresses).mockResolvedValue([]);

            await useCases.getOriginAddresses(5);

            expect(repo.getOriginAddresses).toHaveBeenCalledWith(5);
        });

        it('debería propagar el error cuando falla la consulta de direcciones', async () => {
            vi.mocked(repo.getOriginAddresses).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getOriginAddresses()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // createOriginAddress
    // ---------------------------------------------------------------
    describe('createOriginAddress', () => {
        const createReq: CreateOriginAddressRequest = {
            alias: 'Nueva Bodega',
            company: 'Mi Empresa',
            first_name: 'Carlos',
            last_name: 'Lopez',
            email: 'carlos@test.com',
            phone: '3009876543',
            street: 'Carrera 5 #12-34',
            city_dane_code: '11001',
            city: 'Bogota',
            state: 'Cundinamarca',
        };

        it('debería crear una dirección de origen y retornar la respuesta', async () => {
            const newAddress = makeOriginAddress({ id: 2, alias: 'Nueva Bodega' });
            vi.mocked(repo.createOriginAddress).mockResolvedValue(newAddress);

            const result = await useCases.createOriginAddress(createReq);

            expect(result).toEqual(newAddress);
            expect(repo.createOriginAddress).toHaveBeenCalledWith(createReq, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.createOriginAddress).mockResolvedValue(makeOriginAddress());

            await useCases.createOriginAddress(createReq, 5);

            expect(repo.createOriginAddress).toHaveBeenCalledWith(createReq, 5);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createOriginAddress).mockRejectedValue(new Error('Dirección duplicada'));

            await expect(useCases.createOriginAddress(createReq)).rejects.toThrow('Dirección duplicada');
        });
    });

    // ---------------------------------------------------------------
    // updateOriginAddress
    // ---------------------------------------------------------------
    describe('updateOriginAddress', () => {
        const updateReq: UpdateOriginAddressRequest = { alias: 'Bodega Actualizada' };

        it('debería actualizar una dirección de origen y retornar la respuesta', async () => {
            const updatedAddress = makeOriginAddress({ alias: 'Bodega Actualizada' });
            vi.mocked(repo.updateOriginAddress).mockResolvedValue(updatedAddress);

            const result = await useCases.updateOriginAddress(1, updateReq);

            expect(result).toEqual(updatedAddress);
            expect(repo.updateOriginAddress).toHaveBeenCalledWith(1, updateReq, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.updateOriginAddress).mockResolvedValue(makeOriginAddress());

            await useCases.updateOriginAddress(1, updateReq, 5);

            expect(repo.updateOriginAddress).toHaveBeenCalledWith(1, updateReq, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateOriginAddress).mockRejectedValue(new Error('Dirección no encontrada'));

            await expect(useCases.updateOriginAddress(99, updateReq)).rejects.toThrow('Dirección no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteOriginAddress
    // ---------------------------------------------------------------
    describe('deleteOriginAddress', () => {
        it('debería eliminar una dirección de origen y retornar confirmación', async () => {
            const deleteResponse = { message: 'Dirección eliminada' };
            vi.mocked(repo.deleteOriginAddress).mockResolvedValue(deleteResponse);

            const result = await useCases.deleteOriginAddress(1);

            expect(result).toEqual(deleteResponse);
            expect(repo.deleteOriginAddress).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteOriginAddress).mockResolvedValue({ message: 'OK' });

            await useCases.deleteOriginAddress(1, 5);

            expect(repo.deleteOriginAddress).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteOriginAddress).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteOriginAddress(1)).rejects.toThrow('Network error');
        });
    });
});
