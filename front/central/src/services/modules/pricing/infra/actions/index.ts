'use server';

import { cookies } from 'next/headers';
import { PricingApiRepository } from '../repository/api-repository';
import {
    ClientGroup,
    ClientSummary,
    CatalogPriceRow,
    Paginated,
    SaveClientGroupInput,
    CatalogPriceTarget,
    CatalogPriceItem,
    ActionResult,
    EffectivePrice,
} from '../../domain/types';

async function getRepository(): Promise<PricingApiRepository> {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    return new PricingApiRepository(token);
}

function emptyPage<T>(): Paginated<T> {
    return { data: [], total: 0, page: 1, page_size: 50, total_pages: 0 };
}

export async function listClientGroupsAction(businessId: number | undefined, search = '', page = 1): Promise<Paginated<ClientGroup>> {
    try {
        return await (await getRepository()).listClientGroups(businessId, search, page);
    } catch {
        return emptyPage<ClientGroup>();
    }
}

export async function saveClientGroupAction(businessId: number | undefined, input: SaveClientGroupInput): Promise<ActionResult<ClientGroup>> {
    try {
        const repo = await getRepository();
        const data = input.id
            ? await repo.updateClientGroup(businessId, input.id, input)
            : await repo.createClientGroup(businessId, input);
        return { success: true, data };
    } catch (error) {
        return { success: false, message: error instanceof Error ? error.message : 'Error al guardar el grupo' };
    }
}

export async function deleteClientGroupAction(businessId: number | undefined, groupId: number): Promise<ActionResult> {
    try {
        await (await getRepository()).deleteClientGroup(businessId, groupId);
        return { success: true };
    } catch (error) {
        return { success: false, message: error instanceof Error ? error.message : 'Error al eliminar el grupo' };
    }
}

export async function listGroupMembersAction(businessId: number | undefined, groupId: number, search = '', page = 1): Promise<Paginated<ClientSummary>> {
    try {
        return await (await getRepository()).listGroupMembers(businessId, groupId, search, page);
    } catch {
        return emptyPage<ClientSummary>();
    }
}

export async function addGroupMembersAction(businessId: number | undefined, groupId: number, clientIds: number[]): Promise<ActionResult> {
    try {
        await (await getRepository()).addGroupMembers(businessId, groupId, clientIds);
        return { success: true };
    } catch (error) {
        return { success: false, message: error instanceof Error ? error.message : 'Error al agregar clientes' };
    }
}

export async function removeGroupMemberAction(businessId: number | undefined, groupId: number, clientId: number): Promise<ActionResult> {
    try {
        await (await getRepository()).removeGroupMember(businessId, groupId, clientId);
        return { success: true };
    } catch (error) {
        return { success: false, message: error instanceof Error ? error.message : 'Error al quitar el cliente' };
    }
}

export async function listAvailableClientsAction(businessId: number | undefined, search = '', onlyUngrouped = false, page = 1): Promise<Paginated<ClientSummary>> {
    try {
        return await (await getRepository()).listAvailableClients(businessId, search, onlyUngrouped, page);
    } catch {
        return emptyPage<ClientSummary>();
    }
}

export async function getCatalogPricesAction(businessId: number | undefined, target: CatalogPriceTarget, search = '', page = 1): Promise<Paginated<CatalogPriceRow>> {
    try {
        return await (await getRepository()).getCatalogPrices(businessId, target, search, page);
    } catch {
        return emptyPage<CatalogPriceRow>();
    }
}

export async function saveCatalogPricesAction(businessId: number | undefined, target: CatalogPriceTarget, items: CatalogPriceItem[]): Promise<ActionResult> {
    try {
        await (await getRepository()).saveCatalogPrices(businessId, target, items);
        return { success: true };
    } catch (error) {
        return { success: false, message: error instanceof Error ? error.message : 'Error al guardar los precios' };
    }
}

export async function getEffectivePriceAction(businessId: number | undefined, productId: string, clientId: number): Promise<EffectivePrice | null> {
    try {
        return await (await getRepository()).getEffectivePrice(businessId, productId, clientId);
    } catch {
        return null;
    }
}
