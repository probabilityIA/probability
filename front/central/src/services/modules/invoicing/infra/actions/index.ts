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
  PaginatedInvoiceableOrders,
  BulkCreateInvoicesDTO,
  BulkCreateResult,
  SyncLog,
  CompareResponseData,
  ItemCompareResponseData,
  BankAccountsResponseData,
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

  // Log del request
  console.log(`[INVOICING API Request] ${options.method || 'GET'} ${url}`, {
    headers: Object.keys(headers),
    body: options.body ? JSON.parse(options.body as string) : undefined,
  });

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorText = await response.text();
    console.error(`[INVOICING API Error] ${response.status} ${url}`, errorText);

    let errorMessage = `Error ${response.status}`;
    try {
      const errorBody = JSON.parse(errorText);
      errorMessage = errorBody.message || errorBody.error || errorMessage;
    } catch {
      errorMessage = errorText || errorMessage;
    }

    throw new Error(errorMessage);
  }

  const data = await response.json();

  // Log de la respuesta
  console.log(`[INVOICING API Response] ${response.status} ${url}`, data);

  return data;
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
  if (filters?.invoice_number) params.append('invoice_number', filters.invoice_number);
  if (filters?.order_number) params.append('order_number', filters.order_number);
  if (filters?.customer_name) params.append('customer_name', filters.customer_name);
  if (filters?.start_date) params.append('start_date', filters.start_date);
  if (filters?.end_date) params.append('end_date', filters.end_date);
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.page_size) params.append('page_size', filters.page_size.toString());

  const queryString = params.toString();
  const url = `${API_BASE_URL}/invoicing/invoices${queryString ? '?' + queryString : ''}`;

  const response = await fetchWithAuth(url);
  // Mapear respuesta del backend al formato esperado
  return {
    data: response.items || [],
    total: response.total_count || 0,
    page: response.page || 1,
    page_size: response.page_size || 20,
  };
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

export async function cancelRetryAction(id: number): Promise<void> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/${id}/retry`, {
    method: 'DELETE',
  });
}

export async function enableRetryAction(id: number): Promise<void> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/${id}/retry`, {
    method: 'PUT',
  });
}

export async function deletePendingInvoiceAction(id: number): Promise<void> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/${id}`, {
    method: 'DELETE',
  });
}

export async function getInvoiceSyncLogsAction(id: number): Promise<SyncLog[]> {
  const response = await fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/${id}/sync-logs`);
  return response.sync_logs || [];
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

  const response = await fetchWithAuth(url);
  // Mapear respuesta del backend al formato esperado
  return {
    data: response.items || [],
    total: response.total_count || 0,
    page: response.page || 1,
    page_size: response.page_size || 20,
  };
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

  const response = await fetchWithAuth(url);
  // Mapear respuesta del backend al formato esperado
  return {
    data: response.items || [],
    total: response.total_count || 0,
    page: response.page || 1,
    page_size: response.page_size || 20,
  };
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

export async function updateConfigAction(id: number, data: UpdateConfigDTO, businessId?: number): Promise<InvoicingConfig> {
  const url = businessId
    ? `${API_BASE_URL}/invoicing/configs/${id}?business_id=${businessId}`
    : `${API_BASE_URL}/invoicing/configs/${id}`;
  return fetchWithAuth(url, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteConfigAction(id: number): Promise<void> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}`, {
    method: 'DELETE',
  });
}

export async function enableConfigAction(id: number): Promise<InvoicingConfig> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}/enable`, {
    method: 'POST',
  });
}

export async function disableConfigAction(id: number): Promise<InvoicingConfig> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}/disable`, {
    method: 'POST',
  });
}

export async function enableAutoInvoiceAction(id: number): Promise<InvoicingConfig> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}/enable-auto-invoice`, {
    method: 'POST',
  });
}

export async function disableAutoInvoiceAction(id: number): Promise<InvoicingConfig> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/configs/${id}/disable-auto-invoice`, {
    method: 'POST',
  });
}

// ============================================
// BULK INVOICES
// ============================================

/**
 * Obtiene órdenes facturables (invoiceable=true, invoice_id IS NULL)
 *
 * @param page - Número de página (default: 1)
 * @param pageSize - Tamaño de página (default: 100)
 * @param businessId - Filtro por business (opcional). Solo aplica para super admin (business_id=0)
 *
 * Super admin (business_id = 0):
 *   - Sin businessId: lista órdenes de TODOS los businesses
 *   - Con businessId: filtra solo ese business específico
 * Usuario normal:
 *   - Ignora businessId, siempre filtra por su business_id del JWT
 */
export async function getInvoiceableOrdersAction(
  page: number = 1,
  pageSize: number = 100,
  businessId?: number | null
): Promise<PaginatedInvoiceableOrders> {
  const params = new URLSearchParams({
    page: page.toString(),
    page_size: pageSize.toString(),
  });

  // Si se especifica un businessId, agregarlo al query string
  if (businessId !== null && businessId !== undefined) {
    params.append('business_id', businessId.toString());
  }

  return fetchWithAuth(
    `${API_BASE_URL}/invoicing/invoices/invoiceable-orders?${params.toString()}`,
    { cache: 'no-store' }
  );
}

/**
 * Crea facturas masivamente (async job via RabbitMQ)
 * Retorna resultado o error sin lanzar excepciones,
 * ya que Next.js en producción enmascara errores de Server Actions.
 */
export async function createBulkInvoicesAction(
  dto: BulkCreateInvoicesDTO
): Promise<{ success: true; data: BulkCreateResult } | { success: false; error: string }> {
  try {
    const data = await fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/bulk`, {
      method: 'POST',
      body: JSON.stringify(dto),
    });
    return { success: true, data };
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Error desconocido al crear facturas';
    console.error('[createBulkInvoicesAction] Error:', message);
    return { success: false, error: message };
  }
}

// ============================================
// FACTURA MANUAL (Registro externo)
// ============================================

/**
 * Registra una factura externa (hecha por fuera del sistema).
 * Crea el registro con status "issued" y marca la orden como facturada.
 */
export async function registerManualInvoiceAction(
  orderID: string,
  invoiceNumber: string,
  businessId?: number
): Promise<{ success: true; data: Invoice } | { success: false; error: string }> {
  try {
    const params = businessId ? `?business_id=${businessId}` : '';
    const data = await fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/manual${params}`, {
      method: 'POST',
      body: JSON.stringify({
        order_id: orderID,
        invoice_number: invoiceNumber,
      }),
    });
    return { success: true, data };
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Error desconocido al registrar factura manual';
    console.error('[registerManualInvoiceAction] Error:', message);
    return { success: false, error: message };
  }
}

// ============================================
// COMPARACIÓN DE FACTURAS (Auditoría)
// ============================================

/**
 * Inicia una comparación asíncrona entre facturas del sistema y el proveedor.
 * El resultado llega via SSE con el evento "invoice.compare_ready".
 *
 * @param businessId - ID del negocio (solo requerido para super admin)
 */
export async function requestInvoiceComparisonAction(
  dateFrom: string,
  dateTo: string,
  businessId?: number
): Promise<{ correlation_id: string; message: string }> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/compare`, {
    method: 'POST',
    body: JSON.stringify({
      date_from: dateFrom,
      date_to: dateTo,
      ...(businessId ? { business_id: businessId } : {}),
    }),
  });
}

/**
 * Polls for a compare result by correlation ID.
 * Returns the data if ready, null if not yet available (404).
 *
 * @param correlationId - The correlation_id returned by requestInvoiceComparisonAction
 * @param businessId - ID del negocio (solo requerido para super admin)
 */
export async function getCompareResultAction(
  correlationId: string,
  businessId?: number
): Promise<CompareResponseData | null> {
  const params = businessId ? `?business_id=${businessId}` : '';
  const url = `${API_BASE_URL}/invoicing/invoices/compare/${correlationId}${params}`;

  const cookieStore = await cookies();
  const sessionToken = cookieStore.get('session_token')?.value;
  const businessToken = cookieStore.get('business_token')?.value;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
  if (businessToken) headers['X-Business-Token'] = businessToken;

  const response = await fetch(url, { headers });

  // 404 means the result is not ready yet
  if (response.status === 404) {
    return null;
  }

  if (!response.ok) {
    const errorText = await response.text();
    let errorMessage = `Error ${response.status}`;
    try {
      const errorBody = JSON.parse(errorText);
      errorMessage = errorBody.message || errorBody.error || errorMessage;
    } catch {
      errorMessage = errorText || errorMessage;
    }
    throw new Error(errorMessage);
  }

  return response.json();
}

// ============================================
// COMPARACIÓN DE ÍTEMS/PRODUCTOS (Auditoría)
// ============================================

/**
 * Inicia una comparación asíncrona de ítems del proveedor vs productos del sistema.
 * El resultado llega via SSE con el evento "invoice.list_items_ready".
 */
export async function requestListItemsComparisonAction(
  businessId?: number
): Promise<{ correlation_id: string; message: string }> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/items`, {
    method: 'POST',
    body: JSON.stringify({
      ...(businessId ? { business_id: businessId } : {}),
    }),
  });
}

/**
 * Polls for an items compare result by correlation ID.
 * Returns the data if ready, null if not yet available (404).
 */
export async function getListItemsResultAction(
  correlationId: string,
  businessId?: number
): Promise<ItemCompareResponseData | null> {
  const params = businessId ? `?business_id=${businessId}` : '';
  const url = `${API_BASE_URL}/invoicing/invoices/items/${correlationId}${params}`;

  const cookieStore = await cookies();
  const sessionToken = cookieStore.get('session_token')?.value;
  const businessToken = cookieStore.get('business_token')?.value;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
  if (businessToken) headers['X-Business-Token'] = businessToken;

  const response = await fetch(url, { headers });

  if (response.status === 404) {
    return null;
  }

  if (!response.ok) {
    const errorText = await response.text();
    let errorMessage = `Error ${response.status}`;
    try {
      const errorBody = JSON.parse(errorText);
      errorMessage = errorBody.message || errorBody.error || errorMessage;
    } catch {
      errorMessage = errorText || errorMessage;
    }
    throw new Error(errorMessage);
  }

  return response.json();
}

// ============================================
// CUENTAS BANCARIAS (Softpymes)
// ============================================

/**
 * Inicia una consulta asíncrona de cuentas bancarias del proveedor.
 * El resultado llega via SSE con el evento "invoice.list_bank_accounts_ready".
 */
export async function requestListBankAccountsAction(
  businessId?: number
): Promise<{ correlation_id: string; message: string }> {
  return fetchWithAuth(`${API_BASE_URL}/invoicing/invoices/bank-accounts`, {
    method: 'POST',
    body: JSON.stringify({
      ...(businessId ? { business_id: businessId } : {}),
    }),
  });
}

/**
 * Polls for a bank accounts result by correlation ID.
 * Returns the data if ready, null if not yet available (404).
 */
export async function getListBankAccountsResultAction(
  correlationId: string,
  businessId?: number
): Promise<BankAccountsResponseData | null> {
  const params = businessId ? `?business_id=${businessId}` : '';
  const url = `${API_BASE_URL}/invoicing/invoices/bank-accounts/${correlationId}${params}`;

  const cookieStore = await cookies();
  const sessionToken = cookieStore.get('session_token')?.value;
  const businessToken = cookieStore.get('business_token')?.value;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
  if (businessToken) headers['X-Business-Token'] = businessToken;

  const response = await fetch(url, { headers });

  if (response.status === 404) {
    return null;
  }

  if (!response.ok) {
    const errorText = await response.text();
    let errorMessage = `Error ${response.status}`;
    try {
      const errorBody = JSON.parse(errorText);
      errorMessage = errorBody.message || errorBody.error || errorMessage;
    } catch {
      errorMessage = errorText || errorMessage;
    }
    throw new Error(errorMessage);
  }

  return response.json();
}
