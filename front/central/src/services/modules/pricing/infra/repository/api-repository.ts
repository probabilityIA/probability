import { env } from '@/shared/config/env';
import {
    ClientGroup,
    ClientSummary,
    CatalogPriceRow,
    Paginated,
    SaveClientGroupInput,
    CatalogPriceTarget,
    CatalogPriceItem,
    EffectivePrice,
} from '../../domain/types';

export class PricingApiRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };
        if (this.token) headers['Authorization'] = `Bearer ${this.token}`;
        const res = await fetch(`${this.baseUrl}${path}`, { ...options, headers, cache: 'no-store' });
        const data = await res.json().catch(() => ({}));
        if (!res.ok) throw new Error(data.message || data.error || 'Error en la solicitud');
        return data as T;
    }

    private withBusiness(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
    }

    listClientGroups(businessId: number | undefined, search: string, page: number): Promise<Paginated<ClientGroup>> {
        const params = new URLSearchParams({ page: String(page), page_size: '50' });
        if (search) params.set('search', search);
        return this.request(this.withBusiness(`/pricing/client-groups?${params.toString()}`, businessId));
    }

    createClientGroup(businessId: number | undefined, input: SaveClientGroupInput): Promise<ClientGroup> {
        return this.request(this.withBusiness('/pricing/client-groups', businessId), {
            method: 'POST',
            body: JSON.stringify(input),
        });
    }

    updateClientGroup(businessId: number | undefined, id: number, input: SaveClientGroupInput): Promise<ClientGroup> {
        return this.request(this.withBusiness(`/pricing/client-groups/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(input),
        });
    }

    deleteClientGroup(businessId: number | undefined, id: number): Promise<unknown> {
        return this.request(this.withBusiness(`/pricing/client-groups/${id}`, businessId), {
            method: 'DELETE',
        });
    }

    listGroupMembers(businessId: number | undefined, groupId: number, search: string, page: number): Promise<Paginated<ClientSummary>> {
        const params = new URLSearchParams({ page: String(page), page_size: '50' });
        if (search) params.set('search', search);
        return this.request(this.withBusiness(`/pricing/client-groups/${groupId}/members?${params.toString()}`, businessId));
    }

    addGroupMembers(businessId: number | undefined, groupId: number, clientIds: number[]): Promise<unknown> {
        return this.request(this.withBusiness(`/pricing/client-groups/${groupId}/members`, businessId), {
            method: 'POST',
            body: JSON.stringify({ client_ids: clientIds }),
        });
    }

    removeGroupMember(businessId: number | undefined, groupId: number, clientId: number): Promise<unknown> {
        return this.request(this.withBusiness(`/pricing/client-groups/${groupId}/members/${clientId}`, businessId), {
            method: 'DELETE',
        });
    }

    listAvailableClients(businessId: number | undefined, search: string, onlyUngrouped: boolean, page: number): Promise<Paginated<ClientSummary>> {
        const params = new URLSearchParams({ page: String(page), page_size: '50' });
        if (search) params.set('search', search);
        if (onlyUngrouped) params.set('only_ungrouped', 'true');
        return this.request(this.withBusiness(`/pricing/clients?${params.toString()}`, businessId));
    }

    getCatalogPrices(businessId: number | undefined, target: CatalogPriceTarget, search: string, page: number): Promise<Paginated<CatalogPriceRow>> {
        const params = new URLSearchParams({ page: String(page), page_size: '50' });
        if (search) params.set('search', search);
        if (target.client_group_id) params.set('client_group_id', String(target.client_group_id));
        if (target.client_id) params.set('client_id', String(target.client_id));
        return this.request(this.withBusiness(`/pricing/catalog-prices?${params.toString()}`, businessId));
    }

    saveCatalogPrices(businessId: number | undefined, target: CatalogPriceTarget, items: CatalogPriceItem[]): Promise<unknown> {
        return this.request(this.withBusiness('/pricing/catalog-prices', businessId), {
            method: 'PUT',
            body: JSON.stringify({
                client_group_id: target.client_group_id ?? null,
                client_id: target.client_id ?? null,
                items,
            }),
        });
    }

    getEffectivePrice(businessId: number | undefined, productId: string, clientId: number): Promise<EffectivePrice> {
        const params = new URLSearchParams({ product_id: productId, client_id: String(clientId) });
        return this.request(this.withBusiness(`/pricing/effective-price?${params.toString()}`, businessId));
    }
}
