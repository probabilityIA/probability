/**
 * Tipos del dominio de Facturación
 */

// ===================================
// FACTURAS
// ===================================

export interface Invoice {
  id: number;
  order_id: string;
  business_id: number;
  integration_id: number;
  invoicing_provider_id: number;
  invoice_number: string;
  external_id?: string;
  status: 'pending' | 'issued' | 'cancelled' | 'failed';
  total_amount: number;
  subtotal: number;
  tax: number;
  discount: number;
  currency: string;
  customer_name: string;
  customer_email?: string;
  customer_dni?: string;
  invoice_url?: string;
  pdf_url?: string;
  xml_url?: string;
  cufe?: string;
  issued_at?: string;
  cancelled_at?: string;
  error_message?: string;
  metadata?: Record<string, any>;
  provider_response?: Record<string, any>; // Respuesta completa del proveedor (incluye full_document)
  created_at: string;
  updated_at: string;
  items?: InvoiceItem[];
}

export interface InvoiceItem {
  id: number;
  invoice_id: number;
  product_sku?: string;
  product_name: string;
  description?: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  tax: number;
  tax_rate?: number;
  discount: number;
}

// ===================================
// PROVEEDORES
// ===================================

export interface InvoicingProvider {
  id: number;
  business_id: number;
  provider_type_code: string;
  name: string;
  description?: string;
  config: Record<string, any>;
  credentials: Record<string, any>;
  is_active: boolean;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface InvoicingProviderType {
  code: string;
  name: string;
  description?: string;
  is_active: boolean;
}

// ===================================
// CONFIGURACIONES
// ===================================

export interface InvoicingConfig {
  id: number;
  business_id: number;
  integration_id: number;
  invoicing_integration_id?: number; // Nueva columna FK a integrations (categoria invoicing)
  invoicing_provider_id?: number; // Deprecado - usar invoicing_integration_id
  enabled: boolean;
  auto_invoice: boolean;
  filters?: InvoicingFilters;
  config?: InvoicingSettings;
  description?: string;
  created_at: string;
  updated_at: string;
  integration_name?: string;
  provider_name?: string;
  provider_image_url?: string; // URL del logo del proveedor
}

/**
 * Filtros para determinar qué órdenes se facturan automáticamente
 * Basado en los 15 filtros del backend
 */
export interface InvoicingFilters {
  // Filtros de monto
  min_amount?: number;
  max_amount?: number;

  // Filtros de pago
  payment_status?: 'paid' | 'unpaid' | 'partial';
  payment_methods?: number[];

  // Filtros de orden
  order_types?: string[];
  exclude_statuses?: string[];

  // Filtros de productos
  exclude_products?: string[];
  include_products_only?: string[];
  min_items_count?: number;
  max_items_count?: number;

  // Filtros de cliente
  customer_types?: string[];
  exclude_customer_ids?: string[];

  // Filtros de ubicación
  shipping_regions?: string[];

  // Filtros de fecha
  date_range?: {
    start_date?: string;
    end_date?: string;
  };
}

/**
 * Configuración adicional del proveedor de facturación
 */
export interface InvoicingSettings {
  include_shipping?: boolean;
  apply_discount?: boolean;
  default_tax_rate?: number;
  invoice_type?: string;
  notes?: string;
  provider_config?: Record<string, any>;

  // Softpymes-specific fields
  default_customer_nit?: string;  // NIT por defecto cuando el cliente no tiene DNI
  resolution_id?: number;          // ID de resolución de Softpymes (requerido)
  branch_code?: string;            // Código de sucursal (default: "001")
  seller_nit?: string;             // NIT del vendedor (opcional)
}

// ===================================
// NOTAS DE CRÉDITO
// ===================================

export interface CreditNote {
  id: number;
  invoice_id: number;
  credit_note_number: string;
  external_id?: string;
  amount: number;
  reason: string;
  note_type: string;
  status: string;
  note_url?: string;
  pdf_url?: string;
  xml_url?: string;
  cufe?: string;
  issued_at?: string;
  created_at: string;
}

// ===================================
// DTOs
// ===================================

export interface CreateInvoiceDTO {
  order_id: string;
  business_id: number;
  integration_id: number;
}

export interface CreateProviderDTO {
  business_id: number;
  provider_type_code: string;
  name: string;
  description?: string;
  config: Record<string, any>;
  credentials: Record<string, any>;
  is_active: boolean;
  is_default: boolean;
}

export interface UpdateProviderDTO {
  name?: string;
  description?: string;
  config?: Record<string, any>;
  credentials?: Record<string, any>;
  is_active?: boolean;
  is_default?: boolean;
}

export interface CreateConfigDTO {
  business_id: number;
  integration_id: number;
  invoicing_integration_id: number; // Nueva columna FK a integrations (categoria invoicing)
  enabled?: boolean;
  auto_invoice?: boolean;
  filters?: InvoicingFilters;
  config?: InvoicingSettings;
  description?: string;
}

export interface UpdateConfigDTO {
  enabled?: boolean;
  auto_invoice?: boolean;
  filters?: InvoicingFilters;
  config?: InvoicingSettings;
  invoicing_integration_id?: number;
}

export interface CreateCreditNoteDTO {
  invoice_id: number;
  amount: number;
  reason: string;
  note_type: string;
}

export interface TestProviderResult {
  success: boolean;
  message: string;
}

// ===================================
// RESPUESTAS PAGINADAS
// ===================================

export interface PaginatedInvoices {
  data: Invoice[];
  total: number;
  page: number;
  page_size: number;
}

export interface PaginatedProviders {
  data: InvoicingProvider[];
  total: number;
  page: number;
  page_size: number;
}

export interface PaginatedConfigs {
  data: InvoicingConfig[];
  total: number;
  page: number;
  page_size: number;
}

// ===================================
// FILTROS
// ===================================

export interface InvoiceFilters {
  business_id?: number;
  order_id?: string;
  integration_id?: number;
  status?: string;
  provider_id?: number;
  page?: number;
  page_size?: number;
}

export interface ProviderFilters {
  business_id?: number;
  provider_type?: string;
  is_active?: boolean;
  page?: number;
  page_size?: number;
}

export interface ConfigFilters {
  business_id?: number;
  integration_id?: number;
  provider_id?: number;
  enabled?: boolean;
  page?: number;
  page_size?: number;
}

// ===================================
// RESPUESTAS ADICIONALES
// ===================================

export interface InvoicingConfigsResponse {
  success: boolean;
  message: string;
  data: InvoicingConfig[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface InvoicingConfigResponse {
  success: boolean;
  message: string;
  data: InvoicingConfig;
}

export interface InvoicingProvidersResponse {
  success: boolean;
  message: string;
  data: InvoicingProvider[];
}

export interface InvoicingActionResponse {
  success: boolean;
  message: string;
  error?: string;
}

export interface FilterOption {
  value: string | number;
  label: string;
}

export interface InvoicingStats {
  total_invoices: number;
  total_amount: number;
  pending_invoices: number;
  failed_invoices: number;
  success_rate: number;
  last_invoice_date?: string;
}

// ===================================
// SYNC LOGS (Historial de reintentos)
// ===================================

export interface SyncLog {
  id: number;
  invoice_id: number;
  operation_type: string;
  status: 'pending' | 'processing' | 'success' | 'failed' | 'cancelled';
  error_message?: string;
  error_code?: string;
  retry_count: number;
  max_retries: number;
  next_retry_at?: string;
  triggered_by: string;
  duration_ms?: number;
  started_at: string;
  completed_at?: string;
  created_at: string;
  request_payload?: Record<string, unknown>;
  request_url?: string;
  response_status?: number;
  response_body?: Record<string, unknown>;
}

// ===================================
// CREACIÓN MASIVA DE FACTURAS
// ===================================

export interface InvoiceableOrder {
  id: string;
  business_id: number;
  order_number: string;
  customer_name: string;
  total_amount: number;
  currency: string;
  created_at: string;
}

export interface PaginatedInvoiceableOrders {
  data: InvoiceableOrder[];
  total: number;
  page: number;
  page_size: number;
}

export interface BulkCreateInvoicesDTO {
  order_ids: string[];
  business_id?: number;
}

export interface BulkCreateResult {
  created: number;
  failed: number;
  results: BulkInvoiceResult[];
}

export interface BulkInvoiceResult {
  order_id: string;
  success: boolean;
  invoice_id?: number;
  error?: string;
}

// ===================================
// SSE EVENTS (Tiempo Real)
// ===================================

export type InvoiceSSEEventType =
  | 'invoice.created'
  | 'invoice.failed'
  | 'invoice.cancelled'
  | 'credit_note.created'
  | 'bulk_job.progress'
  | 'bulk_job.completed';

export interface InvoiceSSEEvent {
  id: string;
  type: string;
  business_id: string;
  timestamp: string;
  data: InvoiceSSEEventData;
  metadata: Record<string, any>;
}

export interface InvoiceSSEEventData {
  invoice_id?: number;
  order_id?: string;
  invoice_number?: string;
  total_amount?: number;
  currency?: string;
  status?: string;
  customer_name?: string;
  error_message?: string;
  external_url?: string;
  credit_note_id?: number;
  credit_note_number?: string;
  amount?: number;
  reason?: string;
  job_id?: string;
  total_orders?: number;
  processed?: number;
  successful?: number;
  failed?: number;
  progress?: number;
}
