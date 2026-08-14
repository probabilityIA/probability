import { getSyncProvider } from '../../ui/providers';
import type { CompareInventoryOptions } from '../../domain/types';

export type CompareAction = 'update' | 'unchanged' | 'skip';

export interface CompareRow {
    product_id: string;
    sku: string;
    name: string;
    image_url?: string;
    external_item_id: string;
    external_variant_id?: string;
    probability_qty: number | null;
    channel_qty: number | null;
    delta: number | null;
    action: CompareAction;
    reason?: string;
}

export interface CompareTotals {
    total: number;
    to_update: number;
    unchanged: number;
    skipped: number;
}

export interface ComparePage {
    rows: CompareRow[];
    totals: CompareTotals;
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
    checked_at: string;
    newest_at?: string;
    from_cache?: boolean;
}

export const GRUPO_INVENTARIO = 100;

export async function compararInventario(
    integrationTypeId: number,
    integrationId: number,
    businessId: number | undefined,
    page: number,
    skus?: string[],
    pageSize: number = GRUPO_INVENTARIO,
    opciones?: CompareInventoryOptions,
): Promise<ComparePage> {
    const provider = getSyncProvider(integrationTypeId);
    if (!provider?.compareInventory) {
        throw new Error('Este canal todavia no permite comparar inventario');
    }

    const respuesta = await provider.compareInventory(
        integrationId,
        businessId,
        page,
        pageSize,
        skus,
        opciones,
    ) as { success?: boolean; data?: ComparePage; message?: string; error?: string };

    if (!respuesta?.success || !respuesta.data) {
        throw new Error(respuesta?.message || respuesta?.error || 'No se pudo leer el stock del canal');
    }
    return respuesta.data;
}
