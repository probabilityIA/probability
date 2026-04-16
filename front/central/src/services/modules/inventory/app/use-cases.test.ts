import { describe, it, expect, vi, beforeEach } from 'vitest';
import { InventoryUseCases } from './use-cases';
import { IInventoryRepository } from '../domain/ports';
import {
    InventoryLevel,
    StockMovement,
    MovementType,
    InventoryListResponse,
    MovementListResponse,
    MovementTypeListResponse,
    AdjustStockDTO,
    TransferStockDTO,
    GetInventoryParams,
    GetMovementsParams,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeInventoryLevel = (overrides: Partial<InventoryLevel> = {}): InventoryLevel => ({
    id: 1,
    product_id: 'prod-001',
    warehouse_id: 1,
    location_id: null,
    business_id: 1,
    quantity: 100,
    reserved_qty: 10,
    available_qty: 90,
    min_stock: 5,
    max_stock: 500,
    reorder_point: 20,
    product_name: 'Producto de prueba',
    product_sku: 'SKU-001',
    warehouse_name: 'Bodega principal',
    warehouse_code: 'WH-01',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeStockMovement = (overrides: Partial<StockMovement> = {}): StockMovement => ({
    id: 1,
    product_id: 'prod-001',
    warehouse_id: 1,
    location_id: null,
    business_id: 1,
    movement_type_id: 1,
    movement_type_code: 'adjustment',
    movement_type_name: 'Ajuste',
    reason: 'Ajuste de inventario',
    quantity: 10,
    previous_qty: 90,
    new_qty: 100,
    reference_type: null,
    reference_id: null,
    integration_id: null,
    notes: 'Ajuste manual',
    created_by_id: 1,
    product_name: 'Producto de prueba',
    product_sku: 'SKU-001',
    warehouse_name: 'Bodega principal',
    created_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeMovementType = (overrides: Partial<MovementType> = {}): MovementType => ({
    id: 1,
    code: 'adjustment',
    name: 'Ajuste',
    description: 'Ajuste de inventario',
    is_active: true,
    direction: 'in',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeInventoryListResponse = (overrides: Partial<InventoryListResponse> = {}): InventoryListResponse => ({
    data: [makeInventoryLevel()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
    ...overrides,
});

const makeMovementListResponse = (overrides: Partial<MovementListResponse> = {}): MovementListResponse => ({
    data: [makeStockMovement()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
    ...overrides,
});

const makeMovementTypeListResponse = (overrides: Partial<MovementTypeListResponse> = {}): MovementTypeListResponse => ({
    data: [makeMovementType()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
    ...overrides,
});

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IInventoryRepository {
    return {
        getProductInventory: vi.fn(),
        getWarehouseInventory: vi.fn(),
        adjustStock: vi.fn(),
        transferStock: vi.fn(),
        getMovements: vi.fn(),
        getMovementTypes: vi.fn(),
    } as unknown as IInventoryRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('InventoryUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: InventoryUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new InventoryUseCases(repo as unknown as IInventoryRepository);
    });

    // ---------------------------------------------------------------
    // getProductInventory
    // ---------------------------------------------------------------
    describe('getProductInventory', () => {
        it('debería retornar el inventario de un producto', async () => {
            const inventoryLevels = [makeInventoryLevel()];
            vi.mocked(repo.getProductInventory).mockResolvedValue(inventoryLevels);

            const result = await useCases.getProductInventory('prod-001');

            expect(result).toEqual(inventoryLevels);
            expect(repo.getProductInventory).toHaveBeenCalledOnce();
            expect(repo.getProductInventory).toHaveBeenCalledWith('prod-001', undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const inventoryLevels = [makeInventoryLevel()];
            vi.mocked(repo.getProductInventory).mockResolvedValue(inventoryLevels);

            await useCases.getProductInventory('prod-001', 5);

            expect(repo.getProductInventory).toHaveBeenCalledWith('prod-001', 5);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Producto no encontrado');
            vi.mocked(repo.getProductInventory).mockRejectedValue(expectedError);

            await expect(useCases.getProductInventory('prod-999')).rejects.toThrow('Producto no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getWarehouseInventory
    // ---------------------------------------------------------------
    describe('getWarehouseInventory', () => {
        it('debería retornar el inventario de una bodega con parámetros', async () => {
            const response = makeInventoryListResponse();
            vi.mocked(repo.getWarehouseInventory).mockResolvedValue(response);

            const params: GetInventoryParams = { page: 1, page_size: 10 };
            const result = await useCases.getWarehouseInventory(1, params);

            expect(result).toEqual(response);
            expect(repo.getWarehouseInventory).toHaveBeenCalledWith(1, params);
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan', async () => {
            const response = makeInventoryListResponse();
            vi.mocked(repo.getWarehouseInventory).mockResolvedValue(response);

            await useCases.getWarehouseInventory(1);

            expect(repo.getWarehouseInventory).toHaveBeenCalledWith(1, undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Bodega no encontrada');
            vi.mocked(repo.getWarehouseInventory).mockRejectedValue(expectedError);

            await expect(useCases.getWarehouseInventory(99)).rejects.toThrow('Bodega no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // adjustStock
    // ---------------------------------------------------------------
    describe('adjustStock', () => {
        const adjustDto: AdjustStockDTO = {
            product_id: 'prod-001',
            warehouse_id: 1,
            quantity: 10,
            reason: 'Ajuste manual',
            notes: 'Corrección de conteo',
        };

        it('debería ajustar el stock y retornar el movimiento creado', async () => {
            const movement = makeStockMovement();
            vi.mocked(repo.adjustStock).mockResolvedValue(movement);

            const result = await useCases.adjustStock(adjustDto);

            expect(result).toEqual(movement);
            expect(repo.adjustStock).toHaveBeenCalledOnce();
            expect(repo.adjustStock).toHaveBeenCalledWith(adjustDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const movement = makeStockMovement();
            vi.mocked(repo.adjustStock).mockResolvedValue(movement);

            await useCases.adjustStock(adjustDto, 5);

            expect(repo.adjustStock).toHaveBeenCalledWith(adjustDto, 5);
        });

        it('debería propagar el error cuando el ajuste falla', async () => {
            const expectedError = new Error('Stock insuficiente');
            vi.mocked(repo.adjustStock).mockRejectedValue(expectedError);

            await expect(useCases.adjustStock(adjustDto)).rejects.toThrow('Stock insuficiente');
        });
    });

    // ---------------------------------------------------------------
    // transferStock
    // ---------------------------------------------------------------
    describe('transferStock', () => {
        const transferDto: TransferStockDTO = {
            product_id: 'prod-001',
            from_warehouse_id: 1,
            to_warehouse_id: 2,
            quantity: 5,
            reason: 'Reabastecimiento',
        };

        it('debería transferir stock y retornar el mensaje de confirmación', async () => {
            const response = { message: 'Transferencia exitosa' };
            vi.mocked(repo.transferStock).mockResolvedValue(response);

            const result = await useCases.transferStock(transferDto);

            expect(result).toEqual(response);
            expect(repo.transferStock).toHaveBeenCalledOnce();
            expect(repo.transferStock).toHaveBeenCalledWith(transferDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const response = { message: 'Transferencia exitosa' };
            vi.mocked(repo.transferStock).mockResolvedValue(response);

            await useCases.transferStock(transferDto, 3);

            expect(repo.transferStock).toHaveBeenCalledWith(transferDto, 3);
        });

        it('debería propagar el error cuando la transferencia falla', async () => {
            const expectedError = new Error('Cantidad excede el stock disponible');
            vi.mocked(repo.transferStock).mockRejectedValue(expectedError);

            await expect(useCases.transferStock(transferDto)).rejects.toThrow('Cantidad excede el stock disponible');
        });
    });

    // ---------------------------------------------------------------
    // getMovements
    // ---------------------------------------------------------------
    describe('getMovements', () => {
        it('debería retornar la lista paginada de movimientos cuando se pasan parámetros', async () => {
            const response = makeMovementListResponse();
            vi.mocked(repo.getMovements).mockResolvedValue(response);

            const params: GetMovementsParams = { page: 1, page_size: 10, type: 'adjustment' };
            const result = await useCases.getMovements(params);

            expect(result).toEqual(response);
            expect(repo.getMovements).toHaveBeenCalledOnce();
            expect(repo.getMovements).toHaveBeenCalledWith(params);
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan', async () => {
            const response = makeMovementListResponse();
            vi.mocked(repo.getMovements).mockResolvedValue(response);

            await useCases.getMovements();

            expect(repo.getMovements).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Error de base de datos');
            vi.mocked(repo.getMovements).mockRejectedValue(expectedError);

            await expect(useCases.getMovements()).rejects.toThrow('Error de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getMovementTypes
    // ---------------------------------------------------------------
    describe('getMovementTypes', () => {
        it('debería retornar la lista de tipos de movimiento con parámetros', async () => {
            const response = makeMovementTypeListResponse();
            vi.mocked(repo.getMovementTypes).mockResolvedValue(response);

            const params = { page: 1, page_size: 10, active_only: true };
            const result = await useCases.getMovementTypes(params);

            expect(result).toEqual(response);
            expect(repo.getMovementTypes).toHaveBeenCalledOnce();
            expect(repo.getMovementTypes).toHaveBeenCalledWith(params);
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan', async () => {
            const response = makeMovementTypeListResponse();
            vi.mocked(repo.getMovementTypes).mockResolvedValue(response);

            await useCases.getMovementTypes();

            expect(repo.getMovementTypes).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Servicio no disponible');
            vi.mocked(repo.getMovementTypes).mockRejectedValue(expectedError);

            await expect(useCases.getMovementTypes()).rejects.toThrow('Servicio no disponible');
        });
    });
});
