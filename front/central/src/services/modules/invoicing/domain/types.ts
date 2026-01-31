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
  invoicing_provider_id: number;
  enabled: boolean;
  auto_invoice: boolean;
  filters?: Record<string, any>;
  invoice_config?: Record<string, any>;
  description?: string;
  created_at: string;
  updated_at: string;
  integration_name?: string;
  provider_name?: string;
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
  invoicing_provider_id: number;
  enabled: boolean;
  auto_invoice: boolean;
  filters?: Record<string, any>;
  invoice_config?: Record<string, any>;
  description?: string;
}

export interface UpdateConfigDTO {
  enabled?: boolean;
  auto_invoice?: boolean;
  filters?: Record<string, any>;
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
