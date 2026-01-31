/**
 * Server Actions para el módulo de facturación
 * Estas funciones se ejecutan en el servidor
 */

'use server';

import { cookies } from 'next/headers';
import type {
  Invoice,
  InvoicingProvider,
  InvoicingProviderType,
  InvoicingConfig,
  CreateProviderDTO,
  UpdateProviderDTO,
  CreateConfigDTO,
  UpdateConfigDTO,
  PaginatedInvoices,
  PaginatedProviders,
  PaginatedConfigs,
  InvoiceFilters,
  ProviderFilters,
  ConfigFilters,
} from '../../domain/types';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:3050/api/v1';

/**
 * Función auxiliar para hacer fetch con autenticación desde el servidor
 */
async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const cookieStore = await cookies();
  const sessionToken = cookieStore.get('session_token')?.value;
  const businessToken = cookieStore.get('business_token')?.value;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...options.headers as Record<string, string>,
  };

  if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
  if (businessToken) headers['X-Business-Token'] = businessToken;

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`HTTP ${response.status}: ${errorText}`);
  }

  return response.json();
}

// ============================================
// INVOICES
// ============================================

export async function getInvoicesAction(filters?: InvoiceFilters): Promise<PaginatedInvoices> {
  const params = new URLSearchParams();
  if (filters?.business_id) params.append('business_id', filters.business_id.toString());
  if (filters?.order_id) params.append('order_id', filters.order_id);
  if (filters?.integration_id) params.append('integration_id', filters.integration_id.toString());
  if (filters?.status) params.append('status', filters.status);
  if (filters?.provider_id) params.append('provider_id', filters.provider_id.toString());
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.page_size) params.append('page_size', filters.page_size.toString());

  const queryString = params.toString();
  const url = `${API_BASE_URL}/invoicing/invoices${queryString ? '?' + queryString : ''}`;

  return fetchWithAuth(url);
}

export async function getInvoiceByIdAction(id: number): Promise<Invoice> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/${id}`);
}

export async function cancelInvoiceAction(id: number): Promise<Invoice> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/${id}/cancel`, {
    method: 'POST',
  });
}

export async function retryInvoiceAction(id: number): Promise<Invoice> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/${id}/retry`, {
    method: 'POST',
  });
}

// ============================================
// PROVIDERS
// ============================================

export async function getProvidersAction(filters?: ProviderFilters): Promise<PaginatedProviders> {
  const params = new URLSearchParams();
  if (filters?.business_id) params.append('business_id', filters.business_id.toString());
  if (filters?.provider_type) params.append('provider_type', filters.provider_type);
  if (filters?.is_active !== undefined) params.append('is_active', filters.is_active.toString());
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.page_size) params.append('page_size', filters.page_size.toString());

  const queryString = params.toString();
  const url = `${API_BASE_URL}/invoicing/providers${queryString ? '?' + queryString : ''}`;

  return fetchWithAuth(url);
}

export async function getProviderByIdAction(id: number): Promise<InvoicingProvider> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/providers/${id}`);
}

export async function createProviderAction(data: CreateProviderDTO): Promise<InvoicingProvider> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/providers`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateProviderAction(id: number, data: UpdateProviderDTO): Promise<InvoicingProvider> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/providers/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function testProviderConnectionAction(id: number): Promise<{ success: boolean; message: string }> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/providers/${id}/test`, {
    method: 'POST',
  });
}

export async function getProviderTypesAction(): Promise<InvoicingProviderType[]> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/provider-types`);
}

// ============================================
// CONFIGS
// ============================================

export async function getConfigsAction(filters?: ConfigFilters): Promise<PaginatedConfigs> {
  const params = new URLSearchParams();
  if (filters?.business_id) params.append('business_id', filters.business_id.toString());
  if (filters?.integration_id) params.append('integration_id', filters.integration_id.toString());
  if (filters?.provider_id) params.append('provider_id', filters.provider_id.toString());
  if (filters?.enabled !== undefined) params.append('enabled', filters.enabled.toString());
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.page_size) params.append('page_size', filters.page_size.toString());

  const queryString = params.toString();
  const url = `${API_BASE_URL}/invoicing/configs${queryString ? '?' + queryString : ''}`;

  return fetchWithAuth(url);
}

export async function getConfigByIdAction(id: number): Promise<InvoicingConfig> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}`);
}

export async function createConfigAction(data: CreateConfigDTO): Promise<InvoicingConfig> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateConfigAction(id: number, data: UpdateConfigDTO): Promise<InvoicingConfig> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteConfigAction(id: number): Promise<void> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}`, {
    method: 'DELETE',
  });
}
