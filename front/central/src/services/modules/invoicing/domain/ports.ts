import type {
  Invoice,
  InvoiceFilters,
  PaginatedInvoices,
  CreateInvoiceDTO,
  CreateCreditNoteDTO,
  CreditNote,
  InvoicingProvider,
  ProviderFilters,
  PaginatedProviders,
  CreateProviderDTO,
  UpdateProviderDTO,
  TestProviderResult,
  InvoicingConfig,
  ConfigFilters,
  PaginatedConfigs,
  CreateConfigDTO,
  UpdateConfigDTO,
} from './types';

/**
 * Repository interface for invoicing operations
 * Following hexagonal architecture principles
 */
export interface IInvoicingRepository {
  // ============================================
  // INVOICES
  // ============================================
  getInvoices(filters: InvoiceFilters): Promise<PaginatedInvoices>;
  getInvoiceById(id: number): Promise<Invoice>;
  createInvoice(data: CreateInvoiceDTO): Promise<Invoice>;
  cancelInvoice(id: number): Promise<Invoice>;
  retryInvoice(id: number): Promise<Invoice>;
  createCreditNote(data: CreateCreditNoteDTO): Promise<CreditNote>;

  // ============================================
  // PROVIDERS
  // ============================================
  getProviders(filters: ProviderFilters): Promise<PaginatedProviders>;
  getProviderById(id: number): Promise<InvoicingProvider>;
  createProvider(data: CreateProviderDTO): Promise<InvoicingProvider>;
  updateProvider(id: number, data: UpdateProviderDTO): Promise<InvoicingProvider>;
  testProvider(id: number): Promise<TestProviderResult>;

  // ============================================
  // CONFIGS
  // ============================================
  getConfigs(filters: ConfigFilters): Promise<PaginatedConfigs>;
  getConfigById(id: number): Promise<InvoicingConfig>;
  createConfig(data: CreateConfigDTO): Promise<InvoicingConfig>;
  updateConfig(id: number, data: UpdateConfigDTO): Promise<InvoicingConfig>;
  deleteConfig(id: number): Promise<void>;
}
